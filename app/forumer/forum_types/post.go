package forum_types

type PostID string
type PostContent string
type PostPermamentLink Url

type Post struct {
	*LatestThread
	PostID            PostID
	PostContent       PostContent
	PostPermamentLink PostPermamentLink
}
