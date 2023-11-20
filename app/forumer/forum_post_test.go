package forumer

import (
	"darkbot/app/forumer/forum_types"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/utils"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FixtureLatestThread() *forum_types.LatestThread {
	return &forum_types.LatestThread{
		ThreadLink:      "https://discoverygc.com/forums/showthread.php?tid=200175&action=lastpost",
		ThreadShortName: "To: NNroute.../(BDM-Direk...",
		LastUpdated:     "11-20-2023, 09:35 AM",
		PostAuthorLink:  "https://discoverygc.com/forums/member.php?action=profile&uid=54754",
		PostAuthorName:  "Civil Servant",
	}
}

func TestGetDetailedPost(t *testing.T) {
	thread := FixtureLatestThread()

	detailed_post_content_filepath := filepath.Join(utils.GetCurrrentFolder(), "test_data", "detailed_post_content.html")
	if _, err := os.Stat(detailed_post_content_filepath); err != nil {
		query, err := NewQuery("GET", "https://discoverygc.com/forums/showthread.php?tid=200175&action=lastpost")
		logus.CheckFatal(err, "failed to create mock data")
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
	detailed_post, err := NewDetailedPostRequester(WithMockedRequester(mocked_requester)).GetDetailedPost(thread)
	_ = detailed_post
	fmt.Println("err=", err)
	assert.Nil(t, err, "expected error to be nil")
}
