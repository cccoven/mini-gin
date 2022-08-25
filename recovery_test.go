package mini_gin

import (
	"fmt"
	"testing"
)

func testRecover() {
	// panic 导致的程序退出前，会先处理完当前协程上已经 defer 的任务
	defer func() {
		fmt.Println("defer func")
	
		// Go 提供了 recover 函数，可以避免因为 panic 发生而导致整个程序终止，recover 函数只在 defer 中生效
		if err := recover(); err != nil {
			fmt.Println("recover success")
		}
	}()

	arr := []int{1, 2, 3}
	fmt.Println(arr[3]) // 数组越界，触发 panic
	fmt.Println("after panic")
}

func TestPanicRecovery(t *testing.T) {
	testRecover()
	fmt.Println("after recover")
}
