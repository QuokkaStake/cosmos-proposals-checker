package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type PagerDutyReporter struct {
	PagerDutyURL string
	APIKey       string
	Logger       zerolog.Logger
}

type PagerDutyAlertPayload struct {
	Summary       string            `json:"summary"`
	Timestamp     string            `json:"timestamp"`
	Severity      string            `json:"severity"`
	Source        string            `json:"source"`
	CustomDetails map[string]string `json:"custom_details"`
}

type PagerDutyLink struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

type PagerDutyAlert struct {
	Payload     PagerDutyAlertPayload `json:"payload"`
	RoutingKey  string                `json:"routing_key"`
	EventAction string                `json:"event_action"`
	DedupKey    string                `json:"dedup_key"`
	Client      string                `json:"client"`
	Links       []PagerDutyLink       `json:"links"`
	ClientURL   string                `json:"client_url"`
}

type PagerDutyResponse struct {
	Status  string
	Message string
}

func (r *PagerDutyReporter) NewPagerDutyAlertFromReportEntry(e ReportEntry) PagerDutyAlert {
	eventAction := "trigger"
	if e.HasVoted() {
		eventAction = "resolve"
	}

	dedupKey := fmt.Sprintf(
		"cosmos-proposals-checker alert chain=%s proposal=%s wallet=%s",
		e.Chain.Name,
		e.ProposalID,
		e.Wallet,
	)

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
	}

	links := []PagerDutyLink{}
	if e.Chain.KeplrName != "" {
		links = append(links, PagerDutyLink{
			Href: e.Chain.GetKeplrLink(e.ProposalID),
			Text: "Keplr",
		})
	}

	explorerLinks := e.Chain.GetExplorerProposalsLinks(e.ProposalID)
	for _, link := range explorerLinks {
		links = append(links, PagerDutyLink{
			Href: link.Link,
			Text: link.Name,
		})
	}

	return PagerDutyAlert{
		Payload: PagerDutyAlertPayload{
			Summary: fmt.Sprintf(
				"Wallet %s hasn't voted on proposal %s on %s: %s",
				e.Wallet,
				e.ProposalID,
				e.Chain.GetName(),
				e.ProposalTitle,
			),
			Timestamp: time.Now().Format(time.RFC3339),
			Severity:  "error",
			Source:    hostname,
			CustomDetails: map[string]string{
				"Wallet":               e.Wallet,
				"Chain":                e.Chain.GetName(),
				"Proposal ID":          e.ProposalID,
				"Proposal title":       e.ProposalTitle,
				"Proposal description": e.ProposalDescription,
			},
		},
		Links:       links,
		RoutingKey:  r.APIKey,
		EventAction: eventAction,
		DedupKey:    dedupKey,
		Client:      "cosmos-proposals-checker",
		ClientURL:   "https://github.com/freak12techno/cosmos-proposals-checker",
	}
}

func NewPagerDutyReporter(config PagerDutyConfig, logger *zerolog.Logger) PagerDutyReporter {
	return PagerDutyReporter{
		PagerDutyURL: config.PagerDutyURL,
		APIKey:       config.APIKey,
		Logger:       logger.With().Str("component", "pagerduty_reporter").Logger(),
	}
}

func (r PagerDutyReporter) Init() {
}

func (r PagerDutyReporter) Enabled() bool {
	return r.APIKey != ""
}

func (r PagerDutyReporter) Name() string {
	return "pagerduty-reporter"
}

func (r PagerDutyReporter) SendReport(report Report) error {
	var err error

	for _, entry := range report.Entries {
		if !entry.IsVoteOrNotVoted() {
			continue
		}

		alert := r.NewPagerDutyAlertFromReportEntry(entry)

		if alertErr := r.SendAlert(alert); alertErr != nil {
			err = alertErr
		}
	}

	return err
}

func (r PagerDutyReporter) SendAlert(alert PagerDutyAlert) error {
	var response PagerDutyResponse
	err := r.DoRequest(r.PagerDutyURL+"/v2/enqueue", alert, &response)
	if err != nil {
		return err
	}

	if response.Status != "success" {
		return fmt.Errorf("expected 'success' status, got '" + response.Status + "'. Error: " + response.Message)
	}

	return nil
}

func (r PagerDutyReporter) DoRequest(url string, body interface{}, target interface{}) error {
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
