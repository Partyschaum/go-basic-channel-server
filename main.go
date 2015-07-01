package main

import (
	"fmt"
	"time"
)

type Request struct {
	question        string
	responseChannel chan *Response
}

type Response struct {
	responseCode int
	answer       string
}

type HandlerFunc func(*Request)

func startServer(handler HandlerFunc, requestChannel <-chan *Request) {
	for {
		go handler(<-requestChannel)
	}
}

func main() {
	handler := func(request *Request) {
		request.responseChannel <- &Response{200, "The answer to '" + request.question + "' is 42!"}
	}

	requestChannel := make(chan *Request)
	responseChannel := make(chan *Response)

	go startServer(handler, requestChannel)

	go func() {
		time.Sleep(1 * time.Second)
		requestChannel <- &Request{"What is the meaning of life?", responseChannel}
		time.Sleep(1 * time.Second)
		requestChannel <- &Request{"Why is Go so fun?", responseChannel}
	}()

	for {
		select {
		case response := <-responseChannel:
			fmt.Println(response)
		case <-time.After(3 * time.Second):
			fmt.Println("timed out")
			return
		}
	}
}
