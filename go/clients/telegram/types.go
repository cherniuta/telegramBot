package telegram

// в getUpdate будет еще прочая информация кроме update
// поэтому нужно находить ok и поле result с апдейтами
type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}
type Update struct {
	ID      int `json:"update_id"` //теги json для того,чтобы правильно парсить и находит нужный кусок
	Message int `json:"message"`
}
