package farcaster

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	. "github.com/Hubmakerlabs/hoover/pkg"
	ao "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	pb "github.com/juiceworks/hubble-grpc"
)

func BundlerKind(msg *pb.MessageData) (s string) {
	// Get the type of the message
	switch msg.GetType() {
	case pb.MessageType_MESSAGE_TYPE_CAST_ADD:
		return Post
	case pb.MessageType_MESSAGE_TYPE_REACTION_ADD:
		if msg.GetReactionBody().GetType() == pb.ReactionType_REACTION_TYPE_RECAST {
			return Repost
		} else {
			return Like
		}
	case pb.MessageType_MESSAGE_TYPE_LINK_ADD:
		return Follow
	case pb.MessageType_MESSAGE_TYPE_USER_DATA_ADD:
		return Profile
	}
	return
}

func MessageToBundleItem(msg *pb.Message) (bundle *types.BundleItem, err error) {
	// Get the kind of the message
	var k string
	if k = BundlerKind(msg.GetData()); k == "" {
		return
	}
	bundle = &types.BundleItem{}
	var data *ao.EventData
	bundle.Tags = []types.Tag{
		{Name: Protocol, Value: Farcaster},
		{Name: Kind, Value: k},
		{Name: J(Event, Id), Value: hex.EncodeToString(msg.GetHash())},
		{Name: J(User, Id), Value: fmt.Sprintf("%d", msg.GetData().GetFid())},
		{Name: Timestamp, Value: fmt.Sprintf("%d", msg.GetData().GetTimestamp())},
		{Name: Signature, Value: hex.EncodeToString(msg.GetSignature())},
	}
	switch k {
	case Post:
		data = ao.NewEventData(msg.GetData().GetCastAddBody().GetText())
		embeds_deprecated := msg.GetData().GetCastAddBody().GetEmbedsDeprecated()
		for i := range embeds_deprecated {
			embed := embeds_deprecated[i]
			ao.AppendTag(bundle, J(Embed, Uri), embed)
		}
		mentions := msg.GetData().GetCastAddBody().GetMentions()
		if mentions != nil {
			if len(mentions) == 1 {
				data.Append(J(Mention), strconv.FormatUint(mentions[0], 10))
			} else {
				for i := range mentions {
					mention := mentions[i]
					data.Append(J(Mention, i), strconv.FormatUint(mention, 10))
				}
			}

		}
		parent := msg.GetData().GetCastAddBody().GetParent()
		if parent != nil {
			if x, ok := parent.(*pb.CastAddBody_ParentCastId); ok {
				data.Append(J(Reply, Parent, Id), x.ParentCastId.String())
			} else {
				data.Append(J(Reply, Parent, Uri), parent.(*pb.CastAddBody_ParentUrl).ParentUrl)
			}
		}
		mentions_positions := msg.GetData().GetCastAddBody().GetMentionsPositions()
		if mentions_positions != nil {
			if len(mentions_positions) == 1 {
				data.Append(J(Mention, Start), strconv.FormatUint(uint64(mentions_positions[0]), 10))
			} else {
				for i := range mentions_positions {
					mention_position := mentions_positions[i]
					data.Append(J(Mention, i, Start), strconv.FormatUint(uint64(mention_position), 10))
				}
			}

		}
		embeds := msg.GetData().GetCastAddBody().GetEmbeds()
		for i := range embeds {
			if embeds[i].GetUrl() != "" {
				data.Append(J(Embed, Uri), embeds[i].GetUrl())
			} else {
				data.Append(J(Embed, Id), embeds[i].GetCastId().String())
			}
		}

	case Repost:
		target := msg.GetData().GetReactionBody().GetTarget()
		if x, ok := target.(*pb.ReactionBody_TargetCastId); ok {
			target_id := x.TargetCastId.String()
			data = ao.NewEventData(target_id)
			ao.AppendTag(bundle, J(Repost, Event, Id), target_id)
		} else {
			target_url := target.(*pb.ReactionBody_TargetUrl).TargetUrl
			data = ao.NewEventData(target_url)
			ao.AppendTag(bundle, J(Repost, Event, Uri), target_url)
		}
	case Like:
		target := msg.GetData().GetReactionBody().GetTarget()
		if x, ok := target.(*pb.ReactionBody_TargetCastId); ok {
			target_id := x.TargetCastId.String()
			data = ao.NewEventData(target_id)
			ao.AppendTag(bundle, J(Like, Event, Id), target_id)
		} else {
			target_url := target.(*pb.ReactionBody_TargetUrl).TargetUrl
			data = ao.NewEventData(target_url)
			ao.AppendTag(bundle, J(Like, Event, Uri), target_url)
		}

	case Follow:
		follow_id := strconv.FormatUint(msg.GetData().GetLinkBody().GetTargetFid(), 10)
		data = ao.NewEventData(follow_id)
		ao.AppendTag(bundle, J(Follow, User, Id), follow_id)
	case Profile:
		data = ao.NewEventData(msg.GetData().GetUserDataBody().GetValue())
	}
	// put the ao.EventData into JSON form and place in the bundle.Data field
	var b []byte
	b, err = json.Marshal(data)
	if err != nil {
		return
	}
	bundle.Data = string(b)
	return
}
