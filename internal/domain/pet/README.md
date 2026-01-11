# 宠物领域模块 (Pet Domain)

本模块实现了网络版拓麻歌子游戏的核心宠物系统，包括基因系统、物种系统、性别系统和繁衍系统。

## 目录结构

```
internal/domain/pet/
├── gene.go              # 基因系统（40位十六进制）
├── species.go           # 物种定义与注册表
├── gender.go            # 性别系统
├── breeding.go          # 繁衍系统
├── entity.go            # Pet 实体（聚合根）
├── appearance.go        # 外观值对象
├── personality.go       # 性格值对象
├── skill.go             # 技能值对象
├── service.go           # 领域服务
├── event.go             # 领域事件
├── repository.go        # 仓储接口
└── interpreter/         # 物种基因解释器
    ├── registry.go      # 物种注册表初始化
    ├── feline.go        # 猫科/犬科解释器
    ├── avian.go         # 鸟类解释器
    ├── aquatic.go       # 水生类解释器
    ├── fantasy.go       # 幻想类解释器
    └── hidden.go        # 隐藏物种解释器
```

## 核心设计理念

**"物种决定模板，基因决定表现"**

- **物种 (Species)**: 定义特征的"结构"（有哪些部位、如何解读基因）
- **基因 (Gene)**: 定义特征的"值"（每个部位的具体样式）

同一段基因在不同物种下会产生不同的外观。例如基因片段 `0x7` 在猫身上是"立耳"，在鸟身上是"羽冠"。

---

## 基因系统

### 基因结构（40位十六进制）

```
位置分布:
┌─────────┬─────────┬─────────┬─────────┬─────────┐
│ 0-7     │ 8-15    │ 16-23   │ 24-31   │ 32-39   │
│ 通用外观 │ 物种特征 │ 性格    │ 技能/能力│ 遗传隐藏│
└─────────┴─────────┴─────────┴─────────┴─────────┘
```

### 各区域详细说明

| 区域 | 位置 | 内容 |
|-----|------|------|
| 通用外观 | 0-7 | 主色、副色、体型、眼睛、花纹 |
| 物种特征 | 8-15 | 4个特征槽位 + 4个修饰槽位（由物种解释器解读） |
| 性格 | 16-23 | 活跃度、贪吃度、社交度、好奇度、脾气、忠诚度、智力、玩乐度 |
| 技能/能力 | 24-31 | 主技能、技能强度、副技能、特殊能力、成长速率 |
| 遗传隐藏 | 32-39 | 突变因子、进化倾向、隐藏物种触发、隐性基因 |

### 使用示例

```go
// 生成随机基因
gene := pet.GenerateGene()

// 从字符串创建基因
gene := pet.NewGene("a1b2c3d4e5f6789012345678901234567890abcd")

// 读取基因值
primaryHue := gene.PrimaryColorHue()    // 主色色相 0-360
activity := gene.ActivityTrait()        // 活跃度 0-100
specialA := gene.SpecialA()             // 物种特征A 0-15
```

---

## 物种系统

### 物种分类

| 分类 | Category | 示例物种 |
|-----|----------|---------|
| 哺乳类 | CategoryMammal | 猫、狗、兔 |
| 鸟类 | CategoryAvian | 鹦鹉、猫头鹰 |
| 鱼类 | CategoryFish | 金鱼、热带鱼 |
| 爬行类 | CategoryReptile | 蜥蜴、龟 |
| 幻想类 | CategoryFantasy | 龙、凤凰、史莱姆 |
| 元素类 | CategoryElemental | 火元素、水元素 |

### 物种特征映射

每个物种有 4 个特征槽位，由物种解释器解读：

| 物种 | 特征A | 特征B | 特征C | 特征D |
|-----|-------|-------|-------|-------|
| 猫科 | 耳朵 | 尾巴 | 毛纹 | 胡须 |
| 鸟类 | 翅膀 | 喙 | 羽冠 | 尾羽 |
| 鱼类 | 背鳍 | 鳞片 | 尾鳍 | 触须 |
| 龙类 | 翅膀 | 角 | 鳞甲 | 尾巴 |

### 使用示例

```go
import "pets-server/internal/domain/pet/interpreter"

// 初始化物种注册表
speciesRegistry := interpreter.InitDefaultRegistry()
fusionRegistry := interpreter.InitDefaultFusionRegistry()

// 获取物种定义
species, ok := speciesRegistry.Get(pet.SpeciesCat)

// 解析物种特有外观
if interpreter, ok := speciesRegistry.GetInterpreter(pet.SpeciesCat); ok {
    specialAppearance := interpreter.InterpretSpecialFeatures(gene)
}

// 获取可用物种列表
availableSpecies := speciesRegistry.GetAvailableSpecies()
```

---

## 性别系统

### 性别类型

| 类型 | 值 | 说明 |
|-----|---|------|
| GenderNone | 0 | 无性别（元素生物、史莱姆） |
| GenderMale | 1 | 雄性 |
| GenderFemale | 2 | 雌性 |
| GenderHermaphrodite | 3 | 雌雄同体（蜗牛、凤凰） |

### 性别配对规则

```
雄性 ↔ 雌性          ✅ 可繁殖
雄性 ↔ 雌雄同体      ✅ 可繁殖
雌性 ↔ 雌雄同体      ✅ 可繁殖
雌雄同体 ↔ 雌雄同体  ✅ 可繁殖
无性别               ✅ 仅可分裂繁殖
```

### 物种性别规则

```go
// 默认性别规则（雄性/雌性各50%）
rule := pet.DefaultGenderRule()

// 无性别规则（可分裂繁殖）
rule := pet.AsexualGenderRule()

// 雌雄同体规则
rule := pet.HermaphroditeGenderRule()

// 混合性别规则
rule := pet.MixedGenderRule(40, 40, 20) // 40%雄 40%雌 20%雌雄同体
```

---

## 繁衍系统

### 繁殖类型

1. **有性繁殖**: 需要两只性别兼容的宠物
2. **分裂繁殖**: 无性别物种可自我繁殖

### 繁殖条件

```go
type BreedingRules struct {
    MinStage              Stage // 最低成长阶段（需成年）
    MinLevel              int   // 最低等级
    MinHappiness          int   // 最低快乐度
    CooldownHours         int   // 有性繁殖冷却（小时）
    SelfBreedCooldownHours int  // 分裂繁殖冷却（小时）
}
```

### 物种遗传规则

**同物种繁殖**:
- 100% 产生相同物种

**跨物种繁殖**:
- 45% 继承父方物种
- 45% 继承母方物种
- 10% 触发隐藏物种检查

### 基因遗传算法

```
每个基因位点:
├─ 45% 继承父方
├─ 45% 继承母方
├─ 7%  混合（取平均值）
└─ 3%  突变（完全随机）
```

### 使用示例

```go
// 创建领域服务
domainService := pet.NewDomainService(repo, speciesRegistry, fusionRegistry)

// 有性繁殖
result, err := domainService.BreedPets(parent1, parent2, "小宝宝", ownerID)
if err != nil {
    // 处理错误
}
child := result.Child
isHidden := result.IsHidden // 是否触发隐藏物种

// 分裂繁殖
result, err := domainService.SelfBreedPet(slime, "小史莱姆", ownerID)

// 检查是否可繁殖
err := domainService.CanBreedPair(parent1, parent2)

// 预测后代物种概率
predictions := domainService.PredictOffspringSpecies(parent1, parent2)
```

---

## 隐藏物种

### 融合配置

| 父母组合 | 隐藏物种 | 触发阈值 |
|---------|---------|---------|
| 猫 + 鹦鹉 | 格里芬 | 200 |
| 猫 + 猫头鹰 | 格里芬 | 190 |
| 金鱼 + 蜥蜴 | 龙 | 220 |
| 凤凰 + 蜥蜴 | 龙 | 180 |
| 兔子 + 凤凰 | 独角兽 | 200 |
| 鹦鹉 + 火元素 | 凤凰 | 190 |

### 触发条件

1. 父母物种组合在融合表中
2. 子代基因的隐藏物种触发值（位置34-35）> 触发阈值
3. 随机命中 10% 的隐藏物种检查

---

## 创建宠物

### 方式一：指定物种

```go
pet, err := domainService.CreatePet(userID, "小花", pet.SpeciesCat)
```

### 方式二：随机物种（按稀有度加权）

```go
pet := domainService.CreateRandomPet(userID, "神秘宠物")
```

### 方式三：通过繁殖

```go
result, err := domainService.BreedPets(parent1, parent2, "宝宝", ownerID)
child := result.Child
```

---

## 评分计算

```go
score := domainService.CalculatePetScore(pet)

// 评分因素:
// - 等级 × 10
// - 技能强度 × 50
// - 物种稀有度 × 100
// - 状态平均值
// - 进化阶段 × 100
// - 代数 × 20
```

---

## 扩展指南

### 添加新物种

1. 在 `species.go` 中添加物种ID常量
2. 创建物种解释器（实现 `GeneInterpreter` 接口）
3. 在 `interpreter/registry.go` 中注册物种

```go
// 1. 添加物种ID
const SpeciesNewPet SpeciesID = 701

// 2. 创建解释器
type NewPetInterpreter struct{}

func (n *NewPetInterpreter) GetSpeciesID() pet.SpeciesID {
    return pet.SpeciesNewPet
}

func (n *NewPetInterpreter) InterpretSpecialFeatures(gene pet.Gene) pet.SpecialAppearance {
    // 实现特征解析
}

// 3. 注册物种
registry.Register(&pet.Species{
    ID:          pet.SpeciesNewPet,
    Name:        "新物种",
    Category:    pet.CategoryFantasy,
    Interpreter: NewNewPetInterpreter(),
    // ...
})
```

### 添加新融合

```go
fusionRegistry.Register(pet.SpeciesA, pet.SpeciesB, pet.HiddenSpeciesConfig{
    ResultSpecies:    pet.SpeciesHidden,
    TriggerThreshold: 200,
    Rarity:           5,
})
```
