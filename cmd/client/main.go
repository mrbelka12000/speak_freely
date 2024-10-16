package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
)

func main() {
	obj := models.UserCU{
		FirstName:  getStrPointer("bekabekabekabekabekabekabekabekabekabeka"),
		LastName:   getStrPointer("teka"),
		Nickname:   getStrPointer("naka"),
		Email:      getStrPointer("nakaf@gmail.com"),
		Password:   getStrPointer("1111b#"),
		AuthMethod: getIntPointer(1),
	}

	body, _ := json.Marshal(&obj)
	resp, err := http.Post("http://localhost:8081/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody), resp.StatusCode)

	objl := models.UserLogin{
		Login:    "naka",
		Password: "1111b",
	}
	body, _ = json.Marshal(&objl)

	resp, err = http.Post("http://localhost:8081/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, _ = io.ReadAll(resp.Body)
	fmt.Println(string(respBody), resp.StatusCode)
}

func getIntPointer(i int) *int {
	return &i
}

func getStrPointer(s string) *string {
	return &s
}
