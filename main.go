package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Types for json files

// Basic element of copypasta
type Copypasta struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// Array of copypastas
type Copypastas struct {
	Copypastas []Copypasta `json:"copypastas"`
}

// Basic element of config
type Config struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Array of config
type Configs struct {
	Configs []Config `json:"configs"`
}

// end of types

// returns an array of Copypasta structs
func getCopyPastaJSON(filename string) []Copypasta {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	var copypasta Copypastas
	json.NewDecoder(jsonFile).Decode(&copypasta)
	if copypasta.Copypastas == nil {
		panic("No copypastas found")
	}
	return copypasta.Copypastas
}

// returns an array of Config structs
func getSettings(filename string) []Config {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	var configs Configs
	json.NewDecoder(jsonFile).Decode(&configs)
	return configs.Configs
}

// main function
func main() {
	var basePath = "/root/scripts/copypastabot"
	// get basePath from path of executable
	// get bot token from json file
	configs := getSettings(basePath+"/config.json")
	bot, err := tgbotapi.NewBotAPI(configs[0].Value)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	// get copypasta from json file
	copypastas := getCopyPastaJSON(basePath+"/data.json")
	var result []string
	// get all titles from copypasta and add them to a string array
	for i := 0; i < len(copypastas); i++ {
		result = append(result, copypastas[i].Title)
	}
	// join all titles with a new line
	cmdlist := strings.Join(result, "\n")

	log.Printf("Authorized on bot \"%s\"", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates { // polling cycle
		if update.Message == nil {
			continue
		}

		for i := 0; i < len(copypastas); i++ {
			if strings.ToLower(update.Message.Text) == copypastas[i].Title || strings.ToLower(update.Message.Text) == "/"+copypastas[i].Title {
				// if text is too long, send it in multiple messages
				if len(copypastas[i].Text) > 4096 {
					for j := 0; j < len(copypastas[i].Text); j += 4096 {
						end := j + 4096
						if end > len(copypastas[i].Text) {
							end = len(copypastas[i].Text)
						}
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, copypastas[i].Text[j:end])
						bot.Send(msg)
					}
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, copypastas[i].Text)
					bot.Send(msg)
				}				
			}
		}

		if strings.ToLower(update.Message.Text) == "help" || strings.ToLower(update.Message.Text) == "/help" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "*Copypasta disponibili:*\n"+cmdlist)
			msg.ParseMode = "MarkdownV2"
			bot.Send(msg)
		}

		if strings.ToLower(update.Message.Text) == "start" || strings.ToLower(update.Message.Text) == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Benvenuto su @"+bot.Self.UserName+"!\nDigita '/help' o 'help' per vedere la lista dei copypasta disponibili.")
			bot.Send(msg)
		}
	}
}
