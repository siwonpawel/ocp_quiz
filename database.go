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

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbConn *pgxpool.Pool = createDBPool()

func createDBPool() *pgxpool.Pool {
	connConfig, err := pgxpool.ParseConfig(os.Getenv("DB_CONN"))
	if err != nil {
		log.Panic("Error parsing configuration", err)
	}

	ctx, c := context.WithTimeout(context.Background(), time.Duration(500*time.Millisecond))

	pool, err := pgxpool.ConnectConfig(ctx, connConfig)
	c()
	if err != nil {
		log.Panic("Cannot establish connection to database: ", err)
	}

	ctx2, c := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	if err := pool.Ping(ctx2); err != nil {
		log.Panic(fmt.Errorf("cant ping database: %s", err))
	}
	c()

	return pool
}

func getRandomId(ctx context.Context, except int) (int, error) {
	timeoutContext, c := context.WithTimeout(ctx, time.Duration(500*time.Millisecond))

	vv := dbConn

	result, err := vv.Query(timeoutContext, `SELECT id FROM question WHERE reviewed IS NOT NULL AND id != $1 ORDER BY random() LIMIT 1`, except)
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

func getReviewedQuestionById(ctx context.Context, questionId int) (*Question, error) {
	return queryQuestion(ctx, questionId, `SELECT id, explanation, problem_statement, toughness, type FROM question WHERE id = $1 AND reviewed IS NOT NULL`, questionId)
}

func getAnyQuestionById(ctx context.Context, questionId int) (*Question, error) {
	return queryQuestion(ctx, questionId, `SELECT id, explanation, problem_statement, toughness, type FROM question WHERE id = $1`, questionId)
}

func queryQuestion(ctx context.Context, questionId int, query string, data ...any) (*Question, error) {
	timeoutContext, c := context.WithTimeout(ctx, time.Duration(500*time.Millisecond))

	result, err := dbConn.Query(timeoutContext, query, data...)
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

		ans, err := getAnswers(ctx, questionId)
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

func getAnswers(ctx context.Context, questionId int) ([]Answer, error) {
	timeoutContext, c := context.WithTimeout(ctx, time.Duration(500*time.Millisecond))

	result, err := dbConn.Query(timeoutContext, "SELECT id, correct, short_comment, statement FROM answer WHERE question_id = $1", questionId)
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
