package responses

import (
	"main/pkg/types"
	"strconv"
	"time"

	"cosmossdk.io/math"
)

type ProposalWithID struct {
	ID       int      `json:"id"`
	Proposal Proposal `json:"proposal"`
}
type Proposal struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Expiration  struct {
		AtTime int64 `json:"at_time,string"`
	} `json:"expiration"`
	Status     string `json:"status"`
	TotalPower string `json:"total_power"`

	Votes struct {
		Yes     string `json:"yes"`
		No      string `json:"no"`
		Abstain string `json:"abstain"`
	} `json:"votes"`
}

type ProposalsResponse struct {
	Data struct {
		Proposals []ProposalWithID `json:"proposals"`
	} `json:"data"`
}

func (p ProposalsResponse) ToProposals() []types.Proposal {
	proposals := make([]types.Proposal, len(p.Data.Proposals))

	for index, proposal := range p.Data.Proposals {
		proposals[index] = proposal.ToProposal()
	}

	return proposals
}

func (p ProposalWithID) ToProposal() types.Proposal {
	return types.Proposal{
		ID:          strconv.Itoa(p.ID),
		Title:       p.Proposal.Title,
		Description: p.Proposal.Description,
		EndTime:     time.Unix(0, p.Proposal.Expiration.AtTime),
		Status:      ParseProposalStatus(p.Proposal.Status),
	}
}

func (p ProposalsResponse) ToTally() []types.TallyInfo {
	tallyInfos := make([]types.TallyInfo, 0)

	for _, proposal := range p.Data.Proposals {
		proposalParsed := proposal.ToProposal()

		if !proposalParsed.IsInVoting() {
			continue
		}

		yesVotes := math.LegacyMustNewDecFromStr(proposal.Proposal.Votes.Yes)
		noVotes := math.LegacyMustNewDecFromStr(proposal.Proposal.Votes.No)
		abstainVotes := math.LegacyMustNewDecFromStr(proposal.Proposal.Votes.Abstain)
		totalVotes := math.LegacyMustNewDecFromStr(proposal.Proposal.TotalPower)

		tallyInfos = append(tallyInfos, types.TallyInfo{
			Proposal: proposalParsed,
			Tally: types.Tally{
				{Option: "Yes", Voted: yesVotes},
				{Option: "No", Voted: noVotes},
				{Option: "Abstain", Voted: abstainVotes},
			},
			TotalVotingPower: totalVotes,
		})
	}

	return tallyInfos
}
