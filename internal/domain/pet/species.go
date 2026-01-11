// Package pet 宠物领域
// Species 物种定义 - 定义物种的特征模板和基因解释规则
package pet

// SpeciesID 物种ID类型
type SpeciesID int

// 预定义物种ID
const (
	SpeciesUnknown SpeciesID = 0

	// 哺乳类
	SpeciesCat    SpeciesID = 101 // 猫
	SpeciesDog    SpeciesID = 102 // 狗
	SpeciesRabbit SpeciesID = 103 // 兔
	SpeciesHamster SpeciesID = 104 // 仓鼠

	// 鸟类
	SpeciesParrot SpeciesID = 201 // 鹦鹉
	SpeciesOwl    SpeciesID = 202 // 猫头鹰
	SpeciesCanary SpeciesID = 203 // 金丝雀

	// 鱼类
	SpeciesGoldfish    SpeciesID = 301 // 金鱼
	SpeciesTropicalFish SpeciesID = 302 // 热带鱼
	SpeciesBetta       SpeciesID = 303 // 斗鱼

	// 爬行类
	SpeciesLizard  SpeciesID = 401 // 蜥蜴
	SpeciesTurtle  SpeciesID = 402 // 龟
	SpeciesGecko   SpeciesID = 403 // 壁虎

	// 幻想类/特殊类
	SpeciesSlime   SpeciesID = 501 // 史莱姆
	SpeciesPhoenix SpeciesID = 502 // 凤凰
	SpeciesDragon  SpeciesID = 503 // 龙
	SpeciesGriffin SpeciesID = 504 // 格里芬
	SpeciesUnicorn SpeciesID = 505 // 独角兽

	// 元素类（无性别）
	SpeciesFireSpirit  SpeciesID = 601 // 火元素
	SpeciesWaterSpirit SpeciesID = 602 // 水元素
)

// SpeciesCategory 物种分类
type SpeciesCategory int

const (
	CategoryUnknown   SpeciesCategory = 0
	CategoryMammal    SpeciesCategory = 1 // 哺乳类
	CategoryAvian     SpeciesCategory = 2 // 鸟类
	CategoryFish      SpeciesCategory = 3 // 鱼类
	CategoryReptile   SpeciesCategory = 4 // 爬行类
	CategoryFantasy   SpeciesCategory = 5 // 幻想类
	CategoryElemental SpeciesCategory = 6 // 元素类
)

// CategoryName 分类名称
func (c SpeciesCategory) Name() string {
	names := map[SpeciesCategory]string{
		CategoryUnknown:   "未知",
		CategoryMammal:    "哺乳类",
		CategoryAvian:     "鸟类",
		CategoryFish:      "鱼类",
		CategoryReptile:   "爬行类",
		CategoryFantasy:   "幻想类",
		CategoryElemental: "元素类",
	}
	if name, ok := names[c]; ok {
		return name
	}
	return "未知"
}

// Species 物种定义
type Species struct {
	ID           SpeciesID          // 物种ID
	Name         string             // 显示名称
	Category     SpeciesCategory    // 分类
	BaseParts    []PartType         // 基础部位（所有物种共有）
	SpecialParts []PartType         // 特殊部位（物种特有）
	Rarity       int                // 稀有度 1-5
	IsHidden     bool               // 是否为隐藏物种
	GenderRule   GenderRule         // 性别规则
	BreedRules   BreedingRules      // 繁衍规则
	Interpreter  GeneInterpreter    // 基因解释器
}

// GenderRule 性别规则
type GenderRule struct {
	AllowedGenders []Gender       // 该物种允许的性别
	DefaultRatio   map[Gender]int // 各性别的默认比例 (百分比)
	CanSelfBreed   bool           // 是否可自我繁殖（分裂）
}

// BreedingRules 繁衍规则
type BreedingRules struct {
	MinStage              Stage // 最低成长阶段
	MinLevel              int   // 最低等级
	MinHappiness          int   // 最低快乐度
	CooldownHours         int   // 有性繁殖冷却时间（小时）
	SelfBreedCooldownHours int   // 分裂繁殖冷却时间（小时）
}

// GeneInterpreter 物种基因解释器接口
// 不同物种实现此接口来解释基因物种特征区域（位置8-15）
type GeneInterpreter interface {
	// GetSpeciesID 返回此解释器对应的物种ID
	GetSpeciesID() SpeciesID

	// InterpretSpecialFeatures 解释物种特有特征
	// 从基因中提取物种特有的外观特征
	InterpretSpecialFeatures(gene Gene) SpecialAppearance

	// GetFeatureNames 获取特征名称映射
	GetFeatureNames() map[PartType][]string
}

// SpeciesRegistry 物种注册表
type SpeciesRegistry struct {
	species      map[SpeciesID]*Species
	interpreters map[SpeciesID]GeneInterpreter
}

// NewSpeciesRegistry 创建物种注册表
func NewSpeciesRegistry() *SpeciesRegistry {
	return &SpeciesRegistry{
		species:      make(map[SpeciesID]*Species),
		interpreters: make(map[SpeciesID]GeneInterpreter),
	}
}

// Register 注册物种
func (r *SpeciesRegistry) Register(species *Species) {
	r.species[species.ID] = species
	if species.Interpreter != nil {
		r.interpreters[species.ID] = species.Interpreter
	}
}

// Get 获取物种定义
func (r *SpeciesRegistry) Get(id SpeciesID) (*Species, bool) {
	species, ok := r.species[id]
	return species, ok
}

// GetInterpreter 获取物种基因解释器
func (r *SpeciesRegistry) GetInterpreter(id SpeciesID) (GeneInterpreter, bool) {
	interpreter, ok := r.interpreters[id]
	return interpreter, ok
}

// GetByCategory 按分类获取物种列表
func (r *SpeciesRegistry) GetByCategory(category SpeciesCategory) []*Species {
	var result []*Species
	for _, s := range r.species {
		if s.Category == category {
			result = append(result, s)
		}
	}
	return result
}

// GetAvailableSpecies 获取可用物种（非隐藏）
func (r *SpeciesRegistry) GetAvailableSpecies() []*Species {
	var result []*Species
	for _, s := range r.species {
		if !s.IsHidden {
			result = append(result, s)
		}
	}
	return result
}

// GetHiddenSpecies 获取隐藏物种
func (r *SpeciesRegistry) GetHiddenSpecies() []*Species {
	var result []*Species
	for _, s := range r.species {
		if s.IsHidden {
			result = append(result, s)
		}
	}
	return result
}

// All 获取所有物种
func (r *SpeciesRegistry) All() []*Species {
	result := make([]*Species, 0, len(r.species))
	for _, s := range r.species {
		result = append(result, s)
	}
	return result
}

// InterpretGene 使用指定物种的解释器解析基因
func (r *SpeciesRegistry) InterpretGene(speciesID SpeciesID, gene Gene) SpecialAppearance {
	if interpreter, ok := r.interpreters[speciesID]; ok {
		return interpreter.InterpretSpecialFeatures(gene)
	}
	return NewSpecialAppearance()
}

// SpeciesPair 物种配对（用于融合表）
type SpeciesPair struct {
	SpeciesA SpeciesID
	SpeciesB SpeciesID
}

// NewSpeciesPair 创建物种配对（自动排序，保证唯一性）
func NewSpeciesPair(a, b SpeciesID) SpeciesPair {
	if a > b {
		a, b = b, a
	}
	return SpeciesPair{SpeciesA: a, SpeciesB: b}
}

// HiddenSpeciesConfig 隐藏物种配置
type HiddenSpeciesConfig struct {
	ResultSpecies    SpeciesID // 产生的隐藏物种
	TriggerThreshold int       // 触发阈值（隐藏物种基因值需大于此值）
	Rarity           int       // 稀有度
}

// SpeciesFusionRegistry 物种融合注册表
type SpeciesFusionRegistry struct {
	fusions map[SpeciesPair]HiddenSpeciesConfig
}

// NewSpeciesFusionRegistry 创建物种融合注册表
func NewSpeciesFusionRegistry() *SpeciesFusionRegistry {
	return &SpeciesFusionRegistry{
		fusions: make(map[SpeciesPair]HiddenSpeciesConfig),
	}
}

// Register 注册融合配置
func (r *SpeciesFusionRegistry) Register(speciesA, speciesB SpeciesID, config HiddenSpeciesConfig) {
	pair := NewSpeciesPair(speciesA, speciesB)
	r.fusions[pair] = config
}

// GetFusion 获取融合配置
func (r *SpeciesFusionRegistry) GetFusion(speciesA, speciesB SpeciesID) (HiddenSpeciesConfig, bool) {
	pair := NewSpeciesPair(speciesA, speciesB)
	config, ok := r.fusions[pair]
	return config, ok
}

// CanFuse 检查两个物种是否可以融合
func (r *SpeciesFusionRegistry) CanFuse(speciesA, speciesB SpeciesID) bool {
	pair := NewSpeciesPair(speciesA, speciesB)
	_, ok := r.fusions[pair]
	return ok
}

// CheckFusionTrigger 检查是否触发融合
func (r *SpeciesFusionRegistry) CheckFusionTrigger(speciesA, speciesB SpeciesID, childGene Gene) (SpeciesID, bool) {
	config, ok := r.GetFusion(speciesA, speciesB)
	if !ok {
		return SpeciesUnknown, false
	}

	// 检查隐藏物种基因值是否超过阈值
	hiddenTrigger := childGene.HiddenSpeciesTrigger()
	if hiddenTrigger > config.TriggerThreshold {
		return config.ResultSpecies, true
	}

	return SpeciesUnknown, false
}

// DefaultGenderRule 默认性别规则（雄性/雌性各50%）
func DefaultGenderRule() GenderRule {
	return GenderRule{
		AllowedGenders: []Gender{GenderMale, GenderFemale},
		DefaultRatio:   map[Gender]int{GenderMale: 50, GenderFemale: 50},
		CanSelfBreed:   false,
	}
}

// AsexualGenderRule 无性别规则（可分裂繁殖）
func AsexualGenderRule() GenderRule {
	return GenderRule{
		AllowedGenders: []Gender{GenderNone},
		DefaultRatio:   map[Gender]int{GenderNone: 100},
		CanSelfBreed:   true,
	}
}

// HermaphroditeGenderRule 雌雄同体规则
func HermaphroditeGenderRule() GenderRule {
	return GenderRule{
		AllowedGenders: []Gender{GenderHermaphrodite},
		DefaultRatio:   map[Gender]int{GenderHermaphrodite: 100},
		CanSelfBreed:   false,
	}
}

// MixedGenderRule 混合性别规则（支持雄性、雌性和雌雄同体）
func MixedGenderRule(maleRatio, femaleRatio, hermaphroditeRatio int) GenderRule {
	return GenderRule{
		AllowedGenders: []Gender{GenderMale, GenderFemale, GenderHermaphrodite},
		DefaultRatio: map[Gender]int{
			GenderMale:          maleRatio,
			GenderFemale:        femaleRatio,
			GenderHermaphrodite: hermaphroditeRatio,
		},
		CanSelfBreed: false,
	}
}

// DefaultBreedingRules 默认繁衍规则
func DefaultBreedingRules() BreedingRules {
	return BreedingRules{
		MinStage:              StageAdult,
		MinLevel:              10,
		MinHappiness:          50,
		CooldownHours:         24,
		SelfBreedCooldownHours: 48,
	}
}
