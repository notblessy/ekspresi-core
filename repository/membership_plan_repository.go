package repository

import (
	"context"

	"github.com/notblessy/ekspresi-core/model"
	"github.com/notblessy/ekspresi-core/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type membershipPlanRepository struct {
	db *gorm.DB
}

func NewMembershipPlanRepository(db *gorm.DB) model.MembershipPlanRepository {
	return &membershipPlanRepository{db}
}

func (r *membershipPlanRepository) FindAll(ctx context.Context, query model.MembershipPlanQueryInput) ([]model.MembershipPlan, int64, error) {
	logger := logrus.WithField("query", utils.Dump(query))

	qb := r.db.Model(&model.MembershipPlan{}).WithContext(ctx)

	if query.Keyword != "" {
		qb = qb.Where("name ILIKE ?", "%"+query.Keyword+"%")
	}

	var total int64

	if err := qb.Count(&total).Error; err != nil {
		logger.WithError(err).Error("failed to count membership plans")
		return nil, 0, err
	}

	var plans []model.MembershipPlan

	if err := qb.
		Scopes(query.Paginated()).
		Order(query.Sorted()).
		Find(&plans).Error; err != nil {
		logger.WithError(err).Error("failed to find membership plans")
		return nil, 0, err
	}

	return plans, total, nil
}

func (r *membershipPlanRepository) FindByID(ctx context.Context, id string) (model.MembershipPlan, error) {
	logger := logrus.WithField("id", id)

	var plan model.MembershipPlan

	if err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&plan).Error; err != nil {
		logger.WithError(err).Error("failed to find membership plan")
		return model.MembershipPlan{}, err
	}

	return plan, nil
}

func (r *membershipPlanRepository) Create(ctx context.Context, input model.MembershipPlan) error {
	logger := logrus.WithField("input", utils.Dump(input))

	if err := r.db.
		WithContext(ctx).
		Create(&input).Error; err != nil {
		logger.WithError(err).Error("failed to create membership plan")
		return err
	}

	return nil
}

func (r *membershipPlanRepository) Update(ctx context.Context, id string, input model.MembershipPlan) error {
	logger := logrus.WithField("id", id).WithField("input", utils.Dump(input))

	var toUpdate map[string]interface{}

	if input.Name != "" {
		toUpdate["name"] = input.Name
	}

	if input.Price.String() != "" {
		toUpdate["price"] = input.Price
	}

	if input.BillingCycle != "" {
		toUpdate["billing_cycle"] = input.BillingCycle
	}

	if len(input.Features) > 0 {
		toUpdate["features"] = input.Features
	}

	if input.IsPopular {
		toUpdate["is_popular"] = input.IsPopular
	}

	if input.MaxFolders != 0 {
		toUpdate["max_folders"] = input.MaxFolders
	}

	if input.CustomDomain {
		toUpdate["custom_domain"] = input.CustomDomain
	}

	if input.AdvancedAnalytics {
		toUpdate["advanced_analytics"] = input.AdvancedAnalytics
	}

	if input.StripeProductID != "" {
		toUpdate["stripe_product_id"] = input.StripeProductID
	}

	if err := r.db.
		WithContext(ctx).
		Model(&model.MembershipPlan{}).
		Where("id = ?", id).
		Updates(toUpdate).Error; err != nil {
		logger.WithError(err).Error("failed to update membership plan")
		return err
	}

	return nil
}

func (r *membershipPlanRepository) Delete(ctx context.Context, id string) error {
	logger := logrus.WithField("id", id)

	if err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.MembershipPlan{}).Error; err != nil {
		logger.WithError(err).Error("failed to delete membership plan")
		return err
	}

	return nil
}
