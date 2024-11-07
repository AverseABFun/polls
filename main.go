package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const url = "https://static01.nyt.com/elections-assets/pages/data/2024-11-05/results-president.json"

type Data struct {
	Races []struct {
		Outcome struct {
			Won   []string
			Votes int `json:"electoral_votes"`
		}
	}
}

func getJSON() (Data, error) {
	res, err := http.Get(url)
	if err != nil {
		return Data{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return Data{}, err
		}
		var out = Data{}
		err = json.Unmarshal(bodyBytes, &out)
		return out, err
	}
	return Data{}, errors.New("got non-200 status code")
}

func main() {
	var old = map[string]int{
		"trump":  0,
		"harris": 0,
	}
	for {
		data, err := getJSON()
		if err != nil {
			panic(err)
		}

		var votes = map[string]int{
			"trump":  0,
			"harris": 0,
		}

		for _, race := range data.Races {
			if len(race.Outcome.Won) == 0 {
				continue
			}
			if race.Outcome.Won[0] == "trump-d" {
				votes["trump"] += race.Outcome.Votes
			} else if race.Outcome.Won[0] == "harris-k" {
				votes["harris"] += race.Outcome.Votes
			}
		}

		if old["trump"] != votes["trump"] || old["harris"] != votes["harris"] {
			fmt.Println("\bNew data!")
		}

		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

		fmt.Printf("Trump: %d\nHarris: %d\n", votes["trump"], votes["harris"])

		old = votes

		time.Sleep(5 * time.Second)
	}
}
