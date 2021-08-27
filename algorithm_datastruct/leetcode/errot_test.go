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



func TestError(t *testing.T)  {
	err := &TimeoutErr{err: "timeout for 4ms"}
	fmt.Println(err)
	fmt.Println(IsTimeout(err))
}
