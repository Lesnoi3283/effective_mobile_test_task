package httphandlers

import (
	"encoding/json"
	"errors"
	"io"
	"musiclib/internal/app/entities"
	"musiclib/pkg/databases/dberrors"
	"net/http"
)

// PutSong godoc
// @Summary Updates the song
// @Description Updates current song by id.
// @Accept  json
// @Produce plain
// @Param song body entities.Song true "JSON song data"
// @Success 201 {nil} "Success"
// @Failure 400 {nil} "Bad request"
// @Failure 404 {nil} "Not found"
// @Failure 500 {nil} "Internal server error"
// @Router /song [put]
func (h *handler) PutSong(w http.ResponseWriter, r *http.Request) {
	//get song from request
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Debugf("failed to read body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	song := entities.Song{}
	err = json.Unmarshal(bodyBytes, &song)
	if err != nil {
		h.logger.Debugf("failed to unmarshal body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if song.ID == 0 {
		h.logger.Debugf("song ID is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//update
	err = h.storage.UpdateSong(r.Context(), song)
	if errors.Is(err, dberrors.NewNotFoundErr()) {
		h.logger.Debugf("song was`nt found in db, err: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Debugf("failed to update song in database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//answer
	w.WriteHeader(http.StatusOK)
	return
}
