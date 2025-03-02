package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {

	cotacaoDolar := consultaCotacaoDolar()

	gravaCotacaoDolarArquivo(cotacaoDolar)

}

func consultaCotacaoDolar() Cotacao {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {

		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {

		log.Println("tempo maximo da request excedido")

		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic(errors.New("O server retornou um erro"))
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {

		panic(err)
	}

	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)

	if err != nil {
		panic(err)
	}

	return cotacao

}

func gravaCotacaoDolarArquivo(cotacao Cotacao) {

	file, err := os.Create("cotacao.txt")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	file.WriteString(fmt.Sprintf("Dolar: %s", cotacao.Bid))
}
