package telegram

import "go/clients/telegram"

//реализовывать оба интерфейса будет один тип данных
type Processor struct {
	//телеграм клиент
	tg     *telegram.Client
	offset int
}

//слздаем экземпляр процессора
func New(client *telegram.Client) {
	
}
