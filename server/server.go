package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cotacao struct {
	ID     int `gorm:"PrimaryKey"`
	USDBRL struct {
		Ask        string `json:"ask"`
		Bid        string `json:"bid"`
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		CreateDate string `json:"create_date"`
		High       string `json:"high"`
		Low        string `json:"low"`
		Name       string `json:"name"`
		PctChange  string `json:"pctChange"`
		Timestamp  string `json:"timestamp"`
		VarBid     string `json:"varBid"`
	} `json:"USDBRL"`
}

type Dolar struct {
	ID  int `gorm:"PrimaryKey"`
	Bid string
}

func main() {
	http.HandleFunc("/", QuotationHandler)
	http.ListenAndServe(":8080", nil)
}

func QuotationHandler(w http.ResponseWriter, r *http.Request) {

	quotation, err := SearchQuotation()
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	SaveQuotation(*quotation)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotation.USDBRL.Bid)
}

func SearchQuotation() (*Cotacao, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var data Cotacao

	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	return &data, nil
}

func SaveQuotation(bid Cotacao) {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Dolar{})

	gormCtx, gormCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer gormCancel()
	db.WithContext(gormCtx).Create(&Dolar{
		Bid: bid.USDBRL.Bid,
	})
}
