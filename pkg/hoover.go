package pkg

// Constants defining the standard tag field names

// Common

const (
	Protocol  = "Protocol"
	EventId   = "Event-ID"
	UserId    = "User-ID"
	Timestamp = "Timestamp"
	Kind      = "Kind"
	Path      = "Path"
	Signature = "Signature"
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
	PostText = "Post-Text"
	Richtext = "Richtext"
	Image    = "Image"
	Media    = "Media"
	Links    = "Links"
	Embed    = "Embed"
)

// Profile

const (
	UserName     = "User-Name"
	DisplayName  = "Display-Name"
	AvatarImage  = "Avatar-Image"
	BannerImage  = "Banner-Image"
	Bio          = "Bio"
	Website      = "Website"
	Verification = "Verification"
)

// The rest

const (
	RepostEventId = "Repost-Event-ID"
	LikeEventId   = "Like-Event-ID"
	Add           = "Add"
	FollowUserId  = "Follow-User-ID"
	BlockUserId   = "Block-User-ID"
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
