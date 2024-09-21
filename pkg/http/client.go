package http

import (
	"context"
	"encoding/json"
	"main/pkg/types"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	Hosts  []string
	Logger zerolog.Logger
	Tracer trace.Tracer
}

func NewClient(
	chainName string,
	hosts []string,
	logger *zerolog.Logger,
	tracer trace.Tracer,
) *Client {
	return &Client{
		Hosts: hosts,
		Logger: logger.With().
			Str("component", "http").
			Str("chain", chainName).
			Logger(),
		Tracer: tracer,
	}
}

func (client *Client) Get(
	url string,
	target interface{},
	ctx context.Context,
) []types.NodeError {
	errs, _ := client.GetWithPredicate(url, target, types.HTTPPredicateAlwaysPass(), ctx)
	return errs
}

func (client *Client) GetWithPredicate(
	url string,
	target interface{},
	predicate types.HTTPPredicate,
	ctx context.Context,
) ([]types.NodeError, http.Header) {
	childCtx, span := client.Tracer.Start(ctx, "HTTP request on all nodes")
	defer span.End()

	nodeErrors := make([]types.NodeError, len(client.Hosts))

	for index, lcd := range client.Hosts {
		fullURL := lcd + url
		client.Logger.Trace().Str("url", fullURL).Msg("Trying making request to LCD")

		header, err := client.GetFull(
			fullURL,
			target,
			predicate,
			childCtx,
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
	ctx context.Context,
) (http.Header, error) {
	childCtx, span := client.Tracer.Start(ctx, "HTTP request")
	defer span.End()

	var transport http.RoundTripper

	transportRaw, ok := http.DefaultTransport.(*http.Transport)
	if ok {
		transport = transportRaw.Clone()
	} else {
		transport = http.DefaultTransport
	}

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: otelhttp.NewTransport(transport),
	}

	start := time.Now()

	req, err := http.NewRequestWithContext(childCtx, http.MethodGet, url, nil)
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
