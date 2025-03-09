package repository

import (
	"context"

	"github.com/notblessy/ekspresi-core/model"
	"github.com/notblessy/ekspresi-core/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type membershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) model.MembershipRepository {
	return &membershipRepository{db}
}

func (r *membershipRepository) FindAll(ctx context.Context, query model.MembershipQueryInput) ([]model.Membership, int64, error) {
	logger := logrus.WithField("query", utils.Dump(query))

	qb := r.db.Model(&model.Membership{}).WithContext(ctx)

	if query.UserID != "" {
		qb = qb.Where("user_id = ?", query.UserID)
	}

	if query.MembershipPlanID != "" {
		qb = qb.Where("membership_plan_id = ?", query.MembershipPlanID)
	}

	if query.Status != "" {
		qb = qb.Where("status = ?", query.Status)
	}

	var total int64

	if err := qb.Count(&total).Error; err != nil {
		logger.WithError(err).Error("failed to count memberships")
		return nil, 0, err
	}

	var memberships []model.Membership

	if err := qb.
		Scopes(query.Paginated()).
		Order(query.Sorted()).
		Find(&memberships).Error; err != nil {
		logger.WithError(err).Error("failed to find memberships")
		return nil, 0, err
	}

	return memberships, total, nil
}

func (r *membershipRepository) FindByID(ctx context.Context, id string) (model.Membership, error) {
	logger := logrus.WithField("id", id)

	var membership model.Membership

	if err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&membership).Error; err != nil {
		logger.WithError(err).Error("failed to find membership")
		return model.Membership{}, err
	}

	return membership, nil
}

func (r *membershipRepository) Create(ctx context.Context, input model.Membership) error {
	logger := logrus.WithField("input", utils.Dump(input))

	if err := r.db.
		WithContext(ctx).
		Create(&input).Error; err != nil {
		logger.WithError(err).Error("failed to create membership")
		return err
	}

	return nil
}

func (r *membershipRepository) Update(ctx context.Context, id string, input model.Membership) error {
	logger := logrus.WithField("id", id).WithField("input", utils.Dump(input))

	var membership map[string]interface{}

	if input.MembershipPlanID != "" {
		membership["membership_plan_id"] = input.MembershipPlanID
	}

	if input.Status != "" {
		membership["status"] = input.Status
	}

	if !input.StartDate.IsZero() {
		membership["start_date"] = input.StartDate
	}

	if input.EndDate.Valid {
		membership["end_date"] = input.EndDate
	}

	if input.StripeSubscriptionID != "" {
		membership["stripe_subscription_id"] = input.StripeSubscriptionID
	}

	if input.StripeSubscriptionInterval != "" {
		membership["stripe_subscription_interval"] = input.StripeSubscriptionInterval
	}

	if err := r.db.
		WithContext(ctx).
		Model(&model.Membership{}).
		Where("id = ?", id).
		Updates(membership).Error; err != nil {
		logger.WithError(err).Error("failed to update membership")
		return err
	}

	return nil
}

func (r *membershipRepository) Delete(ctx context.Context, id string) error {
	logger := logrus.WithField("id", id)

	if err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Membership{}).Error; err != nil {
		logger.WithError(err).Error("failed to delete membership")
		return err
	}

	return nil
}
