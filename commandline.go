package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"
)

const (
	issueTracker = "https://github.com/Matt3o12/dict-cc/issues"

	// DictBaseURL The base url of dict.cc
	DictBaseURL = "http://dict.cc/"

	// AllLangaugesGet URL where all available langauge pairs can be found.
	AllLangaugesGet = "http://browse.dict.cc/"

	allAvaiableLangsCSSPath = "#maincontent form[name='langbarchooser'] " +
		"table td a"

	// Version current version
	Version = "0.0.0"

	// Author of the app
	Author = "Matteo Kloiber"

	// Email my email
	Email = "info@matt3o12.de"
)

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Printf("If you believe that is a bug, please open a ticket at %v.\n",
			issueTracker)
		os.Exit(1)

	}
}

func getLangSaveFile() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", errors.New("Could not find home direcotry.")
	}

	langsFile := path.Join(home, ".dict_cc", "languages.json")
	err = os.Mkdir(path.Dir(langsFile), 0755)
	if err != nil && !os.IsExist(err) {
		return "", err
	}

	return langsFile, nil
}

func updateLanguages() {
	fmt.Println("Available languages are being updated.")
	fmt.Println("This may take a while ...")
	fmt.Println()

	response, err := http.Get(AllLangaugesGet)
	handleErr(err)

	langs, err := GetLanguagesFromRemote(response)
	handleErr(err)

	fmt.Printf("Total langs: %v\n", len(langs))

	fileName, err := getLangSaveFile()
	handleErr(err)

	file, err := os.Create(fileName)
	handleErr(err)
	defer file.Close()

	SaveLanguagesToDisk(langs, file)
}

func updateLanguagesCommand(c *cli.Context) {
	if len(c.Args()) > 0 {
		fmt.Println("You may not set additional commands.")
		fmt.Println()

		cli.ShowSubcommandHelp(c)
		os.Exit(1)
	}

	updateLanguages()
}

func lookupCommand(c *cli.Context) {
	langFile, err := getLangSaveFile()
	handleErr(err)

	if stats, err := os.Stat(langFile); err != nil || !stats.Mode().IsRegular() {
		fmt.Println("This is your first usage.")
		updateLanguages()
	}

	// Code for looking up words goes here.
}

func main() {
	app := cli.NewApp()
	app.Name = "Dict.cc client"

	// Usage probably means description
	app.Usage = "Look up any word in many languages!"
	app.Author = Author
	app.Version = Version
	app.Email = Email
	app.Action = lookupCommand
	app.Commands = []cli.Command{
		{
			Name:   "update-langs",
			Usage:  "Updates all languages",
			Action: updateLanguagesCommand,
			Flags:  make([]cli.Flag, 0),
		},
	}

	app.Run(os.Args)
	os.Exit(0)
}
