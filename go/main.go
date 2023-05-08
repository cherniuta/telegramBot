package main

import (
	"flag"
	"go/clients/telegram"
	event_consumer "go/consumer/event-consumer"
	telegram2 "go/events/telegram"
	"go/storage/files"
	"log"
)

// но лучше сделать так жк с флагом,как и с токеном
const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram2.New(
		telegram.New(tgBotHost, mustToken()),
		files.New(storagePath))
	//сообщение, что сервер запущен
	log.Print("servic starte")
	//запускаем консьюмера
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		//ошибка может быть если консюмер по какой-то причине аварийно остановился
		//запишет сообщение об ошибке и остановит программу
		log.Fatal("service is stopped", err)
	}

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
