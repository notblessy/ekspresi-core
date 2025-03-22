package model

import (
	"context"
	"io"
)

type UploaderRepository interface {
	Upload(ctx context.Context, file io.Reader, path string) (string, string, error)
	DeleteByPublicIDs(ctx context.Context, publicID []string) error
	Flush(ctx context.Context) error
	FindByPublicIDs(ctx context.Context, publicIDs []string) ([]Photo, error)
	SavePhoto(ctx context.Context, photo Photo) error
}

type DeleteRequest struct {
	PublicIDs []string `json:"public_ids"`
}
