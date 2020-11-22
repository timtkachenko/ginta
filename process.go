package main

import (
	"encoding/json"
	"fmt"
	"github.com/twinj/uuid"
	"sync"
	"time"
)

type Parsed struct {
	//	id     string
	//	data string
	Number int `json:"number"`
	//	count  int
}

var db = createDB()
var qStore queueSync

func worker(id int, jobs <-chan [][]byte) {
	for {
		chunk := <-jobs
		if len(chunk) > 0 {
			fmt.Println("worker", id, "started  job", chunk[0], "-", chunk[len(chunk)-1])
			receive(chunk)
		}
	}
}

func process() {
	jobs := make(chan [][]byte, 100)
	qStore = queueSync{
		mx:    &sync.RWMutex{},
		store: [][]byte{},
	}
	for w := 1; w <= 4; w++ {
		go worker(w, jobs)
	}
	for {
		<-time.After(5 * time.Millisecond)
		queueCopy := qStore.DrainAll()
		curQueueLength := len(queueCopy)
		if curQueueLength > 0 {
			chunkSize := curQueueLength / 2
			chunk := queueCopy[0:chunkSize]
			chunk2 := queueCopy[chunkSize:curQueueLength]
			jobs <- chunk
			jobs <- chunk2
		}
	}
}

func receive(chunk [][]byte) {
	bulkDoc := db.NewBulkDocument()
	for _, data := range chunk {
		var doc Parsed
		if err := json.Unmarshal(data, &doc); err != nil {
			fmt.Println("data: ", data)
			fmt.Println("Error: ", err)
		}
		//		time.Sleep(time.Second)
		err := bulkDoc.Save(doc, uuid.NewV4().String(), "")
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
	_, err := bulkDoc.Commit()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
