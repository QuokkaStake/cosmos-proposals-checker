package fetchers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestFetcherGetProposals(t *testing.T) {
	t.Parallel()

	fetcher1 := TestFetcher{WithProposalsError: true}
	proposals1, height1, err1 := fetcher1.GetAllProposals(0, context.Background())
	assert.Empty(t, proposals1)
	assert.Equal(t, int64(123), height1)
	require.Error(t, err1)

	fetcher2 := TestFetcher{WithPassedProposals: true}
	proposals2, height2, err2 := fetcher2.GetAllProposals(0, context.Background())
	assert.NotEmpty(t, proposals2)
	assert.Equal(t, int64(123), height2)
	require.Nil(t, err2)

	fetcher3 := TestFetcher{}
	proposals3, height3, err3 := fetcher3.GetAllProposals(0, context.Background())
	assert.NotEmpty(t, proposals3)
	assert.Equal(t, int64(123), height3)
	require.Nil(t, err3)
}

func TestTestFetcherGetVote(t *testing.T) {
	t.Parallel()

	fetcher1 := TestFetcher{WithVoteError: true}
	vote1, height1, err1 := fetcher1.GetVote("proposal", "vote", 0, context.Background())
	assert.Nil(t, vote1)
	assert.Equal(t, int64(456), height1)
	require.Error(t, err1)

	fetcher2 := TestFetcher{WithVote: true}
	vote2, height2, err2 := fetcher2.GetVote("proposal", "vote", 0, context.Background())
	assert.NotNil(t, vote2)
	assert.Equal(t, int64(456), height2)
	require.Nil(t, err2)

	fetcher3 := TestFetcher{}
	vote3, height3, err3 := fetcher3.GetVote("proposal", "vote", 0, context.Background())
	assert.Nil(t, vote3)
	assert.Equal(t, int64(456), height3)
	require.Nil(t, err3)
}

func TestTestFetcherTally(t *testing.T) {
	t.Parallel()

	fetcher1 := TestFetcher{WithTallyError: true}
	tally1, err1 := fetcher1.GetTallies(context.Background())
	assert.NotNil(t, tally1)
	assert.Empty(t, tally1.TallyInfos)
	assert.Nil(t, tally1.Chain)
	require.Error(t, err1)

	fetcher2 := TestFetcher{WithTallyNotEmpty: true}
	tally2, err2 := fetcher2.GetTallies(context.Background())
	assert.NotNil(t, tally2)
	assert.NotEmpty(t, tally2.TallyInfos)
	assert.NotNil(t, tally2.Chain)
	require.Nil(t, err2) //nolint:testifylint // not working

	fetcher3 := TestFetcher{}
	tally3, err3 := fetcher3.GetTallies(context.Background())
	assert.NotNil(t, tally3)
	assert.Empty(t, tally3.TallyInfos)
	assert.Nil(t, tally3.Chain)
	require.Nil(t, err3) //nolint:testifylint // not working
}

func TestTestFetcherParams(t *testing.T) {
	t.Parallel()

	fetcher1 := TestFetcher{WithParamsError: true}
	params1, errs1 := fetcher1.GetChainParams(context.Background())
	assert.NotNil(t, params1)
	assert.Empty(t, params1.Params)
	assert.Nil(t, params1.Chain)
	require.NotEmpty(t, errs1)

	fetcher2 := TestFetcher{}
	params2, errs2 := fetcher2.GetChainParams(context.Background())
	assert.NotNil(t, params2)
	assert.NotEmpty(t, params2.Params)
	assert.NotNil(t, params2.Chain)
	require.Empty(t, errs2)
}
