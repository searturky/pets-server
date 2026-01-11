// Package pet 宠物领域
// Pet 实体 - 宠物聚合根
package pet

import (
	"errors"
	"time"
)

// Stage 成长阶段
type Stage int

const (
	StageEgg     Stage = 0 // 蛋
	StageChild   Stage = 1 // 幼年
	StageTeen    Stage = 2 // 成长期
	StageAdult   Stage = 3 // 成熟期
	StageElderly Stage = 4 // 老年期
)

// StageName 获取阶段名称
func (s Stage) Name() string {
	names := []string{"蛋", "幼年期", "成长期", "成熟期", "老年期"}
	if int(s) < len(names) {
		return names[s]
	}
	return "未知"
}

// FoodType 食物类型
type FoodType int

const (
	FoodTypeBasic   FoodType = 1 // 普通食物
	FoodTypePremium FoodType = 2 // 高级食物
	FoodTypeSpecial FoodType = 3 // 特殊食物
)

// GetHungerRestore 获取食物恢复的饱食度
func (f FoodType) GetHungerRestore() int {
	switch f {
	case FoodTypeBasic:
		return 20
	case FoodTypePremium:
		return 40
	case FoodTypeSpecial:
		return 60
	default:
		return 10
	}
}

// Pet 宠物实体（聚合根）
type Pet struct {
	ID     int64
	UserID int64
	Name   string

	// 物种和性别
	SpeciesID SpeciesID // 物种ID
	Gender    Gender    // 性别

	// 基因与衍生属性
	Gene              Gene             // 基因
	Appearance        Appearance       // 通用外观
	SpecialAppearance SpecialAppearance // 物种特有外观
	Personality       Personality      // 性格
	Skill             Skill            // 技能

	// 成长状态
	Stage Stage
	Exp   int
	Level int

	// 实时状态 (0-100)
	Hunger      int
	Happiness   int
	Cleanliness int
	Energy      int

	// 繁衍相关
	Parent1ID   *int64     // 父方ID (可为空)
	Parent2ID   *int64     // 母方ID (可为空)
	Generation  int        // 代数
	LastBreedAt *time.Time // 上次繁殖时间

	// 时间记录
	LastFedAt     time.Time
	LastPlayedAt  time.Time
	LastCleanedAt time.Time
	BornAt        time.Time
	CreatedAt     time.Time

	// 领域事件收集
	events []any
}

// NewPet 创建新宠物（从蛋开始）
// 使用默认物种（猫）
func NewPet(userID int64, name string) *Pet {
	return NewPetWithSpecies(userID, name, SpeciesCat, nil)
}

// NewPetWithSpecies 创建指定物种的新宠物
func NewPetWithSpecies(userID int64, name string, speciesID SpeciesID, genderRule *GenderRule) *Pet {
	gene := GenerateGene()
	now := time.Now()

	// 确定性别
	var gender Gender
	if genderRule != nil {
		gender = DetermineGender(gene, *genderRule)
	} else {
		gender = DetermineGender(gene, DefaultGenderRule())
	}

	return &Pet{
		UserID:            userID,
		Name:              name,
		SpeciesID:         speciesID,
		Gender:            gender,
		Gene:              gene,
		Appearance:        NewAppearanceFromGene(gene),
		SpecialAppearance: NewSpecialAppearance(), // 由物种解释器填充
		Personality:       NewPersonalityFromGene(gene),
		Skill:             NewSkillFromGene(gene),
		Stage:             StageEgg,
		Exp:               0,
		Level:             1,
		Hunger:            50,
		Happiness:         50,
		Cleanliness:       50,
		Energy:            100,
		Generation:        0,
		BornAt:            now,
		CreatedAt:         now,
	}
}

// NewPetFromBreeding 从繁殖创建新宠物
func NewPetFromBreeding(userID int64, name string, speciesID SpeciesID, gene Gene, gender Gender, parent1ID, parent2ID int64, generation int) *Pet {
	now := time.Now()

	return &Pet{
		UserID:            userID,
		Name:              name,
		SpeciesID:         speciesID,
		Gender:            gender,
		Gene:              gene,
		Appearance:        NewAppearanceFromGene(gene),
		SpecialAppearance: NewSpecialAppearance(),
		Personality:       NewPersonalityFromGene(gene),
		Skill:             NewSkillFromGene(gene),
		Stage:             StageEgg,
		Exp:               0,
		Level:             1,
		Hunger:            50,
		Happiness:         50,
		Cleanliness:       50,
		Energy:            100,
		Parent1ID:         &parent1ID,
		Parent2ID:         &parent2ID,
		Generation:        generation,
		BornAt:            now,
		CreatedAt:         now,
	}
}

// SetSpecialAppearance 设置物种特有外观（由物种解释器调用）
func (p *Pet) SetSpecialAppearance(special SpecialAppearance) {
	p.SpecialAppearance = special
}

// --- 核心业务方法 ---

// Feed 喂食
func (p *Pet) Feed(foodType FoodType) error {
	if p.Stage == StageEgg {
		return ErrPetIsEgg
	}
	if p.Hunger >= 100 {
		return ErrPetIsFull
	}

	// 计算恢复量
	restore := foodType.GetHungerRestore()

	// 性格加成
	if p.Personality.Appetite > 70 {
		restore = int(float64(restore) * 1.2)
	}

	// 技能加成
	if p.Skill.Type == SkillTypeGluttony {
		restore = int(float64(restore) * p.Skill.EffectMultiplier())
	}

	// 更新状态
	p.Hunger = minInt(p.Hunger+restore, 100)
	p.Happiness = minInt(p.Happiness+5, 100)
	p.LastFedAt = time.Now()

	// 获得经验
	exp := int(10 * p.Personality.FeedExpBonus())
	p.addExp(exp)

	// 记录事件
	p.addEvent(PetFedEvent{
		PetID:     p.ID,
		UserID:    p.UserID,
		FoodType:  int(foodType),
		ExpGained: exp,
	})

	return nil
}

// Play 玩耍
func (p *Pet) Play() error {
	if p.Stage == StageEgg {
		return ErrPetIsEgg
	}
	if p.Happiness >= 100 {
		return ErrPetIsHappy
	}
	if p.Energy < 10 {
		return ErrPetIsTired
	}

	// 计算恢复量
	restore := 20
	if p.Personality.Activity > 70 {
		restore = int(float64(restore) * 1.2)
	}
	if p.Skill.Type == SkillTypePlayful {
		restore = int(float64(restore) * p.Skill.EffectMultiplier())
	}

	// 更新状态
	p.Happiness = minInt(p.Happiness+restore, 100)
	p.Energy = maxInt(p.Energy-15, 0)
	p.LastPlayedAt = time.Now()

	// 获得经验
	exp := int(10 * p.Personality.PlayExpBonus())
	p.addExp(exp)

	return nil
}

// Clean 清洁
func (p *Pet) Clean() error {
	if p.Stage == StageEgg {
		return ErrPetIsEgg
	}
	if p.Cleanliness >= 100 {
		return ErrPetIsClean
	}

	restore := 30
	if p.Skill.Type == SkillTypeCleanlover {
		restore = int(float64(restore) * p.Skill.EffectMultiplier())
	}

	p.Cleanliness = minInt(p.Cleanliness+restore, 100)
	p.LastCleanedAt = time.Now()

	// 清洁也给少量经验
	p.addExp(5)

	return nil
}

// Rest 休息恢复能量
func (p *Pet) Rest() {
	p.Energy = minInt(p.Energy+30, 100)
}

// --- 繁殖相关 ---

// CanBreed 检查是否可以繁殖
func (p *Pet) CanBreed(breedRules BreedingRules) error {
	if p.Stage < breedRules.MinStage {
		return ErrPetNotMature
	}
	if p.Level < breedRules.MinLevel {
		return ErrPetLevelTooLow
	}
	if p.Happiness < breedRules.MinHappiness {
		return ErrPetUnhappy
	}
	if p.LastBreedAt != nil {
		cooldown := time.Duration(breedRules.CooldownHours) * time.Hour
		if p.Gender == GenderNone {
			cooldown = time.Duration(breedRules.SelfBreedCooldownHours) * time.Hour
		}
		if time.Since(*p.LastBreedAt) < cooldown {
			return ErrBreedCooldown
		}
	}
	return nil
}

// CanBreedWith 检查是否可以与另一只宠物繁殖
func (p *Pet) CanBreedWith(other *Pet) error {
	if !CanBreedWith(p.Gender, other.Gender) {
		return ErrIncompatibleGender
	}
	return nil
}

// MarkBred 标记已繁殖
func (p *Pet) MarkBred() {
	now := time.Now()
	p.LastBreedAt = &now
}

// --- 状态衰减（由定时任务调用） ---

// DecayStatus 状态衰减
func (p *Pet) DecayStatus(hours float64) {
	if p.Stage == StageEgg {
		return
	}

	// 基础衰减率（每小时）
	baseDecay := 5.0

	// 耐力技能减缓衰减
	if p.Skill.Type == SkillTypeEndurance {
		baseDecay = baseDecay / p.Skill.EffectMultiplier()
	}

	// 饥饿衰减
	hungerDecay := int(baseDecay * p.Personality.HungerDecayRate() * hours)
	p.Hunger = maxInt(p.Hunger-hungerDecay, 0)

	// 快乐衰减
	happinessDecay := int(baseDecay * p.Personality.HappinessDecayRate() * hours)
	p.Happiness = maxInt(p.Happiness-happinessDecay, 0)

	// 清洁衰减
	cleanlinessDecay := int(baseDecay * hours)
	p.Cleanliness = maxInt(p.Cleanliness-cleanlinessDecay, 0)

	// 能量恢复（休息时）
	p.Energy = minInt(p.Energy+int(10*hours), 100)
}

// --- 成长与进化 ---

// addExp 增加经验
func (p *Pet) addExp(exp int) {
	p.Exp += exp
	p.checkLevelUp()
	p.checkEvolution()
}

// checkLevelUp 检查升级
func (p *Pet) checkLevelUp() {
	for {
		required := p.requiredExpForLevel()
		if p.Exp < required {
			break
		}
		p.Exp -= required
		p.Level++

		p.addEvent(PetLevelUpEvent{
			PetID:    p.ID,
			UserID:   p.UserID,
			NewLevel: p.Level,
		})
	}
}

// requiredExpForLevel 当前等级所需经验
func (p *Pet) requiredExpForLevel() int {
	return p.Level * 100
}

// checkEvolution 检查进化
func (p *Pet) checkEvolution() {
	var evolved bool
	switch p.Stage {
	case StageEgg:
		if p.Level >= 3 {
			p.Stage = StageChild
			evolved = true
		}
	case StageChild:
		if p.Level >= 10 {
			p.Stage = StageTeen
			evolved = true
		}
	case StageTeen:
		if p.Level >= 25 {
			p.Stage = StageAdult
			evolved = true
		}
	case StageAdult:
		if p.Level >= 50 {
			p.Stage = StageElderly
			evolved = true
		}
	}

	if evolved {
		p.addEvent(PetEvolvedEvent{
			PetID:    p.ID,
			UserID:   p.UserID,
			NewStage: int(p.Stage),
		})
	}
}

// --- 状态查询 ---

// IsHungry 是否饥饿
func (p *Pet) IsHungry() bool {
	return p.Hunger < 30
}

// IsUnhappy 是否不开心
func (p *Pet) IsUnhappy() bool {
	return p.Happiness < 30
}

// IsDirty 是否脏了
func (p *Pet) IsDirty() bool {
	return p.Cleanliness < 30
}

// IsTired 是否疲劳
func (p *Pet) IsTired() bool {
	return p.Energy < 20
}

// StageName 阶段名称
func (p *Pet) StageName() string {
	return p.Stage.Name()
}

// GenderName 性别名称
func (p *Pet) GenderName() string {
	return p.Gender.Name()
}

// --- 领域事件 ---

func (p *Pet) addEvent(event any) {
	p.events = append(p.events, event)
}

// Events 获取并清空事件
func (p *Pet) Events() []any {
	events := p.events
	p.events = nil
	return events
}

// 领域错误
var (
	ErrPetIsEgg            = errors.New("宠物还在蛋里")
	ErrPetIsFull           = errors.New("宠物已经很饱了")
	ErrPetIsHappy          = errors.New("宠物已经很开心了")
	ErrPetIsClean          = errors.New("宠物已经很干净了")
	ErrPetIsTired          = errors.New("宠物太累了需要休息")
	ErrPetNotFound         = errors.New("宠物不存在")
	ErrPetNotMature        = errors.New("宠物还未成年")
	ErrPetLevelTooLow      = errors.New("宠物等级不足")
	ErrPetUnhappy          = errors.New("宠物不够开心")
	ErrBreedCooldown       = errors.New("繁殖冷却中")
	ErrIncompatibleGender  = errors.New("性别不兼容")
	ErrCannotSelfBreed     = errors.New("该物种不能自我繁殖")
	ErrSpeciesNotFound     = errors.New("物种不存在")
)
