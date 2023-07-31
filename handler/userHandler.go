package handler

import (
	"chadgpt-api/app"
	"chadgpt-api/repository"
	"chadgpt-api/types"
	"encoding/json"
	"fmt"
	"github.com/uptrace/bunrouter"
	"net/http"
)

type UserHandler struct {
	app        *app.App
	repository *repository.UserRepository
}

func NewUserHandler(app *app.App) *UserHandler {
	return &UserHandler{
		app:        app,
		repository: repository.NewUserRepository(app.Database()),
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, req bunrouter.Request) error {
	var data types.LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		return err
	}

	user, err := h.repository.GetUserByEmail(data.Email, req.Context())
	if err != nil {
		return err
	}

	if !user.MatchPassword(data.Password) {
		return fmt.Errorf("authentication failed for user %s", data.Email)
	}

	token, err := user.CreateJwt()
	if err != nil {
		return err
	}

	res := types.LoginResponse{
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
