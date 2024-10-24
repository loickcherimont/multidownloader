package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	M "example.com/downloader/models"
	"github.com/gorilla/mux"
)

// Read data and store them in a variable type []byte
func readData(fileName string, channel chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := os.ReadFile("./data/" + fileName + ".jpg")
	if err != nil {
		log.Fatal(err)
		channel <- nil
		return
	}
	channel <- data
}

// Write data into a folder /out
func writeData(channel chan []byte) {
	defer fmt.Println("Download OK")
	var i int = 0

	for data := range channel {
		path := fmt.Sprintf("./out/imageD%d.jpg", i+1)
		file, err := os.Create(path)
		if err != nil {
			log.Println("Error creating file:", err)
		}
		if _, err := file.Write(data); err != nil {
			log.Println("Error writing file:", err)
		}
		file.Close()
		i++
	}

}
func main() {
	var wg sync.WaitGroup
	channel := make(chan []byte, 2)

	tmpl := template.Must(template.ParseFiles("./index.html"))

	router := mux.NewRouter()

	// Main handler
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		tmpl.Execute(w, M.DownloadData{
			Images:  []string{"image1", "image2"},
			Message: "",
		})
	})

	// Handler for multi-download
	router.HandleFunc("/download", func(w http.ResponseWriter, req *http.Request) {
		/* Multi-download code */
		wg.Add(2)
		go readData("image1", channel, &wg)
		go readData("image2", channel, &wg)

		go func() {
			wg.Wait()
			close(channel)
		}()
		writeData(channel)
		/* End-Multi-download code */
		tmpl.Execute(w, M.DownloadData{
			Images:  []string{"image1", "image2"},
			Message: "Download OK",
		})
	})

	// Server settings
	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server up and running on 8080")
	log.Fatal(server.ListenAndServe())

}
