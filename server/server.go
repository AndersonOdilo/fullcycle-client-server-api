package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

var db *sql.DB

func main() {

	db = criaConexaoBD()
	defer db.Close()

	criaTabelaBD(db)

	http.HandleFunc("/cotacao", consultaCotacaoDolarHandler)
	http.ListenAndServe(":8080", nil)

}

func consultaCotacaoDolarHandler(w http.ResponseWriter, r *http.Request) {

	cotacaoDolar, err := consultaCotacaoDolarInsereBaseDados()

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cotacaoDolar)

}

func consultaCotacaoDolarInsereBaseDados() (Cotacao, error) {

	cotacao, err := consultaCotacaoDolar()

	if err != nil {

		return Cotacao{}, err
	}

	err = insertCotacaoBD(db, cotacao)

	if err != nil {

		return Cotacao{}, err
	}

	return cotacao, nil

}

func consultaCotacaoDolar() (Cotacao, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {

		return Cotacao{}, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {

		log.Println("Tempo maximo da request para consultar a cotação excedido")

		return Cotacao{}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {

		return Cotacao{}, err
	}

	var cotacoes map[string]Cotacao
	err = json.Unmarshal(body, &cotacoes)

	if err != nil {
		return Cotacao{}, err
	}

	return cotacoes["USDBRL"], nil

}

func criaConexaoBD() *sql.DB {
	db, err := sql.Open("sqlite3", "cotacao.db")

	if err != nil {
		panic(err)
	}

	return db
}

func criaTabelaBD(db *sql.DB) {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS cotacao (id INTEGER PRIMARY KEY, code TEXT,codein TEXT, name TEXT, high TEXT, low TEXT, varbid TEXT, pctchange TEXT, bid TEXT, ask TEXT, timestamp TEXT, createdate TEXT)")
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec()

	if err != nil {
		panic(err)
	}
}

func insertCotacaoBD(db *sql.DB, dolar Cotacao) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	stmt, err := db.Prepare("insert into cotacao(code, codein,name,high,low,varBid,pctChange,bid,ask,timestamp,createDate) values(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, dolar.Code, dolar.Codein, dolar.Name, dolar.High, dolar.Low, dolar.VarBid, dolar.PctChange, dolar.Bid, dolar.Ask, dolar.Timestamp, dolar.CreateDate)
	if err != nil {
		log.Print("Tempo maximo para realizar o insert no banco excedido")

		return err
	}

	return nil

}
