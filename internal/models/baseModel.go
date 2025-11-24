package models

import "time"

type BaseModel struct {
	Version   int
	Deleted   bool
	CreatedAt time.Time
	UpdatedAt *time.Time
	UpdatedBy *int64
}
