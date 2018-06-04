package cryptomkt

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"
)

// Client represents a connection to the CryptoMKT API
type Client struct {
	Key    string
	Secret string
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
type OrderType int

// OrderType possible values
const (
	Buy OrderType = iota
	Sell
)

func (ot OrderType) String() string {
	return orderTypesID[ot]
}

var orderTypesID = map[OrderType]string{
	Buy:  "buy",
	Sell: "sell",
}

var orderTypesName = map[string]OrderType{
	"buy":  Buy,
	"sell": Sell,
}

// MarshalJSON converts an OrderType to string
func (ot *OrderType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(orderTypesID[*ot])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON converts a string to an OrderType
func (ot *OrderType) UnmarshalJSON(b []byte) error {
	// unmarshal as string
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	// lookup value
	*ot = orderTypesName[s]
	return nil
}

// WalletType represents a CryptoMKT currency
type WalletType int

func (wt WalletType) String() string {
	return walletTypesID[wt]
}

// WalletType possible values
const (
	ARS WalletType = iota
	BRL
	CPL
	EUR
	ETH
	XLM
	BTC
)

var walletTypesID = map[WalletType]string{
	ARS: "ARS",
	BRL: "BRL",
	CPL: "CPL",
	EUR: "EUR",
	ETH: "ETH",
	XLM: "XLM",
	BTC: "BTC",
}

var walletTypesName = map[string]WalletType{
	"ARS": ARS,
	"BRL": BRL,
	"CPL": CPL,
	"EUR": EUR,
	"ETH": ETH,
	"XLM": XLM,
	"BTC": BTC,
}

// MarshalJSON converts a WalletType to string
func (wt *WalletType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(walletTypesID[*wt])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON converts a string to a WalletType
func (wt *WalletType) UnmarshalJSON(b []byte) error {
	// unmarshal as string
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	// lookup value
	*wt = walletTypesName[s]
	return nil
}

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
type Market int

// OrderType possible values
const (
	ETHARS Market = iota
	ETHEUR
	ETHBRL
	XLMARS
	XLMEUR
	XLMBRL
	BTCARS
	BTCEUR
	BTCBRL
	ETHCLP
	XLMCLP
	BTCCLP
)

func (m Market) String() string {
	return marketsID[m]
}

var marketsID = map[Market]string{
	ETHARS: "ETHARS",
	ETHEUR: "ETHEUR",
	ETHBRL: "ETHBRL",
	XLMARS: "XLMARS",
	XLMEUR: "XLMEUR",
	XLMBRL: "XLMBRL",
	BTCARS: "BTCARS",
	BTCEUR: "BTCEUR",
	BTCBRL: "BTCBRL",
	ETHCLP: "ETHCLP",
	XLMCLP: "XLMCLP",
	BTCCLP: "BTCCLP",
}

var marketsName = map[string]Market{
	"ETHARS": ETHARS,
	"ETHEUR": ETHEUR,
	"ETHBRL": ETHBRL,
	"XLMARS": XLMARS,
	"XLMEUR": XLMEUR,
	"XLMBRL": XLMBRL,
	"BTCARS": BTCARS,
	"BTCEUR": BTCEUR,
	"BTCBRL": BTCBRL,
	"ETHCLP": ETHCLP,
	"XLMCLP": XLMCLP,
	"BTCCLP": BTCCLP,
}

// MarshalJSON converts a Market to string
func (m *Market) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(marketsID[*m])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON converts a string to a Market
func (m *Market) UnmarshalJSON(b []byte) error {
	// unmarshal as string
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	// lookup value
	*m = marketsName[s]
	return nil
}

// MarketResponse is the response of the Markets endpoint
type MarketResponse struct {
	Status string
	Data   []Market
}

// Ticker represents a Ticker in the CryptoMKT API
type Ticker struct {
	High      float64 `json:",string"`
	Volume    float64 `json:",string"`
	Low       float64 `json:",string"`
	Ask       float64 `json:",string"`
	Timestamp Time
	Bid       float64 `json:",string"`
	LastPrice float64 `json:"last_price,string"`
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
	Price     float64 `json:",string"`
	Amount    float64 `json:",string"`
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
	Price       float64 `json:",string"`
	Amount      float64 `json:",string"`
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
	Original  float64 `json:",string"`
	Remaining float64 `json:",string,omitempty"`
	Executed  float64 `json:",string,omitempty"`
}

// Order is the representation of an Order in the CryptoMKT API
type Order struct {
	Status            string
	CreatedAt         Time `json:"created_at"`
	Amount            Amount
	ExecutionPrice    float64 `json:"execution_price,string,omitempty"`
	AvgExecutionPrice float64 `json:"avg_execution_price,string,omitempty"`
	Price             float64 `json:",string"`
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

// OrderResponse is the response of the endpoints Create, Status, and Cancel
type OrderResponse struct {
	Status string
	Data   Order
}

// Wallet represents a Wallet in the CryptoMKT API
type Wallet struct {
	Available float64 `json:",string"`
	Wallet    WalletType
	Balance   float64 `json:",string"`
}

// BalanceResponse is the response of the Balance endpoint
type BalanceResponse struct {
	Status string
	Data   []Wallet
}
