package factory

import (
	"context"
	"wedding-app/assembler"
	"wedding-app/domain/store"
)

func (f *Factory) Assembler(ctx context.Context) (*assembler.Assembler, error) {
	f.assemblerOnce.Do(func() {
		var quizStore store.QuizStore
		quizStore, f.assemblerError = f.QuizStore(ctx)
		if f.assemblerError != nil {
			return
		}
		var userStore store.UserStore
		userStore, f.assemblerError = f.UserStore(ctx)
		if f.assemblerError != nil {
			return
		}
		var questionStore store.QuestionStore
		questionStore, f.assemblerError = f.QuestionStore(ctx)
		if f.assemblerError != nil {
			return
		}
		var answerStore store.AnswerStore
		answerStore, f.assemblerError = f.AnswerStore(ctx)
		if f.assemblerError != nil {
			return
		}

		f.assembler = assembler.NewAssembler(quizStore, userStore, questionStore, answerStore)
	})
	return f.assembler, f.assemblerError
}
