package specs

import (
	"github.com/splitio/go-toolkit/v5/sdk/specs/matchers"
)

type SplitVersionFilter struct {
	excluded map[string]*map[string]struct{}
}

func NewSplitVersionFilter() SplitVersionFilter {
	matchersToExclude := map[string]*map[string]struct{}{
		"1.0": {
			matchers.MatcherEqualToSemver:                  {},
			matchers.MatcherTypeLessThanOrEqualToSemver:    {},
			matchers.MatcherTypeGreaterThanOrEqualToSemver: {},
			matchers.MatcherTypeBetweenSemver:              {},
			matchers.MatcherTypeInListSemver:               {},
		},
		"1.1": {},
	}
	return SplitVersionFilter{
		excluded: matchersToExclude,
	}
}

func (s *SplitVersionFilter) ShouldFilter(matcher string, apiVersion string) bool {
	forVersion, ok := s.excluded[apiVersion]
	if !ok {
		return false
	}
	if _, ok := (*forVersion)[matcher]; ok {
		return true
	}
	return false
}
