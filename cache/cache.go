package cache

import (
	"fmt"
)

type Cache interface {
	//通过主键,获取值
	Get(key string) interface{}
	//通过主键，存放值，并设置过期时间最小单位秒
	Put(key string, v interface{}, expired int64) error
	//删除主键
	Delete(key string) error
	//判定主键是否存在
	IsExists(key string) bool
	//清空户缓存中所有的键值
	ClearAll() error
	//启动缓存服务，并启动缓存的GC服务
	StartAndGC(config string) error
}

//存放允许的Cache类型
var adpaters = make(map[string]Cache)

//注册已经实现的cache类
func RegisterCache(name string, c Cache) {
	if c == nil {
		panic("Register cache is nil")
	}

	if _, ok := adpaters[name]; ok {
		panic("Register cache name repulicated,name = " + name)
	}
	adpaters[name] = c
}

//生成一个cache,输入cache类型和配置文件
func NewCache(cacheTypeName, config string) (adpater Cache, err error) {
	adpater, ok := adpaters[cacheTypeName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name [" + cacheTypeName + "]")
		return
	}
	err = adpater.StartAndGC(config)
	if err != nil {
		adpater = nil
	}
	return
}
