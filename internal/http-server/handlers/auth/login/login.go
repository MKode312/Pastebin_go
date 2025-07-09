package login

import (
	"errors"
	"log/slog"
	"net/http"
	resp "text_sharing/internal/lib/api/response"
	"text_sharing/internal/lib/logger/sl"
	"text_sharing/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UserLogin interface {
	LoginUser(username string, password string) (string, error)
}

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Token string `json:"token"`
	resp.Response
}

func New(log *slog.Logger, userLogin UserLogin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.login.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("requset body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		token, err := userLogin.LoginUser(req.Username, req.Password)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Error("user not found", slog.String("username", req.Username), slog.String("password", req.Password))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("user not found"))
				return
			}

			log.Error("failed to login user", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to login user"))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		log.Info("successful login", slog.String("token", token))

		w.WriteHeader(http.StatusAccepted)

		render.JSON(w, r, Response{
			Token: token,
			Response: resp.OK(),
		})
	}
}
