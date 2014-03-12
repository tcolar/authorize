// History: Mar 12 14 tcolar Creation

package authorize

import (
	"log"
	"testing"
)

func TestCapture(t *testing.T) {
	auth := AuthorizeNet{
		Login:     "<YourLogin>",
		Key:       "<YourKey>",
		DupWindow: 120,
		TestMode:  true,
	}

	card := CardInfo{
		Number: "4111111111111111",
		Cvv:    "555",
		Month:  "11",
		Year:   "2018",
		Method: METHOD_VISA,
	}

	data := AuthData{
		InvoiceNumber: "123",
		Amount:        "5.56",
		Description:   "My Test transaction",
		// ....
		// Fill in the rest of AuthData: adress etc ...
		// ....
	}

	// Authorize a payment
	response := auth.Authorize(card, data, false)
	if !response.Approved() {
		log.Print(response)
		return
	}
	log.Print(response)
	log.Printf("Successful Authorization with id: %s ", response.TransId)

	// Example of capture the preious authorization (using the transactionId)
	response = auth.CapturePreauth(response.TransId, "5.56")
	if !response.Approved() {
		log.Print("Capture failed : ")
		log.Print(response)
		return
	}
	log.Print(response)
	log.Print("Successful Capture !")

}
