package configReader

import (
	"encoding/json"
	"io/ioutil"
)

type RssConfig struct {
	Rss           []string
	RequestPeriod int `json:"request_period"`
}

// Читает конфигурацию из файла
func Read(fileName string) RssConfig {
	fileCont, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("Файл конфигурации не существует")
	}

	var res RssConfig

	err = json.Unmarshal(fileCont, &res)
	if err != nil {
		panic("Файл конфигурации содержит не правильную JSON структуру")
	}

	return res
}
