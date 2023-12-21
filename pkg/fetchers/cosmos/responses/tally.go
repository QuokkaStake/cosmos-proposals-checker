package responses

import (
	"main/pkg/types"

	"cosmossdk.io/math"
)

type TallyRPCResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Tally   *Tally `json:"tally"`
}

type Tally struct {
	Yes        math.LegacyDec `json:"yes"`
	No         math.LegacyDec `json:"no"`
	NoWithVeto math.LegacyDec `json:"no_with_veto"`
	Abstain    math.LegacyDec `json:"abstain"`
}

func (t Tally) ToTally() *types.Tally {
	return &types.Tally{
		{Option: "Yes", Voted: t.Yes},
		{Option: "No", Voted: t.No},
		{Option: "Abstain", Voted: t.Abstain},
		{Option: "No with veto", Voted: t.NoWithVeto},
	}
}
