package baseattack

import (
	"fmt"
	"testing"
)

func TestAPI(t *testing.T) {
	api := basesattackAPI{}.New()
	data := string(api.GetData())
	fmt.Println(data)
}
