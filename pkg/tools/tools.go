package tools

import (
	"io"
	"net/http"
	"time"
	"vacancies/pkg/telegram"
)

var PlatformHeaders = map[string]map[string]string{
	"headhunter": {
		"User-Agent": "Mozilla/5.0 (iPad; CPU OS 7_2_1 like Mac OS X; en-US) AppleWebKit/533.14.6 (KHTML, like Gecko) Version/3.0.5 Mobile/8B116 Safari/6533.14.6",
		// "Authorization": "Bearer " + os.Getenv("HEADHUNTER_TOKEN"),
		"Authorization": "Bearer QQAVSIBVU4B0JCR296THKB22JP05A92H329U49TDD9CRIS8DT9BRPPT7M9OLQ6HD",
	},
	"trudvsem": {},

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