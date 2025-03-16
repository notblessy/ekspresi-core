package repository

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/notblessy/ekspresi-core/model"
)

type uploaderRepository struct {
	cloudinary *cloudinary.Cloudinary
}

// NewUploaderRepository creates a new instance of uploader.
func NewUploaderRepository(cloudinary *cloudinary.Cloudinary) model.UploaderRepository {
	return &uploaderRepository{
		cloudinary: cloudinary,
	}
}

// Upload uploads a file to cloudinary.
func (u *uploaderRepository) Upload(ctx context.Context, file io.Reader, path string) (string, string, error) {
	uploadResult, err := u.cloudinary.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: path,
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
