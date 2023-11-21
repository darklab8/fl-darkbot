package forumer

import (
	"crypto/md5"
	"darkbot/app/configurator"
	"darkbot/app/discorder"
	"darkbot/app/forumer/forum_types"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"fmt"
	"strings"
	"time"

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
			for _, channel := range channelIDs {
				watch_tags, err := v.Forum.Watch.TagsList(channel)
				if logus.CheckDebug(err, "failed to get watch tags") {
					continue
				}

				ignore_tags, err := v.Forum.Ignore.TagsList(channel)
				logus.CheckDebug(err, "failed to get ignore tags")

				do_we_show_this_post := false
				for _, watch_tag := range watch_tags {
					if strings.Contains(string(new_post.ThreadFullName), string(watch_tag)) {
						do_we_show_this_post = true
						break
					}
				}

				for _, ignore_tag := range ignore_tags {
					if strings.Contains(string(new_post.ThreadFullName), string(ignore_tag)) {
						do_we_show_this_post = false
						break
					}
				}

				if !do_we_show_this_post {
					continue
				}

				pingMessage := configurator.GetPingingMessage(channel, v.Configurators, v.Discorder)

				duplication_checker := discorder.NewDeduplicator(func(msgs []*discorder.DiscordMessage) bool {
					for _, msg := range msgs {
						content := msg.Content
						for _, embed := range msg.Embeds {
							content += embed.Description
						}

						if strings.Contains(content, string(new_post.PostID)) &&
							strings.Contains(content, string(new_post.ThreadID)) {
							return true
						}
					}
					return false
				})
				v.Discorder.SendDeduplicatedMsg(
					duplication_checker, channel, func(channel types.DiscordChannelID, dg *discordgo.Session) error {

						embed := &discordgo.MessageEmbed{}
						embed.Title = `✉️  You've got mail`

						// embed.Timestamp = string()
						var content strings.Builder
						content.WriteString(fmt.Sprintf("%s\n", pingMessage))
						content.WriteString(fmt.Sprintf("New post in [%s](<%s>)\n", new_post.ThreadFullName, new_post.PostPermamentLink))

						embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
							Name:   "Topic started by",
							Value:  fmt.Sprintf("[%s](<%s>)", new_post.PostAuthorName, new_post.PostAuthorLink),
							Inline: true,
						})
						embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
							Name:   "Timestamp",
							Value:  string(new_post.LastUpdated),
							Inline: true,
						})
						content.WriteString(fmt.Sprintf("```%s```\n", new_post.PostContent[:600]))
						embed.Description = content.String()

						purple_color := 10181046
						embed.Color = purple_color
						msg, err := dg.ChannelMessageSendEmbed(string(channel), embed)
						logus.CheckError(err, "failed sending msg")
						_ = msg
						return nil
					})
			}
		})
	}
}

func (v *Forumer) Run() {
	for {
		v.update()
		time.Sleep(time.Minute)
	}
}
