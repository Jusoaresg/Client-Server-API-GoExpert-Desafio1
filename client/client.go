package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const REQUEST_MAX_DURATION = 300 * time.Millisecond

type Dollar struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), REQUEST_MAX_DURATION)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	select {
	case <-ctx.Done():
		log.Panic("Request Timeout")
	default:
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		f, err := os.Create("cotacao.txt")
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		var d Dollar
		err = json.Unmarshal(body, &d)
		if err != nil {
			log.Fatal(err)
		}
		f.Write([]byte("Dólar: {" + d.Bid + "}"))
		f.Close()
	}

}
