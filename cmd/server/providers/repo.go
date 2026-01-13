package providers

import (
	"gorm.io/gorm"

	"pets-server/internal/infrastructure/persistence/postgres/repo"
)

// RepoSet 仓储集合
type RepoSet struct {
	User   *repo.UserRepository
	Pet    *repo.PetRepository
	Item   *repo.ItemRepository
	Friend *repo.FriendRepository
	Gift   *repo.GiftRepository
	Trade  *repo.TradeRepository
	Visit  *repo.VisitRepository
}

// ProvideRepoSet 提供所有仓储
func ProvideRepoSet(db *gorm.DB) *RepoSet {
	return &RepoSet{
		User:   repo.NewUserRepository(db),
		Pet:    repo.NewPetRepository(db),
		Item:   repo.NewItemRepository(db),
		Friend: repo.NewFriendRepository(db),
		Gift:   repo.NewGiftRepository(db),
		Trade:  repo.NewTradeRepository(db),
		Visit:  repo.NewVisitRepository(db),
	}
}
