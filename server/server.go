package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 18506})
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("本地地址: <%s> \n", listener.LocalAddr().String())
	peers := make([]net.UDPAddr, 0, 2)
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		log.Printf("<%s> %s\n", remoteAddr.String(), data[:n])
		peers = append(peers, *remoteAddr)
		if len(peers) == 2 {
			log.Printf("进行UDP打洞,建立 %s <--> %s 的连接\n", peers[0].String(), peers[1].String())
			listener.WriteToUDP([]byte(peers[1].String()), &peers[0])
			listener.WriteToUDP([]byte(peers[0].String()), &peers[1])
			time.Sleep(time.Second * 8)
			log.Println("中转服务器退出,仍不影响peers间通信")
			return
		}
	}
}

func main2() {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4zero, Port: 18503})
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("本地地址: <%s> \n", listener.Addr().String())
	peers := make([]net.Conn, 0, 2)
	data := make([]byte, 1024)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accept: %s", err)
		}

		//n, remoteAddr, err := listener.ReadFromUDP(data)
		//if err != nil {
		//	fmt.Printf("error during read: %s", err)
		//}

		n,err:=	conn.Read(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}

		conn.RemoteAddr().String()
		log.Printf("<%s> %s\n", conn.RemoteAddr().String(), data[:n])
		peers = append(peers, conn)
		if len(peers) == 2 {
			log.Printf("进行UDP打洞,建立 %s <--> %s 的连接\n", peers[0].RemoteAddr().String(), peers[1].RemoteAddr().String())
			peers[1].Write([]byte(peers[0].RemoteAddr().String()))
			peers[0].Write([]byte(peers[1].RemoteAddr().String()))
			time.Sleep(time.Second * 8)
			log.Println("中转服务器退出,仍不影响peers间通信")
			return
		}
	}
}