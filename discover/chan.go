package main

import "fmt"

type WatchResponse struct {
	// 响应数据结构
}

func main() {
	// 创建请求通道
	requestChan := make(chan chan WatchResponse, 1)

	// 启动工作goroutine
	go worker(requestChan)

	// 创建响应通道
	responseChan := make(chan WatchResponse)

	// 发送请求
	requestChan <- responseChan
	// 等待响应
	response := <-responseChan
	fmt.Println(response)

	requestChan <- responseChan
	response = <-responseChan
	fmt.Println(response)
}

func worker(requestChan chan chan WatchResponse) {
	for responseChan := range requestChan {
		// 处理请求...
		response := WatchResponse{ /*...*/ }

		// 通过提供的responseChan返回结果
		responseChan <- response
	}
}
