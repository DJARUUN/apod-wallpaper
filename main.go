package main

import (
	"bufio"
	"embed"
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
	URL         string `json:"url"`
	Title       string `json:"title"`
	MediaType   string `json:"media_type"`
}

var directoryPath, _ = os.Getwd()

//go:embed .env
var file embed.FS

func setWallpaper(image string) {
	// Converts file path to something that is usable by the setter
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
	archivedPath := path.Join(directoryPath, "archived")

	// Creates 'archived' directory if it doesnt exist
	if _, err := os.Stat(archivedPath); os.IsNotExist(err) {
		err := os.Mkdir(archivedPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Created directory 'archived'")
	} else if err != nil {
		log.Fatal(err)
	}

	// Goes trough all files and moves any .jpg files that arent todays into the archived folder
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

func downloadImage(Response Response) string {
	fileName := fmt.Sprintf("[%s] %s.jpg", Response.Date, Response.Title)

	// Skips downloading todays image if its already downloaded
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.Name() == fileName {
			fmt.Println("Image already downloaded, skipped")
			return fileName
		}
	}

	// Not a video so do the regular stuff
	response, err := http.Get(Response.HDURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Creates file to write to
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Starts writing to file and links a progressbar to the progress
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
	// Fetches API and converts it into the wanted format
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

func getAPIKey() string {
	// Reads embedded file .env and gets apiKey from it
	data, err := file.ReadFile(".env")
	if err != nil {
		log.Fatal(err)
	}
	apiKey := strings.TrimPrefix(string(data), "API_KEY=")
	return apiKey
}

func main() {
	// Gets apiKey and splices together the full URL
	apiKey := getAPIKey()
	url := "https://api.nasa.gov/planetary/apod?api_key=" + apiKey

	// Executes all functions and skips image related ones if the media is a video
	response := fetchAPI(url)
	if response.MediaType == "video" {
		fmt.Println("Media is a video, skipped")
		fmt.Println(response.URL)
	} else {
		image := downloadImage(response)
		archiveOldImages(image)
		setWallpaper(image)
	}
	fmt.Println("\n", response.Explanation)

	// Keeps program running by waiting for user input
	fmt.Print("\nPress any key to exit... ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	fmt.Println(input)
}
