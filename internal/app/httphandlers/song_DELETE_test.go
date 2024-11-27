package httphandlers

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"musiclib/internal/app/entities"
	"musiclib/internal/app/requiredinterfaces"
	"musiclib/internal/app/requiredinterfaces/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_DeleteSong(t *testing.T) {
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
	}{
		{
			name: "Normal",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					song := entities.Song{
						ID: 1,
					}
					storage.EXPECT().RemoveSong(gomock.Any(), song.ID).Return(nil)
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("DELETE", "/song", bytes.NewBufferString(`{
					"id": 1
				}`)),
			},
			expectedStatus: http.StatusOK,
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
				r: httptest.NewRequest("DELETE", "/song", bytes.NewBufferString(`invalid json`)),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty JSON",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					return mocks.NewMockSongStorage(c)
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("DELETE", "/song", bytes.NewBufferString(`{}`)),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Storage Error",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					storage.EXPECT().RemoveSong(gomock.Any(), gomock.Any()).Return(errors.New("test error"))
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("DELETE", "/song", bytes.NewBufferString(`{
					"id": 2
				}`)),
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Empty Request Body",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					return mocks.NewMockSongStorage(c)
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("DELETE", "/song", nil),
			},
			expectedStatus: http.StatusBadRequest,
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

			h.DeleteSong(tt.args.w, tt.args.r)

			assert.Equal(t, tt.expectedStatus, tt.args.w.Code)
		})
	}
}
