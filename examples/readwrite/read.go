package main

import (
	"fmt"
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
			return err
		}
		fmt.Printf("%+v\n", msg)
	}
	return nil
}
