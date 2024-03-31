package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

var client = &http.Client{}

func getJoke() Joke {
	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	var joke Joke
	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		if err := json.Unmarshal(bodyBytes, &joke); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}
	}
	return joke
}

func main() {
	settingRoutes()
}

func settingRoutes() {
	http.HandleFunc("/v1/joke", func(w http.ResponseWriter, r *http.Request) {
		joke := getJoke()
		payload := "{\n  \"Status\": " + fmt.Sprint(joke.Status) + ",\n" +
			"  \"Joke\": \"" + joke.Joke + "\",\n" +
			"  \"ID\": \"" + joke.ID + "\"\n}"
		fmt.Fprintf(w, "%s", payload)
	})

	count := 0
	http.HandleFunc("/v1/echo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", "echoing "+strconv.Itoa(count))
		count += 2
	})

	arr := []string{"Tito", "Cyrus", "Jen", "Scott", "Joke", "Go"}
	http.HandleFunc("/v1/random", func(w http.ResponseWriter, r *http.Request) {
		index := rand.Intn(len(arr))
		fmt.Fprintf(w, "%s", arr[index])
	})

	port := "80"

	fmt.Print("Starting server at port " + port + "...\n")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Print("Error starting server: ", err)
	}
}
