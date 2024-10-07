package pkg

import (
	"fmt"
	"strings"
)

// Constants defining the standard tag field names

// J joins two strings in HTTP-Header-Case, aka Train-Case
func J(s ...any) string {
	var ss []string
	for _, v := range s {
		if vs, ok := v.(string); ok {
			ss = append(ss, vs)
		} else if vi, ok := v.(int); ok {
			ss = append(ss, fmt.Sprintf("%d", vi))
		} else {
			ss = append(ss, fmt.Sprintf("%v", v))
		}
	}
	return strings.Join(ss, "-")
}

// Common

var (
	App          = "App"
	Version      = "Version"
	AppNameValue = "com.hubmakerlabs.hoover"
	AppVersion   = "v1.0.0"
	Un           = "Un" // todo for future use for delete/unfollow/unblock/unlike
	Id           = "Id"
	Event        = "Event"
	Protocol     = "Protocol"
	User         = "User"
	Unix         = "Unix"
	Time         = "Time"
	Timestamp    = "Timestamp"
	Kind         = "Kind"
	Repository   = "Repository"
	Path         = "Path"
	Signature    = "Signature"
	Label        = "Label"
	Namespace    = "Namespace"
)

var (
	// Protocols
	Nostr     = "Nostr"
	Bsky      = "Bluesky"
	Farcaster = "Farcaster"

	// Kinds
	Post    = "Post"
	Repost  = "Repost"
	Like    = "Like"
	Follow  = "Follow"
	Block   = "Block"
	Profile = "Profile"

	// Posts
	Text        = "Text"
	Richtext    = "Richtext"
	Image       = "Image"
	Embed       = "Embed"
	Alt         = "Alt"
	Ref         = "Ref"
	Facet       = "Facet"
	Mimetype    = "Mimetype"
	Aspect      = "Aspect"
	Title       = "Title"
	Description = "Description"
	External    = "External"
	Record      = "Record"
	Entities    = "Entities"
	Language    = "Language"
	Index       = "Index"
	Type        = "Type"
	Value       = "Value"
	Link        = "Link"
	Start       = "Start"
	End         = "End"
	Media       = "Media"
	Links       = "Links"
	EmbedCid    = "EmbedCid"
	EmbedURI    = "EmbedURI"
	Mention     = "Mention"
	Source      = "Source"
	Hashtag     = "Hashtag"
	Emoji       = "Emoji"
	Content     = "Content"
	Warning     = "Warning"
	Reply       = "Reply"
	Parent      = "Parent"
	Root        = "Root"
	Tag         = "Tag"
	Thumb       = "Thumb"

	// Profiles
	Bio          = "Bio"
	Name         = "Name"
	Display      = "Display"
	Avatar       = "Avatar"
	Banner       = "Banner"
	Website      = "Website"
	Verification = "Verification"
	Payment      = "Payment"
	Address      = "Address"

	// The rest
	Add = "Add"

	// Embeds
	Reference  = "Reference"
	Origin     = "Origin"
	Uri        = "Uri"
	Size       = "Size"
	Dimensions = "Dimensions"
	Duration   = "Duration"
)
