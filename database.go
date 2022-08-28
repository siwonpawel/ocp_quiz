package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
)

var dbConn *pgx.Conn = connect()

func connect() *pgx.Conn {
	ctx, c := context.WithTimeout(context.Background(), time.Duration(500*time.Millisecond))

	conn, err := pgx.Connect(ctx, os.Getenv("DB_CONN"))
	c()
	if err != nil {
		log.Panic(err)
	}

	ctx2, c := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer c()
	if err := conn.Ping(ctx2); err != nil {
		log.Panic(fmt.Errorf("cant ping database: %s", err))
	}

	return conn
}

func getRandomId(except int) (int, error) {
	ctx, c := context.WithTimeout(context.Background(), time.Duration(500*time.Millisecond))

	result, err := dbConn.Query(ctx, `SELECT id FROM question WHERE reviewed IS NOT NULL AND id != $1 ORDER BY random() LIMIT 1`, except)
	c()
	if err != nil {
		return 0, err
	}
	defer result.Close()

	if result.Next() {
		var cnt int
		err = result.Scan(&cnt)
		if err != nil {
			return 0, err
		}

		return cnt, nil
	}

	return 0, errors.New("can't fetch number of questions")
}

func getReviewedQuestionById(questionId int) (*Question, error) {
	return queryQuestion(questionId, `SELECT id, explanation, problem_statement, toughness, type FROM question WHERE id = $1 AND reviewed IS NOT NULL`, questionId)
}

func getAnyQuestionById(questionId int) (*Question, error) {
	return queryQuestion(questionId, `SELECT id, explanation, problem_statement, toughness, type FROM question WHERE id = $1`, questionId)
}

func queryQuestion(questionId int, query string, data ...any) (*Question, error) {
	ctx, c := context.WithTimeout(context.Background(), time.Duration(500*time.Millisecond))

	result, err := dbConn.Query(ctx, query, data...)
	c()
	if err != nil {
		return &Question{}, err
	}

	if result.Next() {
		var id int
		var explanation string
		var problemStatement string
		var toughness string
		var questionType string

		err := result.Scan(&id, &explanation, &problemStatement, &toughness, &questionType)
		result.Close()
		if err != nil {
			return nil, err
		}

		ans, err := getAnswers(questionId)
		if err != nil {
			return &Question{}, err
		}

		return &Question{
			Id:               id,
			ProblemStatement: problemStatement,
			Explanation:      explanation,
			Toughness:        toughness,
			Type:             questionType,
			Answers:          ans,
		}, nil
	} else {
		return &Question{}, fmt.Errorf("can't fetch question with id=%d", questionId)
	}
}

func getAnswers(questionId int) ([]Answer, error) {
	ctx, c := context.WithTimeout(context.Background(), time.Duration(500*time.Millisecond))

	result, err := dbConn.Query(ctx, "SELECT id, correct, short_comment, statement FROM answer WHERE question_id = $1", questionId)
	c()

	if err != nil {
		return []Answer{}, err
	}

	ans := []Answer{}
	for result.Next() {
		var id int
		var correct string
		var shortComment string
		var statement string

		if err := result.Scan(&id, &correct, &shortComment, &statement); err != nil {
			log.Printf("Error parsing Answer object: %s", err)
			return []Answer{}, err
		}

		c, _ := strconv.ParseBool(correct)

		ans = append(ans, Answer{Id: id, Statement: statement, ShortComment: shortComment, Correct: c})
	}

	rand.Shuffle(len(ans), func(i, j int) {
		ans[i], ans[j] = ans[j], ans[i]
	})

	return ans, nil
}
