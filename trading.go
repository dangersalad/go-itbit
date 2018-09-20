package itbit

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

// OrderSide is the side of a trade.
type OrderSide string

// OrderType is the type of the order.
type OrderType string

// Currency is a crypto currency.
type Currency string

const (
	// OrderSideBuy indicates a buy trade.
	OrderSideBuy = "buy"
	// OrderSideSell indicates a sell trade.
	OrderSideSell = "sell"
	// OrderTypeLimit indicates a limit order.
	OrderTypeLimit = "limit"
	// CurrencyXBT is bitcoin.
	CurrencyXBT = "XBT"
	// CurrencyBTC is an alias for CurrencyXBT.
	CurrencyBTC = CurrencyXBT
)

// OrderRequest is a request to place an order on the market.
type OrderRequest struct {
	Side                  OrderSide              `json:"side"`
	Type                  OrderType              `json:"type"`
	Currency              Currency               `json:"currency"`
	Amount                float64                `json:"amount,string"`
	Display               float64                `json:"display,string"`
	Price                 float64                `json:"price,string"`
	Instrument            Symbol                 `json:"instrument"`
	Metadata              map[string]interface{} `json:"metadata"`
	ClientOrderIdentifier string                 `json:"clientOrderIdentifier"`
}

// OrderResponse is a response to an order query or placement
type OrderResponse struct {
	*OrderRequest
	ID                         string    `json:"id"`
	WalletID                   string    `json:"walletId"`
	AmountFilled               float64   `json:"amountFilled,string"`
	VolumeWeightedAveragePrice float64   `json:"volumeWeightedAveragePrice,string"`
	CreatedTime                time.Time `json:"createdTime"`
	Status                     string    `json:"status"`
}

// NewOrder makes a new order.
func NewOrder(conf Config, walletID string, order *OrderRequest) (*OrderResponse, error) {
	if err := validateSymbol(order.Instrument); err != nil {
		return nil, err
	}
	if err := validateOrderSide(order.Side); err != nil {
		return nil, err
	}
	if err := validateOrderType(order.Type); err != nil {
		return nil, err
	}
	if err := validateCurrency(order.Currency); err != nil {
		return nil, err
	}
	resp, err := doReq(conf, "POST", fmt.Sprintf("/wallets/%s/order", walletID), true, order)
	if err != nil {
		return nil, errors.Wrap(err, "requesting order creation")
	}
	dec := json.NewDecoder(resp)
	r := &OrderResponse{
		OrderRequest: &OrderRequest{},
	}
	if err := dec.Decode(r); err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}
	return r, nil
}

func validateCurrency (c Currency) error {
	if c != CurrencyBTC {
		return errors.Errorf("invalid currency %s", c)
	}
	return nil
}
func validateOrderType (c OrderType) error {
	if c != OrderTypeLimit {
		return errors.Errorf("invalid order type %s", c)
	}
	return nil
}

func validateOrderSide (c OrderSide) error {
	if c != OrderSideBuy && c != OrderSideSell {
		return errors.Errorf("invalid side %s", c)
	}
	return nil
}