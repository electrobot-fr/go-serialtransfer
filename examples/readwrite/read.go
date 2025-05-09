package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"

	"github.com/electrobot-fr/go-serialtransfer"
	"go.bug.st/serial"
)

func runReader(serialPort serial.Port) error {
	decoder := serialtransfer.NewDecoder(serialPort)
	for {
		var msg message
		err := decoder.Decode(&msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Error(err)
			continue
		}
		fmt.Printf("%+v\n", msg)
	}
	return nil
}
