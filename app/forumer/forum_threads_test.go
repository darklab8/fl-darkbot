package forumer

import (
	"darkbot/app/settings/logus"
	"darkbot/app/settings/utils"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPosts(t *testing.T) {
	thread_post_content_filepath := filepath.Join(utils.GetCurrrentFolder(), "test_data", "latest_threads.html")
	if _, err := os.Stat(thread_post_content_filepath); err != nil {
		query, err := NewQuery("GET", ThreadPageURL)
		logus.CheckFatal(err, "failed to create mock data")
		os.WriteFile(thread_post_content_filepath, []byte(query.GetContent()), 0644)
	}
	thread_post_content, _ := os.ReadFile(thread_post_content_filepath)
	mocked_requester := func(mt MethodType, u Url) (*QueryResult, error) {
		return &QueryResult{
			content:          string(thread_post_content),
			ResponseRawQuery: ``,
			ResponseFullUrl:  `https://discoverygc.com/forums/portal.php`,
		}, nil
	}

	threads, err := NewLatestThreads(WithMockedPageRequester(mocked_requester)).GetLatestThreads()
	assert.Nil(t, err, "expected nil as error")
	assert.Greater(t, len(threads), 0)
}
