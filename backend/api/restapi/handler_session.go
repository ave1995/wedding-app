package restapi

import (
	"errors"
	"net/http"

	"wedding-app/domain/service"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionService service.SessionService
}

func NewSessionHandler(ss service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: ss}
}

func (h *SessionHandler) startSession(c *gin.Context) {
	quizID := c.Param("quiz_id")

	userID, err := GetUserIDForQuizFromContext(c, quizID)
	if err != nil {
		if errors.Is(err, ErrUserIsNotAuthorizedForQuizInContext) {
			c.Error(NewAPIError(http.StatusUnauthorized, "unauthorized for quiz!", err))
			return
		}
		c.Error(NewInternalAPIError(err))
		return
	}

	// RADEK: Příjde mi tohle trochu jako business logika. Představi si, že budeš chtít implementovat
	// grpc server, tak tam budeš znovu dělat kontrolu na session is completed, a další načítání
	// lepší by bylo, kdyby ti to naservírovala celý service s parametry (bool, currentQuestion, error)
	session, err := h.sessionService.StartSession(c, userID.String(), quizID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	if session.IsCompleted {
		// RADEK: Já bych zase radši viděl, že ta session má metodu, která se bude umět převést do JSONU, který vrátíš pohodlně.
		// Návrh máš v modelu, nebo dost možná stačí předat objekt, který se umí převíst do JSONu, ale nejsem si jistý, nezkoušel jsem.
		// c.JSON(http.StatusOK, session)
		c.JSON(http.StatusOK, gin.H{
			"session_id":   session.ID.String(),
			"is_completed": true,
		})
		return
	}

	// RADEK: Tady moc nechápu nechápu, proč skoro celá služba vrací pěkně domain model, a najednou tohle vrací nějaké dto objekt
	// což dto je vlastně domain model
	questionDto, err := h.sessionService.GetCurrentQuestion(c, session.ID.String())
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	// RADEK: tady je vidět že můžeš posalt objekt a bude to fungovat.
	c.JSON(http.StatusOK, questionDto)
}

func (h *SessionHandler) submitAnswer(c *gin.Context) {
	// RADEK: znovu tenhle handler nese logiku:
	// Kontrola parametru a requestu - to sem správně patří

	//  getUserIDForQuizFromContext - tady by si iděálně měl ještě než něco začneš dělat načíst user id z kontextu
	// protože to je API domená (cooki) a až potom s tím nějak dělat, ty uděláš zase XYZ akcí, a potom zjistíš, že není přihlášenej
	// tak to zabiješ
	// má to být parse REST API input ===> call service ==> return result from service.
	sessionID := c.Param("session_id")

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewAPIError(http.StatusBadRequest, "invalid request body", err))
		return
	}

	// TODO: tohle by se dalo promyslet lépe, abych nemusel načítat session pro kontrolu, ale co se dá dělat
	session, err := h.sessionService.GetSessionByID(c, sessionID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	// TODO: když je uživatel za hosta a přihlásí se za admina, tak stejně může vyplnit kvíz. Zjistit proč
	_, err = GetUserIDForQuizFromContext(c, session.QuizID.String())
	if err != nil {
		if errors.Is(err, ErrUserIsNotAuthorizedForQuizInContext) {
			c.Error(NewAPIError(http.StatusUnauthorized, "unauthorized for session!", err))
			return
		}
		c.Error(NewInternalAPIError(err))
		return
	}

	// RADEK: tohle je si myslím špatný návrh REST API
	// REST API by mělo být bezstavové a vracet vždycky stejné výsledky, ale ty jakmile je ten dotazník vyplněný tak vrátíš
	// úplně jinou JSON strukturu
	// Správný design můžeme probrat potom na callu.
	isCompleted, err := h.sessionService.SubmitAnswer(c, sessionID, req.QuestionID, req.AnswerIDs)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	if isCompleted {
		result, err := h.sessionService.GetResult(c, sessionID)
		if err != nil {
			c.Error(NewInternalAPIError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"completed": true,
			"result": gin.H{
				"score":      result.Score,
				"total":      result.Total,
				"percentage": result.Percentage,
			},
		})
		return
	}

	nextQuestionDto, err := h.sessionService.GetCurrentQuestion(c, sessionID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, nextQuestionDto)
}

func (h *SessionHandler) getCurrentQuestion(c *gin.Context) {
	sessionID := c.Param("session_id")

	// načtení session a kontrola uživatele
	session, err := h.sessionService.GetSessionByID(c, sessionID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}
	if session.IsCompleted {
		c.JSON(http.StatusOK, gin.H{
			"completed": true,
			"question":  nil,
		})
		return
	}

	_, err = GetUserIDForQuizFromContext(c, session.QuizID.String())
	if err != nil {
		if errors.Is(err, ErrUserIsNotAuthorizedForQuizInContext) {
			c.Error(NewAPIError(http.StatusUnauthorized, "unauthorized for session!", err))
			return
		}
		c.Error(NewInternalAPIError(err))
		return
	}

	questionDto, err := h.sessionService.GetCurrentQuestion(c, sessionID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, questionDto)
}

func (h *SessionHandler) getResult(c *gin.Context) {
	sessionID := c.Param("session_id")

	session, err := h.sessionService.GetSessionByID(c, sessionID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	_, err = GetUserIDForQuizFromContext(c, session.QuizID.String())
	if err != nil {
		if errors.Is(err, ErrUserIsNotAuthorizedForQuizInContext) {
			c.Error(NewAPIError(http.StatusUnauthorized, "unauthorized for session!", err))
			return
		}
		c.Error(NewInternalAPIError(err))
		return
	}

	result, err := h.sessionService.GetResult(c, sessionID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"completed": true,
		"result": gin.H{
			"score":      result.Score,
			"total":      result.Total,
			"percentage": result.Percentage,
		},
	})
}

func (h *SessionHandler) getActiveSessionsByQuizID(c *gin.Context) {
	quizID := c.Param("id")

	active, err := h.sessionService.GetActiveSessionsByQuizID(c, quizID)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": len(active)})
}

func (h *SessionHandler) getResults(c *gin.Context) {
	id := c.Param("id")

	results, err := h.sessionService.GetResultsByQuizID(c, id)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, results)
}
