package farcaster

import (
	"encoding/hex"
	"fmt"

	. "github.com/Hubmakerlabs/hoover/pkg"
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
	case Repost:
	case Like:
	case Follow:
	case Profile:
	}
	return
}
