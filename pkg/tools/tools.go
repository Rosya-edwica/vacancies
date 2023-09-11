package tools

import (
	"io"
	"net/http"
	"os"
	"time"
	"vacancies/pkg/logger"
	"vacancies/pkg/telegram"

	"github.com/joho/godotenv"
)

var HH_TOKEN, SUPERJOB_ID, SUPERJOB_TOKEN, SUPERJOB_SECRET string
var PlatformHeaders map[string]map[string]string

func init() {
	err := godotenv.Load(".env")
	CheckErr(err)

	HH_TOKEN = os.Getenv("HH_TOKEN")
	SUPERJOB_TOKEN = os.Getenv("SUPERJOB_TOKEN")
	SUPERJOB_ID = os.Getenv("SUPERJOB_ID")
	SUPERJOB_SECRET = os.Getenv("SUPERJOB_SECRET")
	PlatformHeaders = map[string]map[string]string{
		"headhunter": {
			"User-Agent":    "Mozilla/5.0 (iPad; CPU OS 7_2_1 like Mac OS X; en-US) AppleWebKit/533.14.6 (KHTML, like Gecko) Version/3.0.5 Mobile/8B116 Safari/6533.14.6",
			"Authorization": "Bearer " + HH_TOKEN,
		},
		"superjob": {
			"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/117.0)",
			"Authorization": "Bearer " + SUPERJOB_TOKEN,
			"X-Api-App-Id":  SUPERJOB_SECRET,
		},
		"trudvsem": {},
		"geekjob":  {},
	}
}

func UniqueNames(names []string) (unique []string) {
	allKeys := make(map[string]bool)
	for _, item := range names {
		if _, value := allKeys[item]; !value {
			if len(item) < 2 {
				continue
			}
			allKeys[item] = true
			unique = append(unique, item)
		}
	}
	return
}

func CheckErr(err error) {
	if err != nil {
		telegram.ErrorMessageMailing(err.Error())
		panic(err)
	}
}

// TODO: Обработать случай, когда заканчивается срок действия токена Superjob
func GetJson(url string, platform string) (json string, err error) {
	client := http.Client{
		Timeout: 120 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	CheckErr(err)
	headers := getHeadersByPlatform(platform)
	if len(headers) != 0 {
		for key, val := range headers {
			req.Header.Set(key, val)
		}
	}
	response, err := client.Do(req)
	logger.Log.Printf("GET JSON status: %d by address: %s", response.StatusCode, url)
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

func getHeadersByPlatform(platform string) (headers map[string]string) {
	if val, ok := PlatformHeaders[platform]; ok {
		return val
	}
	return
}
