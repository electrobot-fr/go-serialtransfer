package main

import (
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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
	if err := runCarte(); err != nil {
		log.Fatal(err)
	}
}

func runTelecommande() error {
	speed, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return err
	}
	mode := &serial.Mode{
		BaudRate: speed,
	}
	port, err = serial.Open(os.Args[1], mode)
	if err != nil {
		return err
	}

	for {
		buff := make([]byte, 100)
		n, err := port.Read(buff)
		if err != nil {
			return err
		}
		if n == 0 {
			return fmt.Errorf("EOF")
		}

		//fmt.Printf("Received: %s", hex.EncodeToString(buff[:n]))

		var foo message
		if err := DecodePacket(buff[:n], &foo); err != nil {
			logrus.Error(fmt.Errorf("decoding error: %v", err))
			continue
		}

		fmt.Printf("%+v\n", foo)
	}

}

var input = message{}
var port serial.Port

func runCarte() error {
	speed, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return err
	}
	mode := &serial.Mode{
		BaudRate: speed,
	}
	port, err = serial.Open(os.Args[1], mode)
	if err != nil {
		return err
	}

	p := prompt.New(
		executor,
		completer,
	)
	p.Run()

	return nil
}

func executor(in string) {
	in = strings.TrimSpace(in)
	if in == "" {
		return
	}

	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "quit", "exit":
		fmt.Println("Au revoir!")
		os.Exit(0)
	case "pompe":
		input.CmdPompe = !input.CmdPompe
	case "vanne":
		input.CmdVanne = !input.CmdVanne
	case "aimants-interieurs":
		input.CmdAimantInt = !input.CmdAimantInt
	case "aimants-exterieurs":
		input.CmdAimantExt = !input.CmdAimantExt
	case "glissieres":
		input.CmdGliss = !input.CmdGliss
	case "servo-planche":
		input.CmdServoPlanche = !input.CmdServoPlanche
	case "servo-banniere":
		input.CmdServoBanniere = !input.CmdServoBanniere
	case "compteur":
		if len(blocks) < 2 {
			fmt.Println("Valeur manquante")
			return
		}
		v, err := strconv.Atoi(blocks[1])
		if err != nil {
			fmt.Println("Valeur invalide:", err)
			return
		}
		input.Compteur = v
	case "asc-planche":
		if len(blocks) < 2 {
			fmt.Println("Valeur manquante")
			return
		}
		v, err := strconv.Atoi(blocks[1])
		if err != nil {
			fmt.Println("Valeur invalide:", err)
			return
		}
		input.AscPlanche = v
	case "asc-boites":
		if len(blocks) < 2 {
			fmt.Println("Valeur manquante")
			return
		}
		v, err := strconv.Atoi(blocks[1])
		if err != nil {
			fmt.Println("Valeur invalide:", err)
			return
		}
		input.AscBoites = v
	case "print":
		bin, err := json.MarshalIndent(input, "", "  ")
		if err != nil {
			return
		}
		fmt.Println(string(bin))
		return
	case "help":
		for _, suggest := range s {
			fmt.Printf("%s:\t%s\n", suggest.Text, suggest.Description)
		}
		return
	case "reset":
		input = message{}
	case "sequence":
		input = message{
			AscBoites: 6600,
		}
		sendMessage()
		time.Sleep(5 * time.Second)

		input.AscBoites = 6600
		input.CmdAimantExt = true
		input.CmdAimantInt = true
		sendMessage()
		time.Sleep(5 * time.Second)

		input.AscPlanche = 5000
		input.CmdPompe = true
		sendMessage()
		time.Sleep(5 * time.Second)

		input.AscPlanche = 0
		input.CmdPompe = true
		sendMessage()
		time.Sleep(5 * time.Second)

		input.CmdServoPlanche = true
		sendMessage()
		time.Sleep(3 * time.Second)

		input.CmdGliss = true
		sendMessage()
		time.Sleep(3 * time.Second)

		input.CmdAimantExt = false
		input.CmdServoPlanche = false
		sendMessage()
		time.Sleep(3 * time.Second)

		input.CmdGliss = false
		sendMessage()
		time.Sleep(3 * time.Second)
	case "sequence2":
		input.AscBoites = 1000
		sendMessage()
		time.Sleep(3 * time.Second)

	case "sequence3":
		input.AscBoites = 1500
		sendMessage()
		time.Sleep(3 * time.Second)

		input.CmdAimantInt = false
		sendMessage()
		time.Sleep(3 * time.Second)

		input.AscPlanche = 500
		sendMessage()
		time.Sleep(3 * time.Second)

		input.CmdPompe = false
		input.CmdVanne = true
		sendMessage()
		time.Sleep(1 * time.Second)

		input.CmdVanne = false
		sendMessage()
		time.Sleep(1 * time.Second)
	}

	if err := sendMessage(); err != nil {
		logrus.Error(err)
	}
}

func sendMessage() error {
	bin, err := EncodePacket(&input)
	if err != nil {
		return err
	}
	if _, err := port.Write(bin); err != nil {
		return err
	}

	//buff := make([]byte, 100)
	//n, err := port.Read(buff)
	//if err != nil {
	//	return err
	//}
	//if n == 0 {
	//	return fmt.Errorf("EOF")
	//}
	//
	//fmt.Printf("Received: %s", hex.EncodeToString(buff[:n]))
	//
	//var foo message2
	//if err := DecodePacket(buff[:n], &foo); err != nil {
	//	return fmt.Errorf("decoding error: %v", err)
	//}
	//
	//fmt.Printf(", decoded: %+v\n", foo)
	return nil
}

var s = []prompt.Suggest{
	{Text: "exit", Description: "Quitter"},
	{Text: "pompe", Description: "Activer/désactiver la pompe"},
	{Text: "vanne", Description: "Activer/désactiver la vanne"},
	{Text: "aimants-interieurs", Description: "Activer/désactiver les aimants intérieurs"},
	{Text: "aimants-exterieurs", Description: "Activer/désactiver les aimants extérieurs"},
	{Text: "glissieres", Description: "Rentrer/sortir les glissières"},
	{Text: "servo-planche", Description: "Monter/descendre les planches"},
	{Text: "servo-banniere", Description: "Ouvrir/fermer le servo de la banniere"},
	{Text: "compteur", Description: "Changer la valeur du compteur"},
	{Text: "asc-planche", Description: "Modifier la position de l'ascenseur des planches"},
	{Text: "asc-boites", Description: "Modifier la position de l'ascenseur des bannières"},
	{Text: "print", Description: "Afficher les valeurs courantes"},
	{Text: "help", Description: "Afficher l'aide"},
	{Text: "reset", Description: "Remettre à zéro les valeurs"},
	{Text: "sequence", Description: "Démo de séquence"},
}

func completer(d prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
