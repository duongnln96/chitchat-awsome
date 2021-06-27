package data

import (
	"database/sql"
	"time"
)

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

func (ss *Session) CreatedAtDate() string {
	return ss.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (d *dataHandler) CreateSession(user *User) (Session, error) {
	ssRet := Session{}
	err := psql.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&ssRet.Id, &ssRet.Uuid, &ssRet.Email, &ssRet.UserId, &ssRet.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		},
		"INSERT INTO sessions (uuid, email, user_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id, uuid, email, user_id, created_at",
		user.Uuid, user.Email, user.Id, time.Now(),
	)

	return ssRet, err
}

func (d *dataHandler) GetSessionByUser(user *User) (Session, error) {
	ss := Session{}
	err := psql.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&ss.Id, &ss.Uuid)
			if err != nil {
				return err
			}
			return nil
		},
		"SELECT id, uuid FROM sessions WHERE user_id = $1",
		user.Id,
	)

	return ss, err
}

func (d *dataHandler) GetSessionByUUID(uuid string) (Session, error) {
	ss := Session{}
	err := psql.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&ss.Id, &ss.Uuid, &ss.Email, &ss.UserId, &ss.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		},
		"SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1",
		uuid,
	)
	return ss, err
}

func (d *dataHandler) DeleteSessionByUUID(ss *Session) error {
	err := psql.DBExec(
		"DELETE FROM sessions WHERE uuid=$1",
		ss.Uuid,
	)
	return err
}
