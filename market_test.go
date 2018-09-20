package itbit

import "testing"

func TestGetTicker(t *testing.T) {
	conf := Config{}
	if _, err := GetTicker(conf, SymbolXBTUSD); err != nil {
		t.Fail()
		return
	}
}

func TestGetOrderBook(t *testing.T) {
	conf := Config{}
	if _, err := GetOrderBook(conf, SymbolXBTUSD); err != nil {
		t.Fail()
		return
	}
}

func TestOrder_Amount(t *testing.T) {
	o := Order{"0.1", "0.1"}
	if a, err := o.Amount(); err != nil {
		t.Fail()
		t.Log("parsing error", err)
		return
	} else if a != 0.1 {
		t.Fail()
		return
	}


	o = Order{"0.1", "foo"}
	if _, err := o.Amount(); err == nil {
		t.Fail()
		t.Log("failed to return error on invalid float")
		return
	}

	o = Order{"0.1"}
	if _, err := o.Amount(); err == nil {
		t.Fail()
		t.Log("failed to return error on missing amount")
		return
	}
}

func TestOrder_Price(t *testing.T) {
	o := Order{"0.1", "0.1"}
	if a, err := o.Price(); err != nil {
		t.Fail()
		t.Log("parsing error", err)
		return
	} else if a != 0.1 {
		t.Fail()
		return
	}


	o = Order{"foo", "0.1"}
	if _, err := o.Price(); err == nil {
		t.Fail()
		t.Log("failed to return error on invalid float")
		return
	}

	o = Order{}
	if _, err := o.Price(); err == nil {
		t.Fail()
		t.Log("failed to return error on missing price")
		return
	}
}
