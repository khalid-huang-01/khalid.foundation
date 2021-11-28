package leetcode

import (
	"fmt"
	"strings"
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
	delete(constructMap, "ki")
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

type Student struct {
	name string
}

func TestMap(t *testing.T) {
	a := "mark."
	b := strings.Split(a, ".")
	fmt.Println(b)
	fmt.Println(len(b))
>>>>>>> fe172aefa6156ca5315526a07ea292cf982afa3c
}
