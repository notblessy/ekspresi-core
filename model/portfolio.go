package model

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
)

const (
	DefaultTitle = "Professional Photographer"
)

type PortfolioRepository interface {
	Patch(ctx context.Context, p PortfolioType) error
}

type Portfolio struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Theme          string    `json:"theme"`
	Columns        int       `json:"columns"`
	Gap            int       `json:"gap"`
	RoundedCorners bool      `json:"rounded_corners"`
	ShowCaptions   bool      `json:"show_captions"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Profile struct {
	ID          string `json:"id"`
	PortfolioID string `json:"portfolio_id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Bio         string `json:"bio"`
	Email       string `json:"email"`
	Instagram   string `json:"instagram"`
	Website     string `json:"website"`
}

func (p *Profile) TableName() string {
	return "profiles"
}

type Folder struct {
	ID             string    `json:"id"`
	PortfolioID    string    `json:"portfolio_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	CoverID        int       `json:"cover_id"`
	Columns        int       `json:"columns"`
	Gap            int       `json:"gap"`
	ShowCaptions   bool      `json:"show_captions"`
	RoundedCorners bool      `json:"rounded_corners"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewInitialFolder(portfolioID string) []Folder {
	return []Folder{
		{
			ID:             ulid.Make().String(),
			PortfolioID:    portfolioID,
			Name:           "Landscapes",
			Columns:        3,
			Gap:            16,
			ShowCaptions:   true,
			RoundedCorners: true,
			Description:    "A collection of landscapes",
		},
		{
			ID:             ulid.Make().String(),
			PortfolioID:    portfolioID,
			Name:           "Portraits",
			Columns:        3,
			Gap:            16,
			ShowCaptions:   true,
			RoundedCorners: true,
			Description:    "A collection of portraits",
		},
		{
			ID:             ulid.Make().String(),
			PortfolioID:    portfolioID,
			Name:           "Events",
			Columns:        3,
			Gap:            16,
			ShowCaptions:   true,
			RoundedCorners: true,
			Description:    "A collection of events",
		},
	}
}

type Photo struct {
	ID        string    `json:"id"`
	FolderID  string    `json:"folder_id"`
	Src       string    `json:"src"`
	Alt       string    `json:"alt"`
	Caption   string    `json:"caption"`
	PublicID  string    `json:"public_id"`
	SortIndex int       `json:"sort_index"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *Photo) TableName() string {
	return "photos"
}

type FolderType struct {
	Folder
	Photos []Photo `json:"photos" gorm:"foreignKey:FolderID;references:ID"`
}

func (ft FolderType) TableName() string {
	return "folders"
}

type PortfolioType struct {
	Portfolio
	Profiles Profile      `json:"profiles" gorm:"foreignKey:PortfolioID;references:ID"`
	Folders  []FolderType `json:"folders" gorm:"foreignKey:PortfolioID;references:ID"`
}

func (p *PortfolioType) TableName() string {
	return "portfolios"
}

func (pt *PortfolioType) GetPortfolio() Portfolio {
	return Portfolio{
		ID:             pt.ID,
		UserID:         pt.UserID,
		Title:          pt.Title,
		Columns:        pt.Columns,
		Gap:            pt.Gap,
		RoundedCorners: pt.RoundedCorners,
		ShowCaptions:   pt.ShowCaptions,
		CreatedAt:      pt.CreatedAt,
		UpdatedAt:      pt.UpdatedAt,
	}
}

func (pt *PortfolioType) GetProfiles() Profile {
	return Profile{
		ID:          pt.Profiles.ID,
		PortfolioID: pt.Profiles.PortfolioID,
		Name:        pt.Profiles.Name,
		Title:       pt.Profiles.Title,
		Bio:         pt.Profiles.Bio,
		Email:       pt.Profiles.Email,
		Instagram:   pt.Profiles.Instagram,
		Website:     pt.Profiles.Website,
	}
}

func (pt *PortfolioType) GetFolders() []Folder {
	var folders []Folder

	for _, f := range pt.Folders {
		folders = append(folders, Folder{
			ID:          f.ID,
			PortfolioID: f.PortfolioID,
			Name:        f.Name,
			Description: f.Description,
			CoverID:     f.CoverID,
			CreatedAt:   f.CreatedAt,
			UpdatedAt:   f.UpdatedAt,
		})
	}

	return folders
}
