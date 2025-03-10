package repository

import (
	"context"

	"github.com/notblessy/ekspresi-core/model"
	"github.com/notblessy/ekspresi-core/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type portfolioRepository struct {
	db *gorm.DB
}

// NewPortfolioRepository :nodoc:
func NewPortfolioRepository(d *gorm.DB) model.PortfolioRepository {
	return &portfolioRepository{
		db: d,
	}
}

func (p *portfolioRepository) Patch(ctx context.Context, input model.PortfolioType) error {
	logger := logrus.WithField("ctx", utils.Dump(ctx))

	tx := p.db.Begin()

	porto := input.GetPortfolio()

	if porto.ID != "" {
		var portfolioToUpdate map[string]interface{}

		if porto.Title != "" {
			portfolioToUpdate["title"] = porto.Title
		}

		if porto.Columns != 0 {
			portfolioToUpdate["columns"] = porto.Columns
		}

		if porto.Gap != 0 {
			portfolioToUpdate["gap"] = porto.Gap
		}

		if porto.RoundedCorners {
			portfolioToUpdate["rounded_corners"] = porto.RoundedCorners
		}

		if porto.ShowCaptions {
			portfolioToUpdate["show_captions"] = porto.ShowCaptions
		}

		if err := tx.Model(&model.Portfolio{}).Where("id = ?", porto.ID).Updates(portfolioToUpdate).Error; err != nil {
			tx.Rollback()
			logger.WithError(err).Error("failed to update portfolio")
			return err
		}
	}

	profile := input.GetProfiles()

	if profile.ID != "" {
		var profileToUpdate map[string]interface{}

		if profile.Name != "" {
			profileToUpdate["name"] = profile.Name
		}

		if profile.Title != "" {
			profileToUpdate["title"] = profile.Title
		}

		if profile.Bio != "" {
			profileToUpdate["bio"] = profile.Bio
		}

		if profile.Email != "" {
			profileToUpdate["email"] = profile.Email
		}

		if profile.Instagram != "" {
			profileToUpdate["instagram"] = profile.Instagram
		}

		if profile.Website != "" {
			profileToUpdate["website"] = profile.Website
		}

		if err := tx.Model(&model.Profile{}).Where("id = ?", profile.ID).Updates(profileToUpdate).Error; err != nil {
			tx.Rollback()
			logger.WithError(err).Error("failed to update profile")

			return err
		}
	}

	folders := input.GetFolders()

	for _, folder := range folders {
		if folder.ID != "" {
			var folderToUpdate map[string]interface{}

			if folder.Name != "" {
				folderToUpdate["name"] = folder.Name
			}

			if folder.Description != "" {
				folderToUpdate["description"] = folder.Description
			}

			if folder.CoverID != 0 {
				folderToUpdate["cover_id"] = folder.CoverID
			}

			if err := tx.Model(&model.Folder{}).Where("id = ?", folder.ID).Updates(folderToUpdate).Error; err != nil {
				tx.Rollback()
				logger.WithError(err).Error("failed to update folder")

				return err
			}
		}
	}

	tx.Commit()
	return nil
}
