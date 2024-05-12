package http

import (
	"encoding/json"
	"main/pkg/types"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Client struct {
	Hosts  []string
	Logger zerolog.Logger
}

func NewClient(chainName string, hosts []string, logger *zerolog.Logger) *Client {
	return &Client{
		Hosts: hosts,
		Logger: logger.With().
			Str("component", "http").
			Str("chain", chainName).
			Logger(),
	}
}

func (client *Client) Get(
	url string,
	target interface{},
) []types.NodeError {
	errs, _ := client.GetWithPredicate(url, target, types.HTTPPredicateAlwaysPass())
	return errs
}

func (client *Client) GetWithPredicate(
	url string,
	target interface{},
	predicate types.HTTPPredicate,
) ([]types.NodeError, http.Header) {
	nodeErrors := make([]types.NodeError, len(client.Hosts))

	for index, lcd := range client.Hosts {
		fullURL := lcd + url
		client.Logger.Trace().Str("url", fullURL).Msg("Trying making request to LCD")

		header, err := client.GetFull(
			fullURL,
			target,
			predicate,
		)

		if err == nil {
			return nil, header
		}

		client.Logger.Warn().Str("url", fullURL).Err(err).Msg("LCD request failed")
		nodeErrors[index] = types.NodeError{
			Node:  lcd,
			Error: types.NewJSONError(err),
		}
	}

	client.Logger.Warn().Str("url", url).Msg("All LCD requests failed")
	return nodeErrors, nil
}

func (client *Client) GetFull(
	url string,
	target interface{},
	predicate types.HTTPPredicate,
) (http.Header, error) {
	httpClient := &http.Client{Timeout: 300 * time.Second}
	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "cosmos-proposals-checker")

	client.Logger.Debug().Str("url", url).Msg("Doing a query...")

	res, err := httpClient.Do(req)
	if err != nil {
		client.Logger.Warn().Str("url", url).Err(err).Msg("Query failed")
		return nil, err
	}
	defer res.Body.Close()

	if err := predicate(res); err != nil {
		return nil, err
	}

	client.Logger.Debug().Str("url", url).Dur("duration", time.Since(start)).Msg("Query is finished")

	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}

	return res.Header, nil
}
