// Package bluesky is a set of functions for pulling events from bluesky and
// sending them to Arweave AO.
//
// <insert more notes here>
package bluesky

var Kinds = map[string]string{
	"like":"app.bsky.feed.like",
	"post":"app.bsky.feed.post",
	"follow":"app.bsky.graph.follow",
	"repost":"app.bsky.feed.repost",
	"block":"app.bsky.graph.block",
	"profile":"app.bsky.actor.profile",
	"list":"app.bsky.graph.list",
	"listitem":"app.bsky.graph.listitem",
	"listblock":"app.bsky.graph.listblock",
}

func IsRelevant(kind S) (is bool) {
	for i := range Kinds {
		if kind == Kinds[i] {
			is = true
			break
		}
	}
	return
}
