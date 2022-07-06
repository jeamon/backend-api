package repository

import (
	"context"

	"github.com/jeamon/sample-rest-api/pkg/domain"
	"github.com/jeamon/sample-rest-api/pkg/infrastructure/mockdb"
	"go.uber.org/zap"
)

type MockDBScanInfosRepository struct {
	logger *zap.Logger
	mock   *mockdb.Handler
	dbname string
}

// NewMongoScanInfosRepository provides an instance of MongoScanInfosRepository structure.
func NewMockDBScanInfosRepository(logger *zap.Logger, h *mockdb.Handler, dbname string) *MockDBScanInfosRepository {
	return &MockDBScanInfosRepository{
		logger: logger,
		mock:   h,
		dbname: dbname,
	}
}

func (repo MockDBScanInfosRepository) Save(ctx context.Context, s domain.ScanInfos) (string, error) {
	return "", nil
}

func (repo MockDBScanInfosRepository) FindByID(ctx context.Context, id string) (domain.ScanInfos, error) {
	return domain.ScanInfos{}, nil
}

func (repo MockDBScanInfosRepository) FindAll(ctx context.Context) ([]domain.ScanInfos, error) {
	return []domain.ScanInfos{}, nil
}

func (repo MockDBScanInfosRepository) UpdateByID(ctx context.Context, id string, s domain.ScanInfos) error {
	return nil
}

func (repo MockDBScanInfosRepository) DeleteByID(ctx context.Context, id string) error {
	return nil
}
