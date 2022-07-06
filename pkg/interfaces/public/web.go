package web

import (
	"github.com/jeamon/sample-rest-api/pkg/application"
	"go.uber.org/zap"
)

type ScanInfosService struct {
	logger      *zap.Logger
	application *application.ScanInfosUsecase
}

func New(logger *zap.Logger, application *application.ScanInfosUsecase) *ScanInfosService {
	return &ScanInfosService{logger: logger, application: application}
}
