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
var db = createDB()
var res chan int = make(chan int)

func worker(id int, jobs <-chan []string) {
	for {
		chunk := <-jobs
		if len(chunk) > 0 {
			fmt.Println("worker", id, "started  job", chunk[0], "-", chunk[len(chunk)-1])
			//			res <- len(chunk)
			receive(chunk)
			fmt.Println("worker", id, "finished job", chunk[0], "-", chunk[len(chunk)-1])
		}
	}
}

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
	jobs := make(chan []string, 100)

	go accum(in)
	for w := 1; w <= 4; w++ {
		go worker(w, jobs)
	}
	go consume(jobs)

	go func() {
		ticker := time.NewTicker(time.Millisecond * 3000)
		var i int = 0
		for {
			select {
			case <-ticker.C:
				fmt.Println("res: ", i)
				fmt.Println("cc: ", cc)
				cc = 0
				i = 0
				break
			case c := <-res:
				i += c
				break
			}

		}
	}()
	for {
		n, err := sConn.Read(buf)
		in <- string(buf[0:n])

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

}

var cc int

func accum(in chan string) {
	for {
		queue = append(queue, <-in)
		cc++
	}
}

func consume(jobs chan []string) {
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		//		select {
		//		case <-ticker.C:
		<-ticker.C

		//		fmt.Println("Tick at", t)
		queueCopy := queue[:len(queue)]
		curQueueLength := len(queueCopy)
		res <- len(queue)
		if curQueueLength > 0 {
			//			fmt.Printf("actual queue %v\n", len(queue))
			queue = queue[curQueueLength:]

			chunkSize := curQueueLength / 2
			chunk := queueCopy[0:chunkSize]
			chunk2 := queueCopy[chunkSize:curQueueLength]
			jobs <- chunk
			jobs <- chunk2
			//			fmt.Printf("curQueueLength %v\n", curQueueLength)
			//			fmt.Printf("chunk %v\n", len(chunk))
			//			fmt.Printf("chunk2 %v\n", len(chunk2))
		}

		//		fmt.Println("tick\n")
		//		case res := <-results:
		//			fmt.Println("result: ", res)
		//		default:
		//			//			fmt.Println("no message received")

		//			if count != 0 && len(queue) == 0 && !reported {
		//				reported = true
		//				count = 0
		//				fmt.Printf("count == %v\n", count)
		//			}
		//		}
	}
}

func receive(chunk []string) {
	bulkDoc := db.NewBulkDocument()

	for _, data := range chunk {

		var doc Parsed

		if err := json.Unmarshal([]byte(data), &doc); err != nil {
			fmt.Println("data: ", data)
			panic(err)
		}
		//		time.Sleep(time.Second)
		err := bulkDoc.Save(doc, uuid.NewV4().String(), "")

		//		results <- rev + " - " + chunk[0]
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
	//	if len(chunk) > 0 {
	//		fmt.Printf("%v - %v\n", chunk[0], chunk[len(chunk)-1])
	//	}
	_, err := bulkDoc.Commit()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	//	fmt.Printf("bulk res: %v\n", res)
}
