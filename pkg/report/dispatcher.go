package report

import (
	mutes "main/pkg/mutes"
	reportersPkg "main/pkg/reporters"

	"github.com/rs/zerolog"
)

type Dispatcher struct {
	Logger       zerolog.Logger
	MutesManager *mutes.Manager
	Reporters    []reportersPkg.Reporter
}

func NewDispatcher(
	logger *zerolog.Logger,
	mutesManager *mutes.Manager,
	reporters []reportersPkg.Reporter,
) *Dispatcher {
	return &Dispatcher{
		Logger:       logger.With().Str("component", "report_dispatcher").Logger(),
		MutesManager: mutesManager,
		Reporters:    reporters,
	}
}

func (d *Dispatcher) Init() {
	d.MutesManager.Load()

	for _, reporter := range d.Reporters {
		if err := reporter.Init(); err != nil {
			d.Logger.Fatal().Err(err).
				Str("name", reporter.Name()).
				Msg("Error initializing reporter")
		}
		if reporter.Enabled() {
			d.Logger.Info().Str("name", reporter.Name()).Msg("Init reporter")
		}
	}
}

func (d *Dispatcher) SendReport(report reportersPkg.Report) {
	if report.Empty() {
		d.Logger.Debug().Msg("Empty report, not sending.")
		return
	}

	d.Logger.Debug().Int("len", len(report.Entries)).Msg("Got non-empty report")

	for _, reporter := range d.Reporters {
		if !reporter.Enabled() {
			d.Logger.Debug().
				Str("name", reporter.Name()).
				Msg("Reporter is disabled, not sending report")
			continue
		}

		d.Logger.Debug().
			Str("name", reporter.Name()).
			Msg("Sending report...")

		for _, reportEntry := range report.Entries {
			if err := reporter.SendReportEntry(reportEntry); err != nil {
				d.Logger.Error().
					Err(err).
					Str("name", reporter.Name()).
					Str("entry", reportEntry.Name()).
					Msg("Failed to send report entry")
			}
		}
	}
}
