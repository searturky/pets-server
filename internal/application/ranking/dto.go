// Package ranking 排行榜应用服务
// DTO 数据传输对象
package ranking

// RankType 排行榜类型
type RankType string

const (
	RankTypePetLevel    RankType = "pet_level"    // 宠物等级排行
	RankTypeAchievement RankType = "achievement"  // 成就数量排行
	RankTypeIntimacy    RankType = "intimacy"     // 亲密度排行
)

// RankingRequest 排行榜请求
type RankingRequest struct {
	Type   RankType `form:"type" binding:"required"`
	Offset int      `form:"offset"`
	Limit  int      `form:"limit"`
}

// RankingResponse 排行榜响应
type RankingResponse struct {
	Type     RankType      `json:"type"`
	Rankings []RankItemDTO `json:"rankings"`
	MyRank   *RankItemDTO  `json:"myRank,omitempty"`
}

// RankItemDTO 排行项DTO
type RankItemDTO struct {
	Rank      int    `json:"rank"`
	UserID    int64  `json:"userId"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
	Score     int    `json:"score"`
	PetName   string `json:"petName,omitempty"`
}

