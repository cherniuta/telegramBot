package telegram

import (
	"errors"
	"go/lib/e"
	"go/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

// все команды, которые сможет отправлять бот
// будем смотреть на текст сообщения и будем понимать что это за команда
func (p *Processor) doCmd(text string, chatID int, username string) error {
	//удалим из тектса сообщения лишние пробелы
	text = strings.TrimSpace(text)
	//пропишем логи для отслеживания того,кто нашему боту что пишет
	log.Printf("got new command '%s' from '%s'", text, username)

	//является ли текст ссылкой
	if isAffCmd(text) {
		return p.savePage(chatID, text, username)

	}
	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)

	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	//подготовим станицу, которую хотим сохранить
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	//проверяем существет ли уже такая старница
	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	//если уже существует, то отправляем сообщение,что ссылка уже сохранина
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	//пытаемся сохранить страницу
	if err := p.storage.Save(page); err != nil {
		return err
	}

	//если страница корректнго сохранилась, то сообщаем об этом пользователю
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	//ищем случайную статью
	page, err := p.storage.PickRandom(username)

	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	//особый тип ошибок, когда нет сохраненых страниц
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	//если же мф что-то нашли, отправляем эту ссылку пользователю
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	//если мф нашли и отправили ссылку, то нужно обязательно ее удалить
	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAffCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	//распарсим текст считая его ссылкой
	u, err := url.Parse(text)

	//текст мф считаем ссылкой если ошибка нулевая и есть указанный хост
	return err == nil && u.Host != ""
}
