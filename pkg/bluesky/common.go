package bluesky

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/util"

	. "github.com/Hubmakerlabs/hoover/pkg"
	ao "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
)

func GetCommon(
	bundle *types.BundleItem,
	rr *repo.Repo,
	createdAt time.Time,
	op Op,
	evt Ev,
	resolv *Resolver,
	c context.Context,
) (userid, protocol, timestamp string, err error) {
	split := strings.Split(op.Path, "/")
	if len(split) < 1 {
		return "", "", "", fmt.Errorf("invalid Op.Path: '%s'", op.Path)
	}
	kind := BskyKinds(split[0])
	if kind == "" {
		return "", "", "", fmt.Errorf("invalid Op.Path kind: '%s'", kind)
	}

	userid = rr.SignedCommit().Did
	var pubkey string
	pubkey, err = resolv.Find(userid, c)
	chk.E(err)
	protocol = Bsky
	timestamp = strconv.FormatInt(createdAt.Unix(), 10)

	ao.AppendTag(bundle, J(App, Name), AppNameValue)
	ao.AppendTag(bundle, J(App, Version), AppVersion)
	ao.AppendTag(bundle, Protocol, protocol)
	ao.AppendTag(bundle, Repository, evt.Repo)
	ao.AppendTag(bundle, Kind, kind)
	ao.AppendTag(bundle, J(Event, Id), op.Cid.String())
	ao.AppendTag(bundle, J(User, Id), userid)
	if pubkey != "" {
		ao.AppendTag(bundle, Signer, pubkey)
	}
	ao.AppendTag(bundle, J(Unix, Time), timestamp)
	ao.AppendTag(bundle, Path, op.Path)
	ao.AppendTag(bundle, Signature, hex.EncodeToString(rr.SignedCommit().Sig))
	ao.AppendTag(bundle, J(Signature, Type), fmt.Sprintf("%d", 4))
	ao.AppendTag(bundle, Topic, Bsky)
	ao.AppendTag(bundle, Topic, kind)
	if kind == Profile {
		ao.AppendTag(bundle, Type, ProfileType)
	} else {
		ao.AppendTag(bundle, Type, PostType)
	}
	return
}

// UnmarshalEvent accepts a bsky commit event and a record type, and the concrete type
func UnmarshalEvent(evt Ev, rec Rec, to any) (decoded any, createdAt time.Time, err error) {
	if rec == nil {
		err = errorf.E("nil record, cannot unmarshal")
		return
	}
	banana := lexutil.LexiconTypeDecoder{
		Val: rec,
	}
	var b B
	if b, err = banana.MarshalJSON(); chk.E(err) {
		return
	}
	if err = json.Unmarshal(b, to); chk.E(err) {
		return
	}
	if createdAt, err = time.Parse(util.ISO8601, evt.Time); chk.E(err) {
		return
	}
	return to, createdAt, nil
}

func AppendLexBlobTagsWithoutRef(data *ao.EventData, prefix string, img *lexutil.LexBlob) {
	data.Append(J(prefix, Mimetype), img.MimeType)
	data.Append(J(prefix, Size), strconv.FormatInt(img.Size, 10))
}

func AppendLexBlobTags(data *ao.EventData, prefix string, img *lexutil.LexBlob) {
	// added some iffies due receiving some null pointer panics
	if img != nil {
		if img.Ref != (lexutil.LexLink{}) && img.Ref.String() != "" {
			data.Append(J(prefix, Ref), img.Ref.String())
		}
		if img.MimeType != "" {
			data.Append(J(prefix, Mimetype), img.MimeType)
		}
		if img.Size != 0 {
			data.Append(J(prefix, Size), strconv.FormatInt(img.Size, 10))
		}

	}

}

func AppendImageTags(data *ao.EventData, prefix string,
	img *bsky.EmbedImages_Image) {
	if img.Alt != "" {
		data.Append(J(prefix, Alt), img.Alt)
	}
	AppendLexBlobTags(data, prefix, img.Image)
	if img.AspectRatio != nil {
		data.Append(J(prefix, Aspect),
			fmt.Sprintf("%dx%d", img.AspectRatio.Width, img.AspectRatio.Height))
	}
}

func EmbedRecord(data *ao.EventData, prefix string, record *bsky.EmbedRecord) {
	if record.Record != nil {
		data.Append(J(prefix, Id), record.Record.Cid)
		data.Append(J(prefix, Uri), record.Record.Uri)
	}
}

func EmbedExternal(data *ao.EventData, prefix string, embed *bsky.EmbedExternal) {
	ext := embed.External
	data.Append(J(prefix, Description), ext.Description)
	if ext.Thumb != nil {
		AppendLexBlobTags(data, J(prefix, Thumb), ext.Thumb)
	}
	data.Append(J(prefix, Title), ext.Title)
	data.Append(J(prefix, Uri), ext.Uri)
}

// EmbedExternalRecordWithMedia creates a tag with all the junk in an EmbedRecordWithMedia into one.
// It makes an extremely long tag field but this is the retardation of the bluesky API.
func EmbedExternalRecordWithMedia(data *ao.EventData, prefix string,
	embed *bsky.EmbedRecordWithMedia) {
	if embed.Record != nil {
		EmbedRecord(data, J(prefix, Record), embed.Record)
	}
	if embed.Media != nil {
		if embed.Media.EmbedImages != nil {
			if embed.Media.EmbedImages.Images != nil {
				if len(embed.Media.EmbedImages.Images) == 1 {
					AppendImageTags(data, J(prefix), embed.Media.EmbedImages.Images[0])
				} else {
					for i, img := range embed.Media.EmbedImages.Images {
						AppendImageTags(data, J(prefix, i), img)
					}
				}
			}
		}
		if embed.Media.EmbedExternal != nil && embed.Media.EmbedExternal.External != nil {
			ext := embed.Media.EmbedExternal.External
			data.Append(J(prefix, Uri), ext.Uri)
			AppendLexBlobTags(data, prefix, ext.Thumb)
			data.Append(J(prefix, Description), ext.Description)
			data.Append(J(prefix, Title), ext.Title)
		}
	}
}
