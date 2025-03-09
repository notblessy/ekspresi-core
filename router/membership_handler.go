package router

import (
	"github.com/labstack/echo/v4"
	"github.com/notblessy/ekspresi-core/model"
	"github.com/sirupsen/logrus"
)

func (h *httpService) findAllMemberships(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	var query model.MembershipQueryInput

	if err := c.Bind(&query); err != nil {
		logger.WithError(err).Error("failed to bind query")
		return c.JSON(400, response{Message: "invalid query"})
	}

	plans, total, err := h.membershipRepo.FindAll(c.Request().Context(), query)
	if err != nil {
		logger.WithError(err).Error("failed to find all membership plans")
		return c.JSON(500, response{Message: err.Error()})
	}

	return c.JSON(200, response{Success: true, Data: withPaging(plans, total, query.PageOrDefault(), query.SizeOrDefault())})
}

func (h *httpService) findMembershipByID(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	id := c.Param("id")

	plan, err := h.membershipRepo.FindByID(c.Request().Context(), id)
	if err != nil {
		logger.WithError(err).Error("failed to find membership plan by id")
		return c.JSON(500, response{Message: err.Error()})
	}

	return c.JSON(200, response{Success: true, Data: plan})
}

func (h *httpService) createMembership(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	var input model.MembershipInput

	if err := c.Bind(&input); err != nil {
		logger.WithError(err).Error("failed to bind input")
		return c.JSON(400, response{Message: "invalid input"})
	}

	plan := input.ToMembership("")

	if err := h.membershipRepo.Create(c.Request().Context(), plan); err != nil {
		logger.WithError(err).Error("failed to create membership plan")
		return c.JSON(500, response{Message: err.Error()})
	}

	return c.JSON(201, response{Success: true, Data: plan})
}

func (h *httpService) updateMembership(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	id := c.Param("id")

	var input model.MembershipInput

	if err := c.Bind(&input); err != nil {
		logger.WithError(err).Error("failed to bind input")
		return c.JSON(400, response{Message: "invalid input"})
	}

	plan := input.ToMembership(id)

	if err := h.membershipRepo.Update(c.Request().Context(), id, plan); err != nil {
		logger.WithError(err).Error("failed to update membership plan")
		return c.JSON(500, response{Message: err.Error()})
	}

	return c.JSON(200, response{Success: true, Data: plan})
}

func (h *httpService) deleteMembership(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	id := c.Param("id")

	if err := h.membershipRepo.Delete(c.Request().Context(), id); err != nil {
		logger.WithError(err).Error("failed to delete membership plan")
		return c.JSON(500, response{Message: err.Error()})
	}

	return c.JSON(200, response{Success: true})
}

func (h *httpService) findAllMembershipPlans(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	var query model.MembershipPlanQueryInput
	if err := c.Bind(&query); err != nil {
		logger.WithError(err).Error("failed to bind query")
		return c.JSON(400, response{Message: "invalid query"})
	}

	session, err := authSession(c)
	if err != nil {
		logger.WithError(err).Error("failed to get auth session")
		return c.JSON(500, response{Message: err.Error()})
	}

	if session.Role != model.RoleAdmin {
		return c.JSON(403, response{Message: "forbidden"})
	}

	plans, total, err := h.membershipPlanRepo.FindAll(c.Request().Context(), query)
	if err != nil {
		logger.WithError(err).Error("failed to find all membership plans")
		return c.JSON(500, response{Message: err.Error()})
	}

	return c.JSON(200, response{Success: true, Data: withPaging(plans, total, query.PageOrDefault(), query.SizeOrDefault())})
}

func (h *httpService) findMembershipPlanByID(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	id := c.Param("id")

	session, err := authSession(c)
	if err != nil {
		logger.WithError(err).Error("failed to get auth session")
		return c.JSON(500, response{Message: err.Error()})
	}

	if session.Role != model.RoleAdmin {
		return c.JSON(403, response{Message: "forbidden"})
	}

	plan, err := h.membershipPlanRepo.FindByID(c.Request().Context(), id)
	if err != nil {
		logger.WithError(err).Error("failed to find membership plan by id")
		return c.JSON(500, response{Message: err.Error()})
	}

	return c.JSON(200, response{Success: true, Data: plan})
}

func (h *httpService) createMembershipPlan(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	var input model.MembershipPlanInput

	if err := c.Bind(&input); err != nil {
		logger.WithError(err).Error("failed to bind input")
		return c.JSON(400, response{Message: "invalid input"})
	}

	session, err := authSession(c)
	if err != nil {
		logger.WithError(err).Error("failed to get auth session")
		return c.JSON(500, response{Message: err.Error()})
	}

	if session.Role != model.RoleAdmin {
		return c.JSON(403, response{Message: "forbidden"})
	}

	plan := input.ToMembershipPlan("")

	if err := h.membershipPlanRepo.Create(c.Request().Context(), plan); err != nil {
		logger.WithError(err).Error("failed to create membership plan")
		return c.JSON(500, response{Message: err.Error()})
	}

	return c.JSON(201, response{Success: true, Data: plan})
}

func (h *httpService) updateMembershipPlan(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	id := c.Param("id")

	var input model.MembershipPlanInput

	if err := c.Bind(&input); err != nil {
		logger.WithError(err).Error("failed to bind input")
		return c.JSON(400, response{Message: "invalid input"})
	}

	session, err := authSession(c)
	if err != nil {
		logger.WithError(err).Error("failed to get auth session")
		return c.JSON(500, response{Message: err.Error()})
	}

	if session.Role != model.RoleAdmin {
		return c.JSON(403, response{Message: "forbidden"})
	}

	plan := input.ToMembershipPlan(id)

	if err := h.membershipPlanRepo.Update(c.Request().Context(), id, plan); err != nil {
		logger.WithError(err).Error("failed to update membership plan")
		return c.JSON(500, response{Message: err.Error()})
	}

	return c.JSON(200, response{Success: true, Data: plan})
}

func (h *httpService) deleteMembershipPlan(c echo.Context) error {
	logger := logrus.WithContext(c.Request().Context())

	id := c.Param("id")

	session, err := authSession(c)
	if err != nil {
		logger.WithError(err).Error("failed to get auth session")
		return c.JSON(500, response{Message: err.Error()})
	}

	if session.Role != model.RoleAdmin {
		return c.JSON(403, response{Message: "forbidden"})
	}

	if err := h.membershipPlanRepo.Delete(c.Request().Context(), id); err != nil {
		logger.WithError(err).Error("failed to delete membership plan")
		return c.JSON(500, response{Message: err.Error()})
	}

	return c.JSON(200, response{Success: true})
}
