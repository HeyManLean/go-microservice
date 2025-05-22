package main

import (
	"fmt"
	"testing"
	"time"
)

func worker(done chan bool) {
	/*
		chan 表示在函数周期中可读可写
	*/
	wait(done)
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Println("done")
	finish(done)
}

func wait(done <-chan bool) {
	/* <-chan 在该函数周期中只读*/
	<-done
}

func finish(done chan<- bool) {
	/*
		chan<- 表示在函数周期中只写
	*/
	done <- true
}

func selectChannels() {
	// 多个channel的select，先获取到值的chan先处理
	w1 := make(chan int)
	w2 := make(chan int)
	go func(c chan<- int) {
		time.Sleep(time.Second)
		c <- 1
	}(w1)
	go func(c chan<- int) {
		time.Sleep(time.Second * 2)
		c <- 2
	}(w2)
	for i := 0; i < 2; i++ {
		select {
		case val := <-w1:
			fmt.Println("received: ", val)
		case val := <-w2:
			fmt.Println("received: ", val)
		}
	}
}

func timeout() {
	// 通过 select 和 channel 实现 timeout
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

func closingChannel() {
	/*
		生产者发送多个消息，发送完关闭通道，并等待消费端消费完成后退出
	*/
	w1 := make(chan int)
	done := make(chan bool)

	go func(w <-chan int, done chan<- bool) {
		// 消费者
		for {
			val, ok := <-w // 阻塞, ok 表示是否获取成功（close 情况下是 false)
			if ok {
				fmt.Println("got", val)
			} else {
				fmt.Println("close")
				done <- true
				break
			}
		}

		// 错误写法
		// for {
		// 	select {
		// 	case val := <-w:  // 不会阻塞，拿不到值就走default了
		// 		fmt.Println("got", val)
		// 	default:
		// 		fmt.Println("close")
		// 		done <- true
		// 		break
		// 	}
		// }

	}(w1, done)

	// 生产者
	for i := 0; i < 3; i++ {
		w1 <- i // 阻塞等待消费者接受
	}
	close(w1)
	<-done

	_, ok := <-w1
	fmt.Println("got more:", ok)
}

func TestChannel(t *testing.T) {
	done := make(chan bool, 1)
	go worker(done)
	go worker(done)
	go worker(done)

	finish(done)
	time.Sleep(time.Second)
	wait(done)

	selectChannels()

	timeout()

	closingChannel()

	// // 必须先有接受者，才能向 chan 发送数据，否则死锁
	w1 := make(chan int)
	go func(w chan int) {
		val := <-w1
		fmt.Println(val)
	}(w1)
	w1 <- 1 // 先声明接收者，再发送数据
	// w1 <- 2 // 没有接收者，再发数据会死锁

	// 通过指定缓冲区为2，避免阻塞
	w2 := make(chan int, 2)
	w2 <- 1
	w2 <- 2
	close(w2) // 先 close，避免 range 阻塞死锁

	for val := range w2 {
		fmt.Println("range got", val)
	}

}
