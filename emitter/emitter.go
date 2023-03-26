package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"gitlab.com/Skinass/hakaton-2023-1-1/common/storage"
)

type Emitter struct {
	Ports     []string
	Ratelimit int
	URL       string
}

type WorkerStruct struct {
	Item storage.Message
	Port string
}

func (emitter *Emitter) EmitterServer(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	wg := &sync.WaitGroup{}

	in := make(chan storage.Message)

	workerCh := make(chan WorkerStruct)
	for i := 0; i < len(emitter.Ports); i++ {
		go func(workerCh chan WorkerStruct) {
			wg.Add(1)
			defer wg.Done()
			for item := range workerCh {
				jsonItem, _ := json.Marshal(item.Item)
				bodyReader := bytes.NewReader(jsonItem)

				req, _ := http.NewRequest("POST", "http://"+emitter.URL+":"+item.Port, bodyReader)
				_, err := client.Do(req)
				if err != nil {
					fmt.Printf("Error sending to writer at port %v: %v", item.Port, err)
				}
			}
		}(workerCh)
	}

	reader := NewDataReader("amazon-meta.txt")
	go func() {
		message, ok := reader.NextMessage()
		for ok {
			func() {
				in <- message
				fmt.Println(message)
				message, ok = reader.NextMessage()
			}()
		}
	}()

	ports := []string{}
	ports = append(ports, emitter.Ports...)

	counter := 0
	var tmpWorker WorkerStruct
	for inData := range in {
		time.Sleep(time.Duration(emitter.Ratelimit) * time.Second)
		tmpWorker = WorkerStruct{Item: inData, Port: ports[counter%len(emitter.Ports)]} // fmt.Sprintf("808%d", counter%len(emitter.Ports))}
		workerCh <- tmpWorker
		counter++
	}
}

func main() {
	var writersPorts string
	var num int

	flag.StringVar(&writersPorts, "writers_ports", "", "comma-separated list of writer ports")
	flag.IntVar(&num, "num", 0, "ratelimit of items sending")

	flag.Parse()

	if writersPorts == "" {
		fmt.Println("Error: writers_ports flag is required")
		return
	}

	ports := strings.Split(writersPorts, ",")
	fmt.Printf("Ports: %v\n", ports)
	fmt.Printf("Number: %d\n", num)

	emitter := Emitter{Ports: ports, Ratelimit: num, URL: "localhost"}
	port := ":9493"

	mux := http.NewServeMux()
	mux.HandleFunc("/",
		emitter.EmitterServer)

	server := http.Server{
		Addr:    port,
		Handler: mux,
	}

	fmt.Println("Starting emitter at:", port)
	server.ListenAndServe()
}
