package cancelsubscription

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/sub"

	stripe "github.com/stripe/stripe-go"
)

// Params - charge http request parameters
type Params struct {
	subID string
}

var errorFormat = "{\"error\": {\"message\": \"%s\"}}"

// CancelSubscription - Request: example.com/:sub_id
func CancelSubscription(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_KEY")

	args := Params{
		subID: r.URL.Query().Get("sub_id"),
	}

	s, err := sub.Cancel(args.subID, nil)
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
