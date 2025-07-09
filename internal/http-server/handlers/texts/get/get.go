package get

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
)

type TextGetterDB interface {
	GetLink(linkID string) (string, error)
	GetDate(linkID string) (string, error)
}

type TextGetterMinio interface {
	GetOne(fileName string, date time.Time) (string, error)
}

type TextGetterRedis interface {
	GetLinkFromCache(linkID string) (string, error)
}

type Response struct {
	TextBlock string `json:"text,omitempty"`
	resp.Response
}

func New(log *slog.Logger, textGetterDB TextGetterDB, textGetterMinio TextGetterMinio, textGetterRedis TextGetterRedis) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.texts.Get"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		pathValues := strings.Split(r.URL.Path, "/")
		if len(pathValues) < 3 || pathValues[2] == "" {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}
		objectID := pathValues[2]

		linkForMinio, err := textGetterRedis.GetLinkFromCache(objectID)
		if err != nil {
			if errors.Is(err, storage.ErrCacheMiss) {
				log.Info("link not found in cache")
				linkForMinio, err = textGetterDB.GetLink(objectID)
				if err != nil {
					if errors.Is(err, storage.ErrLinkNotFound) {
						log.Error("link not found", sl.Err(err))
						w.WriteHeader(http.StatusBadRequest)
						render.JSON(w, r, resp.Error("link for text not found"))
						return
					}
					log.Error("failed to get link", sl.Err(err))
					w.WriteHeader(http.StatusInternalServerError)
					render.JSON(w, r, resp.Error("failed to get link"))
					return
				}
			} else {
				log.Error("failed to get link from cache", sl.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, resp.Error("failed to get link from cache"))
				return
			}
		} else {
			log.Info("link was recieved from cache")
		}

		date, err := textGetterDB.GetDate(objectID)
		if err != nil {
			if errors.Is(err, storage.ErrLinkNotFound) {
				log.Error("link not found", sl.Err(err))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("link for text not found"))
				return
			}

			log.Error("failed to get date", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get date"))
			return
		}

		parsedDate, err := time.Parse("2006/01/02 15:04:05", date)
		if err != nil {
			log.Error("failed to parse the date", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to parse the date"))
			return
		}

		fileName := utils.ExtractFilenameFromURL(linkForMinio)

		log.Info("filename is being used", slog.String("filename", fileName))

		fileName, err = utils.ExtractFilenameFromString(fileName)
		if err != nil {
			log.Error("failed to parse filename", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to parse filename"))
			return
		}

		log.Info("filename after parsing", slog.String("filename", fileName))

		textBlock, err := textGetterMinio.GetOne(fileName, parsedDate)
		if err != nil {
			if errors.Is(err, storage.ErrObjectExpired) {
				log.Error("this link is expired", sl.Err(err))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("these link and text are expired"))
				return
			}
			log.Error("failed to get textblock from the file", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get text block"))
			return
		}

		log.Info("successful get link operation")

		render.JSON(w, r, Response{
			TextBlock: textBlock,
			Response:  resp.OK(),
		})

	}
}
