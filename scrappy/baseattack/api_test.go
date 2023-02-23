package baseattack

import (
	"fmt"
	"testing"
)

func TestAPI(t *testing.T) {
	api := basesattackAPI{}.New()
	result, _ := api.GetData()
	data := string(result)
	fmt.Println(data)
}
