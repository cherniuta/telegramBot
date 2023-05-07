package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"go/lib/e"
	"go/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

//тип, который будет реализовывать интерфейс
type Storage struct {
	//нформация о том в какой папке будем его хранить
	basePath string
}

//параметр доступа-у всех пользователей одинаковые права(чтение и запись)
const defaultPerm = 0774

//ошибка если вайлов по указанному пути нет
var ErrNoSavedPages = errors.New("no saved page")

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	//место куда будет сохранятся наш файл
	//не используем path.join из-за не того слеша на виндовс
	//(c чего начинаетя,все ссылки одного пользователя складываем в его личную папку)
	fPath := filepath.Join(s.basePath, page.UserName)

	//создаем путь
	//создаст все директории входимые в данный путь c учетом параметра доступа
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	//название файла
	//нужно, чтобы все файлы имели уникальное имя с помощью хеша
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	//добавим в путь само имя файла
	fPath = filepath.Join(fPath, fName)

	//создаем файл и передаем путь до файла
	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	//в конце фу-ии закрывем созданный файл
	//игнорируем оштбку от Close
	defer func() { _ = file.Close() }()

	//сереализуем нашу страницу(приводим к формату, который мф запишем в файл и потом сможем восстановить исходную структуру
	//станица будет преобразована в формат gob и записана в указаный файл
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err := e.WrapIfErr("can't pick random page", err) }()

	//получаем путь до директориис файлами
	path := filepath.Join(s.basePath, userName)

	//получаем список файлов по пути
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	//если файлов нет(0), то возвращаем заранее определенную ошибку
	//для того,чтобы ее можно ыбло прверить снаружи, выносим ее в переменную
	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}

	//нужно получить случайное число от 0 и до номера последнего файла(если файлов 10, то 0-9)
	//псевдорандом (всегда будет возвращать одинаковую последовательность),поэтому мы используем не const , а текущее время
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	//получаем файл с ген номером
	file := files[n]

	//декодируем файл и вернем содержимое
	//открфываем файл и декодируем
	return s.decodePage(filepath.Join(path, file.Name()))
}

//Удаление ссылок
func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	//добавляем доп сообщение о том ,какой именно файл мы не смогли удалить
	if err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)
		return e.Wrap(msg, err)
	}

	return nil
}

//Существует ли страница(сохранял ли пользователь ее ранее)
func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	//проверяем существование файла
	//возвращает несколько вариантов ошибок
	switch _, err = os.Stat(path); {
	//проверяем что она вернула ошибку о несуществовании файла
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		//обрабатываем все остальные ошибки
		msg := fmt.Sprintf("can't check if file %s exists", path)
		return false, e.Wrap(msg, err)

	}

	return true, nil
}

//декодируем файл
func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	//открфваем файл по пути
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	//закрываем файл
	defer func() { _ = f.Close() }()

	//создаем переменную, в которую файл будет декодирован
	var p storage.Page

	//декодируем с помощью того же gob
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)

	}

	return &p, nil
}

//Определяем имя файл с помощью хэша
//выносим это все в одну функцию, чтобы если мы захотим изменить способ формирования хэша(например:дописывать расширение)
//мы смогли поменять только одну функцию, а не все места ,где был хэщ
func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}