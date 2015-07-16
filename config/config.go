package config

import "fmt"

const (
	INI_CONFIG = "ini"
)

type ConfigManager interface {
	Set(key, val string) error
	String(key string) string
	Strings(key string) []string
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	SaveConfigFile(filename string, bsort bool) error
}

type Config interface {
	Parse(filename string) (ConfigManager, error)
}

var adapters = make(map[string]Config)

func Register(name string, adapter Config) {
	if adapter == nil {
		panic("config: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("config: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

func NewConfig(adapterName, fileaname string) (ConfigManager, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("config: unknown adaptername %q (forgotten import?)", adapterName)
	}
	return adapter.Parse(fileaname)
}
