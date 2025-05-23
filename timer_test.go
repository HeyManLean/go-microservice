package main

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	/*
		timer 是一次性的
	*/
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

func TestTicker(t *testing.T) {
	/*
		ticker 是定时重复的
	*/
	ticker := time.NewTicker(time.Millisecond * 500)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fmt.Println("ticker", time.Now().Unix())
			}
		}
	}()
	time.Sleep(time.Millisecond * 1600)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

func TestTimeout(t *testing.T) {
	/*
		通过 select 和 channel 实现 timeout
	*/
	w1 := make(chan int)
	go func(c chan<- int) {
		time.Sleep(time.Second * 2)
		c <- 1
	}(w1)
	select {
	case val := <-w1:
		fmt.Println("result: ", val)
	case <-time.After(time.Second):
		fmt.Println("timeout 1")
	}

	// 使用 default 避免阻塞，否则报错：fatal error: all goroutines are asleep - deadlock!
	w2 := make(chan int)
	select {
	case val := <-w2:
		fmt.Println("result: ", val)
	default:
		fmt.Println("no result")
	}
}
