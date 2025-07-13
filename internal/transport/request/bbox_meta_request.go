package request

import (
	"gorm.io/datatypes"
)

type BBoxMetaDB struct {
	Page        int            `json:"page"`
	BBoxPercent datatypes.JSON `gorm:"type:jsonb;not null" json:"bbox_percent" `
}
