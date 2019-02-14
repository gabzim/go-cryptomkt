package cryptomkt

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const apiURL = "https://api.cryptomkt.com/"
const version = "v1/"
const limit = 100

// MarketAssetMapping simplifies the obtention of the asset of a market
var MarketAssetMapping = map[Market]WalletType{
	ETHARS: ETH,
	ETHBRL: ETH,
	ETHCLP: ETH,
	ETHEUR: ETH,
	XLMARS: XLM,
	XLMBRL: XLM,
	XLMCLP: XLM,
	XLMEUR: XLM,
	BTCARS: BTC,
	BTCBRL: BTC,
	BTCCLP: BTC,
	BTCEUR: BTC,
	EOSARS: EOS,
	EOSBRL: EOS,
	EOSCLP: EOS,
	EOSEUR: EOS,
}

// MarketCurrencyMapping simplifies the obtention of the currency of a market
var MarketCurrencyMapping = map[Market]WalletType{
	ETHARS: ARS,
	ETHBRL: BRL,
	ETHCLP: CLP,
	ETHEUR: EUR,
	XLMARS: ARS,
	XLMBRL: BRL,
	XLMCLP: CLP,
	XLMEUR: EUR,
	BTCARS: ARS,
	BTCBRL: BRL,
	BTCCLP: CLP,
	BTCEUR: EUR,
	EOSARS: ARS,
	EOSBRL: BRL,
	EOSCLP: CLP,
	EOSEUR: EUR,
}

func NewClient(key, secret string, timeout time.Duration) *Client {
	return &Client{
		key:    key,
		secret: secret,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c Client) formURL(initialURL string, paramsMap map[string]string) (string, error) {
	baseURL, err := url.Parse(initialURL)
	if err != nil {
		return "", fmt.Errorf("Parse error: %s", err)
	}

	params := url.Values{}
	for k, v := range paramsMap {
		params.Add(k, v)
	}

	baseURL.RawQuery = params.Encode()
	return baseURL.String(), nil
}

func (c Client) formHeaders(req *http.Request, path string, data url.Values) {
	req.Header.Add("X-MKT-APIKEY", c.key)

	t := time.Now().Unix()
	body := strconv.FormatInt(t, 10) + "/" + version + path
	if data != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		for _, values := range data {
			for _, value := range values {
				body += value
			}
		}
	}

	h := hmac.New(sha512.New384, []byte(c.secret))
	h.Write([]byte(body))

	req.Header.Add("X-MKT-SIGNATURE", hex.EncodeToString(h.Sum(nil)))
	req.Header.Add("X-MKT-TIMESTAMP", strconv.FormatInt(t, 10))
}

func (c Client) get(path string, params map[string]string, auth bool) (*http.Response, error) {
	var err error

	// First, create the request url with the params map
	requestURL, err := c.formURL(apiURL+version+path, params)
	if err != nil {
		return nil, err
	}

	// Then, create the http Client and set the headers if needed
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %s", err)
	}

	if auth == true {
		c.formHeaders(req, path, nil)
	}

	// Make the request
	return c.client.Do(req)
}

func (c Client) post(path string, data map[string]string) (*http.Response, error) {
	var err error

	// First, create the request url with the params map
	requestURL, err := c.formURL(apiURL+version+path, nil)
	if err != nil {
		return nil, err
	}

	payload := url.Values{}
	for k, v := range data {
		payload.Add(k, v)
	}

	// Then, create the http Client and set the headers
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Request failed: %s", err)
	}
	c.formHeaders(req, path, payload)

	// Make the request
	return c.client.Do(req)
}

// Markets returns a *MarketResponse with an array of Markets
func (c Client) Markets() (*MarketResponse, error) {
	path := "market"

	res, err := c.get(path, nil, false)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result MarketResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding: %s", err)
	}
	return &result, nil
}

// Ticker returns a *TickerResponse with the status of a Market
func (c Client) Ticker(market Market) (*TickerResponse, error) {
	params := map[string]string{"market": string(market)}
	path := "ticker"

	res, err := c.get(path, params, false)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result TickerResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding: %s", err)
	}
	return &result, nil
}

// Book returns an *OrderBookResponse with an array of OrderBookOrders
func (c Client) Book(market Market, ot OrderType, page int) (*OrderBookResponse, error) {
	params := map[string]string{"market": string(market), "type": string(ot), "page": strconv.Itoa(page), "limit": strconv.Itoa(limit)}
	path := "book"

	res, err := c.get(path, params, false)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result OrderBookResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding: %s", err)
	}
	return &result, nil
}

// BuyBook returns an *OrderBookResponse with an array of BUY OrderBookOrders
func (c Client) BuyBook(market Market, page int) (*OrderBookResponse, error) {
	return c.Book(market, BUY, page)
}

// SellBook returns an *OrderBookResponse with an array of SELL OrderBookOrders
func (c Client) SellBook(market Market, page int) (*OrderBookResponse, error) {
	return c.Book(market, SELL, page)
}

// Trades returns a *TradesResponse with an array of Trades
func (c Client) Trades(market Market, start string, end string, page int) (*TradesResponse, error) {
	params := map[string]string{"market": string(market), "start": start, "end": end, "page": strconv.Itoa(page), "limit": strconv.Itoa(limit)}
	path := "trades"

	res, err := c.get(path, params, false)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result TradesResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding: %s", err)
	}
	return &result, nil
}

// ActiveOrders returns an *OrdersResponse with an array of ActiveOrders
func (c Client) ActiveOrders(market Market, page int) (*OrdersResponse, error) {
	params := map[string]string{"market": string(market), "page": strconv.Itoa(page), "limit": strconv.Itoa(limit)}
	path := "orders/active"

	res, err := c.get(path, params, true)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result OrdersResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding: %s", err)
	}
	return &result, nil
}

// ExecutedOrders returns an *OrdersResponse with an array of ExecutedOrders
func (c Client) ExecutedOrders(market Market, page int) (*OrdersResponse, error) {
	params := map[string]string{"market": string(market), "page": strconv.Itoa(page), "limit": strconv.Itoa(limit)}
	path := "orders/executed"

	res, err := c.get(path, params, true)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result OrdersResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding: %s", err)
	}
	return &result, nil
}

// CreateOrder creates an Order and returns an *OrderResponse with the created Order
func (c Client) CreateOrder(market Market, amount float64, price float64, ot OrderType) (*OrderResponse, error) {
	data := map[string]string{
		"amount": strconv.FormatFloat(amount, 'f', 4, 64),
		"market": string(market),
		"price":  strconv.FormatFloat(price, 'f', 4, 64),
		"type":   string(ot),
	}
	path := "orders/create"

	res, err := c.post(path, data)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result OrderResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding: %s", err)
	}
	return &result, nil
}

// OrderStatus returns an *OrderResponse with the status of an Order
func (c Client) OrderStatus(ID string) (*OrderResponse, error) {
	var params = map[string]string{"id": ID}
	path := "orders/status"

	res, err := c.get(path, params, true)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result OrderResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding: %s", err)
	}
	return &result, nil
}

// CancelOrder cancels an Order and returns an *OrderResponse with the status of the Order
func (c Client) CancelOrder(ID string) (*OrderResponse, error) {
	data := map[string]string{"id": ID}
	path := "orders/cancel"

	res, err := c.post(path, data)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result OrderResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding: %s", err)
	}
	return &result, nil
}

// Balance returns a *BalanceResponse with the status of the Wallets
func (c Client) Balance() (*BalanceResponse, error) {
	path := "balance"

	res, err := c.get(path, nil, true)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result BalanceResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding: %s", err)
	}
	return &result, nil
}
// InstantGet Allows you to Find out how much you would receive/need if you were to sell/buy at market price your crypto.
func (c Client) InstantGet(market Market, ot OrderType, amount string) (*InstantGetResponse, error) {
	params := map[string]string{"market": string(market), "type": string(ot), "amount": amount}
	path := "orders/instant/get"
	res, err := c.get(path, params, true)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result InstantGetResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// InstantCreate Allows you to create an order that will be executed at market price.
func (c Client) InstantCreate(market Market, ot OrderType, amount string) (*InstantCreateResponse, error) {
	params := map[string]string{"market": string(market), "type": string(ot), "amount": amount}
	path := "orders/instant/create"
	res, err := c.post(path, params)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result InstantCreateResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
