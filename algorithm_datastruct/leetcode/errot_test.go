package leetcode

import (
	"fmt"
	"testing"
)

type timeout interface {
	Timeout() bool
}

func IsTimeout(err error) bool {
	t, ok := err.(timeout)
	return ok && t.Timeout()
}

type TimeoutErr struct {
	err string
}

func (t *TimeoutErr) Error() string {
	return t.err
}

func (t *TimeoutErr) Timeout() bool {
	return true
}

var constructMap = make(map[string]string)

func TestError(t *testing.T) {
	//err := &TimeoutErr{err: "timeout for 4ms"}
	//fmt.Println(err)
	//fmt.Println(IsTimeout(err))

	constructMap["ni"] = "ki"
	t.Log(constructMap["ni"])
}

const (
	mutexLocked = 1 << iota
	mutexWoken
	mutexWaiterShift = iota
)

func TestIoTA(t *testing.T) {
	fmt.Println(mutexLocked)
	fmt.Println(mutexWoken)
	fmt.Println(mutexWaiterShift)
	fmt.Println(1 << mutexWaiterShift)
	fmt.Print()
}
