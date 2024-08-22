package bluesky

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/util"
)

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

func AppendTags(bundle *types.BundleItem, name S, values []S) {
	var b B
	var err error
	if b, err = json.Marshal(values); chk.E(err) {
		return
	}
	tag := types.Tag{Name: name, Value: S(b)}
	bundle.Tags = append(bundle.Tags, tag)
}

func GetImageTags(img *bsky.EmbedImages_Image) []string {
	tags := []string{
		img.Alt,
	}
	if img.AspectRatio != nil {
		tags = append(tags, fmt.Sprintf("%dx%d", img.AspectRatio.Width, img.AspectRatio.Height))
	}
	tags = append(tags, GetLexBlobTags(img.Image)...)
	return tags
}

func GetLexBlobTags(img *lexutil.LexBlob) (tags []string) {
	tags = []string{
		img.Ref.String(),
		img.MimeType,
		fmt.Sprintf("%d", img.Size)}
	return
}

func EmbedImages(bundle *types.BundleItem, name string, img *bsky.EmbedImages_Image) {
	imgTags := GetImageTags(img)
	for i := range imgTags {
		if imgTags[i]!=""{
			AppendTag(bundle, fmt.Sprintf("%s-%s-%03d", name, pkg.Tag, i), imgTags[i])
		}
	}
	// AppendTags(bundle, name, GetImageTags(img))
}

func EmbedRecord(bundle *types.BundleItem, name string, record *bsky.EmbedRecord) {
	AppendTags(bundle, name, []string{"record", record.Record.Cid, record.Record.Uri})
}

func EmbedExternal(bundle *types.BundleItem, name string, embed *bsky.EmbedExternal) {
	ext := embed.External
	imgTags := []string{ext.Uri, ext.Title, ext.Description}
	if ext.Thumb != nil {
		thumbTags := GetLexBlobTags(ext.Thumb)
		imgTags = append(imgTags, thumbTags...)
	}
	AppendTags(bundle, fmt.Sprintf("%s-%s", name, pkg.External), imgTags)
}

// EmbedExternalRecordWithMedia creates a tag with all the junk in an EmbedRecordWithMedia into one.
// It makes an extremely long tag field but this is the retardation of the bluesky API.
func EmbedExternalRecordWithMedia(bundle *types.BundleItem, name string,
	embed *bsky.EmbedRecordWithMedia) {
	if embed.Record != nil {
		EmbedRecord(bundle, name+"_record", embed.Record)
	}
	if embed.Media != nil {
		if embed.Media.EmbedImages != nil {
			if embed.Media.EmbedImages.Images != nil {
				for i, img := range embed.Media.EmbedImages.Images {
					if img == nil {
						continue
					}
					var tags []string
					var lbtags []string
					if img.Image != nil {
						lbtags = GetLexBlobTags(img.Image)
					}
					log.I.S(img)
					tags = []string{
						fmt.Sprintf("image%03d", i),
					}
					if lbtags != nil {
						tags = append(tags, lbtags...)
					}
					tags = append(tags, img.Alt)
					if img.AspectRatio != nil {
						tags = append(tags,
							fmt.Sprintf("%dx%d", img.AspectRatio.Width, img.AspectRatio.Height))
					}
					AppendTags(bundle, fmt.Sprintf("%s%s%03d", name, "_image", i), tags)
				}
			}
		}
		if embed.Media.EmbedExternal != nil && embed.Media.EmbedExternal.External != nil {
			ext := embed.Media.EmbedExternal.External
			imgTags := []string{ext.Uri, ext.Title, ext.Description}
			if ext.Thumb != nil {
				thumbTags := GetLexBlobTags(ext.Thumb)
				imgTags = append(imgTags, thumbTags...)
			}
			AppendTags(bundle, fmt.Sprintf("%s%s", name, "_media_external"), imgTags)
		}
	}
}
