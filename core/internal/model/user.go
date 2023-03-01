package model

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID              string
	Email           string
	FirstName       string
	LastName        string
	Phone           string
	Role            string
	Password        string
	IsSuperUser     bool
	IsVerified      bool
	IsEmailVerified bool
	ProfileSummary  *string
	LastLogin       time.Time
	UserImages      []*UserImage `bun:"rel:has-many,join:id=user_id"`
	CreatedBy       string
	CreatedAt       time.Time
	UpdatedBy       *string
	UpdatedAt       *time.Time
	DeletedBy       *string
	DeletedAt       time.Time `bun:",soft_delete"`
}

type UserImage struct {
	bun.BaseModel `bun:"table:user_images,alias:ui"`

	ID        string
	UserID    string
	Email     string
	ImageUrl  string
	Order     int64
	CreatedBy string
	CreatedAt time.Time
	UpdatedBy *string
	UpdatedAt *time.Time
	DeletedBy *string
	DeletedAt time.Time `bun:",soft_delete"`
}
