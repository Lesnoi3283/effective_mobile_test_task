package httphandlers

import (
	"encoding/json"
	"io"
	"musiclib/internal/app/entities"
	"net/http"
)

type returnableID struct {
	ID uint64 `json:"id"`
}

// PostSong godoc
// @Summary Creates new song
// @Description Creates new song, asks different service for an additional data (release date, text and link).
// @Accept  json
// @Produce json
// @Param song body entities.Song true "JSON song data"
// @Success 201 {object} returnableID "Song ID"
// @Failure 400 {object} nil "Bad request"
// @Failure 500 {object} nil "Internal server error"
// @Router /song [post]
func (h *handler) PostSong(w http.ResponseWriter, r *http.Request) {
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
	if len(song.Song) == 0 {
		h.logger.Debugf("song name is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(song.Group) == 0 {
		h.logger.Debugf("song group is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//ask api for an additional data
	song.ReleaseDate, song.Text, song.Link, err = h.extraDataProvider.GetExtraSongData(song)
	if err != nil {
		h.logger.Debugf("failed to get song data: %v", err)
	}

	//save
	song.ID, err = h.storage.SaveSong(r.Context(), song)
	if err != nil {
		h.logger.Debugf("failed to create song in database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//answer
	jsonAnswer, err := json.Marshal(returnableID{
		ID: song.ID,
	})

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonAnswer)
	return
}
