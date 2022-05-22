package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"imanukula/lyra-client/pkg/lyra"
	"imanukula/lyra-client/pkg/lyra/response"
)

type PayRequest struct {
	Pan          string
	Expiry       string
	SecurityCode string
}

type CreatePayRequestCustomer struct {
	Email string `json:"email"`
}

type CreatePayRequest struct {
	Amount   int                      `json:"amount"`
	Currency string                   `json:"currency"`
	OrderId  string                   `json:"orderId"`
	Customer CreatePayRequestCustomer `json:"customer"`
}

type PayAnwser struct {
	OrderStatus string `json:"orderStatus"`
}

func main() {

	epaync := lyra.NewClient()
	epaync.SetContext(context.TODO())

	http.HandleFunc("/epay", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("epaync.html"))

		// if r.Method != http.MethodPost {
		// 	tmpl.Execute(w, nil)
		// 	return
		// }

		request := CreatePayRequest{
			Amount:   25,
			Currency: "EUR",
			OrderId:  "orderId",
			Customer: CreatePayRequestCustomer{
				Email: "foo@bar.nc",
			},
		}

		resp, err := epaync.CreatePayment(request)
		if err != nil {
			log.Panicf("%v", err)
		}

		data := struct {
			PublicKey string
			EndPoint  string
			Response  response.EpayncResponse
		}{
			epaync.GetPublicKey(),
			epaync.GetEndpoint(),
			*resp,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/paid", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("paid.html"))

		r.ParseForm()
		msg := lyra.IPNMessage{
			Hash:          r.FormValue("kr-hash"),
			HashAlgorithm: r.FormValue("kr-hash-algorithm"),
			HashKey:       r.FormValue("kr-hash-key"),
			AnswerType:    r.FormValue("kr-answer-type"),
			Answer:        r.FormValue("kr-answer"),
		}

		isOkCheckHash := true
		err := lyra.CheckHash(msg)
		if err != nil {
			log.Print(err)
			isOkCheckHash = false
		}

		anwser := PayAnwser{}
		json.Unmarshal([]byte(msg.Answer), &anwser)

		data := struct {
			CheckHash bool
			Answer    PayAnwser
		}{
			isOkCheckHash,
			anwser,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":8080", nil)
}
