package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/thumbnail", thumbnailHandler)

	fs := http.FileServer(http.Dir("./frontend/dist"))
	http.Handle("/", fs)

	fmt.Println("Server listing on port 3000")
	log.Panic(
		http.ListenAndServe(":3000", nil),
	)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type thumbnailRequest struct {
	Url string `json:"url"`
}

type screenshotAPIRequest struct {
	Token          string `json:"token"`
	Url            string `json:"url"`
	Output         string `json:"output"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	ThumbnailWidth int    `json:"thumbnail_width"`
}

func thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	var decoded thumbnailRequest

	// Decode the request into the shumbnailRequest struct
	err := json.NewDecoder(r.Body).Decode(&decoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a struct with the parameters needed to call the ScreenshotAPI
	apiRequest := screenshotAPIRequest{
		Token:          "P0RSGH2-467MWZQ-QKEQZA1-FB2VXVM",
		Url:            decoded.Url,
		Width:          1920,
		Height:         1080,
		Output:         "json",
		ThumbnailWidth: 300,
	}

	// Convert the struct to a JSON string
	jsonString, err := json.Marshal(apiRequest)
	checkError(err)

	// Create a HTTP request
	req, err := http.NewRequest("POST", "https://shot.screenshotapi.net/screenshot", bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")
	checkError(err)

	// Execute the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	checkError(err)

	// Tell Go to close the response at the end of the function
	defer response.Body.Close()

	// Read the raw response into a Go struct
	type ScreenShotAPIResponse struct {
		Screenshot string `json"screenshot"`
	}
	var apiResponse ScreenShotAPIResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	checkError(err)

	// Pass back the screenshot URL to the frontend
	_, err = fmt.Fprintf(w, `{"screenshot": "%s"}`, apiResponse.Screenshot)
	checkError(err)

}
