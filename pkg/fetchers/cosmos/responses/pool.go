package responses

import "cosmossdk.io/math"

type PoolRPCResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Pool    *Pool  `json:"pool"`
}

type Pool struct {
	BondedTokens math.LegacyDec `json:"bonded_tokens"`
}
