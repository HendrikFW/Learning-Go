package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Post is a post
type Post struct {
	UserID int
	ID     int
	Title  string
	Body   string
}

func main() {
	fmt.Println("Learning Go: web request json")

	response, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("Response status: %s\n", response.Status)

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	var p Post
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("%+v", p)
}
