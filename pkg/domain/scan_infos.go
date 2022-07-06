package domain

import "time"

type ScanInfos struct {
	ID string `db:"id" json:"id" bson:"_id" binding:"required"`

	CompanyID string `db:"company_id" json:"company_id" bson:"company_id" binding:"required"`
	Username  string `db:"username" json:"username" bson:"username" binding:"required"`

	ClientID string `db:"client_id" json:"client_id" bson:"client_id" binding:"required"`

	RepositoryURL string `db:"repository_url" json:"repository_url" bson:"repository_url" binding:"required"`
	CommitID      string `db:"commit_id" json:"commit_id" bson:"commit_id" binding:"required"`
	TagID         string `db:"tag_id" json:"tag_id" bson:"tag_id" binding:"required"`

	Results []string `db:"results" json:"results" bson:"results" binding:"required"`

	StartedAt   int64     `db:"started_at" json:"started_at" bson:"started_at" binding:"required"`
	CompletedAt int64     `db:"completed_at" json:"completed_at" bson:"completed_at" binding:"required"`
	SentAt      int64     `db:"sent_at" json:"sent_at" bson:"sent_at" binding:"required"`
	CreatedAt   time.Time `db:"created_at" json:"-" bson:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"-" bson:"updated_at"`

	Error    string                 `db:"error" json:"error" bson:"error"`
	Metadata map[string]interface{} `db:"metadata" json:"metadata" bson:"metadata" binding:"required"`
}

type StoreScanInfosRequest struct {
	CompanyID string `json:"company_id" binding:"required"`
	Username  string `json:"username" binding:"required"`

	ClientID string `json:"client_id" binding:"required"`

	RepositoryURL string `json:"repository_url" binding:"required"`
	CommitID      string `json:"commit_id" binding:"required"`
	TagID         string `json:"tag_id" binding:"required"`

	Results []string `json:"results" binding:"required"`

	StartedAt   int64 `json:"started_at" binding:"required"`
	CompletedAt int64 `json:"completed_at" binding:"required"`
	SentAt      int64 `json:"sent_at" binding:"required"`

	Error    string                 `json:"error"`
	Metadata map[string]interface{} `json:"metadata" binding:"required"`
}

func (r StoreScanInfosRequest) ToScanInfos() (s ScanInfos) {
	now := time.Now().UTC()
	return ScanInfos{
		CreatedAt:     now,
		UpdatedAt:     now,
		CompanyID:     r.CompanyID,
		Username:      r.Username,
		ClientID:      r.ClientID,
		RepositoryURL: r.RepositoryURL,
		CommitID:      r.CommitID,
		TagID:         r.TagID,
		Results:       r.Results,
		StartedAt:     r.StartedAt,
		CompletedAt:   r.CompletedAt,
		SentAt:        r.SentAt,
		Error:         r.Error,
		Metadata:      r.Metadata,
	}
}
