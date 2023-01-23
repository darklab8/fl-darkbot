package consoler

import (
	"darkbot/configurator"
	"darkbot/consoler/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGettingOutput(t *testing.T) {
	configurator.ConfiguratorChannel{Configurator: configurator.NewConfigurator().Migrate()}.Add("123")
	assert.Contains(t, Consoler{}.New(". ping").Execute(helper.ChannelInfo{ChannelID: "123"}).String(), "Pong!")
}

func TestGrabStdout(t *testing.T) {
	configurator.ConfiguratorChannel{Configurator: configurator.NewConfigurator().Migrate()}.Add("123")
	c := Consoler{}.New(". ping --help")
	result := c.Execute(helper.ChannelInfo{ChannelID: "123"}).String()

	assert.Contains(t, result, "\nFlags:\n  -h, --help   ")
}

func TestAddBaseTag(t *testing.T) {
	configurator.ConfiguratorChannel{Configurator: configurator.NewConfigurator().Migrate()}.Add("123")
	assert.Contains(t, Consoler{}.New(`. base add "bla bla" sdf`).Execute(helper.ChannelInfo{ChannelID: "123"}).String(), "OK tags are added")
}
