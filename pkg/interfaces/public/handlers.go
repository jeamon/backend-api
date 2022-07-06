package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jeamon/sample-rest-api/pkg/domain"
	"go.uber.org/zap"
)

// RequestLoggerMiddleware logs the metadat and the body of the requests.
func (w *ScanInfosService) RequestLoggerMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		requestID := generateRequestID()
		buf, _ := ioutil.ReadAll(c.Request.Body)
		data, err := readBody(ioutil.NopCloser(bytes.NewBuffer(buf)))

		w.logger.Info(
			"received request on:",
			zap.String("url", c.Request.URL.Path),
			zap.String("authorization", c.GetHeader("Authorization")),
			zap.String("method", c.Request.Method),
			zap.String("requestid", requestID),
			zap.String("ip", getIP(c.Request)),
			zap.String("agent", c.Request.UserAgent()),
			zap.String("referer", c.Request.Referer()),
			zap.String("body", data),
			zap.Error(err),
		)
		c.Set("x-requestid", requestID)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		c.Next()
	}
}

func (w *ScanInfosService) NotFoundHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		w.logger.Error("route does not exist", zap.String("requestid", c.GetString("x-requestid")))
		c.JSON(http.StatusNotImplemented, errResponse{
			RequestID:        c.GetString("x-requestid"),
			Message:          "invalid request. make sure to use the exact endpoint.",
			DeveloperMessage: "endpoint called with that method does not exist.",
		})
	}
}

// StoreScanInfosHandler ...
// @Summary store a scan details and get back the infos
// @Description save a scan information and return the full details
// @Tags ScanInfos
// @Accept  json
// @Produce  json
// @Param scanInfosStore body ScanInfos true "Store Scan Infos"
// @Success 200 {object} genericResponse
// @Success 500 {object} errResponse
// @Router /api/v1/scaninfos [post]
func (w *ScanInfosService) StoreScanInfosHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		var req domain.StoreScanInfosRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			w.logger.Error("unable to bind store request input", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
			c.JSON(http.StatusBadRequest, errResponse{
				RequestID:        c.GetString("x-requestid"),
				Message:          "invalid request. make sure to provide expected data format",
				DeveloperMessage: err.Error(),
			})
			return
		}

		id, err := w.application.Store(c, req)
		if err != nil {
			w.logger.Error("unable to store scan infos", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
			c.JSON(http.StatusInternalServerError, errResponse{
				RequestID:        c.GetString("x-requestid"),
				Message:          "an error occurred while storing the scan infos",
				DeveloperMessage: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, genericResponse{
			RequestID:   c.GetString("x-requestid"),
			Message:     "scan infos saved successfully",
			ScanInfosID: id,
		})
	}
}

// GetScanInfosHandler ...
// @Summary get a scan infos
// @Description get a scan information by its id
// @Tags ScanInfos
// @Accept  json
// @Produce  json
// @Param id path string true "ID string"
// @Success 200 {object} ScanInfos
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /api/v1/scaninfos/{id} [get]
func (w *ScanInfosService) GetScanInfosHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		if _, err := uuid.FromString(id); err != nil {
			w.logger.Error("bad request. invalid id", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
			c.JSON(400, errResponse{
				RequestID:        c.GetString("x-requestid"),
				Message:          "bad request. cannot get scan infos.",
				DeveloperMessage: "expect non empty id as parameter of the scan infosto fetch.",
			})
			return
		}

		infos, err := w.application.Get(c, id)
		if err != nil {
			w.logger.Error("unable to get scan infos", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
			c.JSON(500, errResponse{
				RequestID:        c.GetString("x-requestid"),
				Message:          "an error occurred while fetching the scan infos",
				DeveloperMessage: err.Error(),
			})
			return
		}
		c.JSON(200, infos)
	}
}

// GetAllScanInfosHandler ...
// @Summary get all scan infos
// @Description get a list of all stored scan information
// @Tags ScanInfos
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllScanInfosResponse
// @Failure 500 {object} errResponse
// @Router /api/v1/scaninfos [get]
func (w *ScanInfosService) GetAllScanInfosHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		infos, err := w.application.GetAll(c)
		if err != nil {
			w.logger.Error("unable to fetch all scan infos", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
			c.JSON(500, errResponse{
				RequestID:        c.GetString("x-requestid"),
				Message:          "an error occurred while fetching all scan infos",
				DeveloperMessage: err.Error(),
			})
			return
		}

		c.JSON(200, getAllScanInfosResponse{
			RequestID: c.GetString("x-requestid"),
			Message:   "all scan infos fetched successfully",
			Infos:     infos,
		})
	}
}

// UpdateScanInfosHandler ...
// @Summary update an existing scan details
// @Description update a scan information based on its ID
// @Tags ScanInfos
// @Accept  json
// @Produce  json
// @Param scanInfosStore body ScanInfos true "Update Scan Infos"
// @Success 200 {object} genericResponse
// @Success 500 {object} errResponse
// @Router /api/v1/scaninfos [put]
func (w *ScanInfosService) UpdateScanInfosHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		var req domain.ScanInfos
		if err := c.ShouldBindJSON(&req); err != nil {
			w.logger.Error("unable to bind update request input", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
			c.JSON(http.StatusBadRequest, errResponse{
				RequestID:        c.GetString("x-requestid"),
				Message:          "invalid request. make sure to provide expected data format",
				DeveloperMessage: err.Error(),
			})
			return
		}

		if err := w.application.Update(c, req); err != nil {
			w.logger.Error("unable to store scan infos", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
			c.JSON(http.StatusInternalServerError, errResponse{
				RequestID:        c.GetString("x-requestid"),
				Message:          "an error occurred while updating the scan infos",
				DeveloperMessage: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, genericResponse{
			RequestID:   c.GetString("x-requestid"),
			Message:     "scan infos updated successfully",
			ScanInfosID: req.ID,
		})
	}
}

// DeleteScanInfosHandler ...
// @Summary get a scan infos
// @Description get a scan information by its id
// @Tags ScanInfos
// @Accept  json
// @Produce  json
// @Param id path string true "ID string"
// @Success 200 {object} genericResponse
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /api/v1/scaninfos/{id} [delete]
func (w *ScanInfosService) DeleteScanInfosHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		fmt.Println(id)
		if _, err := uuid.FromString(id); err != nil {
			w.logger.Error("bad request. invalid id", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
			c.JSON(http.StatusBadRequest, errResponse{
				RequestID:        c.GetString("x-requestid"),
				Message:          "bad request. cannot delete scan infos.",
				DeveloperMessage: "expect non empty id as parameter of the scan infos to delete.",
			})
			return
		}

		err := w.application.Delete(c, id)
		if err != nil {
			w.logger.Error("unable to delete scan infos", zap.String("requestid", c.GetString("x-requestid")), zap.Error(err))
			c.JSON(http.StatusInternalServerError, errResponse{
				RequestID:        c.GetString("x-requestid"),
				Message:          "an error occurred while deleting the scan infos",
				DeveloperMessage: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, genericResponse{
			RequestID:   c.GetString("x-requestid"),
			Message:     "scan infos deleted successfully",
			ScanInfosID: id,
		})
	}
}
