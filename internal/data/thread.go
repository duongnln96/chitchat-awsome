package data

import (
	"database/sql"
	"time"

	"github.com/chitchat-awsome/internal/utils"
)

func (d *dataHandler) CreateThread(topic string, user *User) (Thread, error) {
	thread := Thread{}
	err := d.db.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&thread.Id, &thread.Uuid, &thread.Topic, &thread.UserId, &thread.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		},
		"INSERT INTO threads (uuid, topic, user_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id, uuid, topic, user_id, created_at",
		utils.CreateUUID(), topic, user.Id, time.Now(),
	)

	return thread, err
}

func (d *dataHandler) GetAllThreads() ([]Thread, error) {
	threads := make([]Thread, 0)
	err := d.db.DBQueryRows(
		func(r *sql.Rows) bool {
			thread := Thread{}
			err := r.Scan(&thread.Id, &thread.Uuid, &thread.Topic, &thread.UserId, &thread.CreatedAt)
			if err != nil {
				d.log.Errorf("Error while reading row: %s", err)
				return false
			}
			threads = append(threads, thread)
			return true
		},
		"SELECT id, uuid, topic, user_id, created_at FROM threads ORDER BY created_at DESC",
	)

	return threads, err
}

func (d *dataHandler) DeleteThread(thread Thread) error {
	err := d.db.DBExec(
		"DELETE FROM threads WHERE id=$1",
		thread.Id,
	)
	return err
}
