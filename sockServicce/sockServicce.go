package sockServicce

import (
	"bytes"
	"encoding/binary"
	_ "encoding/json"
	"fmt"
	"net"
	"os"
	_ "syscall"
)

var unixCon *net.UnixConn

const unixSocketFile = "/tmp/testunix.sock"

func ConnectToUnixSocket() bool {
	addr, err := net.ResolveUnixAddr("unix", unixSocketFile)
	if err != nil {
		fmt.Println(err)
		return false
	}

	unixCon, err = net.DialUnix("unix", nil, addr)
	if err != nil {
		fmt.Println("net.DialUnix error: ", err.Error())
		unixCon = nil
		return false
	}
	return true
}

func CloseUnixSocket() {
	if unixCon != nil {
		unixCon.Close()
	}
}

func SendRequest(conn *net.UnixConn, data []byte) {
	buf := new(bytes.Buffer)
	msglen := uint32(len(data))

	binary.Write(buf, binary.BigEndian, &msglen)
	data = append(buf.Bytes(), data...)

	conn.Write(data)
}

//读结果
func ReadUnix(listener *net.UnixConn) {
	for {
		buf := make([]byte, 1400)
		size, remote, err := listener.ReadFromUnix(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("recv:", string(buf[:size]), " from ", remote.String())
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func RunUnix() {

	addr, err := net.ResolveUnixAddr("unix", unixSocketFile)
	checkError(err)
	unixCon, err = net.DialUnix("unix", nil, addr)
	//defer unixCon.Close()

	wsocketbuf := []byte{0xAA, 0xAA, 0x10, 0x01, 0xFA, 0x23, 0x9C, 0x00, 0x31, 0x01, 0x05, 0x01, 0x64, 0x00, 0x00, 0x01,
		0x26, 0xA7}
	//send to its subs
	go ReadUnix(unixCon)
	go SendRequest(unixCon, wsocketbuf)

	unixCon.Close()
}

//func onMessageReceived(conn *net.UnixConn) {
//	//for { // io Read will wait here, we don't need for loop to check
//	// Read information from response
//	data, err := parseResponse(conn)
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Printf("%v\tReceived from server: %s\n", time.Now(), string(data))
//	}
//
//	// Exit when receive data from server
//	exitSemaphore <- true
//	//}
//}
