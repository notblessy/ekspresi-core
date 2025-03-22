package repository

import (
	"context"

	"github.com/notblessy/ekspresi-core/model"
	"github.com/notblessy/ekspresi-core/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type portfolioRepository struct {
	db           *gorm.DB
	uploaderRepo model.UploaderRepository
}

// NewPortfolioRepository :nodoc:
func NewPortfolioRepository(d *gorm.DB, uploaderRepo model.UploaderRepository) model.PortfolioRepository {
	return &portfolioRepository{
		db:           d,
		uploaderRepo: uploaderRepo,
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

	if len(input.DeletedPhotos) > 0 {
		if err := tx.Where("public_id IN ?", input.DeletedPhotos).Delete(&model.Photo{}).Error; err != nil {
			tx.Rollback()
			logger.WithError(err).Error("failed to delete photos")

			return err
		}

		go p.uploaderRepo.DeleteByPublicIDs(ctx, input.DeletedPhotos)
	}

	if len(input.DeletedFolders) > 0 {
		if err := tx.Where("id IN ?", input.DeletedFolders).Delete(&model.Folder{}).Error; err != nil {
			tx.Rollback()
			logger.WithError(err).Error("failed to delete folders")

			return err
		}
	}

	folders := input.GetFolders()

	var existingFolders []model.Folder

	if err := tx.Where("portfolio_id = ?", porto.ID).Find(&existingFolders).Error; err != nil {
		tx.Rollback()
		logger.WithError(err).Error("failed to get existing folders")

		return err
	}

	folderDict := folderToDict(existingFolders)

	for _, folder := range folders {
		existingFolder, ok := folderDict[folder.ID]
		if !ok {
			if err := tx.Create(&folder).Error; err != nil {
				tx.Rollback()
				logger.WithError(err).Error("failed to create folder")
				return err
			}

			continue
		}

		if folder.Name != "" {
			existingFolder.Name = folder.Name
		}

		if folder.Description != "" {
			existingFolder.Description = folder.Description
		}

		if folder.CoverID != 0 {
			existingFolder.CoverID = folder.CoverID
		}

		if err := tx.Save(&existingFolder).Error; err != nil {
			tx.Rollback()
			logger.WithError(err).Error("failed to save folder")
			return err
		}

		photos, err := p.uploaderRepo.FindByPublicIDs(ctx, folder.GetPhotoPublicIDs())
		if err != nil {
			tx.Rollback()
			logger.WithError(err).Error("failed to find photos")
			return err
		}

		photoDict := photoToDict(photos)

		for i, photo := range folder.Photos {
			if existing, ok := photoDict[photo.PublicID]; ok {
				photo.FolderID = existingFolder.ID
				photo.SortIndex = i

				if err := tx.Save(&existing).Error; err != nil {
					tx.Rollback()
					logger.WithError(err).Error("failed to save photo")
					return err
				}
			}
		}
	}

	tx.Commit()
	return nil
}

func folderToDict(folders []model.Folder) map[string]model.Folder {
	dict := make(map[string]model.Folder)

	for _, folder := range folders {
		dict[folder.ID] = folder
	}

	return dict
}

func photoToDict(photos []model.Photo) map[string]model.Photo {
	dict := make(map[string]model.Photo)

	for _, photo := range photos {
		dict[photo.PublicID] = photo
	}

	return dict
}
