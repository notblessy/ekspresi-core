package router

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/notblessy/ekspresi-core/model"
	"github.com/sirupsen/logrus"
)

func (h *httpService) uploadPhotoHandler(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	file, err := c.FormFile("file")
	if err != nil {
		logger.WithError(err).Error("failed to get file")
		return c.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}

	var photo model.Photo

	if err := c.Bind(&photo); err != nil {
		logger.WithError(err).Error("failed to bind request")
		return c.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}

	src, err := file.Open()
	if err != nil {
		logger.WithError(err).Error("failed to open file")
		return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
	}
	defer src.Close()

	session, err := authSession(c)
	if err != nil {
		logger.WithError(err).Error("failed to get session")
		return c.JSON(http.StatusUnauthorized, response{Message: err.Error()})
	}

	path := fmt.Sprintf("%s/%s/%s", os.Getenv("UPLOADER_BASE_PATH"), "portfolios", session.ID)

	url, publicID, err := h.uploaderRepo.Upload(c.Request().Context(), src, path)
	if err != nil {
		logger.WithError(err).Error("failed to upload file")
		return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
	}

	newPhoto := model.Photo{
		ID:        photo.ID,
		Src:       url,
		PublicID:  publicID,
		Alt:       photo.Alt,
		Caption:   photo.Caption,
		SortIndex: photo.SortIndex,
		CreatedAt: time.Now(),
	}

	err = h.uploaderRepo.SavePhoto(c.Request().Context(), newPhoto)
	if err != nil {
		logger.WithError(err).Error("failed to save photo")
		return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, response{
		Success: true,
		Data:    newPhoto,
	})
}

func (h *httpService) bulkRemovePhotosHandler(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	var req model.DeleteRequest

	if err := c.Bind(&req); err != nil {
		logger.WithError(err).Error("failed to bind request")
		return c.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}

	err := h.uploaderRepo.DeleteByPublicIDs(c.Request().Context(), req.PublicIDs)
	if err != nil {
		logger.WithError(err).Error("failed to delete files")
		return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, response{Success: true})
}

func (h *httpService) flushHandler(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	err := h.uploaderRepo.Flush(c.Request().Context())
	if err != nil {
		logger.WithError(err).Error("failed to flush files")
		return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, response{Success: true})
}
