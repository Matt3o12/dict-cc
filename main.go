package main

import (
	"fmt"
	"net/http"
)

const (
	issueTracker = "https://github.com/Matt3o12/dict-cc/issues"
)

func handleErr(err error) {
	fmt.Println(err)
	fmt.Printf("If you believe that is a bug, please open a ticket at %v.\n",
		issueTracker)
}

func updateLanguages() {
	fmt.Println("This is your first usage.")
	fmt.Println("Available languages are being updated.")
	fmt.Println("This may take a while ...")
	fmt.Println()

	response, err := http.Get(AllLangaugesGet)
	if err != nil {
		handleErr(err)
		return
	}

	langs, err := GetLanguages(response)
	if err != nil {
		handleErr(err)
		return
	}

	fmt.Printf("Total langs: %v\n", len(langs))
}

func main() {
	updateLanguages()
}
