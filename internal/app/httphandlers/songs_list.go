package httphandlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"musiclib/internal/app/entities"
	"musiclib/pkg/databases/dberrors"
	"net/http"
)

type FilterRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Filter entities.Song
}

// SongsListGet godoc
// @Summary Returns list of songs
// @Description filtration and pagination are supported
// @Accept  json
// @Produce application/json
// @Param filter body FilterRequest false "Filter params"
// @Success 201 {array} FilterRequest "Songs list"
// @Failure 204 {object} nil "No songs found"
// @Failure 400 {object} nil "Bad request"
// @Failure 500 {object} nil "Internal server error"
// @Router /songs [post]
func (h *handler) SongsListGet(w http.ResponseWriter, r *http.Request) {
	//get filter from request
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Debugf("failed to read body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	filter := FilterRequest{}
	if len(bodyBytes) > 0 {
		err = json.Unmarshal(bodyBytes, &filter)
		if err != nil {
			h.logger.Debugf("failed to unmarshal body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	//get songs list
	songs, err := h.storage.GetSongList(r.Context(), filter.Filter, filter.Offset, filter.Limit)
	if errors.Is(err, dberrors.NewNotFoundErr()) {
		h.logger.Debugf("no songs with rerquested params: %v", err)
		w.WriteHeader(http.StatusNoContent)
		return
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.logger.Debugf("failed get songs: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//answer
	jsonAnswer, err := json.Marshal(songs)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonAnswer)
	return
}
