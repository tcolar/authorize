// History: Oct 01 13 tcolar Creation
// Authorize.net transaction support

// AuthorizeNet : Some authorize.net Credit Card processing support
// API  Docs: http://www.authorize.net/support/AIM_guide.pdf
package authorize

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var AUTHORIZE_GW = "https://secure.authorize.net/gateway/transact.dll"

// Various types of Authorize.net transaction types
const (
	AUTH_ONLY          = "AUTH_ONLY"
	AUTH_CAPTURE       = "AUTH_CAPTURE"
	PRIOR_AUTH_CAPTURE = "PRIOR_AUTH_CAPTURE"
	CAPTURE_ONLY       = "CAPTURE_ONLY"
	CREDIT             = "CREDIT"
	VOID               = "VOID"
)

// Authorize.net base structure
type AuthorizeNet struct {
	Login     string // authorize.net login
	Key       string //authorize.net key
	DupWindow int    // duplicate window
	TestMode  bool   // test mode or not
}

// Authorize a transaction (does not charge) returns a response
func (a AuthorizeNet) Authorize(card CardInfo, data AuthData, emailCustomer bool) (response AuthorizeResponse) {
	vals := url.Values{
		"x_type":     {AUTH_ONLY},
		"x_login":    {a.Login},
		"x_tran_key": {a.Key},

		"x_version":          {"3.1"},
		"x_method":           {"CC"},
		"x_delim_data":       {"TRUE"},
		"x_delim_char":       {"|"},
		"x_encap_char":       {`"`},
		"x_relay_response":   {"FALSE"},
		"x_duplicate_window": {fmt.Sprintf("%d", a.DupWindow)},

		"x_email_customer": {a.boolToStr(emailCustomer)},
	}

	card.AddToUrlValues(vals)
	data.AddToUrlValues(vals)

	if a.TestMode {
		vals.Set("x_test_request", "TRUE")
	} else {
		vals.Set("x_test_request", "FALSE")

	}
	response = a.Post(vals)
	return response
}

// CapturePreauth: Captures a previously authorized card
// An Empty ammount string means full ammount
func (a AuthorizeNet) CapturePreauth(transactionId string, ammount string) (response AuthorizeResponse) {
	data := url.Values{
		"x_login":      {a.Login},
		"x_tran_key":   {a.Key},
		"x_version":    {"3.1"},
		"x_method":     {"CC"},
		"x_delim_data": {"TRUE"},
		"x_delim_char": {"|"},
		"x_encap_char": {`"`},
		"x_type":       {PRIOR_AUTH_CAPTURE},
		"x_trans_id":   {transactionId},
	}
	if len(ammount) > 0 {
		data["x_amount"] = []string{ammount}
	}
	if a.TestMode {
		data.Set("x_test_request", "TRUE")
	} else {
		data.Set("x_test_request", "FALSE")
	}
	response = a.Post(data)
	return response
}

// Post: posts a query to atuhorize.net and returns the response
func (a AuthorizeNet) Post(data url.Values) (response AuthorizeResponse) {
	resp, err := http.PostForm(AUTHORIZE_GW, data)
	defer resp.Body.Close()
	if err != nil {
		log.Print(err)
		response.ReasonText = "Failed to connect to payment gateway."
		return response
	}

	if resp.StatusCode == 200 { // OK
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		//log.Print(string(bodyBytes))
		response = a.ParseResponse(string(bodyBytes))
	} else {
		msg := fmt.Sprintf("Authorize GW status code %d", resp.StatusCode)
		log.Print(msg)
		response.ReasonText = "Payment gateway returned an error."
	}
	return response
}

// parseResponse parses the response string into an AuthorizeResponse
func (a AuthorizeNet) ParseResponse(response string) AuthorizeResponse {
	data := strings.Split(response, "|")
	// remove the encap (")
	for i, d := range data {
		data[i] = d[1 : len(d)-1]
	}
	return AuthorizeResponse{
		Code:       data[0],
		ReasonCode: data[2],
		ReasonText: data[3],
		AuthCode:   data[4],
		AvsResp:    data[5],
		TransId:    data[6],
		Amount:     data[9],
		TransType:  data[11],
		Tax:        data[32],
		TransMd5:   data[37],
		CvvResp:    data[38],
	}
}

func (a AuthorizeNet) boolToStr(b bool) string {
	if b {
		return "TRUE"
	} else {
		return "FALSE"
	}
}

// AuthorizeResponse Data structure
type AuthorizeResponse struct {
	Code       string
	ReasonCode string
	ReasonText string
	AuthCode   string
	AvsResp    string
	TransId    string
	Amount     string
	TransType  string
	Tax        string
	TransMd5   string
	CvvResp    string
}

// Approved: whether this transaction was approved
func (r *AuthorizeResponse) Approved() bool {
	return r.Code == "1"
}

// Response descriptive text
func (r *AuthorizeResponse) String() string {
	ok := "failed"
	if r.Approved() {
		ok = "succeeded"
	}
	return fmt.Sprintf("Card transaction %s - %s - Code: %s/%s - Reason: %s -. AVS: %s - CVV: %s.",
		r.TransId, ok, r.Code, r.ReasonCode, r.ReasonText, r.AvsResp, r.CvvResp)
}
