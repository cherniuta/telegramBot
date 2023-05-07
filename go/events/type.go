package events

// создаем того, кто будет получать запросы
// и того, кто будет их обрабатывать
type Fatcher interface {
	//не передаем в этот метод offset, тк 1)разберемся с ним внутри
	//2)в других отличных от tg сервисах может не быть такого параметра
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Processor(e Event)
}

type Type int

// все возможные варианты событий
const (
	//что-то непонятное
	Unknown Type = iota
	//сообщение
	Message
)

// объект события
type Event struct {
	Type Type
	Text string
}
