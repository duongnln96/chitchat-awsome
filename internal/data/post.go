package data

import "database/sql"

func (dh *dataHandler) CreatePost(thread Thread, post Post) error {

	return nil
}

func (th *dataHandler) GetPostsIntoThread(thread Thread) ([]Post, error) {
	posts := make([]Post, 0)
	err := th.db.DBQueryRows(
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
