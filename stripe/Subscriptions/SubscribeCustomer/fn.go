package subscribecustomer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/sub"

	stripe "github.com/stripe/stripe-go"
)

// Params - SubscribeCustomer query parameters
type Params struct {
	customerID string
}

var errorFormat = "{\"error\": {\"message\": \"%s\"}}"

// SubscribeCustomer - Request: example.com/?customer_id=test
func SubscribeCustomer(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_KEY")

	args := Params{
		customerID: r.URL.Query().Get("customer_id"),
	}

	params := &stripe.SubscriptionParams{
		Customer: stripe.String(args.customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String("plan-name"),
			},
		},
	}

	s, err := sub.New(params)
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
