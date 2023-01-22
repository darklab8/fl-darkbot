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

	cg.Add("1", "2", "3")

	channels := cg.List()
	fmt.Println(channels)

	assert.Len(t, channels, 3)
}
