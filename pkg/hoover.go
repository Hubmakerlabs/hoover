package pkg

// Constants defining the standard tag field names

// Common

const (
	Protocol       = "Protocol"
	EventId        = "Event-ID"
	UserId         = "User-ID"
	Timestamp      = "Timestamp"
	Kind           = "Kind"
	Path           = "Path"
	Signature      = "Signature"
	Label          = "Label"
	LabelNamespace = "Label-Namespace"
)

// Protocols

const (
	Nostr     = "Nostr"
	Bsky      = "Bluesky"
	Farcaster = "Farcaster"
)

// Kinds

const (
	Post    = "Post"
	Repost  = "Repost"
	Like    = "Like"
	Follow  = "Follow"
	Block   = "Block"
	Profile = "Profile"
)

// Posts

const (
	PostText       = "Post-Text"
	Richtext       = "Richtext"
	Image          = "Image"
	Media          = "Media"
	Links          = "Links"
	Embed          = "Embed"
	Mention        = "Mention"
	ReplyTo        = "Reply-To"
	Source         = "Source"
	Hashtag        = "Hashtag"
	Emoji          = "Emoji"
	ContentWarning = "Content-Warning"
)

// Profile

// About field of Nostr and Bio of other protocols should go in bundle data to allow large bios.
const (
	UserName       = "User-Name"
	DisplayName    = "Display-Name"
	AvatarImage    = "Avatar-Image"
	BannerImage    = "Banner-Image"
	Website        = "Website"
	Verification   = "Verification"
	PaymentAddress = "Payment-Address"
)

// The rest

const (
	RepostEventId = "Repost-Event-ID"
	LikeEventId   = "Like-Event-ID"
	Add           = "Add"
	FollowUserId  = "Follow-User-ID"
	FollowTag     = "Follow-Tag"
	BlockUserId   = "Block-User-ID"
	BlockTag      = "Block-Tag"
)

// Embeds

const (
	Reference  = "Reference"
	MimeType   = "Mime-Type"
	Origin     = "Origin"
	URI        = "URI"
	Size       = "Size"
	Dimensions = "Dimensions"
	Duration   = "Duration"
)
