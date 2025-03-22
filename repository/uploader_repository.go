package repository

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/notblessy/ekspresi-core/model"
	"github.com/notblessy/ekspresi-core/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type uploaderRepository struct {
	cloudinary *cloudinary.Cloudinary
	db         *gorm.DB
}

// NewUploaderRepository creates a new instance of uploader.
func NewUploaderRepository(cloudinary *cloudinary.Cloudinary, db *gorm.DB) model.UploaderRepository {
	return &uploaderRepository{
		cloudinary: cloudinary,
		db:         db,
	}
}

// Upload uploads a file to cloudinary.
func (u *uploaderRepository) Upload(ctx context.Context, file io.Reader, path string) (string, string, error) {
	uploadResult, err := u.cloudinary.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         path,
		Format:         "webp",
		Transformation: "q_auto:good",
	})
	if err != nil {
		return "", "", err
	}

	return uploadResult.SecureURL, uploadResult.PublicID, nil
}

// DeleteByPublicIDs deletes a file from cloudinary by public IDs.
func (u *uploaderRepository) DeleteByPublicIDs(ctx context.Context, publicIDs []string) error {
	_, err := u.cloudinary.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{
		PublicIDs: publicIDs,
		AssetType: "image",
	})

	return err
}

// Flush deletes all files from cloudinary.
func (u *uploaderRepository) Flush(ctx context.Context) error {
	allParam := true

	_, err := u.cloudinary.Admin.DeleteAllAssets(ctx, admin.DeleteAllAssetsParams{
		All: &allParam,
	})
	if err != nil {
		return err
	}

	return nil
}

// SavePhoto saves a photo to the database.
func (u *uploaderRepository) SavePhoto(ctx context.Context, photo model.Photo) error {
	logger := logrus.WithField("photo", utils.Dump(photo))

	err := u.db.Save(&photo).Error
	if err != nil {
		logger.WithError(err).Error("failed to save photo")
		return err
	}

	return nil
}

// FindByPublicIDs finds a photo by public IDs.
func (u *uploaderRepository) FindByPublicIDs(ctx context.Context, publicIDs []string) ([]model.Photo, error) {
	logger := logrus.WithField("public_ids", publicIDs)

	var photos []model.Photo

	err := u.db.Where("public_id IN (?)", publicIDs).Find(&photos).Error
	if err != nil {
		logger.WithError(err).Error("failed to find photos")
		return nil, err
	}

	return photos, nil
}
