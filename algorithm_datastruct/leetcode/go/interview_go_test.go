package interview_go

import (
	"fmt"
	"sync"
	"testing"
)

//https://github.com/lifei6671/interview-go/blob/master/question/q001.md

func printLetterAndNumber() {
	letter := make(chan struct{})
	number := make(chan struct{})
	wg := sync.WaitGroup{}

	// 打印数字，并通知结束
	go func() {
		i := 1
		for {
			<-number
			fmt.Print(i)
			i+=1
			fmt.Print(i)
			i+=1
			letter <- struct{}{}
		}
	}()

	// 打印字母
	wg.Add(1)
	go func() {
		i := 'A'
		for {
			if i > 'Z' {
				// 如果需要让打印数字的协程页退出的化，需要在这里设置一个ch，然后让上面的协程select一下，select到了就退出
				wg.Done()
				return
			}
			<-letter
			fmt.Print(string(i))
			i+=1
			fmt.Print(string(i))
			i+=1
			number <- struct{}{} //
		}
	}()
	number <- struct{}{}
	wg.Wait()

}

func TestQ001(t *testing.T) {
	printLetterAndNumber()
}