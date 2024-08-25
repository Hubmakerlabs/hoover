package bluesky

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	. "github.com/Hubmakerlabs/hoover/pkg"
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
	AppendTag(bundle, Protocol, Bsky)
	AppendTag(bundle, Kind, k)
	AppendTag(bundle, J(Event, Id), op.Cid.String())
	AppendTag(bundle, J(User, Id), rr.SignedCommit().Did)
	AppendTag(bundle, Timestamp, strconv.FormatInt(createdAt.Unix(), 10))
	AppendTag(bundle, Repository, evt.Repo)
	AppendTag(bundle, Path, op.Path)
	AppendTag(bundle, Signature, hex.EncodeToString(rr.SignedCommit().Sig))
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

func AppendTag(bundle *types.BundleItem, name, value string) {
	bundle.Tags = append(bundle.Tags, types.Tag{Name: name, Value: value})
}

func AppendLexBlobTags(bundle *types.BundleItem, prefix string, img *lexutil.LexBlob) {
	AppendTag(bundle, J(prefix, Ref), img.Ref.String())
	AppendTag(bundle, J(prefix, Mimetype), img.MimeType)
	AppendTag(bundle, J(prefix, Size), strconv.FormatInt(img.Size, 10))
}

func AppendImageTags(bundle *types.BundleItem, prefix string,
	img *bsky.EmbedImages_Image) {
	if img.Alt != "" {
		AppendTag(bundle, J(prefix, Alt), img.Alt)
	}
	AppendLexBlobTags(bundle, prefix, img.Image)
	if img.AspectRatio != nil {
		AppendTag(bundle, J(prefix, Aspect),
			fmt.Sprintf("%dx%d", img.AspectRatio.Width, img.AspectRatio.Height))
	}
}

func EmbedRecord(bundle *types.BundleItem, prefix string, record *bsky.EmbedRecord) {
	if record.Record != nil {
		AppendTag(bundle, J(prefix, Id), record.Record.Cid)
		AppendTag(bundle, J(prefix, Uri), record.Record.Uri)
	}
}

func EmbedExternal(bundle *types.BundleItem, prefix string, embed *bsky.EmbedExternal) {
	ext := embed.External
	AppendTag(bundle, J(prefix, Description), ext.Description)
	if ext.Thumb != nil {
		AppendLexBlobTags(bundle, J(prefix, Thumb), ext.Thumb)
	}
	AppendTag(bundle, J(prefix, Title), ext.Title)
	AppendTag(bundle, J(prefix, Uri), ext.Uri)
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
			AppendTag(bundle, J(prefix, Uri), ext.Uri)
			AppendLexBlobTags(bundle, prefix, ext.Thumb)
			AppendTag(bundle, J(prefix, Description), ext.Description)
			AppendTag(bundle, J(prefix, Title), ext.Title)
		}
	}
}
