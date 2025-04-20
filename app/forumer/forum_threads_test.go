package forumer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/darklab8/fl-darkbot/app/forumer/forum_types"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/go-utils/utils/utils_os"

	"github.com/stretchr/testify/assert"
)

func FixtureMockedThreadsRequester() func(mt MethodType, u forum_types.Url) (*QueryResult, error) {
	thread_post_content_filepath := filepath.Join(utils_os.GetCurrentFolder().ToString(), "test_data", "latest_threads2.html")
	if _, err := os.Stat(thread_post_content_filepath); err != nil {
		query, err := NewQuery("GET", ThreadPageURL)
		logus.Log.CheckFatal(err, "failed to create mock data")
		os.WriteFile(thread_post_content_filepath, []byte(query.GetContent()), 0644)
	}
	thread_post_content, _ := os.ReadFile(thread_post_content_filepath)
	mocked_requester := func(mt MethodType, u forum_types.Url) (*QueryResult, error) {
		return &QueryResult{
			content:          string(thread_post_content),
			ResponseRawQuery: ``,
			ResponseFullUrl:  string(ThreadPageURL),
		}, nil
	}
	return mocked_requester
}
func TestGetPosts(t *testing.T) {
	mocked_threads_requester := FixtureMockedThreadsRequester()
	threads, err := NewLatestThreads(WithMockedPageRequester(mocked_threads_requester)).GetLatestThreads()
	assert.Nil(t, err, "expected nil as error")
	assert.Greater(t, len(threads), 0)
}
