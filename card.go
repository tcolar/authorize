// History: Sep 30 13 tcolar Creation

package authorize

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	METHOD_AMEX       = "amex"
	METHOD_DISCOVER   = "discover"
	METHOD_MASTERCARD = "mastercard"
	METHOD_VISA       = "visa"
)

// CardInfo : Credit card data holder
type CardInfo struct {
	Number      string
	Cvv         string
	Month, Year string
	Method      string // visa etc..
}

// LastFour : Last 4 digits
func (c CardInfo) LastFour() string {
	if len(c.Number) < 4 {
		return ""
	}
	return c.Number[len(c.Number)-4 : len(c.Number)]
}

// Wipe: Overwrites all CardInfo values (to clear from memory)
func (c CardInfo) Wipe() {
	c.Cvv, c.Number, c.Month, c.Year, c.Method = "000", "0000000000000000", "01", "1970", "----------------"
}

// ValidateCard checks that the card data **LOOKS** ok (not expired, right number length etc...)
// Does **NOT** validate / authorize the card with the bank.
func (c CardInfo) ValidateCard() (err error) {
	var y, m int
	y, _ = strconv.Atoi(c.Year)
	m, _ = strconv.Atoi(c.Month)
	if y < time.Now().Year() {
		return errors.New("Expired credit card.")
	}
	if y == time.Now().Year() && m <= int(time.Now().Month()) {
		return errors.New("Expired credit card.")
	}
	if c.Number == "5555555555" {
		return nil // test credit card -> pass trough
	}
	if len(c.Cvv) < 3 || len(c.Cvv) > 4 {
		return errors.New("Invalid CVV.")
	}
	l := len(c.Number)
	if l < 13 { // always at east 13 digits
		return errors.New("Invalid credit card number.")
	}
	t := strings.ToLower(c.Method)
	switch t {
	case METHOD_VISA:
		if !strings.HasPrefix(c.Number, "4") || l != 16 && l != 13 {
			return errors.New("Invalid credit card number.")
		}
	case METHOD_MASTERCARD:
		if !(strings.HasPrefix(c.Number, "5")) || l != 16 {
			return errors.New("Invalid credit card number.")
		}
	case METHOD_DISCOVER:
		if !strings.HasPrefix(c.Number, "6011") || l != 16 {
			return errors.New("Invalid credit card number.")
		}
	case METHOD_AMEX:
		if !(strings.HasPrefix(c.Number, "34") || strings.HasPrefix(c.Number, "37")) || l != 15 {
			return errors.New("Invalid credit card number.")
		}
	default:
		return errors.New("Unsupported Credit Card type : " + t)
	}
	// TODO : Could also check the card checksum:
	// http://en.wikipedia.org/wiki/Luhn_algorithm
	return err
}

// AddToUrlValues : Add the card info to the given url.Values map for authorize.net
func (c CardInfo) AddToUrlValues(vals url.Values) {
	v := map[string][]string{
		"x_card_num":  {c.Number},
		"x_card_code": {c.Cvv},
		"x_exp_date":  {fmt.Sprintf("%s/%s", c.Month, c.Year)},
	}
	for k, v := range v {
		vals[k] = v
	}
}
