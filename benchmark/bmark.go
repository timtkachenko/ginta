package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const payload_sz = 32

func main() {
	addr, err := net.ResolveUDPAddr("", "0.0.0.0:22111")
	chk(err)
	iter := 1000 * 1000
	writeToUDP(addr, iter)
	time.Sleep(time.Second*2)
}
func listen(sConn *net.UDPConn) {
	buf := make([]byte, 8192)
	for {
		n, err := sConn.Read(buf)
		//in <- string(buf[0:n])
		_ = buf[0:n]
		//fmt.Printf("response %v", )
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
func writeToUDP(addr *net.UDPAddr, i int) {
	log.Printf("Start `writeToUDP` test with %d iteration\n", i)

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	chk(err)
	go listen(conn)

	for ; i > 0; i-- {
		payload := []byte(fmt.Sprintf(`{"number": %d}`, i))
		n, err := conn.WriteToUDP(payload, addr)
		chk(err)
		time.Sleep(time.Microsecond*3)
		if i%100000 == 0 {
			fmt.Printf("left %d\n", i)
		}
		if n != len(payload) {
			log.Println("Tried to send", len(payload), "bytes but only sent ", n)
		}
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
