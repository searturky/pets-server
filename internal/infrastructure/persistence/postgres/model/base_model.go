package model

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        int       `gorm:"primarykey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime; comment:创建时间"`
	UpdatedAt time.Time `gorm:"autoUpdateTime; comment:更新时间"`
	UUID      uuid.UUID `gorm:"type:uuid; default:uuidv7(); index:unique"`
}
