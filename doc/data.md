# data formats

## ao

Arweave AO "bundle" data type is a data structure that contains the 
following fields (using Go type notation):

```go
type Bundle struct {
    Items            []BundleItem `json:"items"`
    BundleBinary     []byte
    BundleDataReader *os.File
}
```

The `BundleItem` type is a structured type which contains the following fields:

```go
type BundleItem struct {
	SignatureType int    `json:"signatureType"`
	Signature     string `json:"signature"`
	Owner         string `json:"owner"`  //  utils.Base64Encode(pubkey)
	Target        string `json:"target"` // optional, if exist must length 32, and is base64 str
	Anchor        string `json:"anchor"` // optional, if exist must length 32, and is base64 str
	Tags          []Tag  `json:"tags"`
	Data          string `json:"data"`
	Id            string `json:"id"`
	TagsBy        string `json:"tagsBy"` // utils.Base64Encode(TagsBytes) for retry assemble item

	ItemBinary []byte   `json:"-"`
	DataReader *os.File `json:"-"`
}
```

The `Tag` type referred to in the above is a key/value pair, so a smart contract
can access these fields by key names.

These are taken directly from the Go SDK for reading/writing data and 
transactions to Arweave AO called [goar](https://github.com/everFinance/goar)

## nostr

The AO Bundle type described in the [ao](#ao) section above roughly match up 
to some parts of the event data type that is the primary record type in the 
`nostr` protocol.

Because tags are indexed for searching, these should thus match up partially 
with the tags used in `nostr`. Thus, all tag fields will appear in this position.

Conveniently, the `#` is reserved for use in the filter query for tags, and 
only single alphabetical ciphers `[a-zA-Z]` are searchable, however, the 
initial tag field in `nostr` can be anything else, and numerous already 
defined event kinds use these, so we can use the `#` as a marker to 
represent tag values, and the content of the string would be comma separated 
quoted strings.

We can then use special labels that will match up exactly with their use with 
`nostr` and we will reserve the "data" field for carrying the `content` field of
nostr events for simplicity.

The base event object is as follows:

```
{
  "id": <32-bytes lowercase hex-encoded sha256 of the serialized event data>,
  "pubkey": <32-bytes lowercase hex-encoded public key of the event creator>,
  "created_at": <unix timestamp in seconds>,
  "kind": <integer between 0 and 65535>,
  "tags": [
    [<arbitrary string>...],
    // ...
  ],
  "content": <arbitrary string>,
  "sig": <64-bytes lowercase hex of the signature of the sha256 hash of the serialized event data, which is the same as the "id" field>
}
```

### Kind 0 - Profile Metadata and Kind 1 - Text note

Kind 1 events are the primary most common form of event, but there is 
numerous others that are relevant.

Kind 0 user metadata contains embedded JSON within it and is an important 
event kind that must be aggregated to make more sense of nostr events. These 
events are replaceable and for most purposes the newest version is the one 
that should be shown for any event bearing their pubkey:

```json
{
  "name": "",
  "picture": "",
  "display_name": "",
  "displayName": "",
  "about": "",
  "website": "",
  "nip05": "",
  "lud16": "",
  "banner": "",
  "pubkey": ""
}
```

For simplicity, these fields and the primary event fields are reserved as 
keys for the `hoover`, and all tags must have a `#` prefix to distinguish 
them. The tags should then be formatted as json-quoted list of strings 
omitting the implicit surrounding brackets implied by the prefix `#` of the 
tag key.

The bundle can then be used to reconstruct the original event canonical JSON 
form in order to validate the nostr signature using BIP-340 secp256k1 
schnorr signature verification.

#### Preservation of tags for authenticity reasons

In many event kinds there are tags that may not have direct relevance to 
being aggregated into another event store where the interaction is not 
bidirectional, notably relay references. However, they are required to 
correctly construct a canonical array form of the event.

One of the unique features of `nostr` is that there is no central authority 
for anything, events are authenticated by the signatures of the users 
themselves, and they can run their own relays and direct other users to post 
messages to their chosen relays.

Technically the event ID can be derived if all other fields are available 
except for the signature, but the ID acts as a unique identifier for each 
event making references to them easy to locate, if any relay has a copy of it.

### Other event kinds that have relevance to AO

- Kind 3 - Follow List (note that these can be very long)
- Kind 6, 16 - Repost
- Kind 7 - Reaction (content field is a symbol, emoji or image)
- Kind 1984 - Reporting - potentially marking events or users as abusive
- Kind 30023 - Long form article, usually in Markdown format
- Kind 30030 - Emoji sets - defining a collection of text labels and symbols 
  for custom reactions, of which some can be tagged as "favourites" - these 
  allow renderers to find appropriate media to put in place for a reaction 
  of a user to another event
- Kind 30315 - User Statuses - should be shown by clients on the profile 
  page in addition to the less frequently changed kind 0 user metadata 
  "about" field.

The rest of the events are largely ephemeral or protocol specific types that 
are not relevant to externally referenced events or media. The purpose of 
the Hoover is aggregating so that multiple social clients can generate 
content that ends up together and able to be associated by especially URLs 
and other events.

For reasons of easier filtering, these kinds will match to a special `type` 
field in the tag key/value set of the bundle:

- `profile`

Profile metadata as described above

- `note`

A short form post

- `follow_list`

The state at the `created_at` timestamp of the follow list of the user.

This can be represented as a diff for each update event and snapshotted 
periodically to reduce total storage space utilization.

- `repost`

A simple broadcasting of a note with the attribution of reposting to a user.

- `reaction`

A short symbolic message intended to represent a reaction to an event.

- `report`

A user report against a post or user expressing the opinion of it/they being 
abusive.

- `article`

A long form post that is expected to be formatted in markdown, usually.

- `emojis`

A set of emojis and names that can be interpreted to be associated with the 
user's posts for example `:emoji_name_here:` to be replaced with the content 
field of the associated named emoji in this event.

- `status`

A short text associated with a user's profile that usually will be changed 
relatively often, compared to the user metadata.

## bluesky

Bluesky is a federated protocol, so its design is radically different from 
how it can be done with Nostr.

A tool called Bigsky forms the basis of an aggregation spider, and based on 
this system events on bluesky protocol federations will be gathered, parsed 
and converted into bundle data and pushed up to the Arweave AO network.

Note that 

### Profile Information

```go
// ActorDefs_ProfileViewDetailed is a "profileViewDetailed" in the app.bsky.actor.defs schema.
type ActorDefs_ProfileViewDetailed struct {
  Associated           *ActorDefs_ProfileAssociated       `json:"associated,omitempty" cborgen:"associated,omitempty"`
  Avatar               *string                            `json:"avatar,omitempty" cborgen:"avatar,omitempty"`
  Banner               *string                            `json:"banner,omitempty" cborgen:"banner,omitempty"`
  CreatedAt            *string                            `json:"createdAt,omitempty" cborgen:"createdAt,omitempty"`
  Description          *string                            `json:"description,omitempty" cborgen:"description,omitempty"`
  Did                  string                             `json:"did" cborgen:"did"`
  DisplayName          *string                            `json:"displayName,omitempty" cborgen:"displayName,omitempty"`
  FollowersCount       *int64                             `json:"followersCount,omitempty" cborgen:"followersCount,omitempty"`
  FollowsCount         *int64                             `json:"followsCount,omitempty" cborgen:"followsCount,omitempty"`
  Handle               string                             `json:"handle" cborgen:"handle"`
  IndexedAt            *string                            `json:"indexedAt,omitempty" cborgen:"indexedAt,omitempty"`
  JoinedViaStarterPack *GraphDefs_StarterPackViewBasic    `json:"joinedViaStarterPack,omitempty" cborgen:"joinedViaStarterPack,omitempty"`
  Labels               []*comatprototypes.LabelDefs_Label `json:"labels,omitempty" cborgen:"labels,omitempty"`
  PostsCount           *int64                             `json:"postsCount,omitempty" cborgen:"postsCount,omitempty"`
  Viewer               *ActorDefs_ViewerState             `json:"viewer,omitempty" cborgen:"viewer,omitempty"`
}
```

Mapping these to the `nostr` pattern for uniformity:

- name - `handle`
- picture - `avatar` 
- display_name - `displayName`
- about - `description`
- banner - `banner`
- pubkey - `did` (the `pubkey` field also appears on every `nostr` event anyway)

`created_at` will contain the timestamp of the most recent fetch with changes - 
since AO is immutable all old versions will be searchable as well. this will 
not represent actual user activity timestamp in any way, as regards to the the 
change.

### Posts

```go
// EmbedRecord_ViewRecord is a "viewRecord" in the app.bsky.embed.record schema.
//
// RECORDTYPE: EmbedRecord_ViewRecord
type EmbedRecord_ViewRecord struct {
	LexiconTypeID string                                `json:"$type,const=app.bsky.embed.record#viewRecord" cborgen:"$type,const=app.bsky.embed.record#viewRecord"`
	Author        *ActorDefs_ProfileViewBasic           `json:"author" cborgen:"author"`
	Cid           string                                `json:"cid" cborgen:"cid"`
	Embeds        []*EmbedRecord_ViewRecord_Embeds_Elem `json:"embeds,omitempty" cborgen:"embeds,omitempty"`
	IndexedAt     string                                `json:"indexedAt" cborgen:"indexedAt"`
	Labels        []*comatprototypes.LabelDefs_Label    `json:"labels,omitempty" cborgen:"labels,omitempty"`
	LikeCount     *int64                                `json:"likeCount,omitempty" cborgen:"likeCount,omitempty"`
	ReplyCount    *int64                                `json:"replyCount,omitempty" cborgen:"replyCount,omitempty"`
	RepostCount   *int64                                `json:"repostCount,omitempty" cborgen:"repostCount,omitempty"`
	Uri           string                                `json:"uri" cborgen:"uri"`
	// value: The record data itself.
	Value *util.LexiconTypeDecoder `json:"value" cborgen:"value"`
}
```

These have an "embeds" record associated with them, which in `nostr` are 
simply parsed out by clients from the `content` field, with the exception of 
event, addr (abstract event reference) and pubkeys which usually get a tag.

There is several types, so they need to be listed so they get associated to 
the post objects:

```go
// EmbedImages_View is a "view" in the app.bsky.embed.images schema.
//
// RECORDTYPE: EmbedImages_View
type EmbedImages_View struct {
	LexiconTypeID string                   `json:"$type,const=app.bsky.embed.images#view" cborgen:"$type,const=app.bsky.embed.images#view"`
	Images        []*EmbedImages_ViewImage `json:"images" cborgen:"images"`
}

// EmbedExternal_View is a "view" in the app.bsky.embed.external schema.
//
// RECORDTYPE: EmbedExternal_View
type EmbedExternal_View struct {
	LexiconTypeID string                      `json:"$type,const=app.bsky.embed.external#view" cborgen:"$type,const=app.bsky.embed.external#view"`
	External      *EmbedExternal_ViewExternal `json:"external" cborgen:"external"`
}

// EmbedRecord_View is a "view" in the app.bsky.embed.record schema.
//
// RECORDTYPE: EmbedRecord_View
type EmbedRecord_View struct {
	LexiconTypeID string                   `json:"$type,const=app.bsky.embed.record#view" cborgen:"$type,const=app.bsky.embed.record#view"`
	Record        *EmbedRecord_View_Record `json:"record" cborgen:"record"`
}

// EmbedRecordWithMedia_View is a "view" in the app.bsky.embed.recordWithMedia schema.
//
// RECORDTYPE: EmbedRecordWithMedia_View
type EmbedRecordWithMedia_View struct {
	LexiconTypeID string                           `json:"$type,const=app.bsky.embed.recordWithMedia#view" cborgen:"$type,const=app.bsky.embed.recordWithMedia#view"`
	Media         *EmbedRecordWithMedia_View_Media `json:"media" cborgen:"media"`
	Record        *EmbedRecord_View                `json:"record" cborgen:"record"`
}
```

Being such a relatively complicated scheme compared to `nostr`, these will 
then map to the same form as the `nostr` event where the `e` or `a` tag 
refers to the `cid` field of the generic Record type, which will presumably 
be something vaguely like a post ID in `nostr`.

Thus, to aggregate such events, the record is first found, and then its 
embed data is retrieved, and assembled into a single record for publication 
to AO.

### Follows

A cumulative snapshot view of follow graphs, both inbound and outbound can 
be requested from a Bluesky server. Again like most other events except for 
posts these are snapshots and do not represent direct user activity, so they 
have to be polled regularly to remain up to date.

For this reason, these lists need to be stored as diffs on Arweave AO.

```go

// GraphGetFollowers_Output is the output of a app.bsky.graph.getFollowers call.
type GraphGetFollowers_Output struct {
	Cursor    *string                  `json:"cursor,omitempty" cborgen:"cursor,omitempty"`
	Followers []*ActorDefs_ProfileView `json:"followers" cborgen:"followers"`
	Subject   *ActorDefs_ProfileView   `json:"subject" cborgen:"subject"`
}

// GraphGetFollows_Output is the output of a app.bsky.graph.getFollows call.
type GraphGetFollows_Output struct {
	Cursor  *string                  `json:"cursor,omitempty" cborgen:"cursor,omitempty"`
	Follows []*ActorDefs_ProfileView `json:"follows" cborgen:"follows"`
	Subject *ActorDefs_ProfileView   `json:"subject" cborgen:"subject"`
}
```

### Bluesky Event Data Format Encoding

The primary types of events that this aggregator will gather will be 
constructed in the following ways for storage on the Arweave AO:

## farcaster

Farcaster is also quite distinct from both `nostr` and `bluesky` but 
something in-between.

The data is aggregated, as in, the data is aggregated together like a 
centralised distributed data application, but with a less distinct 
centralisation due to being built on top of a large internationally 
distributed cluster of replicas, instead of federated, and unlike `nostr` 
has a consistency across replicas that is enforced rather than completely 
ad-hoc and pure decentralised.

There is 4 main types of data that need to be mapped to Arweave AO:

- User
- Cast
- Recast
- Follow
- Reaction

### User

```go
type User struct {
	Fid              int            `json:"fid"`
	Username         string         `json:"username"`
	DisplayName      string         `json:"displayName"`
	Pfp              Pfp            `json:"pfp"`
	Profile          Profile        `json:"profile"`
	FollowerCount    int            `json:"followerCount"`
	FollowingCount   int            `json:"followingCount"`
	ReferrerUsername string         `json:"referrerUsername"`
	ViewerContext    *ViewerContext `json:"viewerContext"`
}

type Profile struct {
	Bio Bio `json:"bio"`
}

type Bio struct {
	Text     string   `json:"text"`
	Mentions []string `json:"mentions"`
}

type Pfp struct {
	Url      string `json:"url"`
	Verified bool   `json:"verified"`
}
```

### Cast

```go
type Cast struct {
	Hash          string         `json:"hash"`
	ThreadHash    string         `json:"threadHash"`
	ParentHash    string         `json:"parentHash"`
	ParentAuthor  *users.User    `json:"parentAuthor"`
	Author        *users.User    `json:"author"`
	Text          string         `json:"text"`
	Timestamp     uint64         `json:"timestamp"`
	Replies       *Replies       `json:"replies"`
	Reactions     *Reactions     `json:"reactions"`
	Recasts       *Recasts       `json:"recasts"`
	Watches       *Watches       `json:"watches"`
	Recast        bool           `json:"recast"`
	ViewerContext *ViewerContext `json:"viewerContext"`
}
```

### Recast

```go
type Recasts struct {
    Count   int         `json:"count"`
    Recasts []*Recaster `json:"recasters"`
}

type Recaster struct {
	Fid         uint64 `json:"fid"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	RecastHash  string `json:"recastHash"`
}
```

### Follow

```go
type Follow struct {
	FollowerFid  int64 `json:"follower_fid"`
	FollowingFid int64 `json:"following_fid"`
}
```

### Reaction

```go
type Reaction struct {
	Type      string      `json:"type"`
	Hash      string      `json:"hash"`
	Reactor   *users.User `json:"reactor"`
	Timestamp uint64      `json:"timestamp"`
	CastHash  string      `json:"castHash"`
}
```

