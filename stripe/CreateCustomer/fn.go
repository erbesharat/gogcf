package createcustomer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
)

// Params - charge http request parameters
type Params struct {
	email       string
	description string
	token       string
}

var errorFormat = "{\"error\": {\"message\": \"%s\"}}"

// CreateCustomer - Request: example.com/:email/:desc/:token
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	args := Params{
		email:       r.URL.Query().Get("email"),
		description: r.URL.Query().Get("desc"),
		token:       r.URL.Query().Get("token"),
	}

	params := &stripe.CustomerParams{
		Email:       stripe.String(args.email),
		Description: stripe.String(args.description),
	}
	params.SetSource(args.token)
	cus, err := customer.New(params)
	if err != nil {
		errorResponse(w, err)
		return
	}

	customerJSON, err := json.Marshal(cus)
	if err != nil {
		errorResponse(w, err)
		return
	}
	if _, err := w.Write(customerJSON); err != nil {
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
