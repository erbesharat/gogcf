package sendsms

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sfreiberg/gotwilio"
)

// Params - SendSMS query parameters
type Params struct {
	from    string
	to      string
	message string
}

var errorFormat = "{\"error\": {\"message\": \"%s\"}}"

// SendSMS - Request: example.com/?from=+999&to=+999&message=testing
func SendSMS(w http.ResponseWriter, r *http.Request) {
	accountSid := os.Getenv("TWILIO_SID")
	authToken := os.Getenv("TWILIO_TOKEN")
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	args := Params{
		from:    r.URL.Query().Get("from"),
		to:      r.URL.Query().Get("to"),
		message: r.URL.Query().Get("message"),
	}
	sms, _, err := twilio.SendSMS(args.from, args.to, args.message, "", "")
	if err != nil {
		errorResponse(w, err)
		return
	}
	chargeJSON, err := json.Marshal(sms)
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
