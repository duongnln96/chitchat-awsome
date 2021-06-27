package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/chitchat-awsome/internal/utils"
)

type Thread struct {
	Id        int
	Uuid      string
	Topic     string
	UserId    int
	CreatedAt time.Time
}

func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (thread *Thread) User() User {
	user := User{}
	psql.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		},
		"SELECT id, uuid, name, email, created_at FROM users WHERE id = $1",
		thread.UserId,
	)
	return user
}

func (thread *Thread) Posts() []Post {
	posts := make([]Post, 0)
	psql.DBQueryRows(
		func(rows *sql.Rows) bool {
			post := Post{}
			err := rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
			if err != nil {
				return false
			}
			posts = append(posts, post)
			return true
		},
		"SELECT id, uuid, body, user_id, thread_id, created_at FROM posts WHERE thread_id = $1",
		thread.Id,
	)

	fmt.Printf("posts: %+v", posts)
	return posts
}

func (thread *Thread) NumReplies() int {
	var count int = 0
	psql.DBQueryRows(
		func(r *sql.Rows) bool {
			err := r.Scan(&count)
			if err != nil {
				return false
			}
			return true
		},
		"SELECT count(*) FROM posts WHERE thread_id = $1",
		thread.Id,
	)
	return count
}

func (d *dataHandler) CreateThread(topic string, user *User) (Thread, error) {
	thread := Thread{}
	err := psql.DBQueryRow(
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
	err := psql.DBQueryRows(
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

func (d *dataHandler) GetThreadByUUID(uuid string) (Thread, error) {
	thread := Thread{}
	err := psql.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&thread.Id, &thread.Uuid, &thread.Topic, &thread.UserId, &thread.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		},
		"SELECT id, uuid, topic, user_id, created_at FROM threads WHERE uuid = $1",
		uuid,
	)

	return thread, err
}

func (d *dataHandler) DeleteThread(thread Thread) error {
	err := psql.DBExec(
		"DELETE FROM threads WHERE id=$1",
		thread.Id,
	)
	return err
}
