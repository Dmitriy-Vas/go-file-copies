package main

//#region Header
import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
)

const (
	OS_PERMISSIONS os.FileMode = 0644
)

var (
	config     Config
	duplicates map[uint32][]File
)

type Config struct {
	Directories []string `json:"directories"`
}

type File struct {
	Hash uint32
	Path string
}

//#endregion

func isErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//#region Storage
func getConfig() {
	data, err := ioutil.ReadFile("config.json")
	isErr(err)
	err = json.Unmarshal(data, &config)
	isErr(err)
}

func saveConfig() {
	data, err := json.Marshal(&config)
	isErr(err)
	err = ioutil.WriteFile("config.json", data, OS_PERMISSIONS)
	isErr(err)
	fmt.Println("Config successfully saved!")
}

func saveResult() {
	data, err := json.Marshal(&duplicates)
	isErr(err)
	err = ioutil.WriteFile("duplicates.json", data, OS_PERMISSIONS)
	isErr(err)
}

//#endregion

func getCRCHash(name string) uint32 {
	data, err := ioutil.ReadFile(name)
	isErr(err)
	return crc32.ChecksumIEEE(data)
}

func getFilesToScan() {
	duplicates = make(map[uint32][]File)
	// По очереди берём пути из конфига и проходим по ним, собирая все найденные файлы
	for _, dir := range config.Directories {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			// Проверяем путь на директорию, добавляем только файлы
			if !info.IsDir() {
				// Получаем хэш файла
				hash := getCRCHash(path)
				f := File{
					Path: path,
					Hash: hash,
				}
				duplicates[hash] = append(duplicates[hash], f)
			}
			return nil
		})
		isErr(err)
	}

	fmt.Println("Found", len(duplicates), "hashes")
}

func checkFiles() {
	result := make(map[uint32][]File)

	for h, v := range duplicates {
		if len(v) == 1 {
			continue
		}
		result[h] = duplicates[h]
	}
	duplicates = result

	fmt.Println("Found", len(duplicates), "copies")
}

func main() {
	// Загружаем конфиг с диска
	getConfig()

	// Отлавливаем Ctrl^C для сохранения конфига
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)

		<-signalChan

		saveConfig()
		os.Exit(0)
	}()
	// Сохраняем конфиг при ошибке или если программа закончила выполнение
	defer saveConfig()

	getFilesToScan()

	checkFiles()

	saveResult()
}
