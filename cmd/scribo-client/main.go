package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/tent/hawk-go"
)

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s---------------------------------\n\n", data)
	} else {
		log.Fatalf("%s", err)
	}
}

func main() {
	var body []byte
	var response *http.Response
	var request *http.Request

	client := &http.Client{}

	creds := new(hawk.Credentials)
	creds.ID = "benjamin"
	creds.Key = "werxhqb98rpaxn39848xrunpaw3489ruxnpa98w4rxn"
	creds.Hash = sha256.New

	request, err := http.NewRequest("GET", "http://localhost:8080/nodes", nil)

	if err == nil {
		request.Header.Add("Content-Type", "application/json")
		auth := hawk.NewRequestAuth(request, creds, time.Duration(0))
		request.Header.Add("Authorization", auth.RequestHeader())

		debug(httputil.DumpRequestOut(request, true))
		response, err = client.Do(request)
	}

	if err == nil {
		defer response.Body.Close()
		debug(httputil.DumpResponse(response, true))
		body, err = ioutil.ReadAll(response.Body)
	}

	if err == nil {
		fmt.Printf("%s", body)
	} else {
		log.Fatalf("Error: %s", err)
	}

}
