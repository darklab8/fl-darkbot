package forumer

import (
	"darkbot/app/configurator"
	"darkbot/app/forumer/forum_types"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"os"
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
	logus.CheckFatal(err, "unexpected error from GetLatestThreads")
	one_thread := threads[:1]
	return MockedThreadsQuery{threads: one_thread}
}

func FixtureDevEnv() bool {
	return os.Getenv("DEV_ENV") == "true"
}

func TestForumerSending(t *testing.T) {

	mocked_post_requester := FixtureDetailedRequester()
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		dev_env_channel := types.DiscordChannelID("1079189823098724433")
		cg := configurator.NewConfiguratorForumWatch(configurator.NewConfigurator(dbpath))
		cg.TagsAdd(dev_env_channel, []types.Tag{""}...)

		cg_channel := configurator.NewConfiguratorChannel(configurator.NewConfigurator(dbpath))

		if FixtureDevEnv() {
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
