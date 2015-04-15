package cache

import (
	"testing"
)

func TestMemoryCache(t *testing.T) {
	c, err := NewCache("memorycache", `{"interval":30,"cap":1024}`)
	if err != nil {
		println("create cache error")
	}
	c.Put("yjx", "nihao", 20)
	c.Put("yjx2", "nihao2", 20)
	//yjx := c.Get("yjx")
	println(" get yjx = ", c.Get("yjx").(string))
	println(" get yjx2 = ", c.Get("yjx2").(string))
	c.Delete("yjx")
	println(" get yjx after delete = ", c.Get("yjx").(string))
	c.ClearAll()
	println(" get yjx2 after clear all = ", c.Get("yjx2").(string))
}
