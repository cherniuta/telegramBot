package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"tgBotReadAdviser/lib/e"
)

//здесь как и в случае с event будет храниться только интерфейс
//поэтому мы сможем работать с любым сервисом
//будем работать тут с файлами

type Storage interface {
	//сохраняет старницу на вход
	// передавать сарницу будем по ссылке потому что в теории тип(page) может расширяется
	//и если мы будем передавать по значению,то все поля будут копировать , а это не выгодно
	Save(ctx context.Context, p *Page) error
	//какому именно человеку нужно скинуть ссылку
	PickRandom(ctx context.Context, userName string) (*Page, error)
	//удаление
	Remove(ctx context.Context, p *Page) error
	//проверяет есть ли такая ссылка
	IsExists(ctx context.Context, p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved page")

// наша страница(ссылка)
type Page struct {
	//ссылка
	URL string
	//имя пользователя, который ее скинул,чтобы понимать кому ее отдавать
	UserName string
}

// хешируем данные
func (p Page) Hash() (string, error) {
	h := sha1.New()

	//делаем хеш исходя из страницы и пользователя который ее скинул
	//тк одному пользователю нельзя добавлять одну и ту же ссылку
	//а вот разные пользователи могут добавить одну и ту же ссылку
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	//возвращаем сумму хэшей,но она в байтах, поэтому конверитруем ее в строку
	return fmt.Sprintf("%x", h.Sum(nil)), nil

}
