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
correctly construct a canonical array form of the event, 

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

## bluesky



## farcaster