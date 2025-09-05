package session_test

import (
	"fmt"
	"testing"
	"wedding-app/domain/apperrors"
	"wedding-app/domain/model"
	"wedding-app/mocks/wedding-app/domain/storemock"
	"wedding-app/service/session"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCurrentQuestion(t *testing.T) {
	type GetCurrentQuestionTest struct {
		Name           string
		ArgSessionID   string
		ExpectedError  error
		ExpectedResult *model.Question
		Mock           func(*storemock.MockSessionStore, *storemock.MockQuestionStore)
	}

	sessionID := uuid.MustParse("9C33B04D-084A-41EF-8200-4C1D06DD1CD9")
	quizID := uuid.MustParse("68CEA6F4-87FF-44A3-AF11-EC1BD6D92F1D")
	questionID1 := uuid.MustParse("4CF9D44F-2BDC-4BA1-9EA5-F6CFEF9BFF64")

	expectedQuestion := &model.Question{ID: questionID1, QuizID: quizID, Text: "Je mi 29?"}

	tests := []GetCurrentQuestionTest{
		{
			Name:           "should_return_an_question",
			ArgSessionID:   sessionID.String(),
			ExpectedError:  nil,
			ExpectedResult: expectedQuestion,
			Mock: func(ss *storemock.MockSessionStore, qs *storemock.MockQuestionStore) {
				expectedSession := &model.Session{ID: sessionID, QuizID: quizID}
				ss.EXPECT().FindByID(mock.Anything, sessionID).Return(expectedSession, nil)
				expectedQuestions := []*model.Question{expectedQuestion}
				qs.EXPECT().GetOrderedQuestionsByQuizID(mock.Anything, quizID).Return(expectedQuestions, nil)
			},
		},
		{
			Name:         "should_return_error_on_invalid_session_ID",
			ArgSessionID: "invalid-uuid",
			ExpectedError: func() error {
				_, parseErr := uuid.Parse("invalid-uuid")
				return fmt.Errorf("failed to parse session ID %q: %w", "invalid-uuid", parseErr)
			}(),
			ExpectedResult: nil,
			Mock:           func(ss *storemock.MockSessionStore, qs *storemock.MockQuestionStore) {},
		},
		{
			Name:           "should_return_error_if_session_not_found",
			ArgSessionID:   sessionID.String(),
			ExpectedError:  fmt.Errorf("failed to load session: %w", apperrors.ErrNotFound),
			ExpectedResult: nil,
			Mock: func(ss *storemock.MockSessionStore, qs *storemock.MockQuestionStore) {
				ss.EXPECT().FindByID(mock.Anything, sessionID).Return(nil, apperrors.ErrNotFound)
			},
		},
		{
			Name:           "should_return_no_questions_found",
			ArgSessionID:   sessionID.String(),
			ExpectedError:  fmt.Errorf("no questions found for quiz ID %s", quizID),
			ExpectedResult: nil,
			Mock: func(ss *storemock.MockSessionStore, qs *storemock.MockQuestionStore) {
				expectedSession := &model.Session{ID: sessionID, QuizID: quizID}
				ss.EXPECT().FindByID(mock.Anything, sessionID).Return(expectedSession, nil)
				expectedQuestions := []*model.Question{}
				qs.EXPECT().GetOrderedQuestionsByQuizID(mock.Anything, quizID).Return(expectedQuestions, nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ss := storemock.NewMockSessionStore(t)
			qs := storemock.NewMockQuestionStore(t)

			test.Mock(ss, qs)
			// TODO: opravit testy
			svc := session.NewSessionService(ss, qs, nil, nil, nil, nil, nil, nil)
			actual, err := svc.GetCurrentQuestion(t.Context(), test.ArgSessionID)
			if test.ExpectedResult != nil {
				assert.Equal(t, test.ExpectedResult, actual)
			} else {
				assert.Nil(t, actual)
			}
			assert.Equal(t, test.ExpectedError, err)
		})
	}
}
