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

	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/", mainPageHandler)
	mux.HandleFunc("/question/{id}", questionPageHandler)

	mux.HandleFunc("/question/unreviewed/{id}", unreviewedQuestionPageHandler)

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := getRandomId(r.Context(), 0)
	if err != nil {
		log.Println("Health check: FAIL")
		w.WriteHeader(503)
		w.Write([]byte("FAIL"))
	} else {
		log.Println("Health check: OK")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Requested main page")
	id, err := getRandomId(r.Context(), 0)
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
	log.Printf("Requested question [ id = %d ].\n", questionId)

	question, err := getReviewedQuestionById(r.Context(), questionId)
	if err != nil {
		w.Write([]byte("Error " + err.Error()))
		return
	}

	nextQuestionId, err := getRandomId(r.Context(), questionId)
	if err != nil {
		log.Printf("error fetching next questionId: %s", err)
	}

	problemStatementHtml := ParseString(unescapeWhiteCharacters(question.ProblemStatement))
	explanationHtml := ParseString(unescapeWhiteCharacters(question.Explanation))

	err = templates[QUESTION_PAGE].Execute(w, QuestionResponse{
		NextQuestion:     NextQuestion{nextQuestionId},
		Id:               question.Id,
		ProblemStatement: template.HTML(problemStatementHtml),
		Explanation:      template.HTML(explanationHtml),
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

	question, err := getAnyQuestionById(r.Context(), questionId)
	if err != nil {
		w.Write([]byte("Error " + err.Error()))
		return
	}

	nextQuestionId, err := getRandomId(r.Context(), questionId)
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
