package app

import (
	"fmt"
	"github.com/alpha-omega-corp/bunapp-api/app/httputils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/bunrouterotel"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"net/http"
	"strconv"
)

var (
	tokenHeader = "x-jwt-token"
)

func (app *App) initRouter() {
	app.router = bunrouter.New(
		bunrouter.WithMiddleware(bunrouterotel.NewMiddleware()),
		bunrouter.WithMiddleware(reqlog.NewMiddleware(
			reqlog.WithEnabled(app.IsDebug()),
			reqlog.WithVerbose(true),
			reqlog.FromEnv(""),
		)))

	app.apiRouter = app.router.NewGroup("/api",
		bunrouter.WithMiddleware(corsMiddleware),
		bunrouter.WithMiddleware(app.errorHandler),
	)
}

func (app *App) AuthHandler(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		_, err := GetValidTokenFromReq(w, req)
		if err != nil {
			return err
		}

		return next(w, req)
	}
}

func (app *App) AuthClaimHandler(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		token, err := GetValidTokenFromReq(w, req)
		if err != nil {
			return err
		}

		parseInt, err := strconv.ParseInt(req.Param("id"), 10, 64)
		if err != nil {
			return err
		}

		user, err := app.Repositories().User().GetById(parseInt, req.Context())
		if err != nil {
			return err
		}

		if err := user.Claims(token.Claims.(jwt.MapClaims)); err != nil {
			return err
		}

		return next(w, req)
	}
}

func (app *App) errorHandler(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		err := next(w, req)
		if err == nil {
			return nil
		}

		httpErr := httputils.From(err, app.IsDebug())
		if httpErr.Status != 0 {
			w.WriteHeader(httpErr.Status)
		}
		_ = bunrouter.JSON(w, httpErr)

		return err
	}
}

func corsMiddleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		origin := req.Header.Get("Origin")
		fmt.Print(origin)
		h := w.Header()

		h.Set("Access-Control-Allow-Origin", origin)
		h.Set("Access-Control-Allow-Credentials", "true")
		h.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,HEAD")
		h.Set("Access-Control-Allow-Headers", "authorization,content-type")
		h.Set("Access-Control-Max-Age", "86400")

		return next(w, req)
	}
}
