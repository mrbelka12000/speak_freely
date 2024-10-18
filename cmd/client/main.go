package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mrbelka12000/linguo_sphere_backend/internal/models"
	"github.com/mrbelka12000/linguo_sphere_backend/pkg/pointer"
)

func main() {
	obj := models.UserCU{
		FirstName:  getStrPointer("bekabekabekabekabekabekabekabekabekabeka"),
		LastName:   getStrPointer("teka"),
		Nickname:   getStrPointer("beka"),
		Email:      getStrPointer("karshyga.beknur@gmail.com"),
		Password:   getStrPointer("1111b"),
		AuthMethod: getIntPointer(1),
		LanguageID: pointer.Of(int64(1)),
	}

	body, _ := json.Marshal(&obj)
	resp, err := http.Post("http://localhost:8081/api/v1/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody), resp.StatusCode)

	objl := models.UserLogin{
		Login:    "beka",
		Password: "1111b",
	}
	body, _ = json.Marshal(&objl)

	resp, err = http.Post("http://localhost:8081/api/v1/login", "application/json", bytes.NewBuffer(body))
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
