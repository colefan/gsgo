package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var (
	DefaultEvery int = 60 //1 minutes
)

//内存缓存，不是第三方的memcache
//缓存对象
type MemoryItem struct {
	Data           interface{} //data to store
	LastUpdateTime int64       //last updatetime
	Expired        int64       //expired timestamps
}

//缓存
type MemoryCache struct {
	lock  sync.RWMutex           //读写锁，读写要分离
	items map[string]*MemoryItem //缓存存储对象
	Every int                    //每隔多少时间进行一次缓存清理
	dur   time.Duration          //每隔多少时间
}

func NewMemoryCache() *MemoryCache {
	cache := MemoryCache{items: make(map[string]*MemoryItem)}
	return &cache
}

func (c *MemoryCache) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if item, ok := c.items[key]; ok {
		if (item.Expired > 0) && (time.Now().Unix()-item.LastUpdateTime > item.Expired) {
			go c.Delete(key)
			return nil
		} else {
			return item.Data
		}
	}
	return nil
}

//当expired为0时，表示永不过期，否则按具体的过期时间来处理，单位是S
func (c *MemoryCache) Put(key string, v interface{}, expired int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items[key] = &MemoryItem{
		Data:           v,
		LastUpdateTime: time.Now().Unix(),
		Expired:        expired}
	return nil
}

func (c *MemoryCache) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.items[key]; !ok {
		return nil
	}
	delete(c.items, key)
	if _, ok := c.items[key]; ok {
		return fmt.Errorf("delete key from MemoryCache error,key = " + key)
	}
	return nil
}

func (c *MemoryCache) IsExists(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.items[key]
	return ok
}

func (c *MemoryCache) ClearAll() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items = make(map[string]*MemoryItem)
	return nil
}

//config:`{"interval":60,}`
func (c *MemoryCache) StartAndGC(config string) error {
	var cf = make(map[string]int)
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["interval"]; !ok {
		cf["interval"] = DefaultEvery
	}
	c.Every = cf["interval"]
	dur, err := time.ParseDuration(fmt.Sprintf("%ds", cf["interval"]))
	if err != nil {
		return err
	}
	c.dur = dur
	go c.gc()
	return nil
}

func (c *MemoryCache) gc() {
	if c.Every < 1 {
		return
	}

	for {
		<-time.After(c.dur)
		if c.items == nil {
			return
		}
		clearCount := 0
		for name, _ := range c.items {
			if c.item_expired(name) {
				clearCount++
			}
			//每清除10W对象后，退出循环等待下一次清理
			if clearCount > 100000 {
				break
			}
		}
	}

}

func (c *MemoryCache) item_expired(name string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	it, ok := c.items[name]
	if !ok {
		return false
	}

	if it.Expired > 0 && (time.Now().Unix()-it.LastUpdateTime > it.Expired) {
		delete(c.items, name)
		return true
	} else {
		return false
	}
}

func init() {
	RegisterCache("memorycache", NewMemoryCache())
}
