// Package pet 宠物应用服务
// DTO 数据传输对象
package pet

// CreatePetRequest 创建宠物请求
type CreatePetRequest struct {
	Name      string `json:"name" binding:"required,min=1,max=20"` // 宠物名称
	SpeciesID string `json:"speciesId,omitempty"`                  // 物种ID（可选，不填则随机）
}

// CreatePetResponse 创建宠物响应
type CreatePetResponse struct {
	Pet PetDetailDTO `json:"pet"`
}

// FeedPetRequest 喂食请求
type FeedPetRequest struct {
	FoodItemID int `json:"foodItemId" binding:"required"` // 食物道具ID
}

// FeedPetResponse 喂食响应
type FeedPetResponse struct {
	Hunger    int  `json:"hunger"`    // 当前饱食度
	ExpGained int  `json:"expGained"` // 获得经验
	LevelUp   bool `json:"levelUp"`   // 是否升级
	NewLevel  int  `json:"newLevel"`  // 新等级（如果升级）
}

// PlayPetResponse 玩耍响应
type PlayPetResponse struct {
	Happiness int  `json:"happiness"` // 当前快乐度
	Energy    int  `json:"energy"`    // 当前精力
	ExpGained int  `json:"expGained"` // 获得经验
	LevelUp   bool `json:"levelUp"`   // 是否升级
}

// CleanPetResponse 清洁响应
type CleanPetResponse struct {
	Cleanliness int `json:"cleanliness"` // 当前清洁度
}

// PetDetailDTO 宠物详情DTO
type PetDetailDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	// 外观
	Appearance AppearanceDTO `json:"appearance"`

	// 性格
	Personality PersonalityDTO `json:"personality"`

	// 技能
	Skill SkillDTO `json:"skill"`

	// 成长
	Stage     string `json:"stage"`     // 阶段名称
	Level     int    `json:"level"`     // 等级
	Exp       int    `json:"exp"`       // 当前经验
	ExpToNext int    `json:"expToNext"` // 升级所需经验

	// 状态
	Status StatusDTO `json:"status"`

	// 基因码（可选，用于展示独特性）
	GeneCode string `json:"geneCode,omitempty"`
}

// AppearanceDTO 外观DTO
type AppearanceDTO struct {
	ColorPrimary   string `json:"colorPrimary"`
	ColorSecondary string `json:"colorSecondary"`
	PatternType    string `json:"patternType"`
	BodyType       string `json:"bodyType"`
	Description    string `json:"description"`
}

// PersonalityDTO 性格DTO
type PersonalityDTO struct {
	Activity    int    `json:"activity"`
	Appetite    int    `json:"appetite"`
	Social      int    `json:"social"`
	Curiosity   int    `json:"curiosity"`
	Description string `json:"description"`
}

// SkillDTO 技能DTO
type SkillDTO struct {
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Rarity      string `json:"rarity"`
	Description string `json:"description"`
}

// StatusDTO 状态DTO
type StatusDTO struct {
	Hunger      int  `json:"hunger"`
	Happiness   int  `json:"happiness"`
	Cleanliness int  `json:"cleanliness"`
	Energy      int  `json:"energy"`
	IsHungry    bool `json:"isHungry"`
	IsUnhappy   bool `json:"isUnhappy"`
	IsDirty     bool `json:"isDirty"`
	IsTired     bool `json:"isTired"`
}

// PetSimpleDTO 宠物简要信息（用于列表、拜访等）
type PetSimpleDTO struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Level     int    `json:"level"`
	Stage     string `json:"stage"`
	OwnerName string `json:"ownerName,omitempty"`
}

// --- 繁殖相关 DTO ---

// BreedPetsRequest 繁殖请求
type BreedPetsRequest struct {
	Parent1ID int    `json:"parent1Id" binding:"required"` // 父母1 ID
	Parent2ID int    `json:"parent2Id"`                    // 父母2 ID（可选，不填为分裂繁殖）
	ChildName string `json:"childName" binding:"required"` // 后代名称
}

// BreedPetsResponse 繁殖响应
type BreedPetsResponse struct {
	Offspring      PetDetailDTO          `json:"offspring"`      // 后代
	InheritedGenes []string              `json:"inheritedGenes"` // 继承的基因特征
	Mutations      []string              `json:"mutations"`      // 发生的变异
	Parent1Updated PetBreedingStatusDTO  `json:"parent1Updated"` // 父母1更新后状态
	Parent2Updated *PetBreedingStatusDTO `json:"parent2Updated"` // 父母2更新后状态（分裂繁殖时为nil）
}

// PetBreedingStatusDTO 宠物繁殖状态
type PetBreedingStatusDTO struct {
	ID              int `json:"id"`
	BreedingCount   int `json:"breedingCount"`   // 繁殖次数
	RemainingBreeds int `json:"remainingBreeds"` // 剩余可繁殖次数
}

// CanBreedResponse 是否可繁殖响应
type CanBreedResponse struct {
	CanBreed bool   `json:"canBreed"`
	Reason   string `json:"reason,omitempty"` // 不能繁殖的原因
}

// PredictOffspringRequest 预测后代请求
type PredictOffspringRequest struct {
	Parent1ID int `json:"parent1Id" binding:"required"`
	Parent2ID int `json:"parent2Id"`
}

// SpeciesProbabilityDTO 物种概率
type SpeciesProbabilityDTO struct {
	SpeciesID   string  `json:"speciesId"`
	SpeciesName string  `json:"speciesName"`
	Probability float64 `json:"probability"` // 0-1
}

// PredictOffspringResponse 预测后代响应
type PredictOffspringResponse struct {
	PossibleSpecies []SpeciesProbabilityDTO `json:"possibleSpecies"`
}

// --- 物种相关 DTO ---

// SpeciesDTO 物种信息
type SpeciesDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Rarity   int    `json:"rarity"`
}

// GetSpeciesListResponse 物种列表响应
type GetSpeciesListResponse struct {
	Species []SpeciesDTO `json:"species"`
}

// GetPetScoreResponse 宠物评分响应
type GetPetScoreResponse struct {
	Score       int            `json:"score"`
	Breakdown   ScoreBreakdown `json:"breakdown"`
}

// ScoreBreakdown 评分明细
type ScoreBreakdown struct {
	LevelScore   int `json:"levelScore"`
	SkillScore   int `json:"skillScore"`
	RarityScore  int `json:"rarityScore"`
	StatusScore  int `json:"statusScore"`
	StageScore   int `json:"stageScore"`
	GenerationScore int `json:"generationScore"`
}

