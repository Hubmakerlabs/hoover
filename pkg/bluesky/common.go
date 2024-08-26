package bluesky

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	. "github.com/Hubmakerlabs/hoover/pkg"
	ao "github.com/Hubmakerlabs/hoover/pkg/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/util"
)

func GetCommon(bundle *types.BundleItem, rr *repo.Repo, createdAt Time, op Op,
	evt Ev) (err error) {
	split := strings.Split(op.Path, "/")
	if len(split) < 1 {
		return fmt.Errorf("invalid Op.Path: '%s'", op.Path)
	}
	k := BskyKinds(split[0])
	if k == "" {
		return fmt.Errorf("invalid Op.Path kind: '%s'", k)
	}
	ao.AppendTag(bundle, Protocol, Bsky)
	ao.AppendTag(bundle, Kind, k)
	ao.AppendTag(bundle, J(Event, Id), op.Cid.String())
	ao.AppendTag(bundle, J(User, Id), rr.SignedCommit().Did)
	ao.AppendTag(bundle, Timestamp, strconv.FormatInt(createdAt.Unix(), 10))
	ao.AppendTag(bundle, Repository, evt.Repo)
	ao.AppendTag(bundle, Path, op.Path)
	ao.AppendTag(bundle, Signature, hex.EncodeToString(rr.SignedCommit().Sig))
	return
}

// UnmarshalEvent accepts a bsky commit event and a record type, and the concrete type
func UnmarshalEvent(evt Ev, rec Rec, to any) (decoded any, createdAt Time, err error) {
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

func AppendLexBlobTags(bundle *types.BundleItem, prefix string, img *lexutil.LexBlob) {
	ao.AppendTag(bundle, J(prefix, Ref), img.Ref.String())
	ao.AppendTag(bundle, J(prefix, Mimetype), img.MimeType)
	ao.AppendTag(bundle, J(prefix, Size), strconv.FormatInt(img.Size, 10))
}

func AppendImageTags(bundle *types.BundleItem, prefix string,
	img *bsky.EmbedImages_Image) {
	if img.Alt != "" {
		ao.AppendTag(bundle, J(prefix, Alt), img.Alt)
	}
	AppendLexBlobTags(bundle, prefix, img.Image)
	if img.AspectRatio != nil {
		ao.AppendTag(bundle, J(prefix, Aspect),
			fmt.Sprintf("%dx%d", img.AspectRatio.Width, img.AspectRatio.Height))
	}
}

func EmbedRecord(bundle *types.BundleItem, prefix string, record *bsky.EmbedRecord) {
	if record.Record != nil {
		ao.AppendTag(bundle, J(prefix, Id), record.Record.Cid)
		ao.AppendTag(bundle, J(prefix, Uri), record.Record.Uri)
	}
}

func EmbedExternal(bundle *types.BundleItem, prefix string, embed *bsky.EmbedExternal) {
	ext := embed.External
	ao.AppendTag(bundle, J(prefix, Description), ext.Description)
	if ext.Thumb != nil {
		AppendLexBlobTags(bundle, J(prefix, Thumb), ext.Thumb)
	}
	ao.AppendTag(bundle, J(prefix, Title), ext.Title)
	ao.AppendTag(bundle, J(prefix, Uri), ext.Uri)
}

// EmbedExternalRecordWithMedia creates a tag with all the junk in an EmbedRecordWithMedia into one.
// It makes an extremely long tag field but this is the retardation of the bluesky API.
func EmbedExternalRecordWithMedia(bundle *types.BundleItem, prefix string,
	embed *bsky.EmbedRecordWithMedia) {
	if embed.Record != nil {
		EmbedRecord(bundle, J(prefix, Record), embed.Record)
	}
	if embed.Media != nil {
		if embed.Media.EmbedImages != nil {
			if embed.Media.EmbedImages.Images != nil {
				if len(embed.Media.EmbedImages.Images) == 1 {
					AppendImageTags(bundle, J(prefix), embed.Media.EmbedImages.Images[0])
				} else {
					for i, img := range embed.Media.EmbedImages.Images {
						AppendImageTags(bundle, J(prefix, i), img)
					}
				}
			}
		}
		if embed.Media.EmbedExternal != nil && embed.Media.EmbedExternal.External != nil {
			ext := embed.Media.EmbedExternal.External
			ao.AppendTag(bundle, J(prefix, Uri), ext.Uri)
			AppendLexBlobTags(bundle, prefix, ext.Thumb)
			ao.AppendTag(bundle, J(prefix, Description), ext.Description)
			ao.AppendTag(bundle, J(prefix, Title), ext.Title)
		}
	}
}
