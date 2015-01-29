package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get() {
	resp, err := http.Get("http://127.0.0.1:8888")
	if err != nil {
		// do something
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// do something
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(body))
}

func main() {
	Get()
}
