package main

import "os"

func osExit() {
	os.Exit(0) //want "has os.Exit call in main package"
}
