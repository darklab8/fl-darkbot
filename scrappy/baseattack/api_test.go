package baseattack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	api := NewBaseAttackAPI()
	result, _ := api.GetBaseAttackData()
	data := string(result)
	fmt.Println(data)
}

func TestDetectLPAttack(t *testing.T) {
	api := NewMock("data_lp.json")
	result, _ := api.GetBaseAttackData()
	data := string(result)
	assert.True(t, strings.Contains(data, "LP-7743"))
}
