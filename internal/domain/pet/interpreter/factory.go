// Package interpreter 物种基因解释器
package interpreter

import "pets-server/internal/domain/pet"

// InterpreterFactory 解释器工厂
// 通过字符串类型名称获取对应的基因解释器实例
type InterpreterFactory struct {
	creators map[string]func() pet.GeneInterpreter
}

// NewInterpreterFactory 创建解释器工厂
func NewInterpreterFactory() *InterpreterFactory {
	f := &InterpreterFactory{
		creators: make(map[string]func() pet.GeneInterpreter),
	}

	// 注册所有解释器类型
	// 哺乳类
	f.Register("feline", func() pet.GeneInterpreter { return NewFelineInterpreter() })
	f.Register("canine", func() pet.GeneInterpreter { return NewCanineInterpreter() })

	// 鸟类
	f.Register("parrot", func() pet.GeneInterpreter { return NewParrotInterpreter() })
	f.Register("owl", func() pet.GeneInterpreter { return NewOwlInterpreter() })

	// 水生类
	f.Register("goldfish", func() pet.GeneInterpreter { return NewGoldfishInterpreter() })
	f.Register("tropical_fish", func() pet.GeneInterpreter { return NewTropicalFishInterpreter() })

	// 幻想类
	f.Register("slime", func() pet.GeneInterpreter { return NewSlimeInterpreter() })
	f.Register("phoenix", func() pet.GeneInterpreter { return NewPhoenixInterpreter() })
	f.Register("dragon", func() pet.GeneInterpreter { return NewDragonInterpreter() })
	f.Register("griffin", func() pet.GeneInterpreter { return NewGriffinInterpreter() })
	f.Register("unicorn", func() pet.GeneInterpreter { return NewUnicornInterpreter() })

	return f
}

// Register 注册解释器创建函数
func (f *InterpreterFactory) Register(name string, creator func() pet.GeneInterpreter) {
	f.creators[name] = creator
}

// Get 获取解释器实例
// 如果类型名称不存在，返回 nil 和 false
func (f *InterpreterFactory) Get(typeName string) (pet.GeneInterpreter, bool) {
	creator, ok := f.creators[typeName]
	if !ok {
		return nil, false
	}
	return creator(), true
}

// Has 检查是否存在指定类型的解释器
func (f *InterpreterFactory) Has(typeName string) bool {
	_, ok := f.creators[typeName]
	return ok
}

// Types 获取所有已注册的解释器类型名称
func (f *InterpreterFactory) Types() []string {
	types := make([]string, 0, len(f.creators))
	for t := range f.creators {
		types = append(types, t)
	}
	return types
}
