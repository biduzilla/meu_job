package models

import "time"

type BaseModel struct {
	Version   int
	Deleted   bool
	CreatedAt time.Time
	CreatedBy *int64
	UpdatedAt *time.Time
	UpdatedBy *int64
}
