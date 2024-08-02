// Package bluesky is a set of functions for pulling events from bluesky and
// sending them to Arweave AO.
//
// <insert more notes here>
package bluesky

var Kinds = []string{
	"app.bsky.graph.follow",
	"app.bsky.feed.like",
	"app.bsky.feed.repost",
	"app.bsky.feed.post",
	"app.bsky.graph.block",
	"app.bsky.actor.profile",
	"app.bsky.feed.threadgate",
	"app.bsky.graph.listitem",
	"app.bsky.graph.list",
	"app.bsky.graph.starterpack",
	"app.bsky.graph.listblock",
	"app.bsky.feed.generator",
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
