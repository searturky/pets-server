// Package achievement 成就领域
// 成就实体和定义
package achievement

import "time"

// AchievementCategory 成就分类
type AchievementCategory string

const (
	CategoryPet    AchievementCategory = "pet"    // 宠物相关
	CategorySocial AchievementCategory = "social" // 社交相关
	CategoryItem   AchievementCategory = "item"   // 道具相关
	CategoryLogin  AchievementCategory = "login"  // 登录相关
)

// ConditionType 条件类型
type ConditionType string

const (
	ConditionFeedCount     ConditionType = "feed_count"     // 喂食次数
	ConditionPlayCount     ConditionType = "play_count"     // 玩耍次数
	ConditionPetLevel      ConditionType = "pet_level"      // 宠物等级
	ConditionPetStage      ConditionType = "pet_stage"      // 宠物阶段
	ConditionFriendCount   ConditionType = "friend_count"   // 好友数量
	ConditionGiftSentCount ConditionType = "gift_sent"      // 送礼次数
	ConditionVisitCount    ConditionType = "visit_count"    // 拜访次数
	ConditionLoginDays     ConditionType = "login_days"     // 登录天数
	ConditionItemCollect   ConditionType = "item_collect"   // 收集道具数
)

// AchievementDefinition 成就定义（值对象）
type AchievementDefinition struct {
	ID             int
	Name           string
	Description    string
	Category       AchievementCategory
	ConditionType  ConditionType
	ConditionValue int    // 达成条件的值
	RewardCoins    int    // 奖励金币
	RewardDiamonds int    // 奖励钻石
	Icon           string // 图标
}

// UserAchievement 用户成就（实体）
type UserAchievement struct {
	ID            int64
	UserID        int64
	AchievementID int // 对应 AchievementDefinition.ID
	UnlockedAt    time.Time
}

// NewUserAchievement 创建用户成就
func NewUserAchievement(userID int64, achievementID int) *UserAchievement {
	return &UserAchievement{
		UserID:        userID,
		AchievementID: achievementID,
		UnlockedAt:    time.Now(),
	}
}

// AchievementUnlockedEvent 成就解锁事件
type AchievementUnlockedEvent struct {
	UserID        int64     `json:"user_id"`
	AchievementID int       `json:"achievement_id"`
	Name          string    `json:"name"`
	RewardCoins   int       `json:"reward_coins"`
	RewardDiamonds int      `json:"reward_diamonds"`
	Timestamp     time.Time `json:"timestamp"`
}

func (e AchievementUnlockedEvent) EventName() string { return "achievement.unlocked" }

