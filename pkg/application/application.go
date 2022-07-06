package application

import (
	"context"

	"github.com/jeamon/sample-rest-api/pkg/domain"
	"github.com/jeamon/sample-rest-api/pkg/infrastructure/config"
	"go.uber.org/zap"
)

// ScanInfosUsecase is handling the scan informations business logic.
type ScanInfosUsecase struct {
	Logger        *zap.Logger
	ConfigData    *config.Config
	scanInfosRepo domain.ScanInfosRepository
}

// NewScanInfosUsecase initialises and returns a new use case for scan infos use case.
func NewScanInfosUsecase(logger *zap.Logger, configData *config.Config, repo domain.ScanInfosRepository) *ScanInfosUsecase {
	uc := &ScanInfosUsecase{
		Logger:        logger,
		ConfigData:    configData,
		scanInfosRepo: repo,
	}

	return uc
}

func (uc *ScanInfosUsecase) Store(ctx context.Context, req domain.StoreScanInfosRequest) (string, error) {
	return uc.scanInfosRepo.Save(ctx, req.ToScanInfos())
}

func (uc *ScanInfosUsecase) Get(ctx context.Context, id string) (domain.ScanInfos, error) {
	return uc.scanInfosRepo.FindByID(ctx, id)
}

func (uc *ScanInfosUsecase) GetAll(ctx context.Context) ([]domain.ScanInfos, error) {
	return uc.scanInfosRepo.FindAll(ctx)
}

func (uc *ScanInfosUsecase) Delete(ctx context.Context, id string) error {
	return uc.scanInfosRepo.DeleteByID(ctx, id)
}

func (uc *ScanInfosUsecase) Update(ctx context.Context, infos domain.ScanInfos) error {
	return uc.scanInfosRepo.UpdateByID(ctx, infos.ID, infos)
}
