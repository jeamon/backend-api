package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"github.com/jeamon/backend-api/pkg/application"
	"github.com/jeamon/backend-api/pkg/domain"
	"github.com/jeamon/backend-api/pkg/infrastructure/config"
	"github.com/jeamon/backend-api/pkg/infrastructure/mockdb"
	"github.com/jeamon/backend-api/pkg/interfaces/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	testLogger                = zap.L()
	testConfigData            = &config.Config{}
	testScanInfosID           = "7aec1a3e-f22d-11ec-a1c2-37e6aab6bd2c"
	testStoreScanInfosRequest = domain.StoreScanInfosRequest{
		CompanyID:     "0",
		Username:      "jeamon",
		ClientID:      "v1.0.0",
		RepositoryURL: "https://github.com/jeamon/sample-rest-api",
		CommitID:      "d7b8ff1412ebfcde26f9ddfdf9608d1525647958",
		TagID:         "v1.0.0",
		Results:       []string{"found something"},
		StartedAt:     1655903720,
		CompletedAt:   1655903723,
		SentAt:        1655903725,
		Error:         "got an x exception during execution",
		Metadata: map[string]interface{}{
			"os":        "linux",
			"languages": []string{"go", "bash", "html"},
			"arch":      "amd64",
		},
	}
	testScanInfos = domain.ScanInfos{
		ID:            "5c9828c6-f7f4-11ec-aa0a-d3f9ac7b9396",
		CompanyID:     "0",
		Username:      "jeamon",
		ClientID:      "v1.0.0",
		RepositoryURL: "https://github.com/jeamon/sample-rest-api",
		CommitID:      "d7b8ff1412ebfcde26f9ddfdf9608d1525647958",
		TagID:         "v1.0.0",
		Results:       []string{"found something"},
		StartedAt:     1655903720,
		CompletedAt:   1655903723,
		SentAt:        1655903725,
		Error:         "got an x exception during execution",
		Metadata: map[string]interface{}{
			"os":        "linux",
			"languages": []string{"go", "bash", "html"},
			"arch":      "amd64",
		},
	}
)

func setupTestServer() *httptest.Server {
	mockDB := mockdb.Config{}
	mockdbHandler, _ := mockDB.ConnectAndMigrate(testLogger)
	testRepo := repository.NewMockDBScanInfosRepository(testLogger, mockdbHandler, "mockdb")
	scanInfosUc := application.NewScanInfosUsecase(testLogger, testConfigData, testRepo)
	service := New(testLogger, scanInfosUc)
	gin.SetMode(gin.TestMode)
	ts := httptest.NewServer(service.Router(gin.Default()))
	return ts
}

func TestGetScanInfosHandler(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	t.Run("GetScanInfos endpoint tests", func(t *testing.T) {
		t.Run("should pass: valid param <id> value", func(t *testing.T) {
			res, err := req.Get(ts.URL + "/api/v1/scaninfos/" + testScanInfosID)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.Response().StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", res.Response().Header.Get("Content-Type"))
			assert.NotEmpty(t, res.Bytes())
		})

		t.Run("should fail: invalid param <id> value", func(t *testing.T) {
			res, err := req.Get(ts.URL + "/api/v1/scaninfos/7aec1a3e")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, res.Response().StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", res.Response().Header.Get("Content-Type"))
			assert.NotEmpty(t, res.Bytes())
		})
	})
}

func TestGetAllScanInfosHandler(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	t.Run("GetAllScanInfos endpoint tests", func(t *testing.T) {
		t.Run("should pass: endpoint without slash", func(t *testing.T) {
			res, err := req.Get(ts.URL + "/api/v1/scaninfos")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.Response().StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", res.Response().Header.Get("Content-Type"))
			assert.NotEmpty(t, res.Bytes())
		})

		t.Run("should pass: endpoint with slash", func(t *testing.T) {
			res, err := req.Get(ts.URL + "/api/v1/scaninfos/")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.Response().StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", res.Response().Header.Get("Content-Type"))
			assert.NotEmpty(t, res.Bytes())
		})
	})
}

func TestStoreScanInfosHandler(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	t.Run("StoreScanInfos endpoint tests", func(t *testing.T) {
		t.Run("should pass: valid request body", func(t *testing.T) {
			res, err := req.Post(ts.URL+"/api/v1/scaninfos", req.BodyJSON(&testStoreScanInfosRequest))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.Response().StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", res.Response().Header.Get("Content-Type"))
			assert.NotEmpty(t, res.Bytes())
		})

		t.Run("should fail: request with empty body", func(t *testing.T) {
			res, err := req.Post(ts.URL + "/api/v1/scaninfos")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, res.Response().StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", res.Response().Header.Get("Content-Type"))
			assert.NotEmpty(t, res.Bytes())
		})

		t.Run("should fail: request with invalid json body", func(t *testing.T) {
			res, err := req.Post(ts.URL+"/api/v1/scaninfos", req.BodyJSON(&domain.ScanInfos{}))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, res.Response().StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", res.Response().Header.Get("Content-Type"))
			assert.NotEmpty(t, res.Bytes())
		})
	})
}

func TestUpdateScanInfosHandler(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	t.Run("UpdateScanInfos endpoint tests", func(t *testing.T) {
		t.Run("should pass: valid request body", func(t *testing.T) {
			res, err := req.Put(ts.URL+"/api/v1/scaninfos", req.BodyJSON(&testScanInfos))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.Response().StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", res.Response().Header.Get("Content-Type"))
			assert.NotEmpty(t, res.Bytes())
		})

		t.Run("should fail: request with empty body", func(t *testing.T) {
			res, err := req.Put(ts.URL + "/api/v1/scaninfos")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, res.Response().StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", res.Response().Header.Get("Content-Type"))
			assert.NotEmpty(t, res.Bytes())
		})

		t.Run("should fail: request with invalid json request", func(t *testing.T) {
			res, err := req.Put(ts.URL+"/api/v1/scaninfos", req.BodyJSON(&domain.StoreScanInfosRequest{}))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, res.Response().StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", res.Response().Header.Get("Content-Type"))
			assert.NotEmpty(t, res.Bytes())
		})
	})
}

func TestDeleteScanInfosHandler(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	t.Run("DeleteScanInfos endpoint tests", func(t *testing.T) {
		t.Run("should pass: valid param <id> string", func(t *testing.T) {
			res, err := req.Delete(ts.URL + "/api/v1/scaninfos/" + testScanInfosID)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.Response().StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", res.Response().Header.Get("Content-Type"))
			assert.NotEmpty(t, res.Bytes())
		})

		t.Run("should fail: invalid param <id> value", func(t *testing.T) {
			res, err := req.Delete(ts.URL + "/api/v1/scaninfos/7aec1a3e")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, res.Response().StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", res.Response().Header.Get("Content-Type"))
			assert.NotEmpty(t, res.Bytes())
		})

		t.Run("should fail: no param <id> value", func(t *testing.T) {
			res, err := req.Delete(ts.URL + "/api/v1/scaninfos/")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotImplemented, res.Response().StatusCode)
			assert.Equal(t, "application/json; charset=utf-8", res.Response().Header.Get("Content-Type"))
			assert.NotEmpty(t, res.Bytes())
		})
	})
}
