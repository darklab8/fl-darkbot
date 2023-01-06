package base

import (
	"darkbot/utils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBases(t *testing.T) {

	var storage BaseRecords
	storage.api = APIBasespy{}
	storage.addFromAPI()

	bases, err := storage.getLatestRecord()
	utils.CheckPanic(err, "not found latest base record")

	assert.True(t, len(bases.list) > 0)
	fmt.Println(bases.list)
}
