package forum

import (
	"darkbot/app/settings/logus"
	"fmt"
	"net/url"
	"strings"

	"github.com/anaskhan96/soup"
)

type DetailedPost struct {
	*LatestThread
	PostID            PostID
	PostContent       PostContent
	PostPermamentLink PostPermamentLink

	requester func(MethodType, Url) (*QueryResult, error)
}
type PostID string
type PostContent string
type PostPermamentLink Url

type detailedPostParam func(detailedPost *DetailedPost)

func WithMockedRequester(
	requester func(MethodType, Url) (*QueryResult, error),
) detailedPostParam {
	return func(detailedPost *DetailedPost) {
		detailedPost.requester = requester
	}
}

func NewDetailedPost(thread *LatestThread, opts ...detailedPostParam) (*DetailedPost, error) {
	detailed_post := &DetailedPost{requester: NewQuery}
	for _, opt := range opts {
		opt(detailed_post)
	}

	query, err := detailed_post.requester(GET, thread.ThreadLink.GetUrl())
	if logus.CheckError(err, "failed to query ThreadLink.GetUrl()="+string(thread.ThreadLink)) {
		return nil, err
	}

	doc := soup.HTMLParse(query.GetContent())
	params, _ := url.ParseQuery(query.ResponseRawQuery)
	post_id := PostID(params["pid"][0])

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

	detailed_post.LatestThread = thread
	detailed_post.PostID = post_id
	detailed_post.PostContent = PostContent(post_content)
	detailed_post.PostPermamentLink = PostPermamentLink(query.ResponseFullUrl)

	return detailed_post, nil
}
