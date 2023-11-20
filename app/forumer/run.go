package forumer

import (
	"crypto/md5"
	"darkbot/app/configurator"
	"darkbot/app/discorder"
	"darkbot/app/forumer/forum_types"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"time"
)

type Forumer struct {
	Discorder discorder.Discorder
	*configurator.Configurators
	threads_requester *ThreadsRequester
	post_requester    *PostRequester

	cache map[ThreadCacheKey]*forum_types.Post
	// Keeping track as list, to realize which ones are old ones to delete
	cache_keys []ThreadCacheKey
}

type ThreadCacheKey string

func NewThreadCacheKey(thread *forum_types.LatestThread) ThreadCacheKey {
	future_key_as_str := string(thread.LastUpdated) + string(thread.ThreadID) + string(thread.PostAuthorName)
	hash := md5.Sum([]byte(future_key_as_str))
	return ThreadCacheKey(string(hash[:]))
}

type forumerParam func(forum *Forumer)

func WithThreadsRequester(
	threads_page_requester *ThreadsRequester) forumerParam {
	return func(forum *Forumer) { forum.threads_requester = threads_page_requester }
}
func WithDetailedPostRequest(threads_page_requester *ThreadsRequester) forumerParam {
	return func(forum *Forumer) { forum.threads_requester = threads_page_requester }
}

func NewForumer(dbpath types.Dbpath, opts ...forumerParam) *Forumer {

	forum := &Forumer{
		Discorder:         discorder.NewClient(),
		Configurators:     configurator.NewConfigugurators(dbpath),
		threads_requester: NewLatestThreads(),
		post_requester:    NewDetailedPostRequester(),
	}

	for _, opt := range opts {
		opt(forum)
	}

	return forum
}

func (v *Forumer) GetPost(thread *forum_types.LatestThread, new_post_callback func(*forum_types.Post)) *forum_types.Post {
	var err error
	thread_key := NewThreadCacheKey(thread)

	post, ok := v.cache[thread_key]
	if !ok {
		logus.Debug("cache is not found. requesting new post for thread", logus.Thread(thread))
		post, err = v.post_requester.GetDetailedPost(thread)
		logus.CheckError(err, "failed get detailed post for thread=", logus.Thread(thread))
		v.cache[thread_key] = post
		v.cache_keys = append(v.cache_keys, thread_key)
		new_post_callback(post)
	}

	if len(v.cache_keys) > 100 {
		key_to_delete := v.cache_keys[0]
		logus.Debug("deleting old cached key_to_delete=" + string(key_to_delete))
		v.cache_keys = append(v.cache_keys[1:])
		delete(v.cache, key_to_delete)
	}

	return post
}

func (v *Forumer) update() {

	channelIDs, _ := v.Channels.List()

	threads, err := v.threads_requester.GetLatestThreads()
	if logus.CheckError(err, "failed to get threads") {
		return
	}

	for _, thread := range threads {

		v.GetPost(thread, func(new_post *forum_types.Post) {
			// Insert code to push post to channels
		})
	}

	_ = channelIDs
}

func (v *Forumer) Run() {
	for {
		v.update()
		time.Sleep(time.Minute)
	}
}
