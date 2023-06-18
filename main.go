package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	//"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	
	),
)

func main()  {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	  }
	var token string = os.Getenv("TOKEN")

	log.Printf("token %s",token)
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil{
		log.Printf("error %s",err)
	}
	
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates{
		// if update.Message != nil {
		// 	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// 	msg.ReplyToMessageID = update.Message.MessageID

		// 	bot.Send(msg)
		// }
		if update.Message != nil {
			// Construct a new message from the given chat ID and containing
			// the text that we received.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			// If the message was open, add a copy of our numeric keyboard.
			//text := strings.TrimSpace(update.Message.Text)
			// if len(text) > 0 {
				
				switch update.Message.Text {
				case "open":
					msg.ReplyMarkup = numericKeyboard
	
				}
				// if _, err = bot.Send(msg); err != nil {
				// 	panic(err)
				// }
				if update.Message.Photo != nil {
					// Retrieve the photo information
					photos := update.Message.Photo
					photo := photos[len(photos)-1] // Use the last photo (highest resolution)
		
					// Get the file path
					fileConfig := tgbotapi.FileConfig{
						FileID: photo.FileID,
					}
					file, err := bot.GetFile(fileConfig)
					if err != nil {
						log.Println("Failed to get file:", err)
						continue
					}
		
					// Download the photo
					fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
					resp, err := http.Get(fileURL)
					if err != nil {
						log.Println("Failed to download file:", err)
						continue
					}
					defer resp.Body.Close()
		
					// Save the photo to a file
					filePath := fmt.Sprintf("%s.jpg", photo.FileID)
					fileData, err := os.Create(filePath)
					if err != nil {
						log.Println("Failed to create file:", err)
						continue
					}
					defer fileData.Close()
		
					// Copy the downloaded photo data to the file
					_, err = io.Copy(fileData, resp.Body)
					if err != nil {
						log.Println("Failed to save file:", err)
						continue
					}
		
					fmt.Println("Photo downloaded:", filePath)
				}
			// }else{
			// 	emsg := tgbotapi.NewMessage(msg.ChatID,"Belum ada text !")
			// 	if _, err = bot.Send(emsg); err != nil {
			// 		panic(err)
			// 	}
			// }

			

			// Send the message.
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}