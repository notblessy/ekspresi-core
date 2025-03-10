package router

import (
	"github.com/labstack/echo/v4"
	"github.com/notblessy/ekspresi-core/model"
	"github.com/sirupsen/logrus"
)

func (h *httpService) patchPortfolioHandler(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	var input model.PortfolioType

	if err := c.Bind(&input); err != nil {
		logger.WithError(err).Error("failed to bind input")
		return c.JSON(400, response{Message: "invalid input"})
	}

	session, err := authSession(c)
	if err != nil {
		logger.WithError(err).Error("failed to get session")
		return c.JSON(401, response{Message: "unauthorized"})
	}

	if session.ID == "" {
		logger.Error("session id is empty")
		return c.JSON(401, response{Message: "unauthorized"})
	}

	err = h.portfolioRepo.Patch(c.Request().Context(), input)
	if err != nil {
		logger.WithError(err).Error("failed to patch portfolio")
		return c.JSON(500, response{Message: err.Error()})
	}

	return c.JSON(200, response{Success: true})
}
