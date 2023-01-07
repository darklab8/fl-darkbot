package base

import (
	"darkbot/scrappy/shared/api"
	"darkbot/utils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type BaseStorageMocked struct {
	BaseStorage
}

func (b BaseStorageMocked) API() api.APIinterface {
	return APIBasespy{}
}

func TestGetBases(t *testing.T) {

	// var storage storage.IStorage
	// storage = &BaseStorageMocked{}
	storage := BaseStorage{}
	storage.api = APIBasespy{}
	storage.Update()

	bases, err := storage.GetLatestRecord()
	utils.CheckPanic(err, "not found latest base record")

	assert.True(t, len(bases.list) > 0)
	fmt.Println(bases.list)
}

func TestAddManyRecords(t *testing.T) {

	var storage BaseStorage
	storage.api = APIBasespy{}
	storage.Update()
	storage.Update()

	assert.Equal(t, 2, storage.Length())
}
