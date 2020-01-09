package basemodule

import (
	"fmt"
	"sync"

	"github.com/liuqiqi-Y/qiqiCrawler/util"
)

// register 代表组件注册器的接口。
type RegisterIF interface {
	// Register 用于注册组件实例。
	Register(module BaseModuleIF) (bool, error)
	// Unregister 用于注销组件实例。
	Unregister(mid MID) (bool, error)
	// Get 用于获取一个指定类型的组件的实例。
	// 本函数应该基于负载均衡策略返回实例。
	Get(moduleType string) (BaseModuleIF, error)
	// GetAllByType 用于获取指定类型的所有组件实例。
	GetAllByType(moduleType string) (map[MID]BaseModuleIF, error)
	// GetAll 用于获取所有组件实例。
	GetAll() map[MID]BaseModuleIF
	// Clear 会清除所有的组件注册记录。
	Clear()
}

// myregister 代表组件注册器的实现类型。
type Register struct {
	// moduleTypeMap 代表组件类型与对应组件实例的映射。
	moduleTypeMap map[string]map[MID]BaseModuleIF
	// rwlock 代表组件注册专用读写锁。
	mu sync.Mutex
}

// NewRegister 用于创建一个组件注册器的实例。
func NewRegister() RegisterIF {
	return &Register{
		moduleTypeMap: map[string]map[MID]BaseModuleIF{},
	}
}

func (register *Register) Register(module BaseModuleIF) (bool, error) {
	if module == nil {
		util.Trace.Println("Register module is NULL")
		return false, fmt.Errorf("Module can not be NULL")
	}
	mid := module.ID()
	parts, err := SplitMID(mid)
	if err != nil {
		util.Trace.Printf("Failed to split MID: %s", err.Error())
		return false, fmt.Errorf("Failed to split MID: %s", err.Error())
	}
	moduleType := legalLetterTypeMap[parts[0]]
	register.mu.Lock()
	defer register.mu.Unlock()
	modules := register.moduleTypeMap[moduleType]
	if modules == nil {
		modules = map[MID]BaseModuleIF{}
		modules[mid] = module
		register.moduleTypeMap[moduleType] = modules
		return true, nil
	}
	if _, ok := modules[mid]; ok {
		util.Trace.Printf("This module already exist: %s", mid)
		return false, fmt.Errorf("This module already exist: %s", mid)
	}
	modules[mid] = module
	return true, nil
}

func (register *Register) Unregister(mid MID) (bool, error) {
	parts, err := SplitMID(mid)
	if err != nil {
		util.Trace.Printf("Failed to split MID: %s", err.Error())
		return false, fmt.Errorf("Failed to split MID: %s", err.Error())
	}
	moduleType := legalLetterTypeMap[parts[0]]
	register.mu.Lock()
	defer register.mu.Unlock()
	if modules, ok := register.moduleTypeMap[moduleType]; ok {
		if _, ok := modules[mid]; ok {
			delete(modules, mid)
			return true, nil
		}
	}
	return false, fmt.Errorf("this module has not been registered :%s", mid)
}

func (register *Register) GetAllByType(moduleType string) (map[MID]BaseModuleIF, error) {
	if moduleType == "" || legalTypeLetterMap[moduleType] == "" {
		util.Trace.Printf("moduel type is null or is illegal: %s", moduleType)
		return nil, fmt.Errorf("moduel type is null or is illegal: %s", moduleType)
	}
	register.mu.Lock()
	defer register.mu.Unlock()
	modules := register.moduleTypeMap[moduleType]
	if modules == nil || len(modules) == 0 {
		util.Trace.Printf("this type module has not been registered: %s", moduleType)
		return nil, fmt.Errorf("this type module has not been registered: %s", moduleType)
	}
	result := map[MID]BaseModuleIF{}
	for mid, module := range modules {
		result[mid] = module
	}
	return result, nil
}

func (register *Register) Get(moduleType string) (BaseModuleIF, error) {
	modules, err := register.GetAllByType(moduleType)
	if err != nil {
		util.Trace.Printf("Failed to get module: %s", err.Error())
		return nil, fmt.Errorf("Failed to get module: %s", err.Error())
	}
	minScore := uint64(0)
	var selectedModule BaseModuleIF
	for _, module := range modules {
		SetScore(module)
		score := module.Score()
		if minScore == 0 || score < minScore {
			selectedModule = module
			minScore = score
		}
	}
	return selectedModule, nil
}

func (register *Register) GetAll() map[MID]BaseModuleIF {
	result := map[MID]BaseModuleIF{}
	register.mu.Lock()
	defer register.mu.Unlock()
	for _, modules := range register.moduleTypeMap {
		for mid, module := range modules {
			result[mid] = module
		}
	}
	return result
}

func (register *Register) Clear() {
	register.mu.Lock()
	defer register.mu.Unlock()
	register.moduleTypeMap = map[string]map[MID]BaseModuleIF{}
}
