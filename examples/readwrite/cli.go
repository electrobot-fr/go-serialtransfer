package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/electrobot-fr/go-serialtransfer"
	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

var input = message{}
var serialPort serial.Port

func runWriter(port serial.Port) error {
	serialPort = port

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
	}

	if err := sendMessage(); err != nil {
		logrus.Error(err)
	}
}

func sendMessage() error {
	bin, err := serialtransfer.EncodePacket(&input)
	if err != nil {
		return err
	}
	fmt.Println(hex.EncodeToString(bin))
	if _, err := serialPort.Write(bin); err != nil {
		return err
	}
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
