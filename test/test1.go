package main

import "time"

func main() {
	go test1()
	for i := 0; i < 100; i++ {
		println("hello")
	}
	time.Sleep(10)
}

func test1() {
	for i := 0; i < 10; i++ {
		println("word")
	}
}

//export https_proxy=http://127.0.0.1:33210 http_proxy=http://127.0.0.1:33210 all_proxy=socks5://127.0.0.1:33211
