package core

import (
	"log"
	"time"
)

var store map[string]*Obj

type Obj struct {
	val        interface{}
	ExpieresAt int64
}

func init() {
	store = make(map[string]*Obj)

}

func newObj(val interface{}, expDurationMs int64) *Obj {

	var expAt int64 = -1

	if expDurationMs > 0 {
		expAt = time.Now().UnixMilli() + expDurationMs
	}

	return &Obj{
		val:        val,
		ExpieresAt: expAt,
	}
}

func Put(key string, val interface{}, expDurationMs int64) {
	store[key] = newObj(val, expDurationMs)
}

func Get(key string) (interface{}, bool) {

	obj, exists := store[key]

	if !exists {
		return nil, false
	}

	if obj.ExpieresAt != -1 && time.Now().UnixMilli() > obj.ExpieresAt {
		delete(store, key)
		return nil, false
	}

	return obj.val, true

}

func GetTTL(key string) int64 {

	obj, exists := store[key]

	if !exists {
		return -2
	}

	if obj.ExpieresAt == -1 {
		return -1
	}

	if time.Now().UnixMilli() > obj.ExpieresAt {
		delete(store, key)
		return -2
	}

	log.Println("TTL for key", key, "is", (obj.ExpieresAt-time.Now().UnixMilli())/1000)

	return (obj.ExpieresAt - time.Now().UnixMilli()) / 1000
}

func Delete(key string) int {

	_, exists := store[key]

	if exists {
		delete(store, key)
		return 1
	}

	return 0
}

func ExpireAt(key string, expireDurationMs int64) int {

	obj, exists := store[key]

	if !exists {
		return 0
	}
	
	obj.ExpieresAt = time.Now().UnixMilli() + expireDurationMs
	return 1
}

