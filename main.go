package main

import (
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	app := cliApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
