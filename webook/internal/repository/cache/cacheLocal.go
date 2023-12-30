package cache

import (
	"errors"
	"fmt"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

type CodeCacheLocal struct {
	cache *expirable.LRU[string, codeItem]
}

func (c *CodeCacheLocal) Key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func NewCodeCacheLocal(c *expirable.LRU[string, codeItem]) *CodeCacheLocal {
	return &CodeCacheLocal{
		cache: c,
	}
}

type codeItem struct {
	code string
	cnt  int
}

func (c *CodeCacheLocal) Set(biz, phone, code string) error {
	key := c.Key(biz, phone)

	if !c.cache.Contains(key) {
		return errors.New("key is not exist")
	}
	r, ok := c.cache.Get(key)
	if !ok {
		fmt.Printf("key is expired")
		// 验证码没有或者过期
		c.cache.Add(key, codeItem{
			code: code,
			cnt:  3,
		})
		return nil
	}
	return nil
}

func (c *CodeCacheLocal) Verify(biz, phone, inputCode string) (bool, error) {
	key := c.Key(biz, phone)
	if !c.cache.Contains(key) {
		return false, ErrKeyNotExist
	}
	r, ok := c.cache.Get(key)
	if !ok {
		return false, errors.New("key is expired")
	}
	if r.cnt <= 0 {
		return false, ErrCodeVerifyTooMany
	}
	r.cnt--
	return r.code == inputCode, nil
}
