package data

import (
	"database/sql"
	"time"

	"github.com/chitchat-awsome/internal/utils"
)

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

func (usr *User) CreatedAtDate() string {
	return usr.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (d *dataHandler) CreateUser(user *User) (User, error) {
	userRet := User{}
	err := psql.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&userRet.Id, &userRet.Uuid, &userRet.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		},
		"INSERT INTO users (uuid, name, email, password, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, uuid, created_at",
		utils.CreateUUID(), user.Name, user.Email, utils.Encrypt(user.Password), time.Now(),
	)

	return userRet, err
}

func (d *dataHandler) GetUserByEmail(email string) (User, error) {
	user := User{}
	err := psql.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		},
		"SELECT id, uuid, name, email, password, created_at FROM users WHERE email = $1",
		email,
	)

	return user, err
}

func (d *dataHandler) GetUserByUUID(uuid string) (User, error) {
	user := User{}
	err := psql.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, user.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		},
		"SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1",
		uuid,
	)

	return user, err
}

func (d *dataHandler) GetUserBySession(ss *Session) (User, error) {
	user := User{}
	err := psql.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		},
		"SELECT id, uuid, name, email, created_at FROM users WHERE id = $1",
		ss.UserId,
	)

	return user, err
}
