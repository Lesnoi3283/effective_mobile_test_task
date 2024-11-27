package httphandlers

import (
	"bytes"
	"database/sql"
	"errors"
	"musiclib/internal/app/requiredinterfaces"
	"musiclib/internal/app/requiredinterfaces/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func Test_handler_GetSongLyrics(t *testing.T) {
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	type fields struct {
		storage func(c *gomock.Controller) requiredinterfaces.SongStorage
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Normal",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					songID := uint64(1)
					lyrics := "First couplet.\nStill first couplet.\n\nSecond couplet.\n\nThird couplet."
					storage.EXPECT().GetSongLyrics(gomock.Any(), songID).Return(lyrics, nil)
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/song", bytes.NewBufferString(`{
					"song_id": 1,
					"couplet_num": 1
				}`)),
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Second couplet.",
		},
		{
			name: "Bad JSON",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					return mocks.NewMockSongStorage(c)
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/song", bytes.NewBufferString(`invalid json`)),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name: "No song ID",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					return mocks.NewMockSongStorage(c)
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/song", bytes.NewBufferString(`{
					"couplet_num": 1
				}`)),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name: "Song Not Found",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					storage.EXPECT().GetSongLyrics(gomock.Any(), gomock.Any()).Return("", sql.ErrNoRows)
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/song", bytes.NewBufferString(`{
					"song_id": 2,
					"couplet": 1
				}`)),
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			name: "Storage Error",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					storage.EXPECT().GetSongLyrics(gomock.Any(), gomock.Any()).Return("", errors.New("test error"))
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/song", bytes.NewBufferString(`{
					"song_id": 2,
					"couplet": 1
				}`)),
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "",
		},
		{
			name: "No Content - Couplet Number Too High",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					lyrics := "Only one couplet."
					storage.EXPECT().GetSongLyrics(gomock.Any(), gomock.Any()).Return(lyrics, nil)
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/song", bytes.NewBufferString(`{
					"song_id": 3,
					"couplet_num": 200
				}`)),
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name: "First Couplet (0 Index)",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					lyrics := "First couplet.\nStill first couplet.\n\nSecond couplet."
					storage.EXPECT().GetSongLyrics(gomock.Any(), gomock.Any()).Return(lyrics, nil)
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/song", bytes.NewBufferString(`{
					"song_id": 5,
					"couplet": 0
				}`)),
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "First couplet.\nStill first couplet.",
		},
		{
			name: "Last Valid Couplet",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					lyrics := "First couplet.\n\nSecond couplet.\n\nThird couplet."
					storage.EXPECT().GetSongLyrics(gomock.Any(), gomock.Any()).Return(lyrics, nil)
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/song", bytes.NewBufferString(`{
					"song_id": 6,
					"couplet_num": 2
				}`)),
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Third couplet.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := &handler{
				storage: tt.fields.storage(ctrl),
				logger:  sugar,
			}

			h.GetSongLyrics(tt.args.w, tt.args.r)

			assert.Equal(t, tt.expectedStatus, tt.args.w.Code)
			assert.Equal(t, tt.expectedBody, tt.args.w.Body.String())
		})
	}
}
