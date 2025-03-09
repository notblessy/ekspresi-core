package model

import (
	"context"
	"time"

	"github.com/notblessy/ekspresi-core/utils/nuller"
	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

type MembershipRepository interface {
	FindAll(ctx context.Context, query MembershipQueryInput) ([]Membership, int64, error)
	FindByID(ctx context.Context, id string) (Membership, error)
	Create(ctx context.Context, input Membership) error
	Update(ctx context.Context, id string, input Membership) error
	Delete(ctx context.Context, id string) error
}

type Membership struct {
	ID                         string          `json:"id"`
	UserID                     string          `json:"user_id"`
	MembershipPlanID           string          `json:"membership_plan_id"`
	Status                     string          `json:"status"`
	StartDate                  time.Time       `json:"start_date"`
	EndDate                    nuller.NullTime `json:"end_date"`
	StripeSubscriptionID       string          `json:"stripe_subscription_id"`
	StripeSubscriptionInterval string          `json:"stripe_subscription_interval"`
	CreatedAt                  time.Time       `json:"created_at"`
	UpdatedAt                  time.Time       `json:"updated_at"`
}

type MembershipInput struct {
	UserID                     string          `json:"user_id" validate:"required"`
	MembershipPlanID           string          `json:"membership_plan_id" validate:"required"`
	Status                     string          `json:"status" validate:"required"`
	StartDate                  time.Time       `json:"start_date" validate:"required"`
	EndDate                    nuller.NullTime `json:"end_date"`
	StripeSubscriptionID       string          `json:"stripe_subscription_id" validate:"required"`
	StripeSubscriptionInterval string          `json:"stripe_subscription_interval" validate:"required"`
}

func (input MembershipInput) ToMembership(id string) Membership {
	if id == "" {
		id = ulid.Make().String()
	}

	return Membership{
		ID:                         id,
		UserID:                     input.UserID,
		MembershipPlanID:           input.MembershipPlanID,
		Status:                     input.Status,
		StartDate:                  input.StartDate,
		EndDate:                    input.EndDate,
		StripeSubscriptionID:       input.StripeSubscriptionID,
		StripeSubscriptionInterval: input.StripeSubscriptionInterval,
	}
}

type MembershipQueryInput struct {
	UserID           string `query:"user_id"`
	MembershipPlanID string `query:"membership_plan_id"`
	Status           string `query:"status"`
	PaginatedRequest
}

type MembershipPlanRepository interface {
	FindAll(ctx context.Context, query MembershipPlanQueryInput) ([]MembershipPlan, int64, error)
	FindByID(ctx context.Context, id string) (MembershipPlan, error)
	Create(ctx context.Context, input MembershipPlan) error
	Update(ctx context.Context, id string, input MembershipPlan) error
	Delete(ctx context.Context, id string) error
}

type MembershipPlan struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Price             decimal.Decimal `json:"price"`
	BillingCycle      string          `json:"billing_cycle"`
	Features          []string        `json:"features"`
	IsPopular         bool            `json:"is_popular"`
	MaxFolders        int             `json:"max_folders"`
	CustomDomain      bool            `json:"custom_domain"`
	AdvancedAnalytics bool            `json:"advanced_analytics"`
	StripeProductID   string          `json:"stripe_product_id"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type MembershipPlanInput struct {
	Name              string          `json:"name" validate:"required"`
	Price             decimal.Decimal `json:"price" validate:"required"`
	BillingCycle      string          `json:"billing_cycle" validate:"required"`
	Features          []string        `json:"features"`
	IsPopular         bool            `json:"is_popular"`
	MaxFolders        int             `json:"max_folders"`
	CustomDomain      bool            `json:"custom_domain"`
	AdvancedAnalytics bool            `json:"advanced_analytics"`
	StripeProductID   string          `json:"stripe_product_id" validate:"required"`
}

func (input MembershipPlanInput) ToMembershipPlan(id string) MembershipPlan {
	if id == "" {
		id = ulid.Make().String()
	}

	return MembershipPlan{
		ID:                id,
		Name:              input.Name,
		Price:             input.Price,
		BillingCycle:      input.BillingCycle,
		Features:          input.Features,
		IsPopular:         input.IsPopular,
		MaxFolders:        input.MaxFolders,
		CustomDomain:      input.CustomDomain,
		AdvancedAnalytics: input.AdvancedAnalytics,
		StripeProductID:   input.StripeProductID,
	}
}

type MembershipPlanQueryInput struct {
	Keyword string `query:"keyword"`
	PaginatedRequest
}
