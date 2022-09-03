package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func handlers() *mux.Router {
	mux := mux.NewRouter()

	staticContentFS := http.FileServer(http.FS(staticContent))
	mux.PathPrefix("/static/").Handler(staticContentFS)

	faviconFS := http.FileServer(http.FS(favicon))
	mux.Path("/favicon.ico").Handler(faviconFS)

	mux.HandleFunc("/", mainPageHandler)
	mux.HandleFunc("/question/{id}", questionPageHandler)

	mux.HandleFunc("/question/unreviewed/{id}", unreviewedQuestionPageHandler)

	return mux
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getRandomId(0)
	if err != nil {
		log.Printf("error fetching next questionId: %s", err)
	}

	err = templates[HOME_PAGE].Execute(w, NextQuestion{id})
	if err != nil {
		log.Println(err)
	}
}

func questionPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	questionId, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.Write([]byte("Error " + err.Error()))
		return
	}

	question, err := getReviewedQuestionById(questionId)
	if err != nil {
		w.Write([]byte("Error " + err.Error()))
		return
	}

	nextQuestionId, err := getRandomId(questionId)
	if err != nil {
		log.Printf("error fetching next questionId: %s", err)
	}

	problemStatementHtml := ParseString(unescapeWhiteCharacters(question.ProblemStatement))

	err = templates[QUESTION_PAGE].Execute(w, QuestionResponse{
		NextQuestion:     NextQuestion{nextQuestionId},
		Id:               question.Id,
		ProblemStatement: template.HTML(problemStatementHtml),
		Explanation:      template.HTML(question.Explanation),
		Toughness:        question.Toughness,
		Type:             question.Type,
		Answers:          mapAnswers(question.Answers),
	})
	if err != nil {
		log.Println(err)
	}
}

func unreviewedQuestionPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	questionId, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.Write([]byte("Error " + err.Error()))
		return
	}

	question, err := getAnyQuestionById(questionId)
	if err != nil {
		w.Write([]byte("Error " + err.Error()))
		return
	}

	nextQuestionId, err := getRandomId(questionId)
	if err != nil {
		log.Printf("error fetching next questionId: %s", err)
	}

	problemStatementHtml := ParseString(unescapeWhiteCharacters(question.ProblemStatement))

	err = templates[QUESTION_PAGE].Execute(w, QuestionResponse{
		NextQuestion:     NextQuestion{nextQuestionId},
		Id:               question.Id,
		ProblemStatement: template.HTML(problemStatementHtml),
		Explanation:      template.HTML(question.Explanation),
		Toughness:        question.Toughness,
		Type:             question.Type,
		Answers:          mapAnswers(question.Answers),
	})
	if err != nil {
		log.Println(err)
	}
}

func mapAnswers(ans []Answer) []AnswerResponse {
	ansR := make([]AnswerResponse, len(ans))
	for i, v := range ans {
		ansR[i] = AnswerResponse{
			Id:           v.Id,
			Statement:    template.HTML(ParseString(v.Statement)),
			ShortComment: template.HTML(ParseString(v.ShortComment)),
			Correct:      v.Correct,
		}
	}

	return ansR
}

func unescapeWhiteCharacters(value string) string {
	value = strings.ReplaceAll(value, `\n`, "\n")
	value = strings.ReplaceAll(value, `\t`, "\t")

	return value
}
