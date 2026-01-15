// Package pet 宠物领域
// 领域事件定义
package pet

import "time"

// PetFedEvent 宠物被喂食事件
type PetFedEvent struct {
	PetID     int       `json:"pet_id"`
	UserID    int       `json:"user_id"`
	FoodType  int       `json:"food_type"`
	ExpGained int       `json:"exp_gained"`
	Timestamp time.Time `json:"timestamp"`
}

func (e PetFedEvent) EventName() string { return "pet.fed" }

// PetLevelUpEvent 宠物升级事件
type PetLevelUpEvent struct {
	PetID     int       `json:"pet_id"`
	UserID    int       `json:"user_id"`
	NewLevel  int       `json:"new_level"`
	Timestamp time.Time `json:"timestamp"`
}

func (e PetLevelUpEvent) EventName() string { return "pet.level_up" }

// PetEvolvedEvent 宠物进化事件
type PetEvolvedEvent struct {
	PetID     int       `json:"pet_id"`
	UserID    int       `json:"user_id"`
	NewStage  int       `json:"new_stage"`
	Timestamp time.Time `json:"timestamp"`
}

func (e PetEvolvedEvent) EventName() string { return "pet.evolved" }

// PetStatusWarningEvent 宠物状态警告事件
type PetStatusWarningEvent struct {
	PetID       int       `json:"pet_id"`
	UserID      int       `json:"user_id"`
	WarningType string    `json:"warning_type"` // hungry, unhappy, dirty
	Timestamp   time.Time `json:"timestamp"`
}

func (e PetStatusWarningEvent) EventName() string { return "pet.status_warning" }

// PetCreatedEvent 宠物创建事件
type PetCreatedEvent struct {
	PetID     int       `json:"pet_id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	GeneCode  string    `json:"gene_code"`
	Timestamp time.Time `json:"timestamp"`
}

func (e PetCreatedEvent) EventName() string { return "pet.created" }

