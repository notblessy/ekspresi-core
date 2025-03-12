package repository

import (
	"context"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/notblessy/ekspresi-core/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository :nodoc:
func NewUserRepository(d *gorm.DB) model.UserRepository {
	return &userRepository{
		db: d,
	}
}

func (a *userRepository) Authenticate(ctx context.Context, code, requestOrigin string) (model.User, error) {
	logger := logrus.WithFields(logrus.Fields{
		"code":           code,
		"request_origin": requestOrigin,
	})

	auth, err := a.verifyToken(ctx, code, requestOrigin)
	if err != nil {
		logger.Errorf("Error verifying token: %v", err)
		return model.User{}, err
	}

	id, err := gonanoid.New()
	if err != nil {
		logger.Errorf("Error generating id: %v", err)
		return model.User{}, err
	}

	var authUser model.User

	err = a.db.Where("email = ?", auth.Email).First(&authUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorf("Error querying user: %v", err)
		return model.User{}, err
	}

	if err == gorm.ErrRecordNotFound {
		authUser = model.User{
			ID:      id,
			Name:    auth.Name,
			Role:    model.RoleUser,
			Email:   auth.Email,
			Picture: auth.Picture,
		}

		tx := a.db.Begin()

		err = tx.Create(&authUser).Error
		if err != nil {
			logger.Errorf("Error creating user: %v", err)
			tx.Rollback()
			return model.User{}, err
		}

		membership := authUser.NewFreeMembership()

		err = tx.Create(&membership).Error
		if err != nil {
			logger.Errorf("Error creating membership: %v", err)
			tx.Rollback()
			return model.User{}, err
		}

		defaultPortfolio := authUser.NewInitialPortfolio()

		err = tx.Create(defaultPortfolio).Error
		if err != nil {
			logger.Errorf("Error creating portfolio: %v", err)
			tx.Rollback()
			return model.User{}, err
		}

		defaultProfile := authUser.NewDefaultProfile(defaultPortfolio.ID)

		err = tx.Create(&defaultProfile).Error
		if err != nil {
			logger.Errorf("Error creating profile: %v", err)
			tx.Rollback()
			return model.User{}, err
		}

		defaultFolders := model.NewInitialFolder(defaultPortfolio.ID)

		err = tx.Create(&defaultFolders).Error
		if err != nil {
			logger.Errorf("Error creating folders: %v", err)
			tx.Rollback()
			return model.User{}, err
		}

		tx.Commit()
	}

	return authUser, nil
}

func (a *userRepository) FindByID(ctx context.Context, id string) (model.MeResponse, error) {
	logger := logrus.WithField("id", id)

	var user model.MeResponse
	err := a.db.Where("id = ?", id).
		Preload("Portfolio").
		Preload("Portfolio.Profiles").
		Preload("Portfolio.Folders").
		Preload("Portfolio.Folders.Photos").
		First(&user).Error
	if err != nil {
		logger.Errorf("Error querying user: %v", err)
		return model.MeResponse{}, err
	}

	user.OmitPassword()

	return user, nil
}
