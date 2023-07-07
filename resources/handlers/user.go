package handlers

import (
	"chadgpt-api/app"
	"encoding/json"
	"fmt"
	"github.com/uptrace/bunrouter"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
)

type UserHandler struct {
	app *app.App
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
}

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func NewUserHandler(app *app.App) *UserHandler {
	return &UserHandler{
		app: app,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, req bunrouter.Request) error {
	var data LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		return err
	}

	user := new(app.User)
	if err := h.app.Database().NewSelect().Where("email = ?", data.Email).Model(user).Scan(h.app.Context()); err != nil {
		return err
	}

	if !user.IsValid(data.Password) {
		return fmt.Errorf("authentication failed for user %s", data.Email)
	}

	token, err := user.CreateJwt()
	if err != nil {
		return err
	}

	res := LoginResponse{
		Id:    user.ID,
		Token: token,
	}

	return bunrouter.JSON(w, res)
}

func (h *UserHandler) Create(w http.ResponseWriter, req bunrouter.Request) error {
	data := new(CreateUserRequest)
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return err
	}

	encPw, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	res, err := h.app.Database().NewInsert().Model(&app.User{
		FirstName:         data.FirstName,
		LastName:          data.LastName,
		Age:               data.Age,
		Email:             data.Email,
		EncryptedPassword: string(encPw),
		AccountNumber:     int64(rand.Intn(1000000)),
	}).Exec(req.Context())

	if err != nil {
		return err
	}

	return bunrouter.JSON(w, bunrouter.H{
		"id": res,
	})
}

func (h *UserHandler) Get(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()
	id := req.Param("id")

	var user app.User
	if err := h.app.Database().NewSelect().Where("id = ?", id).Model(&user).Scan(ctx); err != nil {
		return err
	}

	return bunrouter.JSON(w, user)
}

func (h *UserHandler) List(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	var users []app.User
	if err := h.app.Database().NewSelect().Model(&users).Scan(ctx); err != nil {
		return err
	}

	return bunrouter.JSON(w, bunrouter.H{
		"users": users,
	})
}
