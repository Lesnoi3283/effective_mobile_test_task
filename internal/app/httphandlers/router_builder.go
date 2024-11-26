package httphandlers

import (
	"go.uber.org/zap"
	"musiclib/internal/app/requiredinterfaces"
)

type handler struct {
	storage           requiredinterfaces.SongStorage
	logger            *zap.SugaredLogger
	extraDataProvider requiredinterfaces.ExtraDataProvider
}
