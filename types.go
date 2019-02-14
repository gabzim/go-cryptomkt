package cryptomkt

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// Client represents a connection to the CryptoMKT API
type Client struct {
	key    string
	secret string
	client *http.Client
}

// FlexInt is a fix for a wrong return on the API, where "null" is returned instead of null
type FlexInt int

// UnmarshalJSON parses inconsistent int || "null" value from CryptoMKT
func (fi *FlexInt) UnmarshalJSON(b []byte) error {
	if b[0] != '"' {
		return json.Unmarshal(b, (*int)(fi))
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	*fi = FlexInt(i)
	return nil
}

// OrderType represents a buy or sell signal
type OrderType string

// OrderType possible values
const (
	BUY  OrderType = "buy"
	SELL OrderType = "sell"
)

// WalletType represents a CryptoMKT currency
type WalletType string

// WalletType possible values
const (
	ARS WalletType = "ARS"
	BRL WalletType = "BRL"
	CLP WalletType = "CLP"
	EUR WalletType = "EUR"
	ETH WalletType = "ETH"
	XLM WalletType = "XLM"
	BTC WalletType = "BTC"
	EOS WalletType = "EOS"
)

// Time represents the custom time format from CryptoMKT
type Time struct {
	time.Time
}

// UnmarshalJSON parses custom date format from CryptoMKT
func (t *Time) UnmarshalJSON(b []byte) error {
	s := string(b)

	// Get rid of the quotes "" around the value.
	// A second option would be to include them
	// in the date format string instead, like so below:
	//   time.Parse(`"`+time.StampMicro+`"`, s)
	s = s[1 : len(s)-1]

	ts, err := time.Parse(time.StampMicro, s)
	if err != nil {
		ts, err = time.Parse("2006-01-02T15:04:05.999999", s)
	}
	t.Time = ts
	return err
}

// Pagination is the representation of the CryptoMKT pagination section of the API results
type Pagination struct {
	Previous FlexInt
	Limit    int
	Page     int
	Next     FlexInt
}

// Market represents a Market in the CryptoMKT API
type Market string

// OrderType possible values
const (
	ETHARS Market = "ETHARS"
	ETHEUR Market = "ETHEUR"
	ETHBRL Market = "ETHBRL"
	ETHCLP Market = "ETHCLP"
	XLMARS Market = "XLMARS"
	XLMEUR Market = "XLMEUR"
	XLMBRL Market = "XLMBRL"
	XLMCLP Market = "XLMCLP"
	BTCARS Market = "BTCARS"
	BTCEUR Market = "BTCEUR"
	BTCBRL Market = "BTCBRL"
	BTCCLP Market = "BTCCLP"
	EOSARS Market = "EOSARS"
	EOSEUR Market = "EOSEUR"
	EOSBRL Market = "EOSBRL"
	EOSCLP Market = "EOSCLP"
)

// MarketResponse is the response of the Markets endpoint
type MarketResponse struct {
	Status string
	Data   []Market
}

// Ticker represents a Ticker in the CryptoMKT API
type Ticker struct {
	High      string
	Volume    string
	Low       string
	Ask       string
	Timestamp Time
	Bid       string
	LastPrice string `json:"last_price"`
	Market    Market
}

// TickerResponse is the response of the Ticker endpoint
type TickerResponse struct {
	Status string
	Data   []Ticker
}

// OrderBookOrder represents an Order in the OrderBook
type OrderBookOrder struct {
	Timestamp Time
	Price     string
	Amount    string
}

// OrderBookResponse is the response of the Book endpoint
type OrderBookResponse struct {
	Status     string
	Pagination Pagination
	Data       []OrderBookOrder
}

// Trade represents a Trade operation in the CryptoMKT API
type Trade struct {
	MarketTaker OrderType `json:"market_taker"`
	Timestamp   Time
	Price       string
	Amount      string
	Market      Market
}

// TradesResponse is the response of the Trades endpoint
type TradesResponse struct {
	Status     string
	Pagination Pagination
	Data       []Trade
}

// Amount represents the different amounts that compose an Order
type Amount struct {
	Original  string
	Remaining string `json:",omitempty"`
	Executed  string `json:",omitempty"`
}

// Order is the representation of an Order in the CryptoMKT API
type Order struct {
	Status            string
	CreatedAt         Time `json:"created_at"`
	Amount            Amount
	ExecutionPrice    string `json:"execution_price,omitempty"`
	AvgExecutionPrice string `json:"avg_execution_price,omitempty"`
	Price             string
	Type              OrderType
	ID                string
	Market            Market
	UpdatedAt         Time `json:"updated_at"`
}

// OrdersResponse is the response of the endpoints ActiveOrders and ExecutedOrders
type OrdersResponse struct {
	Status     string
	Pagination Pagination
	Data       []Order
}

// OrderResponse is the response of the endpoints CreateOrder, OrderStatus, and CancelOrder
type OrderResponse struct {
	Status string
	Data   Order
}

// Wallet represents a Wallet in the CryptoMKT API
type Wallet struct {
	Available string
	Wallet    WalletType
	Balance   string
}

// BalanceResponse is the response of the Balance endpoint
type BalanceResponse struct {
	Status string
	Data   []Wallet
}

type InstantQuote struct {
	Obtained string
	Required string
}

type InstantGetResponse struct {
	Status string
	Data InstantQuote
}

type InstantCreateResponse struct {
	Status string
	Data string
}