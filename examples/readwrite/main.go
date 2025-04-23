package main

import (
	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"os"
	"strconv"
)

type message struct {
	X int16 `struc:"int16,little"`
	Y int16 `struc:"int16,little"`
	Z int16 `struc:"int16,little"`

	CmdGliss         bool `struc:"bool"`         // Glissière : 0: retracter / 1: deployer
	CmdAimantInt     bool `struc:"bool"`         // Pince aimant interieur 0: détacher / 1: attacher
	CmdAimantExt     bool `struc:"bool"`         // Pince aimant exterieur
	CmdPompe         bool `struc:"bool"`         // Commande Pompe : 0: Off / 1: On
	CmdVanne         bool `struc:"bool"`         // Commande Electrovanne : 0: Off / 1: On
	CmdServoPlanche  bool `struc:"bool"`         // Lever les planches
	CmdServoBanniere bool `struc:"bool"`         // Lacher la banniere
	AscPlanche       int  `struc:"int16,little"` // Position de l'ascenceur des planches
	AscBoites        int  `struc:"int16,little"` // Position de l'ascenceur des boites
	Compteur         int  `struc:"uint8,little"` // Compteur
}

func main() {
	speed, err := strconv.Atoi(os.Args[2])
	if err != nil {
		logrus.Fatal(err)
	}
	mode := &serial.Mode{
		BaudRate: speed,
	}
	serialPort, err = serial.Open(os.Args[1], mode)
	if err != nil {
		logrus.Fatal(err)
	}
	defer serialPort.Close()

	if false {
		if err := runReader(serialPort); err != nil {
			logrus.Fatal(err)
		}
	} else {
		if err := runWriter(serialPort); err != nil {
			logrus.Fatal(err)
		}
	}
}
