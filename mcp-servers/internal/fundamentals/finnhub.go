package fundamentals

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

type Finnhub struct {
	http *http.Client
	apiKey string
	baseURL string
}

var _ FundamentalsProvider = (*Finnhub)(nil)

//  I have named it such since go doesnot have classes hence no constructors, for me to understand that the sole purpose of this function is to expose the attr of the struct, since all are private
func FinnhubInit(httpClient *http.Client, apiKey string) *Finnhub {
	return &Finnhub{http: httpClient, apiKey: apiKey, baseURL: FINNHUB_BASE_URL}
}

type finnhubProfile struct {
	Country              string  `json:"country"`
	Currency             string  `json:"currency"`
	Exchange             string  `json:"exchange"`
	FinnhubIndustry      string  `json:"finnhubIndustry"`
	IPO                  string  `json:"ipo"`
	Logo                 string  `json:"logo"`
	MarketCapitalization float64 `json:"marketCapitalization"`
	Name                 string  `json:"name"`
	ShareOutstanding     float64 `json:"shareOutstanding"`
	Ticker               string  `json:"ticker"`
	WebURL               string  `json:"weburl"`
}

type finnhubMetric struct {
	Metric map[string]any `json:"metric"`
	Symbol string         `json:"symbol"`
}

func (f *Finnhub) Fundamentals(ctx context.Context, ticker string) (Fundamentals, error) {
	
	var (
		profile finnhubProfile
		metric finnhubMetric
		profileErr error
		metricErr error
		wg sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		profileErr = f.get(ctx, "/stock/profile2", url.Values{"symbol": {ticker}}, &profile)	
	}()

	go func() {
		defer wg.Done()
		metricErr = f.get(ctx, "/stock/metric", url.Values{"symbol": {ticker}, "metric": {"all"}}, &metric)
	}()

	wg.Wait()

	if profileErr != nil {
		return Fundamentals{}, profileErr
	}

	if metricErr != nil {
		return Fundamentals{}, metricErr
	}
	
	if profile.Name == "" && len(metric.Metric) == 0 {
		return Fundamentals{}, fmt.Errorf("%w: %s", ErrNotFound, ticker)
	}

	return Fundamentals{
		Profile: Profile{
			Ticker:            profile.Ticker,
			Name:              profile.Name,
			Exchange:          profile.Exchange,
			Country:           profile.Country,
			Currency:          profile.Currency,
			Industry:          profile.FinnhubIndustry,
			IPO:               profile.IPO,
			MarketCap:         profile.MarketCapitalization,
			SharesOutstanding: profile.ShareOutstanding,
			WebURL:            profile.WebURL,
			Logo:              profile.Logo,
		},
		Metrics: Metrics{
			High52Week:     num(metric.Metric, "52WeekHigh"),
			Low52Week:      num(metric.Metric, "52WeekLow"),
			Low52WeekDate:  str(metric.Metric, "52WeekLowDate"),
			Beta:           num(metric.Metric, "beta"),
			AvgVolume10Day: num(metric.Metric, "10DayAverageTradingVolume"),
			All:            metric.Metric,
		},
		Source: "finnhub",
	}, nil
}

func (f *Finnhub) get(ctx context.Context, path string, query url.Values, out any) error {
	url := f.baseURL + path + "?" + query.Encode() 
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return fmt.Errorf("finnhub %s: build request: %w", path, err)
	}

	req.Header.Set("X-Finnhub-Token", f.apiKey)

	resp, err := f.http.Do(req)

	if err != nil {
		return fmt.Errorf("finnhub %s: %w", path, err)
	}

	defer resp.Body.Close()

	// if rate limit
	if resp.StatusCode == http.StatusTooManyRequests {
		return fmt.Errorf("finnhub %s: %w", path, ErrRateLimited)
	}
	
	// internal errors, client errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("finnhub %s: status %d: %s", path, resp.StatusCode, body)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("finnhub %s: decode: %w", path, err)
	}
	return nil
}

func num(m map[string]any, key string) float64 {
	v, _ := m[key].(float64)
	return v
}
func str(m map[string]any, key string) string {
	v, _ := m[key].(string)
	return v
}