package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

const (
	DefaultTitle = "Professional Photographer"
)

type Portfolio struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Title          string    `json:"title"`
	Columns        int       `json:"columns"`
	Gap            int       `json:"gap"`
	RoundedCorners bool      `json:"rounded_corners"`
	ShowCaptions   bool      `json:"show_captions"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Profile struct {
	ID          int    `json:"id"`
	PortfolioID string `json:"portfolio_id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Bio         string `json:"bio"`
	Email       string `json:"email"`
	Instagram   string `json:"instagram"`
	Website     string `json:"website"`
}

type Folder struct {
	ID          string    `json:"id"`
	PortfolioID string    `json:"portfolio_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CoverID     int       `json:"cover_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewInitialFolder(portfolioID string) []Folder {
	return []Folder{
		{
			ID:          ulid.Make().String(),
			PortfolioID: portfolioID,
			Name:        "Landscapes",
			Description: "A collection of landscapes",
		},
		{
			ID:          ulid.Make().String(),
			PortfolioID: portfolioID,
			Name:        "Portraits",
			Description: "A collection of portraits",
		},
		{
			ID:          ulid.Make().String(),
			PortfolioID: portfolioID,
			Name:        "Events",
			Description: "A collection of events",
		},
	}
}

type Photo struct {
	ID       int    `json:"id"`
	FolderID string `json:"folder_id"`
	Src      string `json:"src"`
	Alt      string `json:"alt"`
	Caption  string `json:"caption"`
}

type FolderType struct {
	Folder
	Photos []Photo `json:"photos"`
}

type PortfolioType struct {
	Portfolio
	Profiles Profile      `json:"profiles"`
	Folders  []FolderType `json:"folders"`
}
