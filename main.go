package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("https://www.youtube.com")
	if err != nil {
		fmt.Println("ошыбка при отправке Http запроса", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ОШЫБКА", err)
		return
	}
	fmt.Println(string(body))
}
