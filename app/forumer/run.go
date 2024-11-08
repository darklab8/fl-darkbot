package forumer

import (
	"crypto/md5"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/discorder"
	"github.com/darklab8/fl-darkbot/app/forumer/forum_types"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"

	"github.com/darklab8/go-utils/utils"

	"github.com/bwmarrin/discordgo"
)

type iThreadsRequester interface {
	GetLatestThreads(opts ...threadPageParam) ([]*forum_types.LatestThread, error)
}
type Forumer struct {
	Discorder *discorder.Discorder
	*configurator.Configurators
	threads_requester iThreadsRequester
	post_requester    *PostRequester

	cache map[ThreadCacheKey]*forum_types.Post
	// Keeping track as list, to realize which ones are old ones to delete
	cache_keys []ThreadCacheKey
	cache_mu   sync.Mutex
}

type ThreadCacheKey string

func NewThreadCacheKey(thread *forum_types.LatestThread) ThreadCacheKey {
	future_key_as_str := string(thread.LastUpdated) + string(thread.ThreadID) + string(thread.PostAuthorName)
	hash := md5.Sum([]byte(future_key_as_str))
	return ThreadCacheKey(string(hash[:]))
}

type forumerParam func(forum *Forumer)

func WithThreadsRequester(
	threads_page_requester iThreadsRequester) forumerParam {
	return func(forum *Forumer) { forum.threads_requester = threads_page_requester }
}
func WithDetailedPostRequest(post_requester *PostRequester) forumerParam {
	return func(forum *Forumer) { forum.post_requester = post_requester }
}

func NewForumer(dbpath types.Dbpath, opts ...forumerParam) *Forumer {

	forum := &Forumer{
		Discorder:         discorder.NewClient(),
		Configurators:     configurator.NewConfigugurators(dbpath),
		threads_requester: NewLatestThreads(),
		post_requester:    NewDetailedPostRequester(),
		cache:             make(map[ThreadCacheKey]*forum_types.Post),
	}

	for _, opt := range opts {
		opt(forum)
	}

	return forum
}

func (v *Forumer) GetPost(thread *forum_types.LatestThread, new_post_callback func(*forum_types.Post)) *forum_types.Post {
	var err error
	thread_key := NewThreadCacheKey(thread)

	var (
		post *forum_types.Post
		ok   bool
	)
	v.WithCacheLock(func() {
		post, ok = v.cache[thread_key]
	})
	if !ok {
		logus.Log.Debug("cache is not found. requesting new post for thread", logus.Thread(thread))
		post, err = v.post_requester.GetDetailedPost(thread)
		logus.Log.CheckError(err, "failed get detailed post for thread=", logus.Thread(thread))
		new_post_callback(post)

		v.WithCacheLock(func() {
			v.cache[thread_key] = post
			v.cache_keys = append(v.cache_keys, thread_key)
		})
	}

	v.WithCacheLock(func() {
		if len(v.cache_keys) > 50 {
			key_to_delete := v.cache_keys[0]
			logus.Log.Debug("deleting old cached key_to_delete=" + string(key_to_delete))
			v.cache_keys = append(v.cache_keys[1:])
			delete(v.cache, key_to_delete)
		}
	})
	return post
}

func (v *Forumer) isPostMatchTags(channel types.DiscordChannelID, new_post *forum_types.Post) (bool, []string) {
	do_we_show_this_post := false
	var matched_tags []string

	thread_watch_tags, err := v.Forum.Thread.Watch.TagsList(channel)
	logus.Log.CheckDebug(err, "failed to get watch tags")

	thread_ignore_tags, err := v.Forum.Thread.Ignore.TagsList(channel)
	logus.Log.CheckDebug(err, "failed to get ignore tags")

	subforum_watch_tags, err := v.Forum.Subforum.Watch.TagsList(channel)
	logus.Log.CheckDebug(err, "failed to get watch tags")

	subforum_ignore_tags, err := v.Forum.Subforum.Ignore.TagsList(channel)
	logus.Log.CheckDebug(err, "failed to get ignore tags")

	for _, watch_tag := range thread_watch_tags {
		if strings.Contains(string(new_post.ThreadFullName), string(watch_tag)) ||
			strings.Contains(strings.ToLower(string(new_post.ThreadFullName)), strings.ToLower(string(watch_tag))) {
			do_we_show_this_post = true
			matched_tags = append(matched_tags, string(fmt.Sprintf(`"%s"`, watch_tag)))
		}
	}

	for _, watch_tag := range subforum_watch_tags {
		for _, subforum := range new_post.Subforums {
			if strings.Contains(string(subforum), string(watch_tag)) ||
				strings.Contains(strings.ToLower(string(new_post.ThreadFullName)), strings.ToLower(string(watch_tag))) {
				do_we_show_this_post = true
				matched_tags = append(matched_tags, string(fmt.Sprintf(`"%s"`, watch_tag)))
			}
		}

	}

	for _, ignore_tag := range thread_ignore_tags {
		if strings.Contains(string(new_post.ThreadFullName), string(ignore_tag)) {
			do_we_show_this_post = false
			break
		}
	}

	for _, ignore_tag := range subforum_ignore_tags {
		for _, subforum := range new_post.Subforums {
			if strings.Contains(string(subforum), string(ignore_tag)) {
				do_we_show_this_post = false
				break
			}
		}
	}

	return do_we_show_this_post, matched_tags
}

func CreateDeDuplicator(new_post *forum_types.Post, msgs []*discorder.DiscordMessage) *discorder.Deduplicator {
	return discorder.NewDeduplicator(func() bool {
		for _, msg := range msgs {
			content := msg.Content
			for _, embed := range msg.Embeds {
				content += embed.Description
				content += embed.Title
				content += embed.URL
			}

			if strings.Contains(content, string(new_post.PostPermamentLink)) {
				logus.Log.Debug("Post already exists!", logus.Post(new_post))
				return true
			}
		}
		logus.Log.Debug("Post does not exist like that", logus.Post(new_post))
		return false
	})
}

func (v *Forumer) TrySendMsg(channel types.DiscordChannelID, new_post *forum_types.Post, msgs []*discorder.DiscordMessage) {

	pingMessage := configurator.GetPingingMessage(channel, v.Configurators, v.Discorder)
	is_match, matched_tags := v.isPostMatchTags(channel, new_post)
	if !is_match {
		return
	}

	v.Discorder.SendDeduplicatedMsg(
		CreateDeDuplicator(new_post, msgs), channel, func(channel types.DiscordChannelID, dg *discordgo.Session) error {
			embed := &discordgo.MessageEmbed{}
			embed.Title = string(new_post.ThreadFullName)
			embed.URL = string(new_post.PostPermamentLink)

			// embed.Timestamp = string()
			var content strings.Builder
			content.WriteString(
				fmt.Sprintf("received email from %s\n",
					fmt.Sprintf("[%s](<%s>)", new_post.PostAuthorName, new_post.PostAuthorLink)))
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Matched tags",
				Value:  strings.Join(matched_tags, ", "),
				Inline: true,
			})
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Timestamp",
				Value:  string(new_post.LastUpdated),
				Inline: true,
			})

			subforums := utils.CompL(new_post.Subforums, func(x forum_types.Subforum) string { return string(x) })
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Subforums",
				Value:  strings.Join(subforums, " / "),
				Inline: false,
			})

			var post_content string = string(new_post.PostContent)
			if len(post_content) >= 600 {
				post_content = post_content[:600]
			}
			content.WriteString(fmt.Sprintf("```%s```\n", post_content))
			embed.Description = content.String()

			embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: string(new_post.PostAuthorAvatarLink)}

			purple_color := 10181046
			embed.Color = purple_color

			_, err := dg.ChannelMessageSendComplex(string(channel), &discordgo.MessageSend{Embeds: []*discordgo.MessageEmbed{embed}, Content: fmt.Sprintf("mail for %s", string(pingMessage))})

			// helpful to uncover problems with Embeds sending.
			// There was issue with insufficient granted permission to bot once. Inviting bot as admin helped to fix the issue.
			// _, err := dg.ChannelMessageSendEmbed(string(channel), embed)

			// if u will wish simpler msg
			// var normal_msg strings.Builder
			// normal_msg.WriteString(fmt.Sprintf("%s\n", string(pingMessage)))
			// normal_msg.WriteString(fmt.Sprintf("[%s](%s)\n", string(new_post.ThreadFullName), string(new_post.PostPermamentLink)))
			// normal_msg.WriteString(content.String())
			// normal_msg.WriteString(fmt.Sprintln("**Matched tags**: `", strings.Join(matched_tags, ", "), "`"))
			// normal_msg.WriteString(fmt.Sprintln("**Timestamp**: `", string(new_post.LastUpdated), "`"))
			// normal_msg.WriteString(fmt.Sprintln("**Subforums**: `", strings.Join(subforums, " / "), "`"))
			// _, err := dg.ChannelMessageSend(string(channel), normal_msg.String())

			logus.Log.Debug("sent forumer msg",
				logus.MsgContent(post_content),
				logus.ChannelID(channel),
			)
			logus.Log.CheckError(err, "failed sending msg", logus.ChannelID(channel))

			return nil
		})
}

func (v *Forumer) update() {
	channelIDs, _ := v.Channels.List()

	threads, err := v.threads_requester.GetLatestThreads()
	if logus.Log.CheckError(err, "failed to get threads") {
		return
	}

	for _, thread := range threads {
		v.GetPost(thread, func(new_post *forum_types.Post) {
			for _, channel := range channelIDs {
				msgs, err := v.Discorder.GetLatestMessages(channel)
				if logus.Log.CheckError(err, "failed to get discord latest msgs", logus.ChannelID(channel)) {
					continue
				}
				v.TrySendMsg(channel, new_post, msgs)
			}
		})
	}
}

func (v *Forumer) WithCacheLock(callback func()) {
	v.cache_mu.Lock()
	defer v.cache_mu.Unlock()

	callback()
}

func (v *Forumer) RetryMsgs() {

	var copied_keys []ThreadCacheKey
	v.WithCacheLock(func() {
		copied_keys = append(copied_keys, v.cache_keys...)
		utils.ReverseSlice(copied_keys)
	})

	channelIDs, _ := v.Channels.List()
	for _, channel := range channelIDs {
		msgs, err := v.Discorder.GetLatestMessages(channel)
		if logus.Log.CheckError(err, "failed to get discord latest msgs", logus.ChannelID(channel)) {
			continue
		}

		for _, cache_key := range copied_keys {
			var (
				old_post *forum_types.Post
				ok       bool
			)
			v.WithCacheLock(func() {
				old_post, ok = v.cache[cache_key]
			})

			if !ok {
				continue
			}
			v.TrySendMsg(channel, old_post, msgs)
		}

		time.Sleep(time.Second * 3)
	}

}

func (v *Forumer) Run() {
	delay := time.Second * 60
	go func() {
		for {
			logus.Log.Debug("retrying to send msgs")
			v.RetryMsgs()
			time.Sleep(delay)
		}
	}()

	for {
		logus.Log.Debug("trying new forumer cycle")
		v.update()
		time.Sleep(delay)
	}
}
