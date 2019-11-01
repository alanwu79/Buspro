package transfer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	//"github.com/goburrow/serial"
	"log"
)

type CommandFormat struct {
	sourceId []byte `json:"source_address"`
	targetId string `json:target_address`
	Opcode   int    `json:"operate_code"`
	payload  string `json:payload`
}

func ByteToHex(data []byte) string {
	buffer := new(bytes.Buffer)
	for _, b := range data {

		s := strconv.FormatInt(int64(b&0xff), 16)
		if len(s) == 1 {
			buffer.WriteString("0")
		}
		buffer.WriteString(s)
	}

	return buffer.String()
}

func bytesToIntU(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp uint8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp uint16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp uint32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}

func ReadPipeTransJSON(readpipe <-chan []byte) {
	log.Println("enter ReadPipeTransJSON")

	var m CommandFormat
	command := make([]byte, 256)

	command = <-readpipe
	m.sourceId = command[3:4]
	fmt.Println("sourceId:", m.sourceId)
}
