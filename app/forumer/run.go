package forumer

import (
	"darkbot/app/configurator"
	"darkbot/app/discorder"
	"darkbot/app/settings/types"
	"time"
)

type Forumer struct {
	Discorder discorder.Discorder
	*configurator.Configurators
	threads_requester     *ThreadsRequester
	detailed_post_request *PostRequester
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
		Discorder:             discorder.NewClient(),
		Configurators:         configurator.NewConfigugurators(dbpath),
		threads_requester:     NewLatestThreads(),
		detailed_post_request: NewDetailedPostRequester(),
	}

	for _, opt := range opts {
		opt(forum)
	}

	return forum
}

func (v *Forumer) update() {

	channelIDs, _ := v.Channels.List()
	_ = channelIDs
}

func (v *Forumer) Run() {
	for {
		v.update()
		time.Sleep(time.Minute)
	}
}
