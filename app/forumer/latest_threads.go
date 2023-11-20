package forumer

import (
	"darkbot/app/settings/logus"
	"fmt"
	"net/url"

	"github.com/anaskhan96/soup"
)

type Url string

type ThreadLink Url

func (u ThreadLink) GetUrl() Url { return Url(u) }

type ThreadShortName string
type ThreadID string
type ForumTimestamp string

type PostAuthorLink Url

func (u PostAuthorLink) GetUrl() Url { return Url(u) }

type PostAuthorName string

type LatestThread struct {
	ThreadLink     ThreadLink
	ThreadName     ThreadShortName
	ThreadID       ThreadID
	LastUpdated    ForumTimestamp
	PostAuthorLink PostAuthorLink
	PostAuthorName PostAuthorName
}

type ThreadsPage struct {
	Threads   []LatestThread
	requester func(MethodType, Url) (*QueryResult, error)
}

type threadPageParam func(thread_page *ThreadsPage)

func WithMockedPageRequester(
	requester func(MethodType, Url) (*QueryResult, error),
) threadPageParam {
	return func(thread_page *ThreadsPage) {
		thread_page.requester = requester
	}
}

const ThreadPageURL Url = "https://discoverygc.com/forums/portal.php"

func GetLatestThreads(opts ...threadPageParam) (*ThreadsPage, error) {
	thread_page := &ThreadsPage{
		requester: NewQuery,
	}
	for _, opt := range opts {
		opt(thread_page)
	}

	query, err := thread_page.requester(GET, ThreadPageURL)
	if logus.CheckError(err, "Failed to make query") {
		return nil, err
	}

	content := query.GetContent()
	doc := soup.HTMLParse(content)
	forum_posts := doc.FindAll("tr", "class", "latestthreads_portal")

	for _, forum_post := range forum_posts {
		thread := forum_post.Find("strong").Find("a")
		if logus.CheckError(thread.Error, "failed to get thread object") {
			return nil, thread.Error
		}

		thread_link := thread.Attrs()["href"]
		thread_name := thread.Text()
		span_section := forum_post.Find("span")
		if logus.CheckError(span_section.Error, "failed to get span_section object") {
			return nil, span_section.Error
		}

		forum_timestamp := span_section.Find("span").Attrs()["title"]

		author := span_section.Find("a")
		if logus.CheckError(author.Error, "failed to get author object") {
			return nil, author.Error
		}
		author_link := author.Attrs()["href"]
		author_name := author.Text()

		myUrl, _ := url.Parse(thread_link)
		params, _ := url.ParseQuery(myUrl.RawQuery)

		latest_thread := LatestThread{
			ThreadLink:     ThreadLink(thread_link),
			ThreadName:     ThreadShortName(thread_name),
			LastUpdated:    ForumTimestamp(forum_timestamp),
			PostAuthorLink: PostAuthorLink(author_link),
			PostAuthorName: PostAuthorName(author_name),
			ThreadID:       ThreadID(params["tid"][0]),
		}
		thread_page.Threads = append(thread_page.Threads, latest_thread)

		logus.Debug(fmt.Sprintf("latest_thread=%v", latest_thread))
	}
	return thread_page, nil
}
