package model

import (
	"context"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type UserRepository interface {
	Authenticate(ctx context.Context, code, requestOrigin string) (User, error)
	FindByID(ctx context.Context, id string) (User, error)
}

type User struct {
	ID        string         `json:"id"`
	Email     string         `json:"email"`
	Name      string         `json:"name"`
	Password  string         `json:"password,omitempty"`
	Picture   string         `json:"picture"`
	Role      string         `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (u *User) OmitPassword() {
	u.Password = ""
}

func (u *User) NewFreeMembership() Membership {
	trial := os.Getenv("MEMBERSHIP_PLAN_FREE")

	return Membership{
		UserID:                     u.ID,
		MembershipPlanID:           trial,
		Status:                     MembershipStatusActive,
		StartDate:                  time.Now(),
		StripeSubscriptionID:       trial,
		StripeSubscriptionInterval: trial,
		CreatedAt:                  time.Now(),
	}
}

func (u *User) NewInitialPortfolio() *Portfolio {
	return &Portfolio{
		ID:             ulid.Make().String(),
		UserID:         u.ID,
		Title:          u.Name,
		Columns:        3,
		Gap:            16,
		RoundedCorners: true,
		ShowCaptions:   true,
	}
}

func (u *User) NewDefaultProfile(portfolioID string) Profile {
	return Profile{
		PortfolioID: portfolioID,
		Name:        u.Name,
		Title:       DefaultTitle,
	}
}

type Auth struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type AuthRequest struct {
	Code          string `json:"code"`
	RequestOrigin string `json:"request_origin"`
}

type ChangeUsernameRequest struct {
	Username string `json:"username"`
}

type GoogleAuthInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
