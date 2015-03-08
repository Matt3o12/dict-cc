package main

import "fmt"

const (
	issueTracker = "https://github.com/Matt3o12/dict-cc/issues"
)

func updateLanguages() {
	fmt.Println("This is your first usage.")
	fmt.Println("Available languages are being updated.")
	fmt.Println("This may take a while ...")
	fmt.Println()

	langs, err := GetLanguages()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("If you believe that is a bug, please open a ticket at %v.\n",
			issueTracker)
		return
	}

	fmt.Printf("Total langs: %v\n", len(langs))
}

func main() {
	updateLanguages()
}
