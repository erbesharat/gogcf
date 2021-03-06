package chargecustomer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

// Params - ChargeCustomer query parameters
type Params struct {
	customerID string
	source     string
	currency   string
	amount     int64
}

var errorFormat = "{\"error\": {\"message\": \"%s\"}}"

// ChargeCustomer - Request: example.com/?customer_id=test&token=test&amount=1000
func ChargeCustomer(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	amount, _ := strconv.ParseInt(r.URL.Query().Get("amount"), 10, 64)
	args := Params{
		customerID: r.URL.Query().Get("customer_id"),
		source:     r.URL.Query().Get("token"),
		currency:   "usd",
		amount:     amount,
	}
	sourceParams, _ := stripe.SourceParamsFor(args.source)
	params := &stripe.ChargeParams{
		Amount:   stripe.Int64(args.amount),
		Currency: stripe.String(args.currency),
		Source:   sourceParams,
		Customer: stripe.String(args.customerID),
	}
	ch, err := charge.New(params)
	if err != nil {
		errorResponse(w, err)
		return
	}
	chargeJSON, err := json.Marshal(ch)
	if err != nil {
		errorResponse(w, err)
		return
	}
	if _, err := w.Write(chargeJSON); err != nil {
		errorResponse(w, err)
		return
	}
}

func errorResponse(w http.ResponseWriter, err error) {
	errorJSON, err := json.Marshal(fmt.Sprintf(errorFormat, err.Error()))
	if err != nil {
		log.Fatalf("Couldn't convert error to json: %s", err.Error())
	}
	if _, err = w.Write(errorJSON); err != nil {
		log.Fatalf("Couldn't write the error response: %s", err.Error())
	}
}
