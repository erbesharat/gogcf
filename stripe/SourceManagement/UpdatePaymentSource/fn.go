package updatepaymentsource

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/paymentsource"
)

// Params - UpdatePaymentSource query parameters
type Params struct {
	customerID string
	oldToken   string
	newToken   string
}

var errorFormat = "{\"error\": {\"message\": \"%s\"}}"

// UpdatePaymentSource - Request: example.com/?customer_id=test&old_token=testing&new_token=newtesting
func UpdatePaymentSource(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_KEY")

	args := Params{
		customerID: r.URL.Query().Get("customer_id"),
		oldToken:   r.URL.Query().Get("old_token"),
		newToken:   r.URL.Query().Get("new_token"),
	}

	cardParams := &stripe.CardParams{
		Customer: stripe.String(args.customerID),
	}

	_, err := card.Del(args.oldToken, cardParams)
	if err != nil {
		errorResponse(w, err)
		return
	}

	params := &stripe.CustomerSourceParams{
		Customer: stripe.String(args.customerID),
		Source: &stripe.SourceParams{
			Token: stripe.String(args.newToken),
		},
	}

	if err := params.SetSource(args.newToken); err != nil {
		errorResponse(w, err)
		return
	}

	s, err := paymentsource.New(params)
	if err != nil {
		errorResponse(w, err)
		return
	}

	sourceJSON, err := json.Marshal(s)
	if err != nil {
		errorResponse(w, err)
		return
	}
	if _, err := w.Write(sourceJSON); err != nil {
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
