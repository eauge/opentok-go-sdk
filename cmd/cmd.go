package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/eauge/opentok-go-sdk"
)

func main() {

	apiKeyString := os.Getenv("API_KEY")
	if len(apiKeyString) == 0 {
		fmt.Println("API_KEY must be set in order to use this utility")
		return
	}
	apiKey, err := strconv.Atoi(apiKeyString)
	if err != nil {
		fmt.Println("API_KEY must be an int: ", err)
		return
	}
	apiSecret := os.Getenv("API_SECRET")
	if len(apiSecret) == 0 {
		fmt.Println("API_SECRET must be set in order to use this utility")
		return
	}

	ot := opentok.New(apiKey, apiSecret)
	s, err := ot.Session(nil)
	if err != nil {
		fmt.Println("Session could not be created: err: ", err)
		return
	}

	fmt.Println("Session created: ", s)

	t, err := ot.Token(s, nil)
	if err != nil {
		fmt.Println("Token could not be created: err: ", err)
		return
	}
	fmt.Println("Token created: ", t)
}
