package config

import (
	"github.com/goburrow/serial"
	"time"
)

var Conf = &serial.Config{
Address:  "COM6",
BaudRate: 9600,
DataBits: 8,
StopBits: 1,
Parity:   "E",
//Timeout時間決定port.read()的等待時間上限
Timeout: 10 * time.Second,
}