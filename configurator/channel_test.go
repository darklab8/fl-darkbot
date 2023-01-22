package configurator

import (
	"darkbot/settings"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannels(t *testing.T) {
	os.Remove(settings.Dbpath)
	cg := ConfiguratorChannel{Configurator: NewConfigurator().Migrate()}

	cg.Add("1")
	cg.Add("2")
	cg.Add("3")

	channels, _ := cg.List()
	fmt.Println(channels)
	assert.Len(t, channels, 3)

	cg.Remove("3")

	channels, _ = cg.List()
	fmt.Println(channels)
	assert.Len(t, channels, 2)

	cg.Add("3")

	channels, _ = cg.List()
	fmt.Println(channels)
	assert.Len(t, channels, 3)
}
