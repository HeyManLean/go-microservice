package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type Job struct {
	Id     int
	Result int
}

func jobWorker(jobs <-chan Job, result chan<- Job) {
	for job := range jobs {
		fmt.Println("got job: ", job.Id)
		time.Sleep(time.Microsecond * 100)
		job.Result = job.Id + 1
		result <- job
	}
}

func TestConcurrency(t *testing.T) {
	/*
		并发消费，并支持获取结果
	*/
	jobs := make(chan Job)
	result := make(chan Job)
	// 启动3个worker并发处理消息
	for i := 0; i < 3; i++ {
		go jobWorker(jobs, result)
	}

	// 发送消息
	for i := 0; i < 3; i++ {
		jobs <- Job{Id: i + 1}
	}
	close(jobs)

	// 接受结果, for range 会阻塞，不会退出
	// for job := range result {
	// 	fmt.Println("result: ", job.Id)
	// }
	// 硬编码 3 条消息，可以考虑使用 wait group 避免该问题
	for i := 0; i < 3; i++ {
		job := <-result
		fmt.Println("result: ", job.Id)
	}
}

func TestWaitGroup(t *testing.T) {
	/*
		并发消费消息，并等待全部消息处理完成
	*/
	wg := &sync.WaitGroup{}

	worker := func(id int) {
		fmt.Println("worker start: ", id)
		time.Sleep(time.Millisecond * 100 * time.Duration(id))
		fmt.Println("worker done: ", id)
	}

	for i := 0; i < 3; i++ {
		// Add 添加协程
		wg.Add(1)
		go func() {
			worker(i + 1)
			// Done 结束协程
			wg.Done()
		}()
	}
	// Wait 直到所有协程结束
	wg.Wait()
}
