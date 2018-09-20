package itbit

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

// Symbol represents a trading pair.
type Symbol string

const (
	// SymbolXBTUSD is the USD/BTC trading pair.
	SymbolXBTUSD = "XBTUSD"
	// SymbolXBTSGD is the USD/SGD trading pair.
	SymbolXBTSGD = "XBTSGD"
	// SymbolXBTEUR is the USD/EUR trading pair.
	SymbolXBTEUR = "XBTEUR"
)

// Ticker represents market data.
type Ticker struct {
	Pair          Symbol    `json:"pair"`
	Bid           float64   `json:"bid,string"`
	BidAmt        float64   `json:"bidAmt,string"`
	Ask           float64   `json:"ask,string"`
	AskAmt        float64   `json:"askAmt,string"`
	LastPrice     float64   `json:"lastPrice,string"`
	LastAmt       float64   `json:"lastAmt,string"`
	Volume24h     float64   `json:"volume24h,string"`
	VolumeToday   float64   `json:"volumeToday,string"`
	High24h       float64   `json:"high24h,string"`
	Low24h        float64   `json:"low24h,string"`
	HighToday     float64   `json:"highToday,string"`
	LowToday      float64   `json:"lowToday,string"`
	OpenToday     float64   `json:"openToday,string"`
	VwapToday     float64   `json:"vwapToday,string"`
	Vwap24h       float64   `json:"vwap24h,string"`
	ServerTimeUTC time.Time `json:"serverTimeUTC,string"`
}

func (t *Ticker) String () string {
	return fmt.Sprintf("%0.2f - %0.2f %s", t.Bid, t.Ask, t.Pair)
}

// Order represents an order on the order book.
type Order []string

// Price returns the price from the order.
func (o Order) Price () (float64, error) {
	if len(o) < 1 {
		return 0, errors.New("invalid order")
	}
	p, err := strconv.ParseFloat(o[0], 64)
	if err != nil {
		return 0, errors.Wrap(err, "parsing price")
	}
	return p, nil
}

// Amount returns the price from the order.
func (o Order) Amount () (float64, error) {
	if len(o) < 2 {
		return 0, errors.New("invalid order")
	}
	a, err := strconv.ParseFloat(o[1], 64)
	if err != nil {
		return 0, errors.Wrap(err, "parsing amount")
	}
	return a, nil
}

// OrderBook represents open orders.
type OrderBook struct {
	Asks []Order `json:"asks"`
	Bids []Order `json:"bids"`
}

// GetTicker gets the market data for a symbol.
func GetTicker(conf Config, symbol Symbol) (*Ticker, error) {
	if err := validateSymbol(symbol); err != nil {
		return nil, err
	}
	body, err := doReq(conf, "GET", fmt.Sprintf("/markets/%s/ticker", symbol), false,nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting ticker from itBit")
	}
	defer body.Close()
	dec := json.NewDecoder(body)
	resp := &Ticker{}
	if err := dec.Decode(resp); err!= nil {
		return nil, errors.Wrap(err, "decoding JSON")
	}
	return resp, nil
}

// GetOrderBook gets a list of the open orders.
func GetOrderBook(conf Config, symbol Symbol) (*OrderBook, error) {
	if err := validateSymbol(symbol); err != nil {
		return nil, err
	}
	body, err := doReq(conf, "GET", fmt.Sprintf("/markets/%s/order_book", symbol), false,nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting ticker from itBit")
	}
	defer body.Close()
	dec := json.NewDecoder(body)
	resp := &OrderBook{}
	if err := dec.Decode(resp); err!= nil {
		return nil, errors.Wrap(err, "decoding JSON")
	}
	return resp, nil
}

func validateSymbol (symbol Symbol) error {
	if symbol != SymbolXBTUSD && symbol != SymbolXBTSGD && symbol != SymbolXBTEUR {
		return errors.Errorf("invalid symbol %s", symbol)
	}
	return nil
}