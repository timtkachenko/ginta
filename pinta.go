package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/satori/go.uuid"
)

type Parsed struct {
	//	id     string
	//	data string
	Number int
	//	count  int
}

var queue []string
var in chan string = make(chan string)

//var out chan string = make(chan string, 2)
var db = createDB()

func main() {
	sAddr, err := net.ResolveUDPAddr("udp", ":22111")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	sConn, err := net.ListenUDP("udp", sAddr)
	sConn.SetReadBuffer(212992)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer sConn.Close()

	buf := make([]byte, 8192)

	go accum(in)
	//	go accum(in)
	go consume()
	i := 0
	for {
		n, err := sConn.Read(buf)

		in <- string(buf[0:n])
		if err != nil {
			fmt.Println("Error: ", err)
		}
		i++
		fmt.Printf("m %+v\n", i)
	}

}
func accum(in chan string) {
	for {
		queue = append(queue, <-in)
	}
}

func consume() {
	ticker := time.NewTicker(time.Millisecond * 10)
	i := 0
	for {
		<-ticker.C
		//		fmt.Println("Tick at", t)
		if len(queue) > 0 {
			i++
			//			fmt.Printf("%+v\n", queue)
			var data string
			data, queue = queue[len(queue)-1], queue[:len(queue)-1]
			//			fmt.Printf("%+v\n", string(data))
			var doc Parsed
			if err := json.Unmarshal([]byte(data), &doc); err != nil {
				panic(err)
			}

			_, err := db.Save(doc, uuid.NewV4().String(), "")
			if err != nil {
				fmt.Println("Error: ", err)
			}
			//			fmt.Printf("%s\n", rev)

		} else {
			if i != 0 {
				fmt.Printf("c %v\n", i)
				i = 0
			}
		}
	}
}

func receive(chunk) {
	for data := range chunk {

		var doc Parsed
		if err := json.Unmarshal([]byte(data), &doc); err != nil {
			panic(err)
		}

		_, err := db.Save(doc, uuid.NewV4().String(), "")
		if err != nil {
			fmt.Println("Error: ", err)
		}
		//			fmt.Printf("%s\n", rev)
	}
}
