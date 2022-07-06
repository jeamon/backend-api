package domain

import "context"

// Repository ...
type ScanInfosRepository interface {
	Save(ctx context.Context, scanInfos ScanInfos) (string, error)
	FindByID(ctx context.Context, id string) (ScanInfos, error)
	UpdateByID(ctx context.Context, id string, scanInfos ScanInfos) error
	DeleteByID(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]ScanInfos, error)
}
