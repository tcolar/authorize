// History: Mar 12 14 tcolar Creation

package authorize

import (
	"net/url"
)

// Data structure for authorization request
type AuthData struct {
	InvoiceNumber                                                 string
	Amount                                                        string
	Description                                                   string
	FirstName, LastName                                           string
	Company                                                       string
	Address, City, State, Zip, Country                            string // billing
	Phone, Email                                                  string
	CustomerId, CustomerIp                                        string
	ShipToFirstName, ShipToLastName, ShipToCompany, ShipToAddress string
	ShipToCity, ShipToState, ShipToZip, ShipToCountry             string
}

// AddToUrlValues : Add the AuthData values to the given url.Values map for authorize.net
func (a AuthData) AddToUrlValues(vals url.Values) {
	v := map[string][]string{
		"x_invoice_num":        {a.InvoiceNumber},
		"x_amount":             {a.Amount},
		"x_description":        {a.Description},
		"x_first_name":         {a.FirstName},
		"x_last_name":          {a.LastName},
		"x_company":            {a.Company},
		"x_address":            {a.Address},
		"x_city":               {a.City},
		"x_state":              {a.State},
		"x_zip":                {a.Zip},
		"x_country":            {a.Country},
		"x_phone":              {a.Phone},
		"x_email":              {a.Email},
		"x_cust_id":            {a.CustomerId},
		"x_customer_ip":        {a.CustomerIp},
		"x_ship_to_first_name": {a.ShipToFirstName},
		"x_ship_to_last_name":  {a.ShipToLastName},
		"x_ship_to_company":    {a.ShipToCompany},
		"x_ship_to_address":    {a.ShipToAddress},
		"x_ship_to_city":       {a.ShipToCity},
		"x_ship_to_state":      {a.ShipToState},
		"x_ship_to_zip":        {a.ShipToZip},
		"x_ship_to_country":    {a.ShipToCountry},
	}
	for k, v := range v {
		vals[k] = v
	}
}
