package main

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"time"
)

var (
	BEANS_TUBE = "metrics"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	conn, err := beanstalk.Dial("tcp", "127.0.0.1:11300")

	check(err)

	defer conn.Close()

	remoteSet := &beanstalk.Tube{conn, BEANS_TUBE}

	body := []byte(`{"uri":"action=log&prod=example.com&guid=C4BAE81F3A308FC5FBC690D616B496A4&mid=","ip":"109.48.178.88","ts":1529270575}`)

	for i := 1; i < 499; i++ {
		_, err = remoteSet.Put(body, 1, 0, 7*24*time.Hour)

		check(err)
	}

	stats, err := remoteSet.Stats()

	check(err)

	fmt.Printf("current-jobs-ready: %s\n\n", stats["current-jobs-ready"])
}
