package baseattack

import (
	"strings"
	"testing"

	"github.com/darklab8/fl-darkbot/app/settings/logus"

	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	api := FixtureBaseAttackAPIMock()
	result, _ := api.GetBaseAttackData()
	data := string(result)
	logus.Log.Debug(data)
}

func TestDetectLPAttack(t *testing.T) {
	api := NewMock("data_lp.json")
	result, _ := api.GetBaseAttackData()
	data := string(result)
	assert.True(t, strings.Contains(data, "LP-7743"))
}
