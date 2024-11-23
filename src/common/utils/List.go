package utils

import (
	"errors"
	"fmt"
)

// List 表示一个基于切片的列表 注意 改对象无锁的 因此不支持并发操作
type List[T comparable] struct {
	data []T
}

// NewList 创建一个新的 List 实例
func NewList[T comparable]() List[T] {
	return List[T]{data: make([]T, 0)}
}

// Add 添加一个元素到列表中
func (l *List[T]) Add(element T) {
	for _, v := range l.data {
		if v == element {
			return // 元素已存在，不重复添加
		}
	}
	l.data = append(l.data, element)
}

// Remove 删除列表中的一个元素
func (l *List[T]) Remove(element T) {
	for i, v := range l.data {
		if v == element {
			l.data = append(l.data[:i], l.data[i+1:]...)
			return
		}
	}
}

func (l *List[T]) RemoveIndex(index int) (T, error) {
	var zeroValue T
	if !l.isValidIndex(index) {
		return zeroValue, fmt.Errorf("index out of range")
	}

	// 实现移除元素的逻辑
	value := l.data[index]
	l.data = append(l.data[:index], l.data[index+1:]...)
	return value, nil
}

// Contains 检查列表中是否包含某个元素
func (l *List[T]) Contains(element T) bool {
	for _, v := range l.data {
		if v == element {
			return true
		}
	}
	return false
}

// Size 返回列表的大小
func (l *List[T]) Size() int {
	return len(l.data)
}

// String 返回列表的字符串表示
func (l *List[T]) String() string {
	return fmt.Sprintf("%v", l.data)
}

// Clear 清空列表
func (l *List[T]) Clear() {
	l.data = []T{}
}

// Get retrieves an element from the List at the specified index.
// If the index is out of range, it returns the zero value of type T and an error.
// This function demonstrates how to handle errors and return values in Go.
func (l *List[T]) Get(index int) (T, error) {
	// Initialize a zero value variable of type T to return in case of an error.
	var zeroValue T

	// Check if the index is out of range.
	if !l.isValidIndex(index) {
		// If the index is out of range, return the zero value and an error.
		return zeroValue, errors.New("index out of range")
	}

	// If the index is valid, return the element at the index and nil error.
	return l.data[index], nil
}

// IndexOf 在列表中查找指定元素的索引。
// 如果元素存在于列表中，则返回其索引；如果元素不存在，则返回-1。
// 这个方法通过遍历列表数据来寻找匹配的元素。
// 参数 element 是需要查找的元素。
// 返回值是元素在列表中的索引，如果没有找到，则返回-1。
func (l *List[T]) IndexOf(element T) int {
	for index, v := range l.data {
		if v == element {
			return index
		}
	}
	return -1
}

// 辅助函数，用于检查索引是否有效
func (l *List[T]) isValidIndex(index int) bool {
	return index >= 0 && index < l.Size()
}

func (l *List[T]) RemoveFirst() (T, error) {
	var zeroValue T
	if l.Size() > 0 {
		zeroValue = l.data[0]
		l.data = l.data[1:]
		return zeroValue, nil
	}
	return zeroValue, errors.New("list is empty")
}

// ForEachAndClear 遍历列表中的每个元素，并在遍历结束后清空列表。
// 参数 forEachFunc 是一个接受列表元素类型 T 的函数，用于在遍历过程中对每个元素执行操作。
// 此方法先遍历列表中的每个元素并调用提供的函数，然后清空列表，释放资源。
func (l *List[T]) ForEachAndClear(forEachFunc func(T)) {
	for _, item := range l.data {
		forEachFunc(item)
	}
	l.Clear()

}

func (l *List[T]) ForEach(f func(T)) {
	for _, task := range l.data {
		f(task)
	}
}
