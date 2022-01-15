package main

import (
	"log"
	"strconv"

	"encoding/json"
	"fmt"
	"os"

	uexchange "github.com/Sagleft/uexchange-go"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	TelegramBotToken string
	Keyp             string
	Passp            string
}

func main() {

	//read congig
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(configuration.TelegramBotToken)
	fmt.Println(configuration.Keyp)
	fmt.Println(configuration.Passp)

	// bot-token

	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// ini channel
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updatesChann := bot.GetUpdatesChan(ucfg)

	//utp

	// create client
	client := uexchange.NewClient()

	// auth

	_, err = client.Auth(uexchange.Credentials{
		AccountPublicKey: configuration.Keyp,
		Password:         configuration.Passp,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Authorized on account %s", "CrpEX")

	// get balance

	//json goutpbot

	// update
	for {
		select {
		case update := <-updatesChann:
			// User bot
			UserName := update.Message.From.UserName

			// ID chat.

			ChatID := update.Message.Chat.ID

			// Text massage user
			Text := update.Message.Text

			log.Printf("[%s] %d %s", UserName, ChatID, Text)

			// reply
			reply := "Ok"
			// create massage
			msg := tgbotapi.NewMessage(ChatID, reply)
			// send
			bot.Send(msg)

			switch Text {

			case "/GetBalance":

				balanceData, err := client.GetBalance()

				if err != nil {
					log.Fatalln(err)
				}

				log.Println(balanceData)

				/*jsonBytes, err := json.Marshal(balanceData)

				if err != nil {
					log.Println(err)
					return
				}

				x1 := fmt.Sprintf("% s", jsonBytes)
				*/

				result := ""
				for _, data := range balanceData {
					result += data.Currency.Name + ": " + strconv.FormatFloat(data.Balance, 'f', 8, 32) + "\n"
				}

				log.Println("chek")
				log.Println(client.GetBalance())
				log.Println("chek")

				reply := result
				msg := tgbotapi.NewMessage(ChatID, "\n \n-- BALANCE \n\n"+reply+"\n \n-- BALANCE")

				bot.Send(msg)

			case "/How":

				reply := "Set up a bot in json and use it by adding new elements"
				msg := tgbotapi.NewMessage(ChatID, reply)

				bot.Send(msg)

			default:

				fmt.Println("commands")

				reply := "Commands:\n /GetBalance \n /How"
				msg := tgbotapi.NewMessage(ChatID, reply)

				bot.Send(msg)

			}

		}

	}
}
