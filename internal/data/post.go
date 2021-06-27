package data

import (
	"database/sql"
	"time"

	"github.com/chitchat-awsome/internal/utils"
)

type Post struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	ThreadId  int
	CreatedAt time.Time
}

func (post *Post) CreatedAtDate() string {
	return post.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (post *Post) User() User {
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
		post.UserId,
	)
	return user
}

func (d *dataHandler) CreatePost(user *User, thread *Thread, body string) (Post, error) {
	post := Post{}
	err := psql.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		},
		"insert into posts (uuid, body, user_id, thread_id, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, body, user_id, thread_id, created_at",
		utils.CreateUUID(), body, user.Id, thread.Id, time.Now(),
	)
	return post, err
}

func (th *dataHandler) GetPostsIntoThread(thread Thread) ([]Post, error) {
	posts := make([]Post, 0)
	err := psql.DBQueryRows(
		func(rows *sql.Rows) bool {
			post := Post{}
			err := rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
			if err != nil {
				th.log.Errorf("ERROR While reading Post %+v", err)
				return false
			}
			posts = append(posts, post)
			return true
		},
		"SELECT id, uuid, body, user_id, thread_id, created_at FROM posts WHERE thread_id = $1",
		thread.Id,
	)

	return posts, err
}
