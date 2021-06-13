package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/chitchat-awsome/internal/utils"
	"github.com/chitchat-awsome/pkg/psqlconnector"
	"go.uber.org/zap"
)

type Thread struct {
	Id        int
	Uuid      string
	Topic     string
	UserId    int
	CreatedAt time.Time
}

type Post struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	ThreadId  int
	CreatedAt time.Time
}

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

type DataHandlerI struct {
}

type DataHandlerDeps struct {
	Log *zap.SugaredLogger
	Ctx context.Context
	Db  psqlconnector.PsqlClientI
}

type dataHandler struct {
	log *zap.SugaredLogger
	ctx context.Context
	db  psqlconnector.PsqlClientI
}

// func NewThreadHanler(deps ThreadHandlerDeps) ThreadHandlerI {
// 	return &threadHandler{
// 		log: deps.Log,
// 		ctx: deps.Ctx,
// 		db:  deps.Db,
// 	}
// }

// format the CreatedAt date to display nicely on the screen
func (thread *Thread) ThreadCreatedAtDate() string {
	return thread.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (post *Post) PostCreatedAtDate() string {
	return post.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (usr *User) UserCreatedAtDate() string {
	return usr.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (ss *Session) SessionCreatedAtDate() string {
	return ss.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

// thread functionality
func (th *dataHandler) CreateThread()

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

func (th *dataHandler) DeleteThread(thread Thread) error {
	err := th.db.DBExec(
		"DELETE FROM threads WHERE id=$1",
		thread.Id,
	)
	return err
}

// end of thread

// user functionality
func (th *dataHandler) CreateUser(user User) (User, error) {
	userRet := User{}
	err := th.db.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&userRet.Id, &userRet.Uuid, &userRet.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		},
		"INSERT INTO users (uuid, name, email, password, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, uuid, created_at",
		utils.CreateUUID(), user.Name, utils.Encrypt(user.Password), time.Now(),
	)

	return userRet, err
}

func (th *dataHandler) GetUserByEmail(email string) (User, error) {
	user := User{}
	err := th.db.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, user.CreatedAt)
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

func (th *dataHandler) GetUserByUUID(uuid string) (User, error) {
	user := User{}
	err := th.db.DBQueryRow(
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

func (th *dataHandler) UpdateUserInfo(user User) {

}

// end of user

// session functionality
func (th *dataHandler) CreateSession(user User) (Session, error) {
	ssRet := Session{}
	err := th.db.DBQueryRow(
		func(r *sql.Row) error {
			err := r.Scan(&ssRet.Id, &ssRet.Uuid, &ssRet.Email, &ssRet.UserId, &ssRet.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		},
		"INSERT INTO sessions (uuid, email, user_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id, uuid, email, user_id, created_at",
		user.Uuid, user.Email, user.Id, user.CreatedAt,
	)

	return ssRet, err
}

func (th *dataHandler) GetSessionByUUID(uuid string) (Session, error) {
	ss := Session{}
	err := th.db.DBQueryRow(
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

func (th *dataHandler) DeleteSessionByUUID(ss Session) error {
	err := th.db.DBExec(
		"DELETE FROM sessions WHERE uuid=$1",
		ss.Uuid,
	)
	return err
}

// end of session

// routing handler function

// GET /login
// Show the login page
func (th *dataHandler) Login(w http.ResponseWriter, r *http.Request) {
	t := utils.ParseTemplateFiles("login.layout", "public.navbar", "login")
	t.Execute(w, nil)
}

// GET /signup
// Show the signup page
func (th *dataHandler) Signup(w http.ResponseWriter, r *http.Request) {
	utils.GenerateHTML(w, nil, "login.layout", "public.navbar", "signup")
}

// POST /authenticate
// Authenticate the user given the email and password
func (th *dataHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		th.log.Errorf("Cannot parse form %+v", err)
	}

	user, err := th.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		th.log.Errorf("User not found %+v", err)
	}
	th.log.Debugf("User: %v", user)
	if user.Password == utils.Encrypt(r.PostFormValue("password")) {
		ss, err := th.CreateSession(user)
		if err != nil {
			th.log.Errorf("Cannot create session for this user %+v", ss)
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    ss.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		th.log.Debug("Password is not correct")
		http.Redirect(w, r, "/login", 302)
	}
}

// POST /signup
// Create the user account
func (th *dataHandler) SignupUserAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		th.log.Errorf("Cannot parse form %+v", err)
	}

	user := User{
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}

	usrRet, err := th.CreateUser(user)
	if err != nil {
		th.log.Errorf("Cannot create user: %+v, err: %+v", user, err)
	} else {
		th.log.Debugf("Created user: %+v", usrRet)
	}
	http.Redirect(w, r, "/", 302)
}

// GET /logout
// Logs the user out
func (th *dataHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		ss := Session{
			Uuid: cookie.Value,
		}
		err := th.DeleteSessionByUUID(ss)
		if err != nil {
			th.log.Errorf("Cannot logout, cannot delete session %+v", err)
		}
	}
	http.Redirect(w, r, "/", 302)
}

func (th *dataHandler) GetSession(w http.ResponseWriter, r *http.Request) (Session, error) {
	ss := Session{}
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		ss, err = th.GetSessionByUUID(cookie.Value)
		if err != nil {
			th.log.Errorf("Invalid Session %+v", err)
		}
	} else {
		th.log.Errorf("Cookie not found %+v", err)
	}
	return ss, err
}

func (th *dataHandler) CreatThread(w http.ResponseWriter, r *http.Request) {
	// ss, err :=
}

// end of routing handler
