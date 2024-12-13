package logging

import (
	"testing"

	"github.com/splitio/go-toolkit/v6/common"
	"github.com/stretchr/testify/assert"
)

func TestContextData(t *testing.T) {
	one := NewContext().WithTag("key", "value")
	assert.Equal(t, one.Get("key"), "value")
	assert.Empty(t, one.Get("key2"))
	assert.Equal(t, one.String(), "[key=value]")
	assert.Nil(t, Merge(nil, nil))
	assert.Equal(t, Merge(one, nil).String(), "[key=value]")
	two := NewContext().WithTag("key2", "value2")
	assert.Equal(t, Merge(nil, two).String(), "[key2=value2]")
	assert.Equal(t, Merge(one, two).String(), "[key=value, key2=value2]")
	three := NewContext().WithBaselineContext(Baseline{
		orgID: common.Ref("org"),
		envID: common.Ref("env"),
		name:  common.Ref("name"),
		txIDs: common.Ref("tx"),
	})
	assert.Equal(t, three.String(), "[env_id=env, ls_name=name, org_id=org, tx_ids=tx]")
}
