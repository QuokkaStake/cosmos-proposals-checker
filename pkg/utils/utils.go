package utils

func Filter[T any](slice []T, f func(T) bool) []T {
	var n []T
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func ResolveVote(value string) string {
	votes := map[string]string{
		"VOTE_OPTION_YES":          "Yes",
		"VOTE_OPTION_ABSTAIN":      "Abstain",
		"VOTE_OPTION_NO":           "No",
		"VOTE_OPTION_NO_WITH_VETO": "No with veto",
	}

	if vote, ok := votes[value]; ok && vote != "" {
		return vote
	}

	return value
}
