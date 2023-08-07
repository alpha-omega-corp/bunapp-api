package httphandlers

import (
	"encoding/json"
	"fmt"
	"github.com/alpha-omega-corp/bunapp-api/app"
	"github.com/alpha-omega-corp/bunapp-api/repository"
	"github.com/alpha-omega-corp/bunapp-api/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/uptrace/bunrouter"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

type UserHandler struct {
	app        *app.App
	repository repository.IUserRepository
}

func NewUserHandler(app *app.App) *UserHandler {
	return &UserHandler{
		app:        app,
		repository: app.Repositories().User(),
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, req bunrouter.Request) error {
	data := new(types.LoginRequest)
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return err
	}

	user, err := h.repository.GetByEmail(data.Email, req.Context())
	if err != nil {
		return err
	}

	if !user.Verify(data.Password) {
		return fmt.Errorf("authentication failed for user %s", data.Email)
	}

	token, err := user.CreateToken()
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, types.LoginResponse{
		User:  user,
		Token: token,
	})
}

func (h *UserHandler) Register(w http.ResponseWriter, req bunrouter.Request) error {
	data := new(types.CreateUserRequest)
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return err
	}

	encPw, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user, err := h.repository.CreateUser(&types.UserRaw{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Age:       data.Age,

		EncryptedPassword: string(encPw),
	}, req.Context())

	if err != nil {
		return err
	}

	return bunrouter.JSON(w, user)
}

func (h *UserHandler) Get(w http.ResponseWriter, req bunrouter.Request) error {
	parseInt, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	user, err := h.repository.GetById(parseInt, req.Context())
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, user)
}

func (h *UserHandler) List(w http.ResponseWriter, req bunrouter.Request) error {
	users, err := h.repository.GetAll(req.Context())
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, users)
}

func (h *UserHandler) UserFromToken(w http.ResponseWriter, req bunrouter.Request) error {
	token, err := app.GetValidTokenFromReq(w, req)
	if err != nil {
		return err
	}

	claims := token.Claims.(jwt.MapClaims)
	user, err := h.repository.GetByEmail(claims["email"].(string), req.Context())
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, user)
}
