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
	data := ao.NewEventData("")
	bundle.Tags = []types.Tag{
		{Name: J(App, Name), Value: AppNameValue},
		{Name: J(App, Version), Value: AppVersion},
		{Name: Protocol, Value: Farcaster},
		{Name: Kind, Value: k},
		{Name: J(Event, Id), Value: hex.EncodeToString(msg.GetHash())},
		{Name: J(User, Id), Value: fmt.Sprintf("%d", msg.GetData().GetFid())},
		{Name: J(Unix, Time), Value: fmt.Sprintf("%d",
			int64(msg.GetData().GetTimestamp())+1609459200)}, // offset from 2021-01-10T00:00
		{Name: Signature, Value: hex.EncodeToString(msg.GetSignature())},
		{Name: "Forward-For", Value: hex.EncodeToString(msg.GetSigner())},
		{Name: J(Signature, Type), Value: fmt.Sprintf("%d", msg.GetSignatureScheme())},
	}

	switch k {
	case Post:
		data.Content = msg.GetData().GetCastAddBody().GetText()
		embeds_deprecated := msg.GetData().GetCastAddBody().GetEmbedsDeprecated()
		for i := range embeds_deprecated {
			embed := embeds_deprecated[i]
			data.Append(J(Embed, Uri), embed)
		}
		mentions := msg.GetData().GetCastAddBody().GetMentions()
		if mentions != nil {
			for i := range mentions {
				mention := mentions[i]
				data.Append(Mention, strconv.FormatUint(mention, 10))
			}
		}
		parent := msg.GetData().GetCastAddBody().GetParent()
		if parent != nil {
			if x, ok := parent.(*pb.CastAddBody_ParentCastId); ok {
				fid := x.ParentCastId.Fid
				hash := x.ParentCastId.Hash
				ao.AppendTag(bundle, J(Reply, Parent, Id), fmt.Sprintf("%0x", hash))
				data.Append(J(Reply, Parent, User, Id), fmt.Sprintf("%d", fid))
			}
			if x, ok := parent.(*pb.CastAddBody_ParentUrl); ok {
				ao.AppendTag(bundle, J(Reply, Parent, Uri), x.ParentUrl)
			}
			mentions_positions := msg.GetData().GetCastAddBody().GetMentionsPositions()
			if mentions_positions != nil {
				if len(mentions_positions) == 1 {
					data.Append(J(Mention, Start),
						strconv.FormatUint(uint64(mentions_positions[0]), 10))
				} else {
					for i := range mentions_positions {
						mention_position := mentions_positions[i]
						data.Append(J(Mention, i, Start),
							strconv.FormatUint(uint64(mention_position), 10))
					}
				}

			}
			embeds := msg.GetData().GetCastAddBody().GetEmbeds()
			for i := range embeds {
				if embeds[i].GetUrl() != "" {
					data.Append(J(Embed, i, Uri), embeds[i].GetUrl())
				} else {
					fid := fmt.Sprintf("%d", embeds[i].GetCastId().Fid)
					hash := fmt.Sprintf("%0x", embeds[i].GetCastId().Hash)
					data.Append(J(Embed, i, User, Id), fid)
					data.Append(J(Embed, i, Event, Id), hash)
				}
			}
		}

	case Repost:
		target := msg.GetData().GetReactionBody().GetTarget()
		if x, ok := target.(*pb.ReactionBody_TargetCastId); ok {
			targetFid := x.TargetCastId.Fid
			targetHash := x.TargetCastId.Hash
			ao.AppendTag(bundle, J(Repost, Event, Id),
				fmt.Sprintf("%0x", targetHash))
			data.Append(J(Repost, User, Id),
				fmt.Sprintf("%d", targetFid))
		} else {
			target_url := target.(*pb.ReactionBody_TargetUrl).TargetUrl
			data.Append(J(Repost, Event, Uri), target_url)
		}

	case Like:
		target := msg.GetData().GetReactionBody().GetTarget()
		if x, ok := target.(*pb.ReactionBody_TargetCastId); ok {
			targetFid := x.TargetCastId.Fid
			targetHash := x.TargetCastId.Hash
			ao.AppendTag(bundle, J(Like, Event, Id),
				fmt.Sprintf("%0x", targetHash))
			data.Append(J(Like, User, Id), fmt.Sprint(targetFid))
		} else {
			target_url := target.(*pb.ReactionBody_TargetUrl).TargetUrl
			data.Append(J(Like, Event, Uri), target_url)
		}

	case Follow:
		follow_id := fmt.Sprintf("%d", msg.GetData().GetLinkBody().GetTargetFid())
		ao.AppendTag(bundle, J(Follow, User, Id), follow_id)

	case Profile:
		//add data to EventData in all cases because there is a None user data type as well
		data.Content = msg.GetData().GetUserDataBody().GetValue()
		switch msg.GetData().GetUserDataBody().GetType() {
		case pb.UserDataType_USER_DATA_TYPE_PFP:
			ao.AppendTag(bundle, J(Avatar, Image), data.Content)
		case pb.UserDataType_USER_DATA_TYPE_DISPLAY:
			ao.AppendTag(bundle, J(Display, Name), data.Content)
		case pb.UserDataType_USER_DATA_TYPE_BIO:
			data.Append(Bio, data.Content)
		case pb.UserDataType_USER_DATA_TYPE_URL:
			data.Append(Website, data.Content)
		case pb.UserDataType_USER_DATA_TYPE_USERNAME:
			ao.AppendTag(bundle, J(User, Name), data.Content)
		}

	}

	if data != nil && (data.Content != "" || len(data.EventTags) > 0) {
		// put the ao.EventData into JSON form and place in the bundle.Data field
		var b []byte
		b, err = json.Marshal(data)
		if err != nil {
			return
		}
		bundle.Data = string(b)
	}
	return
}
