package handlers

import (
	"chadgpt-api/app"
	"chadgpt-api/types"
	"encoding/json"
	"fmt"
	"github.com/uptrace/bunrouter"
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
	User  types.UserResponse `json:"user"`
	Token string             `json:"token"`
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
		User:  user.ToResponse(),
		Token: token,
	}

	return bunrouter.JSON(w, res)
}

func (h *UserHandler) Register(w http.ResponseWriter, req bunrouter.Request) error {
	data := new(types.CreateUserRequest)
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return err
	}

	user, err := h.app.NewUser(req.Context(), data)
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, user)
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

	return bunrouter.JSON(w, users)
}
