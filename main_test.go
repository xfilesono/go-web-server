package main

import (
	"encoding/json"
	"fmt"
	"go-web-server/config"
	"go-web-server/models"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestCreateMovie(t *testing.T) {

	var (
		requestBody models.Mov
	)

	var err error
	db, err = config.GetDB()
	if err != nil {
		log.Fatalln(http.StatusNotFound)
	}

	requestBody.Cast = "Tom Cruise"
	requestBody.Country = "ABD"
	requestBody.Title = "Top Gun: Maverick"
	requestBody.Director = "Joseph Kosinski"
	requestBody.Genre = "Action, Drama"
	requestBody.ID = "199"
	requestBody.Year = "2022"

	f, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}

	reqBody := strings.NewReader(string(f))

	s := httptest.NewServer(http.HandlerFunc(createMovie))
	request, err := http.NewRequest(http.MethodPost, s.URL, reqBody)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	if string(body) != strconv.Itoa(http.StatusCreated) {
		t.Errorf("Unexpected body returned. Want %q, got %q", strconv.Itoa(http.StatusCreated), string(body))
	} else {
		fmt.Printf("Test passed. want: %s , got %s\n", strconv.Itoa(http.StatusCreated), string(body))
	}

}
