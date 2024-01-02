package forumer

import (
	"darkbot/app/configurator"
	"darkbot/app/forumer/forum_types"
	"darkbot/app/settings/darkbot_logus"
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
	"testing"
)

type MockedThreadsQuery struct {
	threads []*forum_types.LatestThread
}

func (m MockedThreadsQuery) GetLatestThreads(opts ...threadPageParam) ([]*forum_types.LatestThread, error) {
	return m.threads, nil
}

func newMockedThreadsQuery() MockedThreadsQuery {
	mocked_threads_requester := FixtureMockedThreadsRequester()
	threads_requester := NewLatestThreads(WithMockedPageRequester(mocked_threads_requester))
	threads, err := threads_requester.GetLatestThreads()
	darkbot_logus.Log.CheckFatal(err, "unexpected error from GetLatestThreads")
	one_thread := threads[:1]
	return MockedThreadsQuery{threads: one_thread}
}

func TestForumerSending(t *testing.T) {

	mocked_post_requester := FixtureDetailedRequester()
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		dev_env_channel := types.DiscordChannelID("1079189823098724433")
		cg := configurator.NewConfiguratorForumWatch(configurator.NewConfigurator(dbpath))
		cg.TagsAdd(dev_env_channel, []types.Tag{""}...)

		cg_channel := configurator.NewConfiguratorChannel(configurator.NewConfigurator(dbpath))

		if utils.FixtureDevEnv() {
			cg_channel.Add(dev_env_channel)
		}

		forum := NewForumer(
			dbpath,
			WithThreadsRequester(newMockedThreadsQuery()),
			WithDetailedPostRequest(NewDetailedPostRequester(WithMockedRequester(mocked_post_requester))),
		)

		forum.update()
	})
}

func TestSubForumSending(t *testing.T) {

	mocked_post_requester := FixtureDetailedRequester()
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		dev_env_channel := types.DiscordChannelID("1079189823098724433")
		cg := configurator.NewConfiguratorSubForumWatch(configurator.NewConfigurator(dbpath))
		cg.TagsAdd(dev_env_channel, []types.Tag{"Communication Channel"}...)

		cg_channel := configurator.NewConfiguratorChannel(configurator.NewConfigurator(dbpath))

		if utils.FixtureDevEnv() {
			cg_channel.Add(dev_env_channel)
		}

		forum := NewForumer(
			dbpath,
			WithThreadsRequester(newMockedThreadsQuery()),
			WithDetailedPostRequest(NewDetailedPostRequester(WithMockedRequester(mocked_post_requester))),
		)

		forum.update()
	})
}
