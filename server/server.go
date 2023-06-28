package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	uuid "github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DB_TIMEOUT = 10 * time.Millisecond
const REQUEST_MAX_DURATION = 200 * time.Millisecond

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

type Dollar struct {
	ID  uuid.UUID `gorm:"type:uuid;primaryKey;" json:"-"`
	Bid string    `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", QuotationHandler)
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
	bid, err := SaveQuotation(*quotation)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bid)
}

func SearchQuotation() (*Cotacao, error) {

	ctx, cancel := context.WithTimeout(context.Background(), REQUEST_MAX_DURATION)
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

func SaveQuotation(bid Cotacao) (*Dollar, error) {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Dollar{})

	gormCtx, gormCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer gormCancel()

	dollarBid := &Dollar{
		ID:  uuid.New(),
		Bid: bid.USDBRL.Bid,
	}

	db.WithContext(gormCtx).Create(&dollarBid)
	return dollarBid, nil
}
