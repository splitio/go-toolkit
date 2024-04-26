package specs

import (
	"github.com/splitio/go-toolkit/v5/sdk/specs/matchers"
	"testing"
)

func Test_splitVersionFilter_shouldFilter(t *testing.T) {
	filter := NewSplitVersionFilter()
	shouldFilter := filter.ShouldFilter(matchers.MatcherTypeBetweenSemver, "1.0")
	if !shouldFilter {
		t.Error("It should filtered")
	}

	shouldFilter = filter.ShouldFilter(matchers.MatcherTypeEqualTo, "1.0")
	if shouldFilter {
		t.Error("It should not filtered")
	}

	shouldFilter = filter.ShouldFilter(matchers.MatcherTypeBetweenSemver, "1.1")
	if shouldFilter {
		t.Error("It should not filtered")
	}
}
