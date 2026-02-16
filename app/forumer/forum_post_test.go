package forumer

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/darklab8/fl-darkbot/app/forumer/forum_types"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/go-utils/utils/utils_os"

	"github.com/stretchr/testify/assert"
)

func FixtureDetailedRequester() func(mt MethodType, u forum_types.Url) (*QueryResult, error) {
	detailed_post_content_filepath := filepath.Join(utils_os.GetCurrentFolder().ToString(), "test_data", "detailed_post_content.html")
	if _, err := os.Stat(detailed_post_content_filepath); err != nil {
		query, err := NewQuery("GET", "https://discoverygc.com/forums/showthread.php?tid=200175&action=lastpost")
		logus.Log.CheckFatal(err, "failed to create mock data")
		os.WriteFile(detailed_post_content_filepath, []byte(query.GetContent()), 0644)
	}
	detailed_post_content, _ := os.ReadFile(detailed_post_content_filepath)
	mocked_requester := func(mt MethodType, u forum_types.Url) (*QueryResult, error) {
		return &QueryResult{
			content:          string(detailed_post_content),
			ResponseRawQuery: `tid=200175&pid=2315295`,
			ResponseFullUrl:  `https://discoverygc.com/forums/showthread.php?tid=200175&pid=2315295#pid2315295`,
		}, nil
	}
	return mocked_requester
}

func TestGetDetailedPost(t *testing.T) {
	thread := &forum_types.LatestThread{
		ThreadLink:      "https://discoverygc.com/forums/showthread.php?tid=200175&action=lastpost",
		ThreadShortName: "To: NNroute.../(BDM-Direk...",
		LastUpdated:     "11-20-2023, 09:35 AM",
		PostAuthorLink:  "https://discoverygc.com/forums/member.php?action=profile&uid=54754",
		PostAuthorName:  "Civil Servant",
	}

	mocked_requester := FixtureDetailedRequester()
	detailed_post, err := NewDetailedPostRequester(WithMockedRequester(mocked_requester)).GetDetailedPost(thread)
	_ = detailed_post
	fmt.Println("err=", err)
	assert.Nil(t, err, "expected error to be nil")
	assert.Greater(t, len(detailed_post.Subforums), 0)
}

func FixtureThread2() func(mt MethodType, u forum_types.Url) (*QueryResult, error) {
	detailed_post_content_filepath := filepath.Join(utils_os.GetCurrentFolder().ToString(), "test_data", "detailed_post_content2.html")
	if _, err := os.Stat(detailed_post_content_filepath); err != nil {
		query, err := NewQuery("GET", "https://discoverygc.com/forums/showthread.php?tid=188959&action=lastpost")
		logus.Log.CheckFatal(err, "failed to create mock data")
		os.WriteFile(detailed_post_content_filepath, []byte(query.GetContent()), 0644)
	}
	detailed_post_content, _ := os.ReadFile(detailed_post_content_filepath)
	mocked_requester := func(mt MethodType, u forum_types.Url) (*QueryResult, error) {
		return &QueryResult{
			content:          string(detailed_post_content),
			ResponseRawQuery: `tid=188959&pid=2387860`,
			ResponseFullUrl:  `https://discoverygc.com/forums/showthread.php?tid=188959&pid=2387860#pid2387860`,
		}, nil
	}
	return mocked_requester
}

func TestGetDetailedPost2(t *testing.T) {
	thread := &forum_types.LatestThread{
		ThreadLink:      "https://discoverygc.com/forums/showthread.php?tid=188959&action=lastpost",
		ThreadShortName: "[ignore] testing, testing, 123 (17)",
		LastUpdated:     "11-20-2023, 09:35 AM",
		PostAuthorLink:  "https://discoverygc.com/forums/member.php?action=profile&uid=42166",
		PostAuthorName:  "darkwind",
	}

	mocked_requester := FixtureThread2()
	detailed_post, err := NewDetailedPostRequester(WithMockedRequester(mocked_requester)).GetDetailedPost(thread)
	_ = detailed_post
	fmt.Println("err=", err)
	assert.Nil(t, err, "expected error to be nil")
	assert.Greater(t, len(detailed_post.Subforums), 0)
}
