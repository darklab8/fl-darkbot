package forumer

import (
	"darkbot/app/configurator"
	"darkbot/app/forumer/forum_types"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
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

func TestForumerSending(t *testing.T) {

	mocked_post_requester := FixtureDetailedRequester()
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		forum := NewForumer(
			dbpath,
			WithThreadsRequester(newMockedThreadsQuery()),
			WithDetailedPostRequest(NewDetailedPostRequester(WithMockedRequester(mocked_post_requester))),
		)

		forum.update()
	})

}
