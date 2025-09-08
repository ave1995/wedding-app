package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"wedding-app/config"
	"wedding-app/domain/model"
	"wedding-app/factory"
)

type QuizDefinition struct {
	Name      string
	Questions []model.Question
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := slog.Default()
	logger.Info("starting quiz setup")

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("error loading config", "error", err)
		os.Exit(1)
	}

	fc := factory.NewFactory(cfg)

	qs, err := fc.QuizService(ctx)
	if err != nil {
		logger.Error("error creating quiz service", "error", err)
		os.Exit(1)
	}

	quests, err := fc.QuestionService(ctx)
	if err != nil {
		logger.Error("error creating questions service", "error", err)
	}

	answers, err := fc.AnswerService(ctx)
	if err != nil {
		logger.Error("error creating answers service", "error", err)
	}

	definition := getTestDefinition()
	//definition := finalQuizDefinition()

	quiz, err := qs.CreateQuiz(ctx, definition.Name)
	if err != nil {
		logger.Error("error creating quiz", "error", err)
		os.Exit(1)
	}

	for _, question := range definition.Questions {
		q, err := quests.CreateQuestion(ctx, question.Text, quiz.ID.String(), question.Type, question.PhotoPath)
		if err != nil {
			logger.Error("error creating question", "error", err)
			os.Exit(1)
		}

		for _, answer := range question.Answers {
			_, err := answers.CreateAnswer(ctx, answer.Text, q.ID.String(), answer.IsCorrect)
			if err != nil {
				logger.Error("error creating answer", "error", err)
				os.Exit(1)
			}
		}
	}

	//users, err := fc.UserService(ctx)
	//if err != nil {
	//	logger.Error("error creating users service", "error", err)
	//	os.Exit(1)
	//}
	//_, err = users.RegisterUser(ctx, model.RegisterUserParams{
	//	Username: "admin",
	//	Email:    "admin@test.cz",
	//	Password: "secret",
	//	IconURL:  "",
	//})
	//if err != nil {
	//	logger.Error("error registering user", "error", err)
	//	os.Exit(1)
	//}

	logger.Info("quiz setup finished")
	os.Exit(0)
}

func getTestDefinition() QuizDefinition {
	return QuizDefinition{
		Name: "quiz_setup_test_2",
		Questions: []model.Question{
			{
				Text:      "Otázka 1",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("quiz_setup_test_2/q_1.png"),
				Answers: []*model.Answer{
					{Text: "Odpověď 1", IsCorrect: false},
					{Text: "Odpověď 2", IsCorrect: false},
					{Text: "Odpověď 3", IsCorrect: true},
					{Text: "Odpověď 4", IsCorrect: false},
				},
			},
			{
				Text:      "Otázka 2",
				Type:      model.MultipleChoice,
				PhotoPath: Ptr("quiz_setup_test_2/q_2.png"),
				Answers: []*model.Answer{
					{Text: "Odpověď 1", IsCorrect: false},
					{Text: "Odpověď 2", IsCorrect: true},
					{Text: "Odpověď 3", IsCorrect: true},
					{Text: "Odpověď 4", IsCorrect: false},
				},
			},
			{
				Text:      "Otázka 3",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("quiz_setup_test_2/q_3.png"),
				Answers: []*model.Answer{
					{Text: "Odpověď 1", IsCorrect: true},
					{Text: "Odpověď 2", IsCorrect: false},
					{Text: "Odpověď 3", IsCorrect: false},
					{Text: "Odpověď 4", IsCorrect: false},
				},
			},
			{
				Text:      "Otázka 4",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("quiz_setup_test_2/q_4.png"),
				Answers: []*model.Answer{
					{Text: "Odpověď 1", IsCorrect: false},
					{Text: "Odpověď 2", IsCorrect: false},
					{Text: "Odpověď 3", IsCorrect: false},
					{Text: "Odpověď 4", IsCorrect: true},
				},
			},
		},
	}
}

func finalQuizDefinition() QuizDefinition {
	return QuizDefinition{
		Name: "Svatební kvíz",
		Questions: []model.Question{
			{
				Text:      "1. Jaké je nevěstino první slovo?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_1.png"),
				Answers: []*model.Answer{
					{Text: "Ryba", IsCorrect: false},
					{Text: "Kapr", IsCorrect: true},
					{Text: "Prut", IsCorrect: false},
					{Text: "Síť", IsCorrect: false},
				},
			},
			{
				Text:      "2. Jakou měl ženich přezdívku na základní škole?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_2.png"),
				Answers: []*model.Answer{
					{Text: "Bedna", IsCorrect: false},
					{Text: "Béďa", IsCorrect: false},
					{Text: "Čoky", IsCorrect: true},
					{Text: "Koky", IsCorrect: false},
				},
			},
			{
				Text:      "3. Jaké bylo oblíbené pití nevěsty v dětství?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Coca Cola", IsCorrect: false},
					{Text: "Sprite", IsCorrect: true},
					{Text: "Kofola", IsCorrect: false},
					{Text: "Fanta", IsCorrect: false},
				},
			},
			{
				Text:      "4. Jaké vysněné povolání měl ženich jako malý?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "IT specialista", IsCorrect: false},
					{Text: "Kuchař", IsCorrect: true},
					{Text: "Hasič", IsCorrect: false},
					{Text: "Prodavač", IsCorrect: false},
				},
			},
			{
				Text:      "5. Kterou pohádku měla nevěsta nejraději?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Malá Čarodějnice", IsCorrect: true},
					{Text: "Rákosníček", IsCorrect: false},
					{Text: "Pohádky z mechu a kapradí", IsCorrect: false},
					{Text: "Včelí Medvídci", IsCorrect: false},
				},
			},
			{
				Text:      "6. Jaký byl první sport, který zkoušel ženich?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_6.png"),
				Answers: []*model.Answer{
					{Text: "Plavání", IsCorrect: false},
					{Text: "Basketbal", IsCorrect: true},
					{Text: "Florbal", IsCorrect: false},
					{Text: "Atletika", IsCorrect: false},
				},
			},
			{
				Text:      "7. S jakým sportem začala nevěsta jako první?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_7.png"),
				Answers: []*model.Answer{
					{Text: "Plavání", IsCorrect: false},
					{Text: "Volejbal", IsCorrect: false},
					{Text: "Aerobik", IsCorrect: true},
					{Text: "Basketbal", IsCorrect: false},
				},
			},
			{
				Text:      "8. Jaký byl oblíbený školní předmět nevěsty?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_8.png"),
				Answers: []*model.Answer{
					{Text: "Pracovní činnosti", IsCorrect: true},
					{Text: "Čeština", IsCorrect: false},
					{Text: "Matematika", IsCorrect: false},
					{Text: "Žádný", IsCorrect: false},
				},
			},
			{
				Text:      "9. Jaký byl oblíbený školní předmět ženicha?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_9.png"),
				Answers: []*model.Answer{
					{Text: "Výtvarná výchova", IsCorrect: true},
					{Text: "Přírodověda", IsCorrect: false},
					{Text: "Matematika", IsCorrect: false},
					{Text: "Čeština", IsCorrect: false},
				},
			},
			{
				Text:      "10. Jak se jmenuje nevěstin volejbalový tým?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_10.png"),
				Answers: []*model.Answer{
					{Text: "JTNCH", IsCorrect: false},
					{Text: "Rumburk", IsCorrect: false},
					{Text: "Vikingové", IsCorrect: true},
					{Text: "Jen lehké drogy", IsCorrect: false},
				},
			},
			{
				Text:      "11. Jaký byl první koncert, na kterém ženich byl?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_11.png"),
				Answers: []*model.Answer{
					{Text: "Refew", IsCorrect: false},
					{Text: "Prago Union", IsCorrect: false},
					{Text: "Maxim Turbulenc", IsCorrect: true},
					{Text: "Troublegang", IsCorrect: false},
				},
			},
			{
				Text:      "12. Jaký domácí mazlíček byl u nevěsty doma první?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Králík", IsCorrect: false},
					{Text: "Kočka", IsCorrect: false},
					{Text: "Křeček", IsCorrect: true},
					{Text: "Žádný", IsCorrect: false},
				},
			},
			{
				Text:      "13. Kam jezdil ženich na tábor?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_13.png"),
				Answers: []*model.Answer{
					{Text: "Tábor", IsCorrect: false},
					{Text: "Varvažov", IsCorrect: false},
					{Text: "Čeřín", IsCorrect: true},
					{Text: "Ostrovec", IsCorrect: false},
				},
			},
			{
				Text:      "14. Jaký byl první zahraniční výlet nevěsty letadlem?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_14.png"),
				Answers: []*model.Answer{
					{Text: "Chorvatsko", IsCorrect: false},
					{Text: "Turecko", IsCorrect: true},
					{Text: "Egypt", IsCorrect: false},
					{Text: "Bulharsko", IsCorrect: false},
				},
			},
			{
				Text:      "15. Kterou příchuť zmrzliny měla nevěsta nejraději?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_15.png"),
				Answers: []*model.Answer{
					{Text: "Vanilka", IsCorrect: false},
					{Text: "Šmoulová", IsCorrect: true},
					{Text: "Čokoláda", IsCorrect: false},
					{Text: "Pistácie", IsCorrect: false},
				},
			},
			{
				Text:      "16. Na jaký hudební nástroj hrál ženich?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_16.png"),
				Answers: []*model.Answer{
					{Text: "Klavír", IsCorrect: false},
					{Text: "Ozvučná dřívka", IsCorrect: false},
					{Text: "Flétna", IsCorrect: true},
					{Text: "Kytara", IsCorrect: false},
				},
			},
			{
				Text:      "17. Jaký byl první film v kině, který viděla nevěsta?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Harry Potter", IsCorrect: false},
					{Text: "Auta", IsCorrect: false},
					{Text: "Hledá se Nemo", IsCorrect: true},
					{Text: "Příběh žraloka", IsCorrect: false},
				},
			},
			{
				Text:      "18. Jaký byl ženichův oblíbený plyšák?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_18.png"),
				Answers: []*model.Answer{
					{Text: "Králík", IsCorrect: false},
					{Text: "Želva", IsCorrect: false},
					{Text: "Veverka", IsCorrect: true},
					{Text: "Medvěd", IsCorrect: false},
				},
			},
			{
				Text:      "19. Jakou pohádku četli nevěstě nejčastěji před spaním?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Birlibán", IsCorrect: true},
					{Text: "Honzíkova cesta", IsCorrect: false},
					{Text: "Děti z Bullerbynu", IsCorrect: false},
					{Text: "Luisa a Lotka", IsCorrect: false},
				},
			},
			{
				Text:      "20. Jaká byla ženichova oblíbená počítačová hra?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "GTA San Andreas", IsCorrect: false},
					{Text: "Trackmania", IsCorrect: false},
					{Text: "Battlefield 3", IsCorrect: false},
					{Text: "Mafie", IsCorrect: true},
				},
			},
			{
				Text:      "21. Kde jste se poprvé potkali?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_21.png"),
				Answers: []*model.Answer{
					{Text: "Ve škole", IsCorrect: false},
					{Text: "Na diskotéce", IsCorrect: true},
					{Text: "V hospodě", IsCorrect: false},
					{Text: "V metru", IsCorrect: false},
				},
			},
			{
				Text:      "22. Kdo oslovil toho druhého jako první?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Ženich", IsCorrect: false},
					{Text: "Nevěsta", IsCorrect: false},
					{Text: "Oba současně", IsCorrect: false},
					{Text: "Nikdo", IsCorrect: true},
				},
			},
			{
				Text:      "23. Kde jsme spolu poprvé bydleli?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_23.png"),
				Answers: []*model.Answer{
					{Text: "Strašnice", IsCorrect: false},
					{Text: "Bohdalec", IsCorrect: true},
					{Text: "Letňany", IsCorrect: false},
					{Text: "Dašice", IsCorrect: false},
				},
			},
			{
				Text:      "24. Jaký byl první společný film, který jste viděli?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "X-Men", IsCorrect: true},
					{Text: "Někdo to rád blond", IsCorrect: false},
					{Text: "Pán prstenů", IsCorrect: false},
					{Text: "Harry Potter", IsCorrect: false},
				},
			},
			{
				Text:      "25. Kam byl první společný výlet?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_25.png"),
				Answers: []*model.Answer{
					{Text: "Kunětická hora", IsCorrect: false},
					{Text: "Morava", IsCorrect: false},
					{Text: "Kutná Hora", IsCorrect: false},
					{Text: "Drážďany", IsCorrect: true},
				},
			},
			{
				Text:      "26. Kdo je větší spáč?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Ženich", IsCorrect: true},
					{Text: "Nevěsta", IsCorrect: false},
					{Text: "Oba stejně", IsCorrect: false},
					{Text: "Nikdo", IsCorrect: false},
				},
			},
			{
				Text:      "27. Kdo tráví víc času na telefonu?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Ženich", IsCorrect: true},
					{Text: "Nevěsta", IsCorrect: false},
					{Text: "Oba stejně", IsCorrect: false},
					{Text: "Nikdo", IsCorrect: false},
				},
			},
			{
				Text:      "28. Kde byla první společná dovolená v zahraničí?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_28.png"),
				Answers: []*model.Answer{
					{Text: "Chorvatsko", IsCorrect: false},
					{Text: "Rakousko", IsCorrect: true},
					{Text: "Maďarsko", IsCorrect: false},
					{Text: "Tunisko", IsCorrect: false},
				},
			},
			{
				Text:      "29. Jakou sportovní událost jsme spolu jako první sledovaly?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "MS v hokeji", IsCorrect: true},
					{Text: "Olympiáda", IsCorrect: false},
					{Text: "ME ve volejbale", IsCorrect: false},
					{Text: "MS v šipkách", IsCorrect: false},
				},
			},
			{
				Text:      "30. Kdo má větší nepořádek ve skříni?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Ženich", IsCorrect: false},
					{Text: "Nevěsta", IsCorrect: true},
					{Text: "Oba stejně", IsCorrect: false},
					{Text: "Nikdo", IsCorrect: false},
				},
			},
			{
				Text:      "31. Kdo utráci víc peněz?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Ženich", IsCorrect: false},
					{Text: "Nevěsta", IsCorrect: false},
					{Text: "Oba stejně", IsCorrect: true},
					{Text: "Raději nikod", IsCorrect: false},
				},
			},
			{
				Text:      "32. Jak dlouho spolu byli před svatbou?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "5 let", IsCorrect: false},
					{Text: "6 let", IsCorrect: true},
					{Text: "7 let", IsCorrect: false},
					{Text: "8 let", IsCorrect: false},
				},
			},
			{
				Text:      "33. Jaký je největší ženichův zlozvyk?",
				Type:      model.MultipleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Trávení času na PC i mimo práci", IsCorrect: false},
					{Text: "Přehnané uklízení", IsCorrect: true},
					{Text: "Dokupování nepotřebných věcí", IsCorrect: true},
					{Text: "Žádný", IsCorrect: false},
				},
			},
			{
				Text:      "34. Kde jsme strávili první výročí (1 rok)?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_34.png"),
				Answers: []*model.Answer{
					{Text: "Field", IsCorrect: false},
					{Text: "Brasileiro", IsCorrect: true},
					{Text: "Café Imperial", IsCorrect: false},
					{Text: "Doma", IsCorrect: false},
				},
			},
			{
				Text:      "35. Kam byl náš nejkratší let?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_35.png"),
				Answers: []*model.Answer{
					{Text: "Budapešť", IsCorrect: true},
					{Text: "Karlovy Vary", IsCorrect: false},
					{Text: "Bratislava", IsCorrect: false},
					{Text: "Vídeň", IsCorrect: false},
				},
			},
			{
				Text:      "36. Jaký je největší nevěsty zlozvyk?",
				Type:      model.MultipleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Bordelářství", IsCorrect: true},
					{Text: "Usínání na gauči", IsCorrect: true},
					{Text: "Žádné", IsCorrect: false},
					{Text: "Musí mít poslední slovo", IsCorrect: false},
				},
			},
			{
				Text:      "37. Kdo z vás dvou se víc bojí hororů?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Ženich", IsCorrect: true},
					{Text: "Nevěsta", IsCorrect: false},
					{Text: "Oba", IsCorrect: false},
					{Text: "Nikdo", IsCorrect: false},
				},
			},
			{
				Text:      "38. K jakým společným aktivitám jsme se dopracovali?",
				Type:      model.MultipleChoice,
				PhotoPath: Ptr("wedding-quiz/q_38.png"),
				Answers: []*model.Answer{
					{Text: "Bowling", IsCorrect: false},
					{Text: "Hraní na počítači", IsCorrect: false},
					{Text: "Badminton", IsCorrect: true},
					{Text: "Tenis", IsCorrect: true},
				},
			},
			{
				Text:      "39. Kdo z nás dvou více vyhrává?",
				Type:      model.SingleChoice,
				PhotoPath: nil,
				Answers: []*model.Answer{
					{Text: "Ženich", IsCorrect: false},
					{Text: "Nevěsta", IsCorrect: true},
					{Text: "Je to vyrovnané", IsCorrect: false},
					{Text: "Radši nikdo", IsCorrect: false},
				},
			},
			{
				Text:      "40. Co předcházelo zasnoubení?",
				Type:      model.SingleChoice,
				PhotoPath: Ptr("wedding-quiz/q_40.png"),
				Answers: []*model.Answer{
					{Text: "Večeře v restauraci U Matěje", IsCorrect: false},
					{Text: "Večeře v restauraci Field", IsCorrect: true},
					{Text: "Večeře v restauraci Špejle", IsCorrect: false},
					{Text: "Večeře v restauraci Entree", IsCorrect: false},
				},
			},
		},
	}
}

func Ptr[T any](input T) *T {
	return &input
}
