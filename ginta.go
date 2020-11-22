package main

import (
	"fmt"
	"github.com/libp2p/go-reuseport"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

const PORT = 22111
const maxQueueSize = 1000

type OutgoingMessage struct {
	recipient *net.UDPAddr
	data      []byte
}

var sCounter counterSync

func main() {
	log.Println("Starting...")

	maxListeners := runtime.NumCPU() / 2
	runtime.GOMAXPROCS(maxListeners)
	fmt.Printf("maxListeners: %d\n", maxListeners)

	go stats()
	go process()
	quit := make(chan os.Signal, 1)
	for i := 0; i < maxListeners; i++ {
		go beginListen(quit)
	}
	signal.Notify(quit, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-quit
	time.Sleep(time.Second * 2)
	log.Println("Shutdown ...")
}

func beginListen(c chan os.Signal) {
	addr := net.UDPAddr{
		Port: PORT,
		IP:   net.IP{0, 0, 0, 0},
	}

	connection, err := reuseport.ListenPacket("udp", addr.String())
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	outbox := make(chan OutgoingMessage, maxQueueSize)

	sendFromOutbox := func() {
		n, err := 0, error(nil)
		for msg := range outbox {
			n, err = connection.(*net.UDPConn).WriteToUDP(msg.data, msg.recipient)
			if err != nil {
				log.Println(err)
			}
			if n != len(msg.data) {
				log.Println("Tried to send", len(msg.data), "bytes but only sent ", n)
			}
		}
	}

	for i := 1; i <= 4; i++ {
		go sendFromOutbox()
	}

	listen(connection.(*net.UDPConn), c, outbox)

	close(outbox)
}

func listen(udpConn *net.UDPConn, c chan os.Signal, outbox chan OutgoingMessage) {
	for {
		buf := make([]byte, 8192)
		n, addr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		msg := buf[0:n]
		sCounter.Incr()
		qStore.Set(msg)
		outbox <- OutgoingMessage{
			recipient: addr,
			data:      msg,
		}
	}

}
func stats() {
	sCounter = counterSync{
		mx:    &sync.RWMutex{},
		store: 0,
	}
	var history int
	passed := 1
	for {
		select {
		case <-time.After(time.Second):
			fmt.Printf("processed: %d\n", sCounter.Get())
			passed++
			if passed > 5 && history == sCounter.Get() {
				sCounter.Reset()
				passed = 0
			} else if passed > 5 {
				history = sCounter.Get()
				passed = 0
			}
			break
		}

	}
}
