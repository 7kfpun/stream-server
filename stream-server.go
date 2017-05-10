package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	// "strconv"
	"time"

	"github.com/nubunto/tts"
)

var (
	chars        = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func generateSessionID() string {
	id := make([]rune, 6)
	for i := range id {
		id[i] = chars[rand.Intn(len(chars))]
	}
	return string(id)
}

func speechText(text string) string {
	log.Printf("Start speech" + text)
	s, err := tts.Speak(tts.Config{
		Speak:    text,
		Language: "zh-HK",
	})
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	filename := now.Format("2006-01-02_") + generateSessionID() + ".mp3"
	err = ioutil.WriteFile(filename, s.Bytes(), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	return filename
}

func main() {
	fmt.Println("Starting http stream server")
	http.HandleFunc("/", HandleClient)
	err := http.ListenAndServe(GetPort(), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func GetPort() string {
	var port = os.Getenv("PORT")

	if port == "" {
		port = "4747"
		log.Printf("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func HandleClient(writer http.ResponseWriter, request *http.Request) {
	// First of check if Get is set in the URL
	text := request.URL.Query().Get("text")
	if text == "" {
		// Get not set, send a 400 bad request
		defaultMusic(writer, request)
		return
	}

	c := make(chan string)
	go func() {
		c <- speechText(text)
	}()

	var result string
	result = <- c
	Openfile, err := os.Open(result)
	defer Openfile.Close()  // Close after function return
	if err != nil {
		http.Error(writer, "File not found.", 404)
		return
	}

	// Get the Content-Type of the file
	// Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	// Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)

	// Get the file size
	// FileStat, _ := Openfile.Stat()                      // Get info from file
	// FileSize := strconv.FormatInt(FileStat.Size(), 10)  // Get file size as a string

	// Send the headers
	// writer.Header().Set("Content-Length", FileSize)
	writer.Header().Set("Content-Type", "video/mp4")
	writer.Header().Set("Content-Type", "video/mpeg")
	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Header().Set("Content-Disposition", "inline")

	// Send the file
	// We read 512 bytes from the file already so we reset the offset back to 0
	Openfile.Seek(0, 0)
	// 'Copy' the file to the client
	_, err = io.Copy(writer, Openfile)
	if err != nil {
		log.Print(err)
	}
}


func defaultMusic (writer http.ResponseWriter, request *http.Request) {
	// Check if file exists and open
	Openfile, err := os.Open("iloveyou.mp3")
	defer Openfile.Close()  // Close after function return
	if err != nil {
		http.Error(writer, "File not found.", 404)
		return
	}

	// Get the Content-Type of the file
	// Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	// Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)

	// Get the file size
	// FileStat, _ := Openfile.Stat()                      // Get info from file
	// FileSize := strconv.FormatInt(FileStat.Size(), 10)  // Get file size as a string

	// Send the headers
	// writer.Header().Set("Content-Length", FileSize)
	writer.Header().Set("Content-Type", "video/mp4")
	writer.Header().Set("Content-Type", "video/mpeg")
	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Header().Set("Content-Disposition", "inline")

	// Send the file
	// We read 512 bytes from the file already so we reset the offset back to 0
	for {
		Openfile.Seek(0, 0)
		// 'Copy' the file to the client
		_, err = io.Copy(writer, Openfile)
		if err != nil {
			log.Print(err)
			return
		}

		fmt.Println("Finish once")
	}

	return
}
