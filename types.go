package main

import "html/template"

type NextQuestion struct {
	NextQuestionId int
}

type Question struct {
	Id               int      `json:"id"`
	ProblemStatement string   `json:"problemStatement"`
	Explanation      string   `json:"explanation"`
	Toughness        string   `json:"toughness"`
	Type             string   `json:"type"`
	Answers          []Answer `json:"answers"`
}

type Answer struct {
	Id           int    `json:"id"`
	Statement    string `json:"statement"`
	ShortComment string `json:"shortComment"`
	Correct      bool   `json:"correct"`
}

type QuestionResponse struct {
	NextQuestion
	Id               int
	ProblemStatement template.HTML
	Explanation      template.HTML
	Toughness        string
	Type             string
	Answers          []AnswerResponse
}

type AnswerResponse struct {
	Id           int           `json:"id"`
	Statement    template.HTML `json:"statement"`
	ShortComment template.HTML `json:"shortComment"`
	Correct      bool          `json:"correct"`
}
