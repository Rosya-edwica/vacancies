package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

)


func ErrorMessageMailing(text string) {
	Mailing("🛑🛑 ОШИБКА 🛑🛑\n" + text)
}

func SuccessMessageMailing(text string) {
	Mailing("✅✅ УСПЕШНО ✅✅\n" + text)
}

func Mailing(text string) {
	for _, chat := range chats {
		SendMessage(text, chat)
	}
}

func SendMessage(text string, chatId string) (bool, error) {
	url := fmt.Sprintf("%s/sendMessage", getUrl())
	body, _ := json.Marshal(map[string]string{
		"chat_id": chatId,
		"text":    "Парсер вакансий:\n\n" + text,
	})
	response, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	return true, nil
}
