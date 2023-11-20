package forumer

import (
	"darkbot/app/forumer/forum_types"
	"darkbot/app/settings/logus"
	"fmt"
	"net/url"
	"strings"

	"github.com/anaskhan96/soup"
)

type PostRequester struct {
	requester func(MethodType, forum_types.Url) (*QueryResult, error)
}

type detailedPostParam func(p *PostRequester)

func WithMockedRequester(
	requester func(MethodType, forum_types.Url) (*QueryResult, error),
) detailedPostParam {
	return func(p *PostRequester) {
		p.requester = requester
	}
}

func NewDetailedPostRequester(opts ...detailedPostParam) *PostRequester {
	res := &PostRequester{
		requester: NewQuery,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (p *PostRequester) GetDetailedPost(thread *forum_types.LatestThread) (*forum_types.Post, error) {
	query, err := p.requester(GET, thread.ThreadLink.GetUrl())
	if logus.CheckError(err, "failed to query ThreadLink.GetUrl()="+string(thread.ThreadLink)) {
		return nil, err
	}

	doc := soup.HTMLParse(query.GetContent())
	params, _ := url.ParseQuery(query.ResponseRawQuery)
	post_id := forum_types.PostID(params["pid"][0])

	forum := doc.Find("div", "id", "forum")
	if logus.CheckError(forum.Error, "failed to get forum object") {
		return nil, forum.Error
	}
	thread_header := forum.Find("td", "class", "thead")
	if logus.CheckError(thread_header.Error, "failed to get thread_header object") {
		return nil, thread_header.Error
	}
	thread_name := thread_header.FullText()
	thread_name = strings.ReplaceAll(thread_name, "\n", "")
	thread_name = strings.ReplaceAll(thread_name, "\t", "")
	logus.Debug("thread_name=" + thread_name)

	post := doc.Find("table", "id", fmt.Sprintf("post_%s", string(post_id)))
	if logus.CheckError(post.Error, "failed to get post object") {
		return nil, post.Error
	}
	post_body := post.Find("div", "class", "post_body")
	if logus.CheckError(post_body.Error, "failed to get post_body object") {
		return nil, post_body.Error
	}
	post_content := post_body.FullText()

	post_content = strings.ReplaceAll(post_content, "\t", "")
	for i := 0; i < 5; i++ {
		post_content = strings.ReplaceAll(post_content, "\n\n", "\n")
	}

	return &forum_types.Post{
		LatestThread:      thread,
		PostID:            post_id,
		PostContent:       forum_types.PostContent(post_content),
		PostPermamentLink: forum_types.PostPermamentLink(query.ResponseFullUrl),
	}, nil
}
