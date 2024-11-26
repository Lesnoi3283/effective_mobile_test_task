// Code generated by MockGen. DO NOT EDIT.
// Source: required_interfaces.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	entities "musiclib/internal/app/entities"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockExtraDataProvider is a mock of ExtraDataProvider interface.
type MockExtraDataProvider struct {
	ctrl     *gomock.Controller
	recorder *MockExtraDataProviderMockRecorder
}

// MockExtraDataProviderMockRecorder is the mock recorder for MockExtraDataProvider.
type MockExtraDataProviderMockRecorder struct {
	mock *MockExtraDataProvider
}

// NewMockExtraDataProvider creates a new mock instance.
func NewMockExtraDataProvider(ctrl *gomock.Controller) *MockExtraDataProvider {
	mock := &MockExtraDataProvider{ctrl: ctrl}
	mock.recorder = &MockExtraDataProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExtraDataProvider) EXPECT() *MockExtraDataProviderMockRecorder {
	return m.recorder
}

// GetExtraSongData mocks base method.
func (m *MockExtraDataProvider) GetExtraSongData(song entities.Song) (string, string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExtraSongData", song)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetExtraSongData indicates an expected call of GetExtraSongData.
func (mr *MockExtraDataProviderMockRecorder) GetExtraSongData(song interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExtraSongData", reflect.TypeOf((*MockExtraDataProvider)(nil).GetExtraSongData), song)
}

// MockSongStorage is a mock of SongStorage interface.
type MockSongStorage struct {
	ctrl     *gomock.Controller
	recorder *MockSongStorageMockRecorder
}

// MockSongStorageMockRecorder is the mock recorder for MockSongStorage.
type MockSongStorageMockRecorder struct {
	mock *MockSongStorage
}

// NewMockSongStorage creates a new mock instance.
func NewMockSongStorage(ctrl *gomock.Controller) *MockSongStorage {
	mock := &MockSongStorage{ctrl: ctrl}
	mock.recorder = &MockSongStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSongStorage) EXPECT() *MockSongStorageMockRecorder {
	return m.recorder
}

// GetSongList mocks base method.
func (m *MockSongStorage) GetSongList(ctx context.Context, filter entities.Song, offset int) ([]entities.Song, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSongList", ctx, filter, offset)
	ret0, _ := ret[0].([]entities.Song)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSongList indicates an expected call of GetSongList.
func (mr *MockSongStorageMockRecorder) GetSongList(ctx, filter, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSongList", reflect.TypeOf((*MockSongStorage)(nil).GetSongList), ctx, filter, offset)
}

// GetSongLyrics mocks base method.
func (m *MockSongStorage) GetSongLyrics(ctx context.Context, id int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSongLyrics", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSongLyrics indicates an expected call of GetSongLyrics.
func (mr *MockSongStorageMockRecorder) GetSongLyrics(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSongLyrics", reflect.TypeOf((*MockSongStorage)(nil).GetSongLyrics), ctx, id)
}

// RemoveSong mocks base method.
func (m *MockSongStorage) RemoveSong(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSong", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSong indicates an expected call of RemoveSong.
func (mr *MockSongStorageMockRecorder) RemoveSong(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSong", reflect.TypeOf((*MockSongStorage)(nil).RemoveSong), ctx, id)
}

// SaveSong mocks base method.
func (m *MockSongStorage) SaveSong(ctx context.Context, song entities.Song) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSong", ctx, song)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveSong indicates an expected call of SaveSong.
func (mr *MockSongStorageMockRecorder) SaveSong(ctx, song interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSong", reflect.TypeOf((*MockSongStorage)(nil).SaveSong), ctx, song)
}

// UpdateSong mocks base method.
func (m *MockSongStorage) UpdateSong(ctx context.Context, song entities.Song) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSong", ctx, song)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSong indicates an expected call of UpdateSong.
func (mr *MockSongStorageMockRecorder) UpdateSong(ctx, song interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSong", reflect.TypeOf((*MockSongStorage)(nil).UpdateSong), ctx, song)
}