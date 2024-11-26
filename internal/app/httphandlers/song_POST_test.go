package httphandlers

import (
	"bytes"
	"fmt"
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

func Test_handler_PostSong(t *testing.T) {
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	type fields struct {
		storage           func(c *gomock.Controller) requiredinterfaces.SongStorage
		extraDataProvider func(c *gomock.Controller) requiredinterfaces.ExtraDataProvider
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
			name: "Ok",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					storage.EXPECT().SaveSong(gomock.Any(), entities.Song{
						Song:        "some song",
						Group:       "some group",
						ReleaseDate: "10.10.2010",
						Text:        "some text\ntext2\n\ntext3.",
						Link:        "https://example.com/somesong",
					}).Return(uint64(10), nil)
					return storage
				},
				extraDataProvider: func(c *gomock.Controller) requiredinterfaces.ExtraDataProvider {
					provider := mocks.NewMockExtraDataProvider(c)
					provider.EXPECT().GetExtraSongData(entities.Song{
						Song:  "some song",
						Group: "some group",
					}).Return("10.10.2010", "some text\ntext2\n\ntext3.", "https://example.com/somesong", nil)
					return provider
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/song", bytes.NewBufferString(`{"group":"some group","song":"some song"}`)),
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":10}`,
		},
		{
			name: "empty request",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					return storage
				},
				extraDataProvider: func(c *gomock.Controller) requiredinterfaces.ExtraDataProvider {
					provider := mocks.NewMockExtraDataProvider(c)
					return provider
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/song", nil),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name: "Empty JSON",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					return storage
				},
				extraDataProvider: func(c *gomock.Controller) requiredinterfaces.ExtraDataProvider {
					provider := mocks.NewMockExtraDataProvider(c)
					return provider
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/song", bytes.NewBufferString(`{}`)),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name: "Db error",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					storage.EXPECT().SaveSong(gomock.Any(), gomock.Any()).Return(uint64(0), fmt.Errorf("test error"))
					return storage
				},
				extraDataProvider: func(c *gomock.Controller) requiredinterfaces.ExtraDataProvider {
					provider := mocks.NewMockExtraDataProvider(c)
					provider.EXPECT().GetExtraSongData(entities.Song{
						Song:  "some song",
						Group: "some group",
					}).Return("10.10.2010", "some text\ntext2\n\ntext3.", "https://example.com/somesong", nil)
					return provider
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/song", bytes.NewBufferString(`{"group":"some group","song":"some song"}`)),
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "",
		},
		{
			name: "Provider error",
			fields: fields{
				storage: func(c *gomock.Controller) requiredinterfaces.SongStorage {
					storage := mocks.NewMockSongStorage(c)
					storage.EXPECT().SaveSong(gomock.Any(), entities.Song{
						Song:  "some song",
						Group: "some group",
					}).Return(uint64(10), nil)
					return storage
				},
				extraDataProvider: func(c *gomock.Controller) requiredinterfaces.ExtraDataProvider {
					provider := mocks.NewMockExtraDataProvider(c)
					provider.EXPECT().GetExtraSongData(gomock.Any()).Return("", "", "", fmt.Errorf("test error"))
					return provider
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/song", bytes.NewBufferString(`{"group":"some group","song":"some song"}`)),
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":10}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			h := &handler{
				storage:           tt.fields.storage(c),
				logger:            sugar,
				extraDataProvider: tt.fields.extraDataProvider(c),
			}
			h.PostSong(tt.args.w, tt.args.r)

			assert.Equal(t, tt.expectedStatus, tt.args.w.Code)
			assert.Equal(t, tt.expectedBody, tt.args.w.Body.String())
		})
	}
}
