package main

import (
	"bufio"
	"os"

	"github.com/dxbednarczyk/checkers"
)

type Ripper int

var (
	XLD Ripper = 1
	EAC Ripper = 2
)

func main() {
	lf, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		panic(err)
	}

	defer lf.Close()

	log := bufio.NewScanner(lf)
	if !log.Scan() {
		panic("Empty log file")
	}

	var ripper Ripper

	switch {
	case log.Text()[0] == 'X':
		ripper = XLD
	default:
		panic("Unknown log format")
	}

	switch ripper {
	case XLD:
		err = checkers.XLD(log)
	}

	if err != nil {
		panic(err)
	}
}
