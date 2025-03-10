package router

import (
	"github.com/labstack/echo/v4"
	"github.com/notblessy/ekspresi-core/model"
	"gorm.io/gorm"
)

type httpService struct {
	db                 *gorm.DB
	userRepo           model.UserRepository
	membershipPlanRepo model.MembershipPlanRepository
	membershipRepo     model.MembershipRepository
	portfolioRepo      model.PortfolioRepository
}

func NewHTTPService() *httpService {
	return &httpService{}
}

func (h *httpService) RegisterPostgres(db *gorm.DB) {
	h.db = db
}

func (h *httpService) RegisterUserRepository(repo model.UserRepository) {
	h.userRepo = repo
}

func (h *httpService) RegisterMembershipPlanRepository(repo model.MembershipPlanRepository) {
	h.membershipPlanRepo = repo
}

func (h *httpService) RegisterMembershipRepository(repo model.MembershipRepository) {
	h.membershipRepo = repo
}

func (h *httpService) RegisterPortfolioRepository(repo model.PortfolioRepository) {
	h.portfolioRepo = repo
}

func (h *httpService) Router(e *echo.Echo) {
	e.GET("/ping", h.ping)
	e.GET("/health", h.health)

	v1 := e.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.POST("/login/google", h.loginWithGoogleHandler)

	v1.Use(NewJWTMiddleware().ValidateJWT)
	users := v1.Group("/users")
	users.GET("/me", h.profileHandler)

	membershipPlans := v1.Group("/membership-plans")
	membershipPlans.POST("", h.createMembershipPlan)
	membershipPlans.GET("", h.findAllMembershipPlans)
	membershipPlans.GET("/:id", h.findMembershipPlanByID)
	membershipPlans.PUT("/:id", h.updateMembershipPlan)
	membershipPlans.DELETE("/:id", h.deleteMembershipPlan)

	portfolios := v1.Group("/portfolios")
	portfolios.PATCH("", h.patchPortfolioHandler)
}

func (h *httpService) ping(c echo.Context) error {
	return c.JSON(200, response{Data: "pong"})
}

func (h *httpService) health(c echo.Context) error {
	err := h.db.Raw("SELECT 1").Error
	if err != nil {
		return c.JSON(500, response{Message: err.Error()})
	}

	return c.JSON(200, "OK")
}
