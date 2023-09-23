package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const dbTimeout = time.Second * 3

var db *sql.DB

type Job struct {
	Id      string     `json:"id"`
	Payload JobPayload `json:"payload"`
}

type JobPayload struct {
	Service string          `json:"service"`
	Action  string          `json:"action"`
	Data    json.RawMessage `json:"data"`
}

func SetConnection(dbPool *sql.DB) {
	db = dbPool
}

func GetUnhandledJobs() ([]*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, payload
		FROM jobs 
		WHERE reserved_at IS NULL
		ORDER BY created_at DESC`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*Job

	for rows.Next() {
		var job Job
		var payloadData []byte

		err := rows.Scan(
			&job.Id,
			&payloadData,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(payloadData, &job.Payload)
		if err != nil {
			log.Println("Error unmarshalling payload data", err)
			return nil, err
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (j *Job) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `DELETE FROM jobs WHERE id = $1`

	_, err := db.ExecContext(ctx, stmt, j.Id)
	if err != nil {
		return err
	}

	return nil
}

func (j *Job) Insert() (*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `INSERT INTO jobs (payload, created_at, updated_at)
		VALUES ($1, $1, $1) RETURNING id`

	payloadData, err := json.Marshal(j.Payload)
	if err != nil {
		return j, err
	}

	err = db.QueryRowContext(ctx, query,
		payloadData,
		time.Now(),
		time.Now(),
	).Scan(&j.Id)

	if err != nil {
		return j, err
	}

	return j, nil
}
