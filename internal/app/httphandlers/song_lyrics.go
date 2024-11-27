package httphandlers

import (
	"errors"
	"musiclib/pkg/databases/dberrors"
	"net/http"
	"strconv"
	"strings"
)

// GetSongLyrics godoc
// @Summary Returns a specific couplet from a song's lyrics
// @Description Retrieves the requested couplet number from the specified song.
// @Produce text/plain
// @Param song_id query uint64 true "Song ID"
// @Param couplet_num query int true "Couplet Number"
// @Success 200 {string} string "Couplet text"
// @Failure 400 {object} nil "Bad request"
// @Failure 404 {object} nil "Song not found"
// @Failure 204 {object} nil "Requested couplet number is bigger than the number of couplets"
// @Failure 500 {object} nil "Internal server error"
// @Router /lyrics [get]
func (h *handler) GetSongLyrics(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	query := r.URL.Query()
	songIDStr := query.Get("song_id")
	coupletNumStr := query.Get("couplet_num")

	songID, err := strconv.ParseUint(songIDStr, 10, 64)
	if err != nil {
		h.logger.Debugf("invalid song_id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	coupletNum, err := strconv.Atoi(coupletNumStr)
	if err != nil {
		h.logger.Debugf("invalid couplet_num: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get lyrics
	text, err := h.storage.GetSongLyrics(r.Context(), songID)
	if errors.Is(err, dberrors.NewNotFoundErr()) {
		h.logger.Debugf("no song found with id %d", songID)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Debugf("failed get song lyrics from db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get couplet
	couplets := strings.Split(text, "\n\n")
	if coupletNum > len(couplets)-1 {
		h.logger.Debugf("requested couplet num is `%v`, but this song has only `%v` couplets.", coupletNum, len(couplets))
		w.WriteHeader(http.StatusNoContent)
		return
	}
	answer := []byte(couplets[coupletNum])

	// Answer
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
	return
}
