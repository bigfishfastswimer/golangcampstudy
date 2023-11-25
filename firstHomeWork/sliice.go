package slice

import "log"

func SliceDelete[T any](index int, slice []T) []T {
	if index < 0 || index >= len(slice) {
		log.Printf("invalid index input: %d", index)
		return slice // 或者处理错误
	}

	// 删除指定索引的元素
	copy(slice[index:], slice[index+1:])
	slice = slice[:len(slice)-1]

	// 缩容 平衡内存使用和避免频繁的内存分配。
	if cap(slice) > 2*len(slice) {
		newSlice := make([]T, len(slice))
		copy(newSlice, slice)
		return newSlice
	}
	return slice
}
