package web

import (
	"github.com/jeamon/sample-rest-api/pkg/domain"
)

type errResponse struct {
	RequestID        string `json:"request_id"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developer_message"`
}

type getAllScanInfosResponse struct {
	RequestID string             `json:"request_id"`
	Message   string             `json:"message"`
	Infos     []domain.ScanInfos `json:"infos"`
}

type genericResponse struct {
	RequestID   string `json:"request_id"`
	Message     string `json:"message"`
	ScanInfosID string `json:"scan_infos_id"`
}
