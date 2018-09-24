package itbit

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/url"
)

const walletsPerPage = 50

type Wallet struct {
	ID string `json:"id"`
	UserID string `json:"userId"`
	Name string `json:"name"`
	Balances []*WalletBalance `json:"balances"`
}

type WalletBalance struct {
	Currency Currency `json:"currency"`
	Available float64 `json:"availableBalance,string"`
	Total float64 `json:"totalBalance,string"`
}

// GetAllWallets returns the wallets for the user ID in the given config
func GetAllWallets(conf Config) ([]*Wallet, error) {
	if conf.UserID == "" {
		return nil, errors.New("no UserID on config")
	}
	return readAllWallets(conf, 1)
}

// readAllWallets recursively gets all wallets for the user ID in the conf
func readAllWallets(conf Config, page int) ([]*Wallet, error) {
	params := url.Values{
		"userId": []string{conf.UserID},
		"page": []string{fmt.Sprintf("%d", page)},
		"perPage": []string{fmt.Sprintf("%d", walletsPerPage)},
	}

	resp, err := doReq(conf, "GET", fmt.Sprintf("/wallets?%s", params.Encode()), true, nil)
	if err != nil {
		return nil, errors.Wrapf(err,"getting page %d", page)
	}

	dec := json.NewDecoder(resp)
	var wallets []*Wallet
	if err := dec.Decode(wallets); err != nil {
		return nil, errors.Wrapf(err,"decoding page %d", page)
	}

	if len(wallets) >= walletsPerPage {
		more, err := readAllWallets(conf, page+1)
		if err != nil {
			return nil, errors.Wrapf(err, "recursing after page %d", page)
		}
		wallets = append(wallets, more...)
	}

	return wallets, nil
}