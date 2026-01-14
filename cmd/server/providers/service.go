package providers

import (
	authApp "pets-server/internal/application/auth"
	petApp "pets-server/internal/application/pet"
	rankingApp "pets-server/internal/application/ranking"
	socialApp "pets-server/internal/application/social"
	"pets-server/internal/domain/pet"
	"pets-server/internal/domain/shared"
	"pets-server/internal/infrastructure/external/wechat"
	"pets-server/internal/infrastructure/persistence/postgres"
	"pets-server/internal/infrastructure/persistence/redis"
	"pets-server/internal/pkg/config"
)

// ServiceSet 应用服务集合
type ServiceSet struct {
	Auth    *authApp.Service
	Pet     *petApp.Service
	Social  *socialApp.Service
	Ranking *rankingApp.Service
}

// ProvideServiceSet 提供所有应用服务
func ProvideServiceSet(
	cfg *config.Config,
	repos *RepoSet,
	speciesRegistry *pet.SpeciesRegistry,
	fusionRegistry *pet.SpeciesFusionRegistry,
	uow *postgres.UnitOfWork,
	cache *redis.CacheService,
	rankingStore *redis.RankingStore,
	wechatAuth *wechat.AuthService,
	eventPublisher shared.EventPublisher,
) *ServiceSet {
	// 创建领域服务（使用注入的注册表）
	petDomainService := pet.NewDomainService(repos.Pet, speciesRegistry, fusionRegistry)

	return &ServiceSet{
		Auth: authApp.NewService(
			repos.User,
			uow,
			wechatAuth.GetOpenID,
			cfg.JWT.Secret,
			cfg.JWT.ExpireHours,
		),
		Pet: petApp.NewService(
			repos.Pet,
			repos.Item,
			petDomainService,
			uow,
			eventPublisher,
			cache,
		),
		Social: socialApp.NewService(
			repos.Friend,
			repos.Gift,
			repos.Trade,
			repos.Visit,
			uow,
			eventPublisher,
		),
		Ranking: rankingApp.NewService(rankingStore),
	}
}
