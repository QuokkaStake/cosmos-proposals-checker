package report

import (
	"context"
	mutes "main/pkg/mutes"
	reportersPkg "main/pkg/reporters"

	"go.opentelemetry.io/otel/trace"

	"github.com/rs/zerolog"
)

type Dispatcher struct {
	Logger       zerolog.Logger
	MutesManager *mutes.Manager
	Reporters    []reportersPkg.Reporter
	Tracer       trace.Tracer
}

func NewDispatcher(
	logger *zerolog.Logger,
	mutesManager *mutes.Manager,
	reporters []reportersPkg.Reporter,
	tracer trace.Tracer,
) *Dispatcher {
	return &Dispatcher{
		Logger:       logger.With().Str("component", "report_dispatcher").Logger(),
		MutesManager: mutesManager,
		Reporters:    reporters,
		Tracer:       tracer,
	}
}

func (d *Dispatcher) Init() error {
	for _, reporter := range d.Reporters {
		if err := reporter.Init(); err != nil {
			d.Logger.Error().Err(err).
				Str("name", reporter.Name()).
				Msg("Error initializing reporter")
			return err
		}
		if reporter.Enabled() {
			d.Logger.Info().Str("name", reporter.Name()).Msg("Init reporter")
		}
	}

	return nil
}

func (d *Dispatcher) SendReport(report reportersPkg.Report, ctx context.Context) {
	childCtx, span := d.Tracer.Start(ctx, "Sending report")
	defer span.End()

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
			if d.MutesManager.IsEntryMuted(reportEntry) {
				d.Logger.Debug().
					Str("entry", reportEntry.Name()).
					Msg("Notifications are muted, not sending.")
				continue
			}

			if err := reporter.SendReportEntry(reportEntry, childCtx); err != nil {
				d.Logger.Error().
					Err(err).
					Str("name", reporter.Name()).
					Str("entry", reportEntry.Name()).
					Msg("Failed to send report entry")
			}
		}
	}
}
