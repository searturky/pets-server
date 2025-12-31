// Package cron 定时任务调度
package cron

import (
	"context"
	"log"
	"time"

	"pets-server/internal/domain/pet"
	"pets-server/internal/domain/shared"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	petRepo   pet.Repository
	uow       shared.UnitOfWork
	publisher shared.EventPublisher
	stopCh    chan struct{}
}

// NewScheduler 创建调度器
func NewScheduler(
	petRepo pet.Repository,
	uow shared.UnitOfWork,
	publisher shared.EventPublisher,
) *Scheduler {
	return &Scheduler{
		petRepo:   petRepo,
		uow:       uow,
		publisher: publisher,
		stopCh:    make(chan struct{}),
	}
}

// Start 启动定时任务
func (s *Scheduler) Start() {
	go s.runPetStatusDecay()
	log.Println("Scheduler started")
}

// Stop 停止定时任务
func (s *Scheduler) Stop() {
	close(s.stopCh)
	log.Println("Scheduler stopped")
}

// runPetStatusDecay 宠物状态衰减任务
// 每小时执行一次
func (s *Scheduler) runPetStatusDecay() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.decayAllPetStatus()
		}
	}
}

// decayAllPetStatus 衰减所有宠物状态
func (s *Scheduler) decayAllPetStatus() {
	ctx := context.Background()

	// 获取总数
	total, err := s.petRepo.CountAll(ctx)
	if err != nil {
		log.Printf("Failed to count pets: %v", err)
		return
	}

	log.Printf("Starting pet status decay for %d pets", total)

	// 分批处理
	batchSize := 100
	for offset := 0; offset < int(total); offset += batchSize {
		pets, err := s.petRepo.FindAll(ctx, offset, batchSize)
		if err != nil {
			log.Printf("Failed to find pets: %v", err)
			continue
		}

		for _, p := range pets {
			s.decayPetStatus(ctx, p)
		}
	}

	log.Println("Pet status decay completed")
}

// decayPetStatus 衰减单个宠物状态
func (s *Scheduler) decayPetStatus(ctx context.Context, p *pet.Pet) {
	err := s.uow.Do(ctx, func(txCtx context.Context) error {
		// 衰减1小时的状态
		p.DecayStatus(1.0)

		// 检查是否需要发送警告
		var events []shared.Event
		if p.IsHungry() {
			events = append(events, pet.PetStatusWarningEvent{
				PetID:       p.ID,
				UserID:      p.UserID,
				WarningType: "hungry",
				Timestamp:   time.Now(),
			})
		}
		if p.IsUnhappy() {
			events = append(events, pet.PetStatusWarningEvent{
				PetID:       p.ID,
				UserID:      p.UserID,
				WarningType: "unhappy",
				Timestamp:   time.Now(),
			})
		}

		// 保存
		if err := s.petRepo.Save(txCtx, p); err != nil {
			return err
		}

		// 发布事件
		if s.publisher != nil && len(events) > 0 {
			for _, e := range events {
				_ = s.publisher.Publish(txCtx, e)
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Failed to decay pet %d status: %v", p.ID, err)
	}
}

