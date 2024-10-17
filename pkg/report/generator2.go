package report

import (
	"context"
	"fmt"
	databasePkg "main/pkg/database"
	"main/pkg/events"
	fetchersPkg "main/pkg/fetchers"
	"main/pkg/report/entry"
	"main/pkg/reporters"
	"main/pkg/types"
	"sync"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type NewGenerator struct {
	Chains   types.Chains
	Logger   zerolog.Logger
	Database databasePkg.Database
	Fetchers map[string]fetchersPkg.Fetcher
	Tracer   trace.Tracer
}

func NewReportNewGenerator(
	logger *zerolog.Logger,
	chains types.Chains,
	database databasePkg.Database,
	tracer trace.Tracer,
) *NewGenerator {
	fetchers := make(map[string]fetchersPkg.Fetcher, len(chains))

	for _, chain := range chains {
		fetchers[chain.Name] = fetchersPkg.GetFetcher(chain, logger, tracer)
	}

	return &NewGenerator{
		Chains:   chains,
		Logger:   logger.With().Str("component", "report_generator").Logger(),
		Tracer:   tracer,
		Fetchers: fetchers,
		Database: database,
	}
}

func (g *NewGenerator) GenerateReport(ctx context.Context) reporters.Report {
	_, span := g.Tracer.Start(ctx, "Generating report")
	defer span.End()

	entries := []entry.ReportEntry{}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	wg.Add(len(g.Chains))

	for _, chain := range g.Chains {
		go func(chain *types.Chain) {
			defer wg.Done()

			chainEntries := g.ProcessChain(chain, ctx)

			mutex.Lock()
			entries = append(entries, chainEntries...)
			mutex.Unlock()
		}(chain)
	}

	wg.Wait()

	return reporters.Report{Entries: entries}
}

func (g *NewGenerator) ProcessChain(chain *types.Chain, ctx context.Context) []entry.ReportEntry {
	childCtx, span := g.Tracer.Start(ctx, "Processing chain")
	span.SetAttributes(attribute.String("chain", chain.Name))
	defer span.End()

	fetcher := g.Fetchers[chain.Name]

	g.Logger.Trace().Str("chain", chain.Name).Msg("Processing chain...")

	prevHeight, prevHeightErr := g.Database.GetLastBlockHeight(chain, "proposals")
	if prevHeightErr != nil {
		g.Logger.Error().Err(prevHeightErr).Msg("Failed to fetch last block height")
		span.RecordError(prevHeightErr)
		return []entry.ReportEntry{}
	}

	proposals, newHeight, err := fetcher.GetAllProposals(prevHeight, childCtx)
	if err != nil {
		g.Logger.Error().Str("chain", chain.Name).Err(err).Msg("Error fetching proposals")
		span.RecordError(err)
		return []entry.ReportEntry{
			events.ProposalsQueryErrorEvent{Chain: chain, Error: err},
		}
	}

	if insertErr := g.Database.UpsertLastBlockHeight(chain, "proposals", newHeight); insertErr != nil {
		g.Logger.Error().Err(prevHeightErr).Msg("Failed to insert last block height")
		span.RecordError(prevHeightErr)
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	entries := make([]entry.ReportEntry, 0)

	wg.Add(len(proposals))

	for _, proposal := range proposals {
		go func(proposal types.Proposal) {
			defer wg.Done()

			proposalEntries := g.ProcessProposal(chain, proposal, childCtx)

			mutex.Lock()
			entries = append(entries, proposalEntries...)
			mutex.Unlock()
		}(proposal)
	}

	wg.Wait()

	return entries
}

func (g *NewGenerator) ProcessProposal(
	chain *types.Chain,
	proposal types.Proposal,
	ctx context.Context,
) []entry.ReportEntry {
	childCtx, span := g.Tracer.Start(ctx, "Processing proposal")
	span.SetAttributes(attribute.String("chain", chain.Name))
	span.SetAttributes(attribute.String("proposal_id", proposal.ID))
	defer span.End()

	previousProposal, err := g.Database.GetProposal(chain, proposal.ID)
	if err != nil {
		g.Logger.Error().Err(err).Msg("Failed to fetch proposal from DB")
		span.RecordError(err)
		return []entry.ReportEntry{}
	}

	entries := []entry.ReportEntry{}

	if previousProposal != nil && previousProposal.IsInVoting() && !proposal.IsInVoting() {
		g.Logger.Trace().
			Str("chain", chain.Name).
			Str("proposal", proposal.ID).
			Msg("Voting on a proposal has finished")

		entries = append(entries, events.FinishedVotingEvent{
			Chain:    chain,
			Proposal: proposal,
		})
	}

	if previousProposal == nil || !previousProposal.Equals(proposal) {
		if updateErr := g.Database.UpsertProposal(chain, proposal); updateErr != nil {
			g.Logger.Error().Err(updateErr).Msg("Failed to update proposal in DB")
			span.RecordError(err)
		}
	}

	if !proposal.IsInVoting() {
		g.Logger.Trace().
			Str("chain", chain.Name).
			Str("proposal", proposal.ID).
			Msg("Proposal is not in voting period - not fetching votes.")
		return entries
	}

	g.Logger.Trace().
		Str("chain", chain.Name).
		Str("proposal", proposal.ID).
		Msg("Processing proposal...")

	var wg sync.WaitGroup
	var mutex sync.Mutex

	wg.Add(len(chain.Wallets))

	for _, wallet := range chain.Wallets {
		go func(wallet *types.Wallet) {
			defer wg.Done()

			walletEntries := g.ProcessWallet(chain, proposal, wallet, childCtx)

			mutex.Lock()
			entries = append(entries, walletEntries...)
			mutex.Unlock()
		}(wallet)
	}

	wg.Wait()

	return entries
}

func (g *NewGenerator) ProcessWallet(
	chain *types.Chain,
	proposal types.Proposal,
	wallet *types.Wallet,
	ctx context.Context,
) []entry.ReportEntry {
	childCtx, span := g.Tracer.Start(ctx, "Processing wallet")
	span.SetAttributes(attribute.String("chain", chain.Name))
	span.SetAttributes(attribute.String("proposal_id", proposal.ID))
	span.SetAttributes(attribute.String("wallet", wallet.Address))
	defer span.End()

	g.Logger.Trace().
		Str("chain", chain.Name).
		Str("proposal", proposal.ID).
		Str("address", wallet.Address).
		Msg("Processing wallet...")

	fetcher := g.Fetchers[chain.Name]
	storeKey := fmt.Sprintf("proposal_%s_vote_%s", proposal.ID, wallet.Address)

	prevHeight, prevHeightErr := g.Database.GetLastBlockHeight(chain, storeKey)
	if prevHeightErr != nil {
		g.Logger.Error().Err(prevHeightErr).Msg("Failed to fetch last block height")
		span.RecordError(prevHeightErr)
		return []entry.ReportEntry{}
	}

	vote, newHeight, err := fetcher.GetVote(proposal.ID, wallet.Address, prevHeight, childCtx)
	if err != nil {
		g.Logger.Error().Err(err).Msg("Failed to fetch vote from chain")
		span.RecordError(err)
		return []entry.ReportEntry{
			events.VoteQueryError{
				Chain:    chain,
				Proposal: proposal,
				Error:    err,
			},
		}
	}

	if insertErr := g.Database.UpsertLastBlockHeight(chain, storeKey, newHeight); insertErr != nil {
		g.Logger.Error().Err(prevHeightErr).Msg("Failed to insert last block height")
		span.RecordError(prevHeightErr)
	}

	if vote == nil {
		g.Logger.Trace().
			Str("chain", chain.Name).
			Str("proposal", proposal.ID).
			Str("address", wallet.Address).
			Msg("Wallet has not voted - sending an alert.")
		return []entry.ReportEntry{
			events.NotVotedEvent{
				Chain:    chain,
				Proposal: proposal,
				Wallet:   wallet,
			},
		}
	}

	previousVote, dbErr := g.Database.GetVote(chain, proposal, wallet)
	if dbErr != nil {
		g.Logger.Error().Err(err).Msg("Failed to fetch vote from DB")
		span.RecordError(err)
		return []entry.ReportEntry{}
	}

	if previousVote == nil || !previousVote.VotesEquals(vote) {
		if updateErr := g.Database.UpsertVote(chain, proposal, wallet, vote, childCtx); updateErr != nil {
			g.Logger.Error().Err(updateErr).Msg("Failed to update vote in DB")
			span.RecordError(err)
		}
	}

	if previousVote == nil {
		g.Logger.Trace().
			Str("chain", chain.Name).
			Str("proposal", proposal.ID).
			Str("address", wallet.Address).
			Msg("Wallet has voted - sending an alert.")

		return []entry.ReportEntry{
			events.VotedEvent{
				Chain:    chain,
				Proposal: proposal,
				Wallet:   wallet,
				Vote:     vote,
			},
		}
	}

	if !vote.VotesEquals(previousVote) {
		g.Logger.Trace().
			Str("chain", chain.Name).
			Str("proposal", proposal.ID).
			Str("address", wallet.Address).
			Msg("Wallet has changed its vote - sending an alert.")

		return []entry.ReportEntry{
			events.RevotedEvent{
				Chain:    chain,
				Proposal: proposal,
				Wallet:   wallet,
				Vote:     vote,
				OldVote:  previousVote,
			},
		}
	}

	return []entry.ReportEntry{}
}
