package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"syscall"
	"unsafe"

	"github.com/schollz/progressbar/v3"
)

type Response struct {
	Date        string `json:"date"`
	Explanation string `json:"explanation"`
	HDURL       string `json:"hdurl"`
	Title       string `json:"title"`
}

var directoryPath, _ = os.Getwd()

func printExplation(explanation string) {
	fmt.Println("\n", explanation)
}

func setWallpaper(image string) {
	fullImagePath := path.Join(directoryPath, image)
	filenameUTF16, err := syscall.UTF16PtrFromString(fullImagePath)
	if err != nil {
		log.Fatal(err)
	}

	syscall.NewLazyDLL("user32.dll").NewProc("SystemParametersInfoW").Call(
		uintptr(0x0014),
		uintptr(0x0000),
		uintptr(unsafe.Pointer(filenameUTF16)),
		uintptr(0x01|0x02),
	)

	fmt.Printf("Wallpaper set as '%s'\n", fullImagePath)
}

func archiveOldImages(image string) {
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() != image && strings.HasSuffix(file.Name(), ".jpg") {
			oldPath := path.Join(directoryPath, file.Name())
			newPath := path.Join(directoryPath, "archived", file.Name())

			err := os.Rename(oldPath, newPath)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Archived '%s'\n", file.Name())
		}
	}
}

func downloadImage(link string, date string, title string) string {
	response, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	fileName := fmt.Sprintf("[%s] %s.jpg", date, title)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bar := progressbar.DefaultBytes(
		response.ContentLength,
		"Downloading image",
	)

	_, err = io.Copy(io.MultiWriter(file, bar), response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Image downloaded as '%s'\n", fileName)
	return fileName
}

func fetchAPI(api string) Response {
	response, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	fmt.Println("API data fetched")
	return responseObject
}

func main() {
	url := "https://api.nasa.gov/planetary/apod?api_key=T7l7P7NmsfhoV0rhg7HWz8NYXaINQTaBzCTpWfZy"

	response := fetchAPI(url)
	image := downloadImage(response.HDURL, response.Date, response.Title)
	archiveOldImages(image)
	setWallpaper(image)
	printExplation(response.Explanation)

	// Keeps program running by waiting for user input
	fmt.Print("\nPress any key to exit... ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	fmt.Println(input)
}
