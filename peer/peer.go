package main

import (
	"fmt"
	"log"
	"moqikaka.com/goutil/mathUtil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)
var tag string
const HAND_SHAKE_MSG = "我是打洞消息"

func main() {
	// 当前进程标记字符串,便于显示
	if len(os.Args)>1{
		tag = os.Args[1]
	}else {
		tag= strconv.Itoa(mathUtil.GetRand().Intn(10000))
	}

	udp()
}

func udp()  {
	defSrcPort := 9982
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: defSrcPort} // 注意端口必须固定
	//dstAddr := &net.UDPAddr{IP: net.ParseIP("207.148.70.129"), Port: 9981}
	dstAddr := &net.UDPAddr{IP: net.ParseIP("106.52.100.147"), Port: 18506}
	//dstAddr := &net.UDPAddr{IP: net.ParseIP("10.255.0.16"), Port: 18503}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)

		for  {
			defSrcPort++
			srcAddr = &net.UDPAddr{IP: net.IPv4zero, Port: defSrcPort}
			conn, err = net.DialUDP("udp", srcAddr, dstAddr)
			if err==nil{
				break
			}

			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(500*time.Millisecond)
		}
	}

	if _, err = conn.Write([]byte("hello, I'm new peer:" + tag)); err != nil {
		log.Panic(err)
	}
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
	}
	conn.Close()
	anotherPeer := parseAddr(string(data[:n]))
	fmt.Printf("local:%s server:%s another:%s\n", srcAddr, remoteAddr, anotherPeer.String())
	// 开始打洞
	bidirectionHole(srcAddr, &anotherPeer)
}

func tcp()  {
	defSrcPort := 9982
	srcAddr := &net.TCPAddr{IP: net.IPv4zero, Port: defSrcPort} // 注意端口必须固定
	//dstAddr := &net.UDPAddr{IP: net.ParseIP("207.148.70.129"), Port: 9981}
	dstAddr := &net.TCPAddr{IP: net.ParseIP("106.52.100.147"), Port: 18503}
	//dstAddr := &net.TCPAddr{IP: net.ParseIP("10.255.0.16"), Port: 18503}
	conn, err := net.DialTCP("tcp", srcAddr, dstAddr)
	if err != nil {
		log.Println(err)

		for  {
			defSrcPort++
			srcAddr = &net.TCPAddr{IP: net.IPv4zero, Port: defSrcPort}
			conn, err = net.DialTCP("tcp", srcAddr, dstAddr)
			if err==nil{
				break
			}

			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(500*time.Millisecond)
		}
	}

	if _, err = conn.Write([]byte("hello, I'm new peer:" + tag)); err != nil {
		log.Panic(err)
	}
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
	}
	conn.Close()
	anotherPeer := parseAddr(string(data[:n]))
	fmt.Printf("local:%s server:%s another:%s\n", srcAddr, conn.RemoteAddr().String(), anotherPeer.String())
	// 开始打洞
	udpSrcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: defSrcPort}
	bidirectionHole(udpSrcAddr, &anotherPeer)
}

func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

func bidirectionHole(srcAddr *net.UDPAddr, anotherAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	// 向另一个peer发送一条udp消息(对方peer的nat设备会丢弃该消息,非法来源),用意是在自身的nat设备打开一条可进入的通道,这样对方peer就可以发过来udp消息
	if _, err = conn.Write([]byte(HAND_SHAKE_MSG)); err != nil {
		log.Println("send handshake:", err)
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from [" + tag + "]")); err != nil {
				log.Println("send msg fail", err)
			}
		}
	}()
	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		} else {
			log.Printf("收到数据:%s\n", data[:n])
		}
	}
}