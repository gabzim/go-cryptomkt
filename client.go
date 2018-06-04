package cryptomkt

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const apiURL = "https://api.cryptomkt.com/"
const version = "v1/"
const limit = 100
const maxRetries = 5

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
	req.Header.Add("X-MKT-APIKEY", c.Key)

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

	h := hmac.New(sha512.New384, []byte(c.Secret))
	h.Write([]byte(body))

	req.Header.Add("X-MKT-SIGNATURE", hex.EncodeToString(h.Sum(nil)))
	req.Header.Add("X-MKT-TIMESTAMP", strconv.FormatInt(t, 10))
}

func (c Client) get(path string, params map[string]string, auth bool) ([]byte, error) {
	var err error

	for i := 0; i < maxRetries; i++ {
		// First, create the request url with the params map
		requestURL, err := c.formURL(apiURL+version+path, params)
		if err != nil {
			return nil, err
		}

		// Then, create the http Client and set the headers if needed
		client := &http.Client{}
		req, err := http.NewRequest("GET", requestURL, nil)
		if err != nil {
			return nil, fmt.Errorf("Request failed: %s", err)
		}

		if auth == true {
			c.formHeaders(req, path, nil)
		}

		// Make the request
		resp, err := client.Do(req)
		if err != nil {
			err = fmt.Errorf("Request failed: %s", err)
			time.Sleep(2 * time.Second)
			continue
		}
		defer resp.Body.Close()

		// Test the response code
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("Request failed: %s", resp.Status)
			time.Sleep(2 * time.Second)
			continue
		}

		// Extract the body before closing and return it
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Body reading failed: %s", err)
		}

		return body, nil
	}
	return nil, err
}

func (c Client) post(path string, data map[string]string) ([]byte, error) {
	var err error

	for i := 0; i < maxRetries; i++ {
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
		client := &http.Client{}
		req, err := http.NewRequest("POST", requestURL, strings.NewReader(payload.Encode()))
		if err != nil {
			return nil, fmt.Errorf("Request failed: %s", err)
		}
		c.formHeaders(req, path, payload)

		// Make the request
		resp, err := client.Do(req)
		if err != nil {
			err = fmt.Errorf("Request failed: %s", err)
			time.Sleep(2 * time.Second)
			continue
		}
		defer resp.Body.Close()

		// Test the response code
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("Request failed: %s", resp.Status)
			time.Sleep(2 * time.Second)
			continue
		}

		// Extract the body before closing and return it
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Body reading failed: %s", err)
		}

		return body, nil
	}
	return nil, err
}

// Markets returns a *MarketResponse with an array of Markets
func (c Client) Markets() (*MarketResponse, error) {
	path := "market"

	body, err := c.get(path, nil, false)
	if err != nil {
		return nil, err
	}

	var result MarketResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Error decoding: %s", err)
	}
	return &result, nil
}

// Ticker returns a *TickerResponse with the status of a Market
func (c Client) Ticker(market Market) (*TickerResponse, error) {
	params := map[string]string{"market": market.String()}
	path := "ticker"

	body, err := c.get(path, params, false)
	if err != nil {
		return nil, err
	}

	var result TickerResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Error decoding: %s", err)
	}
	return &result, nil
}

// Book returns an *OrderBookResponse with an array of OrderBookOrders
func (c Client) Book(market Market, ot OrderType, page int) (*OrderBookResponse, error) {
	params := map[string]string{"market": market.String(), "type": ot.String(), "page": strconv.Itoa(page), "limit": strconv.Itoa(limit)}
	path := "book"

	body, err := c.get(path, params, false)
	if err != nil {
		return nil, err
	}

	var result OrderBookResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Error decoding: %s", err)
	}
	return &result, nil
}

// BuyBook returns an *OrderBookResponse with an array of Buy OrderBookOrders
func (c Client) BuyBook(market Market, page int) (*OrderBookResponse, error) {
	return c.Book(market, Buy, page)
}

// SellBook returns an *OrderBookResponse with an array of Sell OrderBookOrders
func (c Client) SellBook(market Market, page int) (*OrderBookResponse, error) {
	return c.Book(market, Sell, page)
}

// Trades returns a *TradesResponse with an array of Trades
func (c Client) Trades(market Market, start string, end string, page int) (*TradesResponse, error) {
	params := map[string]string{"market": market.String(), "start": start, "end": end, "page": strconv.Itoa(page), "limit": strconv.Itoa(limit)}
	path := "trades"

	body, err := c.get(path, params, false)
	if err != nil {
		return nil, err
	}

	var result TradesResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Error decoding: %s", err)
	}
	return &result, nil
}

// ActiveOrders returns an *OrdersResponse with an array of ActiveOrders
func (c Client) ActiveOrders(market Market, page int) (*OrdersResponse, error) {
	params := map[string]string{"market": market.String(), "page": strconv.Itoa(page), "limit": strconv.Itoa(limit)}
	path := "orders/active"

	body, err := c.get(path, params, true)
	if err != nil {
		return nil, err
	}

	var result OrdersResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Error decoding: %s", err)
	}
	return &result, nil
}

// ExecutedOrders returns an *OrdersResponse with an array of ExecutedOrders
func (c Client) ExecutedOrders(market Market, page int) (*OrdersResponse, error) {
	params := map[string]string{"market": market.String(), "page": strconv.Itoa(page), "limit": strconv.Itoa(limit)}
	path := "orders/executed"

	body, err := c.get(path, params, true)
	if err != nil {
		return nil, err
	}

	var result OrdersResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Error decoding: %s", err)
	}
	return &result, nil
}

// Create creates an Order and returns an *OrderResponse with the created Order
func (c Client) Create(market Market, amount float64, price float64, ot OrderType) (*OrderResponse, error) {
	data := map[string]string{
		"amount": strconv.FormatFloat(amount, 'f', 4, 64),
		"market": market.String(),
		"price":  strconv.FormatFloat(price, 'f', 4, 64),
		"type":   ot.String(),
	}
	path := "orders/create"

	body, err := c.post(path, data)
	if err != nil {
		return nil, err
	}

	var result OrderResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Error decoding: %s", err)
	}
	return &result, nil
}

// Status returns an *OrderResponse with the status of an Order
func (c Client) Status(ID string) (*OrderResponse, error) {
	var params = map[string]string{"id": ID}
	path := "orders/status"

	body, err := c.get(path, params, true)
	if err != nil {
		return nil, err
	}

	var result OrderResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Error decoding: %s", err)
	}
	return &result, nil
}

// Cancel cancels an Order and returns an *OrderResponse with the status of the Order
func (c Client) Cancel(ID string) (*OrderResponse, error) {
	data := map[string]string{"id": ID}
	path := "orders/cancel"

	body, err := c.post(path, data)
	if err != nil {
		return nil, err
	}

	var result OrderResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Error decoding: %s", err)
	}
	return &result, nil
}

// Balance returns a *BalanceResponse with the status of the Wallets
func (c Client) Balance() (*BalanceResponse, error) {
	path := "balance"

	body, err := c.get(path, nil, true)
	if err != nil {
		return nil, err
	}

	var result BalanceResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Error decoding: %s", err)
	}
	return &result, nil
}
