package forum_types

type PostID string
type PostContent string
type PostPermamentLink Url
type ThreadFullName string

type Post struct {
	*LatestThread
	PostID            PostID
	PostContent       PostContent
	PostPermamentLink PostPermamentLink
	ThreadFullName    ThreadFullName
}
