package main

import (
	"flag"
	"go/clients/telegram"
	"log"
)

// но лучше сделать так жк с флагом,как и с токеном
const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken())

}

// фу-ия аварийно завершает прогамму, если токен оказался пустым (must)
func mustToken() string {
	//токен передаем из командной сторки при запуске программы
	//(имя флага,значение по умолчанию,подсказака для данного флага)
	//bot -tg-bot-token 'my token'
	token := flag.String(
		"token-bot-token",
		"",
		"token for access to telegram bot",
	)
	//значение попадает во время вызова метода парс
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified") //аварийно завершаем
	}

	return *token

}
