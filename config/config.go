package config

import (
	//"github.com/goburrow/serial"
	"github.com/tarm/serial"
	"time"
)

var Conf = &serial.Config{
	Name:  "COM6",
	Baud: 9600,
	Size: 8,
	StopBits: 1,
	Parity:   'E',
//Timeout時間決定port.read()的等待時間上限
	ReadTimeout: 10 * time.Second,
}