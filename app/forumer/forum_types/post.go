package forum_types

type PostID string
type PostContent string
type PostPermamentLink Url
type ThreadFullName string
type Subforum string
type Post struct {
	*LatestThread
	PostID               PostID
	PostContent          PostContent
	PostPermamentLink    PostPermamentLink
	ThreadFullName       ThreadFullName
	PostAuthorAvatarLink Url
	Subforums            []Subforum
}
