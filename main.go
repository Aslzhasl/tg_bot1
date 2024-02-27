package main

import (
	"log"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

var (
	bot *tgbotapi.BotAPI

	conversationPaused bool

	dailyMessages = map[int]struct {
		Message  string
		ImageURL string
	}{
		1: {"Monday's message goes here.", "https://example.com/monday_image.jpg"},
		2: {"Tuesday's message goes here.", "https://example.com/tuesday_image.jpg"},
		3: {"Wednesday's message goes here.", "https://example.com/wednesday_image.jpg"},
		4: {"Thursday's message goes here.", "https://example.com/thursday_image.jpg"},
		5: {"Friday's message goes here.", "https://example.com/friday_image.jpg"},
		6: {"Saturday's message goes here.", "https://example.com/saturday_image.jpg"},
		7: {"Sunday's message goes here.", "https://example.com/sunday_image.jpg"},
	}

	// Variable to store the messages for /give term command
	giveTermMessages = []string{
		"Интерфейс - граница между двумя функциональными объектами, требования к которой определяются стандартом; совокупность средств, методов и правил взаимодействия (управления, контроля и т. д.) между элементами системы",
		"Фронтенд -презентационная часть web приложений",
		"Компилировать -  составление какого-либо текста, произведения путём использования чужих текстов, трудов без самостоятельной обработки источников и без ссылок на авторов ",
		"Тестить- это процесс проверки программного обеспечения на соответствие требованиям, выявление ошибок и дефектов. ",
		"Парсить - означает анализировать и обрабатывать данные в определенном формате или структуре. В компьютерном программировании, это часто относится к разбору строки или текстового файла на составляющие элементы с целью извлечения нужной информации или выполнения определенных действий в зависимости от содержания данных.",
		"Митап – 'собрание специалистов определенной сферы деятельности для обмена опытом, в образовательных целях или просто для общения в неформальной обстановке'",
		"Патч – 'информация, предназначенная для автоматизированного внесения определѐнных изменений в компьютерные файлы'",
		"рефакторинг – процесс изменения внутренней структуры программы",
		"синиор – человек, имеющий опыт работы от 5 лет и более",
		"аутсорс – 'передача организацией, на основании договора, определѐнных видов или функций производственной предпринимательской деятельности другой компании, действующей в нужной области'",
		"session – сеанс, который имеет два значения для компьютерных систем: а) связь пользователя и компьютера; б) последовательность операций по установлению соединения и обмену данными между компьютерами",
	}

	currentTermMessageIndex int // To keep track of the current term message index
	selectedTime            string
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env not loaded")
	}

	botToken := os.Getenv("TG_API_BOT_TOKEN")

	var botErr error
	bot, botErr = tgbotapi.NewBotAPI(botToken)
	if botErr != nil {
		log.Fatal("Error initializing bot:", botErr)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Initialize cron scheduler

	// Start listening for incoming updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("Error getting update channel:", err)
	}

	// Inside the for loop where updates are received
	for update := range updates {
		// Check if the update is a message
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID

		// Handle commands
		switch update.Message.Text {
		case "/start":
			if conversationPaused {
				sendMainMenu(chatID)
				conversationPaused = false
			} else {
				sendStartInfo(chatID)
			}
		case "/give term":
			if !conversationPaused {
				sendNextTermMessage(chatID)
			}
		case "/set time":
			sendTimeSelectionMessage(chatID)
		case "/off":
			conversationPaused = true
			reply := tgbotapi.NewMessage(chatID, "Разговор приостановлен. Пожалуйста, используйте /start, чтобы возобновить.")
			if _, err := bot.Send(reply); err != nil {
				log.Println(err)
			}
		default:
			if !conversationPaused {
				reply := tgbotapi.NewMessage(chatID, "Неверная команда. Пожалуйста, выберите из меню.")
				if _, err := bot.Send(reply); err != nil {
					log.Println(err)
				}
				sendMainMenu(chatID)
			}
		}
	}

}

func sendMainMenu(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выбрать опцию :")
	msg.ReplyMarkup = mainMenuKeyboard()
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func mainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	buttons := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("/start"),
		tgbotapi.NewKeyboardButton("/give term"),
		tgbotapi.NewKeyboardButton("/set time"),
		tgbotapi.NewKeyboardButton("/off"),
	}
	return tgbotapi.NewReplyKeyboard(buttons)
}

func sendStartInfo(chatID int64) {
	startInfo := "Добро пожаловать в бота! Вы можете использовать следующие команды:\n" +
		"/start - Показать это сообщение\n" +
		"/give term - Получить информацию о термине\n" +
		"/set time - Установить время для ежедневных сообщений\n" +
		"/off - Приостановить разговор"
	msg := tgbotapi.NewMessage(chatID, startInfo)
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func sendTimeSelectionMessage(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Please select the time:")
	msg.ReplyMarkup = timeSelectionKeyboard()
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func timeSelectionKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var buttonRows [][]tgbotapi.KeyboardButton
	for i := 0; i < 24; i++ {
		hour := strconv.Itoa(i)
		button := tgbotapi.NewKeyboardButton(hour)
		buttonRow := []tgbotapi.KeyboardButton{button}
		buttonRows = append(buttonRows, buttonRow)
	}
	return tgbotapi.NewReplyKeyboard(buttonRows...)
}

func handleTimeSelection(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	// Assuming the user selects time in the format of an integer from 0 to 23
	selectedTime = text // Store the selected time in the selectedTime variable

	// Confirm the selected time to the user
	sendTimeConfirmationMessage(chatID, selectedTime)
}

func sendDailyMessages(chatID int64) {
	if selectedTime != "" {
		// Get today's message and image URL based on the day of the week
		now := time.Now()
		dayOfWeek := int(now.Weekday())
		messageData, ok := dailyMessages[dayOfWeek]
		if !ok {
			// If there's no message for today, send a default message
			messageData = dailyMessages[1] // Monday's message as default
		}

		// Send the message
		msg := tgbotapi.NewMessage(chatID, messageData.Message)
		msg.ParseMode = tgbotapi.ModeMarkdown
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}

		// Send the image
		imageURL := messageData.ImageURL
		imageMsg := tgbotapi.NewPhotoShare(chatID, imageURL)
		if _, err := bot.Send(imageMsg); err != nil {
			log.Println(err)
		}
	}
}

func handlePausedConversation(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text
	if text == "/start" {
		sendStartInfo(chatID)
		conversationPaused = false
	} else {
		msg := tgbotapi.NewMessage(chatID, "Conversation paused. Please use /start to resume.")
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Error sending pause message:", err)
		}
	}
}

func sendNextTermMessage(chatID int64) {
	if currentTermMessageIndex < len(giveTermMessages) {
		// Get the term message
		message := giveTermMessages[currentTermMessageIndex]

		// Get the corresponding image URL from dailyMessages map
		dayOfWeek := currentTermMessageIndex + 1
		messageData, ok := dailyMessages[dayOfWeek]
		if !ok {
			log.Println("Image URL not found for the current term message.")
			return
		}
		imageURL := messageData.ImageURL

		// Send the message with text and image
		msg := tgbotapi.NewMessage(chatID, message)
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = mainMenuKeyboard() // Optional: add keyboard for additional actions
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
			return
		}

		// Send the image
		imageMsg := tgbotapi.NewPhotoShare(chatID, imageURL)
		if _, err := bot.Send(imageMsg); err != nil {
			log.Println(err)
		}

		currentTermMessageIndex++
	} else {
		msg := tgbotapi.NewMessage(chatID, "No more term messages available.")
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}

func sendTimeConfirmationMessage(chatID int64, selectedTime string) {
	msg := tgbotapi.NewMessage(chatID, "Daily messages will be sent at "+selectedTime+" every day.")
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
	selectedTime = selectedTime // Update selectedTime variable
}
