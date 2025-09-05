package utils

import (
	"wedding-app/domain/model"
)

func IsQuestionCorrect(q *model.Question, userAnswers []*model.Attempt, correctAnswers []*model.Answer) (bool, error) {
	switch q.Type {
	case model.SingleChoice:
		return len(userAnswers) == 1 && userAnswers[0].IsCorrect, nil
	case model.MultipleChoice:
		if len(userAnswers) != len(correctAnswers) {
			return false, nil // musí zvolit správný počet odpovědí, aby někdo nemohl zaklikávat všechno a mít plný počet bodů
		}

		for _, ua := range userAnswers {
			if !ua.IsCorrect {
				return false, nil
			}
		}
		return true, nil
	default:
		return false, nil
	}
}
