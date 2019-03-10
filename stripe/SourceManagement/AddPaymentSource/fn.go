package addpaymentsource

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentsource"
)

// Params - AddPaymentSource query parameters
type Params struct {
	customerID string
	token      string
}

var errorFormat = "{\"error\": {\"message\": \"%s\"}}"

// AddPaymentSource - Request: example.com/?customer_id=test&token=testing
func AddPaymentSource(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_KEY")

	args := Params{
		customerID: r.URL.Query().Get("customer_id"),
		token:      r.URL.Query().Get("token"),
	}

	params := &stripe.CustomerSourceParams{
		Customer: stripe.String(args.customerID),
		Source: &stripe.SourceParams{
			Token: stripe.String(args.token),
		},
	}

	if err := params.SetSource(args.token); err != nil {
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
