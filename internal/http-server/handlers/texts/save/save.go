package save

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	resp "text_sharing/internal/lib/api/response"
	"text_sharing/internal/lib/logger/sl"
	"text_sharing/internal/lib/utils"
	"text_sharing/internal/storage"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type TextSaverDB interface {
	SaveLink(urlToSave string, date string, objectID string) (string, error)
}

type TextSaverMinio interface {
	CreateOne(file utils.FileDataType, expires time.Duration) (string, string, error)
}

type TextSaverRedis interface {
	SetLinksCache(linkID string, link string) error
}

type Request struct {
	Filename    string `json:"filename" validate:"required"`
	TextBlock   string `json:"text" validate:"required"`
	ExpiresMins int64  `json:"expiresAfter(mins)" validate:"required"`
}

type Response struct {
	URLtoGetFromDB string `json:"url"`
	resp.Response
}

func New(log *slog.Logger, textSaverDB TextSaverDB, textSaverMinio TextSaverMinio, textSaverRedis TextSaverRedis) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.texts.Save"

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

		file := utils.FileDataType{
			FileName: req.Filename,
			Data:     []byte(req.TextBlock),
		}

		log.Info("link will be expired after", slog.Duration("expires", time.Duration(req.ExpiresMins) * time.Minute))

		url, objectID, err := textSaverMinio.CreateOne(file, time.Duration(req.ExpiresMins) * time.Minute)
		if err != nil {
			log.Error("failed to create link for the text", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create link for the text"))
			return
		}

		log.Info("link was created", slog.String("url", url), slog.String("objectId", objectID))

		date := time.Now()

		linkID, err := textSaverDB.SaveLink(url, date.Format("2006/01/02 15:04:05"), objectID)
		if err != nil {
			if errors.Is(err, storage.ErrLinkExists) {
				log.Error("link for this text already exists in db", slog.String("link", url))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("link for this text has been already saved"))
				return
			}

			log.Error("failed to save link for the text", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to save link for this text"))
			return
		}

		err = textSaverRedis.SetLinksCache(objectID, url)
		if err != nil {
			log.Error("failed to set link in the cache", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to set link in the cache"))
			return
		}

		log.Info("link was successfully set in the cache")

		parts := strings.Split(r.URL.String(), "/")
		if len(parts) > 0 {
			parts = parts[:len(parts)-1]
		} else {
			log.Error("invalid request url", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("you have invalid url"))
			return
		}

		newURL := "http://localhost:8082" + strings.Join(parts, "/") + "/" + linkID
		log.Info("link saved", slog.String("linkID", linkID))
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, Response{
			Response:       resp.OK(),
			URLtoGetFromDB: newURL,
		})

	}
}
