package telegram

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"vacancies/pkg/logger"

	"github.com/joho/godotenv"
)

var chats = []string{
	"544490770",  // Ярослав
	"1487312575", // Гриша
	"941543716",  // Дима
}

func getJson(url string) (json string, err error) {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	checkErr(err)
	response, err := client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	return string(data), nil
}

func checkErr(err error) {
	if err != nil {
		Mailing(fmt.Sprintf("Программа остановилась: %s", err))
		logger.Log.Fatal(err)
		panic(err)
	}
}

func getUrl() string {
	err := godotenv.Load(".env")
	checkErr(err)
	return fmt.Sprintf("https://api.telegram.org/bot%s", os.Getenv("TELEGRAM_TOKEN"))
}
