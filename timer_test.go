package main

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	start := time.Now().Unix()
	timer := time.NewTimer(time.Second) // 马上开始计时，只能使用一次 <-timer.C
	time.Sleep(time.Second * 3)
	<-timer.C
	fmt.Println("timer 2", time.Now().Unix()-start)

	// stop 可以提前终止 timer，避免执行定时逻辑
	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("timer2 execute") // stop 之后，不会执行该语句
	}()
	stopped := timer2.Stop()
	if stopped {
		fmt.Println("timer2 stopped")
	}
	time.Sleep(time.Second * 2)

}
