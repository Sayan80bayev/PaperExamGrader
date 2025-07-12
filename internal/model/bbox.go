package model

import (
	"time"

	"gorm.io/datatypes"
)

type BBoxMetaDB struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Page        int            `json:"page"`
	ExamID      uint           `json:"exam_id"`
	BBoxPercent datatypes.JSON `gorm:"type:jsonb;not null" json:"bbox_percent"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
