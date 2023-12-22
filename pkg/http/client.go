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

func NewClient(chainName string, hosts []string, logger zerolog.Logger) *Client {
	return &Client{
		Hosts: hosts,
		Logger: logger.With().
			Str("component", "http").
			Str("chain", chainName).
			Logger(),
	}
}

func (client *Client) Get(url string, target interface{}) []types.NodeError {
	nodeErrors := make([]types.NodeError, len(client.Hosts))

	for index, lcd := range client.Hosts {
		fullURL := lcd + url
		client.Logger.Trace().Str("url", fullURL).Msg("Trying making request to LCD")

		err := client.GetFull(
			fullURL,
			target,
		)

		if err == nil {
			return nil
		}

		client.Logger.Warn().Str("url", fullURL).Err(err).Msg("LCD request failed")
		nodeErrors[index] = types.NodeError{
			Node:  lcd,
			Error: types.NewJSONError(err),
		}
	}

	client.Logger.Warn().Str("url", url).Msg("All LCD requests failed")
	return nodeErrors
}

func (client *Client) GetFull(url string, target interface{}) error {
	httpClient := &http.Client{Timeout: 300 * time.Second}
	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "cosmos-proposals-checker")

	client.Logger.Debug().Str("url", url).Msg("Doing a query...")

	res, err := httpClient.Do(req)
	if err != nil {
		client.Logger.Warn().Str("url", url).Err(err).Msg("Query failed")
		return err
	}
	defer res.Body.Close()

	client.Logger.Debug().Str("url", url).Dur("duration", time.Since(start)).Msg("Query is finished")

	return json.NewDecoder(res.Body).Decode(target)
}
