package forumer

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/darklab/fl-darkbot/app/forumer/forum_types"
	"github.com/darklab/fl-darkbot/app/settings/logus"

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
	if logus.Log.CheckError(err, "failed to query ThreadLink.GetUrl()="+string(thread.ThreadLink)) {
		return nil, err
	}

	doc := soup.HTMLParse(query.GetContent())
	params, _ := url.ParseQuery(query.ResponseRawQuery)
	post_id := forum_types.PostID(params["pid"][0])

	forum := doc.Find("div", "id", "forum")
	if logus.Log.CheckError(forum.Error, "failed to get forum object") {
		return nil, forum.Error
	}
	thread_header := forum.Find("td", "class", "thead")
	if logus.Log.CheckError(thread_header.Error, "failed to get thread_header object") {
		return nil, thread_header.Error
	}
	thread_name := thread_header.FullText()
	thread_name = strings.ReplaceAll(thread_name, "\n", "")
	thread_name = strings.ReplaceAll(thread_name, "\t", "")
	logus.Log.Debug("thread_name=" + thread_name)

	post := doc.Find("table", "id", fmt.Sprintf("post_%s", string(post_id)))
	if logus.Log.CheckError(post.Error, "failed to get post object") {
		return nil, post.Error
	}

	post_author_avatar := post.Find("div", "class", "author_avatar")
	post_author_avatar_a := post_author_avatar.Find("img")
	author_avatar_url := post_author_avatar_a.Attrs()["src"]

	// If u wish getting author from here
	// post_author := post.Find("td", "class", "postcat")
	// iflogus.CheckError(post_author.Error, "failed to get post author") {
	// 	return nil, post_author.Error
	// }
	// post_author_a := post_author.Find("a")
	// iflogus.CheckError(post_author_a.Error, "failed to get post_author_a") {
	// 	return nil, post_author_a.Error
	// }
	// post_author_name := post_author_a.Text()
	// post_author_link := post_author_a.Attrs()["href"]

	post_body := post.Find("div", "class", "post_body")
	if logus.Log.CheckError(post_body.Error, "failed to get post_body object") {
		return nil, post_body.Error
	}
	post_content := post_body.FullText()

	post_content = strings.ReplaceAll(post_content, "\t", "")
	for i := 0; i < 5; i++ {
		post_content = strings.ReplaceAll(post_content, "\n\n", "\n")
	}

	navigation := doc.Find("div", "class", "navigation")
	navigation_childs := navigation.Children()

	subforums := []forum_types.Subforum{}
	for _, may_be_subforum := range navigation_childs {
		if may_be_subforum.NodeValue != "a" {
			continue
		}

		subforums = append(subforums, forum_types.Subforum(may_be_subforum.Text()))
	}

	return &forum_types.Post{
		LatestThread:         thread,
		PostID:               post_id,
		PostContent:          forum_types.PostContent(post_content),
		PostPermamentLink:    forum_types.PostPermamentLink(query.ResponseFullUrl),
		ThreadFullName:       forum_types.ThreadFullName(thread_name),
		PostAuthorAvatarLink: forum_types.Url(author_avatar_url),
		Subforums:            subforums[1:], // first subforum is always root
	}, nil
}
