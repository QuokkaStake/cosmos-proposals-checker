package pagerduty

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"main/pkg/events"
	"main/pkg/report/entry"
	"main/pkg/types"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/rs/zerolog"
)

type Reporter struct {
	PagerDutyURL string
	APIKey       string
	Logger       zerolog.Logger
	Tracer       trace.Tracer
}

type AlertPayload struct {
	Summary       string            `json:"summary"`
	Timestamp     string            `json:"timestamp"`
	Severity      string            `json:"severity"`
	Source        string            `json:"source"`
	CustomDetails map[string]string `json:"custom_details"`
}

type Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

type Alert struct {
	Payload     AlertPayload `json:"payload"`
	RoutingKey  string       `json:"routing_key"`
	EventAction string       `json:"event_action"`
	DedupKey    string       `json:"dedup_key"`
	Client      string       `json:"client"`
	Links       []Link       `json:"links"`
	ClientURL   string       `json:"client_url"`
}

type Response struct {
	Status  string
	Message string
}

func (r *Reporter) NewAlertFromReportEntry(eventRaw entry.ReportEntry) (Alert, error) {
	event, ok := eventRaw.(entry.ReportEntryNotError)
	if !ok {
		return Alert{}, errors.New("error converting alert entry")
	}

	eventAction := "trigger"
	if _, ok := event.(events.VotedEvent); ok {
		eventAction = "resolve"
	}

	dedupKey := fmt.Sprintf(
		"cosmos-proposals-checker alert chain=%s proposal=%s wallet=%s",
		event.GetChain().Name,
		event.GetProposal().ID,
		event.GetWallet().AddressOrAlias(),
	)

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
	}

	links := []Link{}
	explorerLinks := event.GetChain().GetExplorerProposalsLinks(event.GetProposal().ID)
	for _, link := range explorerLinks {
		links = append(links, Link{
			Href: link.Href,
			Text: link.Name,
		})
	}

	return Alert{
		Payload: AlertPayload{
			Summary: fmt.Sprintf(
				"Wallet %s hasn't voted on proposal %s on %s: %s",
				event.GetWallet().AddressOrAlias(),
				event.GetProposal().ID,
				event.GetChain().GetName(),
				event.GetProposal().Title,
			),
			Timestamp: time.Now().Format(time.RFC3339),
			Severity:  "error",
			Source:    hostname,
			CustomDetails: map[string]string{
				"Wallet":               event.GetWallet().AddressOrAlias(),
				"Chain":                event.GetChain().GetName(),
				"Proposal ID":          event.GetProposal().ID,
				"Proposal title":       event.GetProposal().Title,
				"Proposal description": event.GetProposal().Description,
			},
		},
		Links:       links,
		RoutingKey:  r.APIKey,
		EventAction: eventAction,
		DedupKey:    dedupKey,
		Client:      "cosmos-proposals-checker",
		ClientURL:   "https://github.com/QuokkaStake/cosmos-proposals-checker",
	}, nil
}

func NewPagerDutyReporter(
	config types.PagerDutyConfig,
	logger *zerolog.Logger,
	tracer trace.Tracer,
) Reporter {
	return Reporter{
		PagerDutyURL: config.PagerDutyURL,
		APIKey:       config.APIKey,
		Logger:       logger.With().Str("component", "pagerduty_reporter").Logger(),
		Tracer:       tracer,
	}
}

func (r Reporter) Init() error {
	return nil
}

func (r Reporter) Enabled() bool {
	return r.APIKey != ""
}

func (r Reporter) Name() string {
	return "pagerduty-reporter"
}

func (r Reporter) SendReportEntry(reportEntry entry.ReportEntry, ctx context.Context) error {
	_, span := r.Tracer.Start(ctx, "Sending Pagerduty report entry")
	defer span.End()

	if !reportEntry.IsAlert() {
		return nil
	}

	alert, alertCreateErr := r.NewAlertFromReportEntry(reportEntry)
	if alertCreateErr != nil {
		return alertCreateErr
	}

	return r.SendAlert(alert)
}

func (r Reporter) SendAlert(alert Alert) error {
	var response Response
	err := r.DoRequest(r.PagerDutyURL+"/v2/enqueue", alert, &response)
	if err != nil {
		return err
	}

	if response.Status != "success" {
		return errors.New("expected 'success' status, got '" + response.Status + "'. Error: " + response.Message)
	}

	return nil
}

func (r Reporter) DoRequest(url string, body interface{}, target interface{}) error {
	client := &http.Client{Timeout: 10 * 1000000000}
	start := time.Now()

	jsonBody, err := json.Marshal(body)
	if err != nil {
		r.Logger.Err(err).Msg("Error marshalling request body")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer((jsonBody)))
	if err != nil {
		r.Logger.Err(err).Msg("Error instantiating request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	r.Logger.Debug().Str("url", url).Msg("Doing a query...")

	res, err := client.Do(req)
	if err != nil {
		r.Logger.Warn().Str("url", url).Err(err).Msg("Query failed")
		return err
	}
	defer res.Body.Close()

	r.Logger.Debug().Str("url", url).Dur("duration", time.Since(start)).Msg("Query is finished")
	return json.NewDecoder(res.Body).Decode(target)
}
