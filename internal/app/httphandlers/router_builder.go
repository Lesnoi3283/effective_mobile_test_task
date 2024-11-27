package httphandlers

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"musiclib/internal/app/requiredinterfaces"
)

type handler struct {
	storage           requiredinterfaces.SongStorage
	logger            *zap.SugaredLogger
	extraDataProvider requiredinterfaces.ExtraDataProvider
}

func NewHTTPRouter(logger *zap.SugaredLogger, storage requiredinterfaces.SongStorage, extraDataProvider requiredinterfaces.ExtraDataProvider) chi.Router {

	h := &handler{
		storage:           storage,
		logger:            logger,
		extraDataProvider: extraDataProvider,
	}

	r := chi.NewRouter()

	r.Post("/song", h.PostSong)
	r.Put("/song", h.PutSong)
	r.Get("/lyrics", h.GetSongLyrics)
	r.Delete("/song", h.DeleteSong)
	r.Get("/songs", h.SongsListGet)

	return r
}
