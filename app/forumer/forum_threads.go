package forumer

import (
	"fmt"
	"net/url"

	"github.com/darklab8/fl-darkbot/app/forumer/forum_types"
	"github.com/darklab8/fl-darkbot/app/settings/logus"

	"github.com/anaskhan96/soup"
)

type ThreadsRequester struct {
	requester func(MethodType, forum_types.Url) (*QueryResult, error)
}

type threadPageParam func(thread_page *ThreadsRequester)

func WithMockedPageRequester(
	requester func(MethodType, forum_types.Url) (*QueryResult, error),
) threadPageParam {
	return func(thread_page *ThreadsRequester) {
		thread_page.requester = requester
	}
}

const ThreadPageURL forum_types.Url = "https://discoverygc.com/forums/portal.php"

func NewLatestThreads(opts ...threadPageParam) *ThreadsRequester {
	thread_page := &ThreadsRequester{
		requester: NewQuery,
	}

	for _, opt := range opts {
		opt(thread_page)
	}
	return thread_page
}

func (p *ThreadsRequester) GetLatestThreads(opts ...threadPageParam) ([]*forum_types.LatestThread, error) {
	records := []*forum_types.LatestThread{}

	query, err := p.requester(GET, ThreadPageURL)
	if logus.Log.CheckError(err, "Failed to make query") {
		return nil, err
	}

	content := query.GetContent()
	doc := soup.HTMLParse(content)
	forum_posts := doc.FindAll("tr", "class", "latestthreads_portal")

	for _, forum_post := range forum_posts {
		thread := forum_post.Find("strong").Find("a")
		if logus.Log.CheckError(thread.Error, "failed to get thread object") {
			return nil, thread.Error
		}

		thread_link := thread.Attrs()["href"]
		thread_name := thread.Text()
		span_section := forum_post.Find("span")
		if logus.Log.CheckError(span_section.Error, "failed to get span_section object") {
			return nil, span_section.Error
		}

		forum_timestamp := span_section.Find("span").Attrs()["title"]

		author := span_section.Find("a")
		if logus.Log.CheckError(author.Error, "failed to get author object") {
			return nil, author.Error
		}
		author_link := author.Attrs()["href"]
		author_name := author.Text()

		myUrl, _ := url.Parse(thread_link)
		params, _ := url.ParseQuery(myUrl.RawQuery)

		latest_thread := &forum_types.LatestThread{
			ThreadLink:      forum_types.ThreadLink(thread_link),
			ThreadShortName: forum_types.ThreadShortName(thread_name),
			LastUpdated:     forum_types.ForumTimestamp(forum_timestamp),
			PostAuthorLink:  forum_types.PostAuthorLink(author_link),
			PostAuthorName:  forum_types.PostAuthorName(author_name),
			ThreadID:        forum_types.ThreadID(params["tid"][0]),
		}
		records = append(records, latest_thread)

		logus.Log.Debug(fmt.Sprintf("latest_thread=%v", latest_thread))
	}
	return records, nil
}
