package register

import (
	"errors"
	"log/slog"
	"net/http"
	resp "text_sharing/internal/lib/api/response"
	"text_sharing/internal/lib/logger/sl"
	"text_sharing/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UserRegister interface {
	SaveUser(username string, password string) (int64, error)
}

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
}

func New(log *slog.Logger, userRegister UserRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.register.New"

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

		id, err := userRegister.SaveUser(req.Username, req.Password)
		if err != nil {
			if errors.Is(err, storage.ErrUserExists) {
				log.Error("user already registered", slog.String("username", req.Username), slog.String("password", req.Password))

				w.WriteHeader(http.StatusBadRequest)

				render.JSON(w, r, resp.Error("user already registered"))

				return
			}

			log.Error("failed to register user", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)

			render.JSON(w, r, resp.Error("failed to register user"))

			return
		}

		log.Info("user regitered", slog.Int64("id", id))

		w.WriteHeader(http.StatusCreated)

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})

	}
}
