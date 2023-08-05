package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jeamon/backend-api/pkg/domain"
	"github.com/jeamon/backend-api/pkg/infrastructure/postgres"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type PostgresScanInfosRepository struct {
	logger *zap.Logger
	pg     *postgres.Handler
}

// NewPostgresScanInfosRepository ...
func NewPostgresScanInfosRepository(logger *zap.Logger, h *postgres.Handler) *PostgresScanInfosRepository {
	return &PostgresScanInfosRepository{
		logger: logger,
		pg:     h,
	}
}

// Save will create a new scan infos and not update existing one.
func (repo PostgresScanInfosRepository) Save(ctx context.Context, s domain.ScanInfos) (string, error) {
	var id string
	sql, args, err := psql.Insert("data.scan_infos").SetMap(
		map[string]interface{}{
			"company_id":     s.CompanyID,
			"client_id":      s.ClientID,
			"username":       s.Username,
			"repository_url": s.RepositoryURL,
			"commit_id":      s.CommitID,
			"tag_id":         s.TagID,
			"results":        s.Results,
			"started_at":     s.StartedAt,
			"completed_at":   s.CompletedAt,
			"sent_at":        s.SentAt,
			"created_at":     s.CreatedAt,
			"updated_at":     s.UpdatedAt,
			"metadata":       s.Metadata,
		}).Suffix("RETURNING id").ToSql()
	if err != nil {
		return id, errors.Wrapf(err, "cannot save scan infos. failed to build query statement")
	}

	err = repo.pg.PGx.QueryRow(ctx, sql, args...).Scan(&id)
	return id, errors.Wrapf(err, "could not save scan infos")
}

func (repo PostgresScanInfosRepository) FindByID(ctx context.Context, id string) (domain.ScanInfos, error) {
	var s domain.ScanInfos
	sql, args, err := psql.Select("*").From("data.scan_infos").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return s, errors.Wrap(err, "cannot get scan infos")
	}

	err = pgxscan.Get(ctx, repo.pg.PGx, &s, sql, args...)
	return s, errors.Wrapf(err, "could not find scan infos with ID: %s", id)
}

func (repo PostgresScanInfosRepository) FindAll(ctx context.Context) ([]domain.ScanInfos, error) {
	sql, args, err := psql.Select("*").From("data.scan_infos").ToSql()
	if err != nil {
		return nil, err
	}

	res := []domain.ScanInfos{}
	err = pgxscan.Select(ctx, repo.pg.PGx, &res, sql, args...)
	return res, errors.Wrapf(err, "could not find all scan infos")
}

// UpdateByID updates a scan infos record by its ID.
func (repo PostgresScanInfosRepository) UpdateByID(ctx context.Context, id string, s domain.ScanInfos) error {
	sql, args, err := psql.Update("data.scan_infos").SetMap(
		map[string]interface{}{
			"company_id":     s.CompanyID,
			"client_id":      s.ClientID,
			"username":       s.Username,
			"repository_url": s.RepositoryURL,
			"commit_id":      s.CommitID,
			"tag_id":         s.TagID,
			"results":        s.Results,
			"started_at":     s.StartedAt,
			"completed_at":   s.CompletedAt,
			"sent_at":        s.SentAt,
			"updated_at":     time.Now().UTC(),
			"error":          s.Error,
			"metadata":       s.Metadata,
		}).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return errors.Wrapf(err, "cannot update scan infos with ID: %s", id)
	}

	_, err = repo.pg.PGx.Exec(ctx, sql, args...)
	return errors.Wrapf(err, "cannot update scan infos with ID: %s", id)
}

// DeleteByID deletes a scan infos record by its ID.
func (repo PostgresScanInfosRepository) DeleteByID(ctx context.Context, id string) error {
	sql, args, err := psql.Delete("data.scan_infos").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return errors.Wrapf(err, "cannot delete scan infos record with ID: %s", id)
	}

	_, err = repo.pg.PGx.Exec(ctx, sql, args...)
	return errors.Wrapf(err, "could not delete scan infos with ID: %s", id)
}
