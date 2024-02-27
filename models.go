package main

type Update struct {
	Message  Message `json:"message"`
	UpdateId int     `json:"update_id"`
}
type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}
type Chat struct {
	Id int `json:"id"`
}
type RestResponse struct {
	Result []Update `json:"result"`
}
type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

/*"Интерфейс - граница между двумя функциональными объектами, требования к которой определяются стандартом; совокупность средств, методов и правил взаимодействия (управления, контроля и т. д.) между элементами системы": "monday_image.jpg",
"Фронтенд -презентационная часть web приложений": "https://example.com/image2.jpg",
"Компилировать -  составление какого-либо текста, произведения путём использования чужих текстов, трудов без самостоятельной обработки источников и без ссылок на авторов ": "https://example.com/image3.jpg",
"Тестить- это процесс проверки программного обеспечения на соответствие требованиям, выявление ошибок и дефектов. ": */
