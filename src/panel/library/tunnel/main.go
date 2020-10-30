package main

// 端口转发

import (
	"fmt"
	"log"
	"net"
)

func main() {
	//源端口，目的端口,转，隧道
	var fromport, toport int = 10086, 22
	fromaddr := fmt.Sprintf(":%d", fromport)
	toaddr := fmt.Sprintf("127.0.0.1:%d", toport)

	fromlistener, err := net.Listen("tcp", fromaddr)

	if err != nil {
		log.Fatal("Unable to listen on: %s, error: %s\n", fromaddr, err.Error())
	}
	defer fromlistener.Close()

	for {
		fromcon, err := fromlistener.Accept()
		if err != nil {
			fmt.Printf("Unable to accept a request, error: %s\n", err.Error())
		} else {
			fmt.Println("new connect:" + fromcon.RemoteAddr().String())
		}

		//这边最好也做个协程，防止阻塞
		toCon, err := net.Dial("tcp", toaddr)
		if err != nil {
			fmt.Printf("can not connect to %s\n", toaddr)
			continue
		}

		go handleConnection(fromcon, toCon)
		go handleConnection(toCon, fromcon)
	}
}

func handleConnection(r, w net.Conn) {
	defer r.Close()
	defer w.Close()

	var buffer = make([]byte, 100000)
	for {
		n, err := r.Read(buffer)
		if err != nil {
			break
		}

		n, err = w.Write(buffer[:n])
		if err != nil {
			break
		}
	}
}
