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
	"musiclib/pkg/databases/dberrors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_PutSong(t *testing.T) {
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	type fields struct {
		storage func(c *gomock.Controller) requiredinterfaces.SongStorage
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
	}{
		{
			name: "OK",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					song := entities.Song{
						ID:          1,
						Song:        "updated song",
						Group:       "updated group",
						ReleaseDate: "10.10.2010",
						Text:        "Updated text",
						Link:        "https://example.com/updatedsong",
					}
					storage.EXPECT().UpdateSong(gomock.Any(), song).Return(nil)
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("PUT", "/song", bytes.NewBufferString(`{
					"id": 1,
					"song": "updated song",
					"group": "updated group",
					"release_date": "10.10.2010",
					"text": "Updated text",
					"link": "https://example.com/updatedsong"
				}`)),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Bag JSON",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					return mocks.NewMockSongStorage(c)
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("PUT", "/song", bytes.NewBufferString(`invalid json`)),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty Song ID",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					return mocks.NewMockSongStorage(c)
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("PUT", "/song", bytes.NewBufferString(`{
					"song": "song without ID",
					"group": "group without ID",
					"release_date": "12.12.2012",
					"text": "Some text",
					"link": "https://example.com/songwithoutid"
				}`)),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Song Not Found",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					song := entities.Song{
						ID:          999,
						Song:        "nonexistent song",
						Group:       "nonexistent group",
						ReleaseDate: "12.12.2012",
						Text:        "No text",
						Link:        "https://example.com/nonexistentsong",
					}
					storage.EXPECT().UpdateSong(gomock.Any(), song).Return(dberrors.NewNotFoundErr())
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("PUT", "/song", bytes.NewBufferString(`{
					"id": 999,
					"song": "nonexistent song",
					"group": "nonexistent group",
					"release_date": "12.12.2012",
					"text": "No text",
					"link": "https://example.com/nonexistentsong"
				}`)),
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Storage Error",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					song := entities.Song{
						ID:          2,
						Song:        "song with error",
						Group:       "group with error",
						ReleaseDate: "09.09.2009",
						Text:        "Error text",
						Link:        "https://example.com/songwitherror",
					}
					storage.EXPECT().UpdateSong(gomock.Any(), song).Return(errors.New("test error"))
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("PUT", "/song", bytes.NewBufferString(`{
					"id": 2,
					"song": "song with error",
					"group": "group with error",
					"release_date": "09.09.2009",
					"text": "Error text",
					"link": "https://example.com/songwitherror"
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
				r: httptest.NewRequest("PUT", "/song", nil),
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
				r: httptest.NewRequest("PUT", "/song", bytes.NewBufferString("{}")),
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := handler{
				storage: tt.fields.storage(ctrl),
				logger:  sugar,
			}

			h.PutSong(tt.args.w, tt.args.r)

			assert.Equal(t, tt.expectedStatus, tt.args.w.(*httptest.ResponseRecorder).Code)
		})
	}
}
