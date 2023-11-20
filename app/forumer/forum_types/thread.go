package forum_types

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
