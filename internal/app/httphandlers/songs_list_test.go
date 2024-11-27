package httphandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func Test_handler_SongsListGet(t *testing.T) {
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
					filter := entities.Song{
						Song:  "some song",
						Group: "some group",
					}
					expectedFilter := filter
					songs := []entities.Song{
						{
							ID:          1,
							Song:        "some song",
							Group:       "some group",
							ReleaseDate: "10.10.2010",
							Text:        "Sample text",
							Link:        "https://example.com/song1",
						},
						{
							ID:          2,
							Song:        "some song",
							Group:       "some group",
							ReleaseDate: "12.12.2012",
							Text:        "Another sample text",
							Link:        "https://example.com/song2",
						},
					}
					storage.EXPECT().GetSongList(gomock.Any(), expectedFilter, 20, 10).Return(songs, nil)
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/songs", bytes.NewBufferString(`{
					"offset": 20,
					"limit": 10,
					"filter": {
						"song": "some song",
						"group": "some group"
					}
				}`)),
			},
			expectedStatus: http.StatusOK,
			expectedBody: `[
				{
					"id":1,
					"song":"some song",
					"group":"some group",
					"release_date":"10.10.2010",
					"text":"Sample text",
					"link":"https://example.com/song1"
				},
				{
					"id":2,
					"song":"some song",
					"group":"some group",
					"release_date":"12.12.2012",
					"text":"Another sample text",
					"link":"https://example.com/song2"
				}
			]`,
		},
		{
			name: "No songs found",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					storage.EXPECT().GetSongList(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dberrors.NewNotFoundErr())
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/songs", bytes.NewBufferString(`{
					"offset": 0,
					"limit": 10,
					"filter": {
						"song": "some song",
						"group": "some group"
					}
				}`)),
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name: "Storage error",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					storage.EXPECT().GetSongList(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("test error"))
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/songs", bytes.NewBufferString(`{
					"offset": 0,
					"limit": 10,
					"filter": {
						"song": "some song",
						"group": "some group"
					}
				}`)),
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "",
		},
		{
			name: "Empty request body",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					songs := []entities.Song{
						{
							ID:          1,
							Song:        "some song",
							Group:       "some group",
							ReleaseDate: "10.10.2010",
							Text:        "Sample text",
							Link:        "https://example.com/song1",
						},
						{
							ID:          2,
							Song:        "some song",
							Group:       "some group",
							ReleaseDate: "12.12.2012",
							Text:        "Another sample text",
							Link:        "https://example.com/song2",
						},
					}
					storage.EXPECT().GetSongList(gomock.Any(), entities.Song{}, 0, 0).Return(songs, nil)
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/songs", nil),
			},
			expectedStatus: http.StatusOK,
			expectedBody: `[
				{
					"id":1,
					"song":"some song",
					"group":"some group",
					"release_date":"10.10.2010",
					"text":"Sample text",
					"link":"https://example.com/song1"
				},
				{
					"id":2,
					"song":"some song",
					"group":"some group",
					"release_date":"12.12.2012",
					"text":"Another sample text",
					"link":"https://example.com/song2"
				}
			]`,
		},
		{
			name: "Empty JSON object",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					songs := []entities.Song{
						{
							ID:          1,
							Song:        "some song",
							Group:       "some group",
							ReleaseDate: "10.10.2010",
							Text:        "Sample text",
							Link:        "https://example.com/song1",
						},
						{
							ID:          2,
							Song:        "some song",
							Group:       "some group",
							ReleaseDate: "12.12.2012",
							Text:        "Another sample text",
							Link:        "https://example.com/song2",
						},
					}
					storage.EXPECT().GetSongList(gomock.Any(), entities.Song{}, 0, 0).Return(songs, nil)
					return storage
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/songs", bytes.NewBufferString(`{}`)),
			},
			expectedStatus: http.StatusOK,
			expectedBody: `[
				{
					"id":1,
					"song":"some song",
					"group":"some group",
					"release_date":"10.10.2010",
					"text":"Sample text",
					"link":"https://example.com/song1"
				},
				{
					"id":2,
					"song":"some song",
					"group":"some group",
					"release_date":"12.12.2012",
					"text":"Another sample text",
					"link":"https://example.com/song2"
				}
			]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			h := handler{
				storage: tt.fields.storage(c),
				logger:  sugar,
			}

			h.SongsListGet(tt.args.w, tt.args.r)

			assert.Equal(t, tt.expectedStatus, tt.args.w.Code)
			//assert.Equal(t, tt.expectedBody, tt.args.w.Body.String())
			if tt.expectedBody != "" {
				var expected, actual interface{}
				err := json.Unmarshal([]byte(tt.expectedBody), &expected)
				assert.NoError(t, err)
				err = json.Unmarshal(tt.args.w.Body.Bytes(), &actual)
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			} else {
				assert.Equal(t, tt.expectedBody, tt.args.w.Body.String())
			}
		})
	}
}
