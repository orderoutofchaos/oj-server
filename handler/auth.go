package handler

import (
	"encoding/hex"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/orderoutofchaos/oj-server/dao"
	"golang.org/x/crypto/scrypt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const salt = ""

func encodeHash(value string) string {
	key, _ := scrypt.Key([]byte(value), []byte(salt), 32768, 8, 1, 32)
	return hex.EncodeToString(key)
}

func V1Login(req *LoginRequest, c echo.Context) Response {
	repo := dao.GetRepo()
	rows, err := repo.Reader().Query("SELECT id, username FROM user WHERE username = ? AND password = ?", req.Username, encodeHash(req.Password))
	if err != nil {
		panic(err)
	}

	var (
		id       int64
		username string
	)
	if rows.Next() {
		err := rows.Scan(&id, &username)
		if err != nil {
			panic(err)
		}

		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		sess.Values["user_id"] = id
		sess.Values["username"] = username
		sess.Save(c.Request(), c.Response())

		return Response{
			Code: http.StatusOK,
		}
	} else {
		return Response{
			Code:  http.StatusUnauthorized,
			Error: ErrInvalidCredential,
		}
	}
}
