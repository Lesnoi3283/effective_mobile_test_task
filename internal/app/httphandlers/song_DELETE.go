package httphandlers

import (
	"encoding/json"
	"errors"
	"io"
	"musiclib/internal/app/entities"
	"musiclib/pkg/databases/dberrors"
	"net/http"
)

// DeleteSong godoc
// @Summary Deletes the song
// @Description Deletes current song by id. Only ID is required.
// @Accept  json
// @Produce plain
// @Param song body entities.Song true "JSON song ID"
// @Success 201 {object} nil "Success"
// @Failure 400 {object} nil "Bad request"
// @Failure 500 {object} nil "Internal server error"
// @Router /song [delete]
func (h *handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
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

	//delete
	err = h.storage.RemoveSong(r.Context(), song.ID)
	if errors.Is(err, dberrors.NewNotFoundErr()) {
		h.logger.Debugf("no songs found with ID `%v`, err: %v", song.ID, err)
		w.WriteHeader(http.StatusNoContent)
		return
	} else if err != nil {
		h.logger.Debugf("failed to remove song in database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//answer
	w.WriteHeader(http.StatusOK)
	return
}
