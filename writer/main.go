package main

import (
	"flag"
	"fmt"
	"gitlab.com/Skinass/hakaton-2023-1-1/writer/networking"
	"net/http"
	"os"
)

var fileName = "storage/worker%d.txt"

func main() {
	var port = flag.Int("port", -1, "port flag")
	var num = flag.Int("num", -1, "writer id flag")

	flag.Parse()

	if *port < 0 {
		panic(fmt.Sprintf("wrong port: %d", *port))
	}

	if *num < 0 {
		panic(fmt.Sprintf("wrong num: %d", *num))
	}

	file, err := os.Create(fmt.Sprintf(fileName, *num))

	if err != nil {
		panic(fmt.Sprintf("err open file: %s", err))
	}

	defer file.Close()

	wHandler := &networking.Handler{File: file}
	http.Handle("/", wHandler)

	fmt.Printf("starting %d server at :%d\n", *num, *port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)

	if err != nil {
		panic(err)
	}
}

//curl -d '{"discounted":false,"id":5,"ASIN":"fsdjkfs","title":"dsjfsd","group":"fsdfsd","salesrank":0,"similarsnum":0,"similars":null,"time":"0001-01-01T00:00:00Z"}' -X POST http://localhost:2222/
//go run ./writer/main.go -port 8081 -num 1
