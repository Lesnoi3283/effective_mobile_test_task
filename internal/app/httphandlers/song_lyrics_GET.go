package httphandlers

import (
	"encoding/json"
	"errors"
	"io"
	"musiclib/pkg/databases/dberrors"
	"net/http"
	"strings"
)

type SongLyricsRequest struct {
	SongID     uint64 `json:"song_id"`
	CoupletNum int    `json:"couplet_num"`
}

// GetSongLyrics godoc
// @Summary Returns song`s lyrics
// @Description
// @Accept  json
// @Produce text/plain
// @Param song body SongLyricsRequest true "JSON song lyrics data"
// @Success 201 {object} returnableID "Song ID"
// @Failure 400 {object} nil "Bad request"
// @Failure 404 {object} nil "Song not found"
// @Failure 204 {object} nil "Requested couplet number is bigger than len of couplets"
// @Failure 500 {object} nil "Internal server error"
// @Router /lyrics [post]
func (h *handler) GetSongLyrics(w http.ResponseWriter, r *http.Request) {
	//get song from request
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Debugf("failed to read body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	lyricsRequest := SongLyricsRequest{}
	err = json.Unmarshal(bodyBytes, &lyricsRequest)
	if err != nil {
		h.logger.Debugf("failed to unmarshal body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if lyricsRequest.SongID == 0 {
		h.logger.Debugf("song ID is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//get lyrics
	text, err := h.storage.GetSongLyrics(r.Context(), lyricsRequest.SongID)
	if errors.Is(err, dberrors.NewNotFoundErr()) {
		h.logger.Debugf("no song found with id %d", lyricsRequest.SongID)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Debugf("failed get song lyrics from db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//get couplet
	couplets := strings.Split(text, "\n\n")
	if lyricsRequest.CoupletNum > len(couplets)-1 {
		h.logger.Debugf("requested couplet num is `%v`, but this song has only `%v` couplets.", lyricsRequest.CoupletNum, len(couplets))
		w.WriteHeader(http.StatusNoContent)
		return
	}
	answer := []byte(couplets[lyricsRequest.CoupletNum])

	//answer
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
	return
}
