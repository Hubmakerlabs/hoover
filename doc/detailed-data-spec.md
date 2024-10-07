# Arweave Social Media Data and Explorer

The Hoover is a social media protocol data aggregator that pulls a small subset of the most important events in the respective social media protocol data sets, and makes them accessible to applications that can query the Arweave permaweb.

The purpose of this application is to aggregate these disparate social media databases and allow them to be integrated into Arweave applications. It is a possible future use case where the data that can be generated to put into the Hoover data set could also generate events to send out to the respective social networks and create a full bridge between them all.

## Explorer

The explorer is a simple one page browser that loads the latest transactions created by the Hoover found in an arweave gateway and displays them, with the ability to step backwards in the history to earlier events.

Its output is extremely rudimentary, and some of the non-searchable fields, contained in JSON data alongside the content field of the event data, is not displayed.

The purpose of the Explorer is to demonstrate a simple GraphQL based search for Hoover data and shows a checkmark symbol on Nostr events that have valid User-Id/Event-Id/Signature data.

Note that as mentioned previously, there are some fields that can be omitted from the event data and full reversal of the conversion cannot currently be achieved to enable full authentication of the data. This could be the subject of later work... for Bluesky especially it is complex. One possibility is to sub-bundle the raw event data, which would be a very simple change, however, this also means a lot more data per event.

## Common Data Format and Conventions

There is 5 types of events that the Hoover extracts from:

- `Post` - text content with possibly links and references to other events and users
- `Repost` - an event that is simply the making visible of a post event to the followers of a user
- `Like` - an event that signifies a response from a user to a post
- `Follow` - publicly advertising that a user is subscribing to the events of another user
- `Profile` - an event containing metadata about the user, display names, contact methods, avatar, banner images

In addition, there is a set of common fields found in all events, stored in the Arweave transaction bundle tags, as follows:

- `App-Name`
  
  this should be the usual reverse structured naming scheme as used in many other app systems like Android - `com.hubmakerlabs.hoover`

- `App-Version`
  
  this should be a semver style version string, the final version for this project release will be `v1.0.0`

- `Protocol`
  
  this field encodes the protocol, in this release we are providing `Bluesky`, `Nostr` and `Farcaster` protocols

- `Repository`
  
  - for `Nostr` this is the relay where the event was found, very often there is many others that also have it but this is the one it was retrieved from
  - for `Bluesky` this is the main bluesky API endpoint provided by the Bluesky project themselves
  - for `Farcaster` this is one of several providing API access

- `Kind`
  
  kind is one of `Post`, `Repost`, `Like`, `Follow`, and `Profile` as described above

- `Event-Id`
  
  this is the protocol specific event identifier, which is in all cases for this release a hash of the canonically structured event data... Bluesky uses base32, Farcaster and Nostr use hexadecimal

- `User-Id`
  
  this is the network identifier for the user, for Nostr it is the x-only public key of secp256k1 curve, for Bluesky it is a hashed version of the NIST-p256 curve public key, and for Farcaster it is a serial identifier as the accounts (public key, ED25519 type) are stored when created on a blockchain and thus can be assigned serial numbers - the Nostr key can be used directly in signature verification, Bluesky ID must be looked up from the user's most recent profile event, and the Ethereum blockchain must be searched for Farcaster... Note that for Farcaster, the identity is tied to a client installation and cannot currently be directly accessed aside from the client, and the public key is found on a [key registry](https://github.com/farcasterxyz/protocol/discussions/103).

- `Unix-Time`
  
  this is the time at which the client created the event as a UNIX timestamp which designates 1 January 1970 at 00:00 as the zero second or epoch. The Farcaster event encodes this instead as an epoch offset from 1 January 2021 00:00 and for this application we add the offset to the Unix epoch to get a consistent relative value

- `Signature`
  
  This is the signature on the `Event-Id` made using the secret counterpart of the public key represented in the `User-Id`. Only Nostr enables this to be done without any further searches, Bluesky requires acquiring the profile event where the key can be extracted, and the Ethereum blockchain is where the key, related to the `User-Id` can be found. (this is the reason why currently the browser app does not verify signatures on either `Farcaster` or `Bluesky`)

There can also be protocol specific fields present in this common set as well, related to each of the event kinds that we handle with Hoover.

These have been treated specific to the protocol in the list below but mainly because there are differences in the encoding used for the relevant fields, such as `Event-Id` and `User-Id`

There are some fields that are not carried across from the original protocol, and some that are only put with the `Content` field inside the arweave bundle data, we only cover the ones that are in the bundle tags and are thus searchable using the GraphQL endpoint.

The reason for this is the fields that are present but not added to the indexable tags is that they can be of indefinite size and there is an informal limitation on total size of tag data of 2kb, and the reason for limiting indexable data is that indexes grow nonlinearly with the amount of indexable data points.

Note also that even though in most cases much of the extra data for each event kind is encoded there can be missing fields that prevent complete reconstruction of the original event date.

This could be fixed for a future version if there is a need for this extra data, or to act secondarily as a backing store for a protocol service using Arweave as a shared data store, but was built this way primarily as a minimally complete data set for Arweave applications presenting the data sourced externally.

## Event Specific Data

In this section, we show under the *Kind* heading the common elements that can appear, and then under each of the protocol subheadings, the specifics of other extra data that may be found.

### Kind `Post`

- `Reply-Root-Id`

  This designates the original post (OP) that this post part of the post thread tree. Older events may not have this field, which makes it more complex to derive a post graph for a given original post.

- `Reply-Parent-Id`

  This designates the immediate previous post that this event is a reply to.

#### Nostr

- `Mention`, `Hashtag` - other user's public keys, hashtags
- `Mention-Event-Id` - when the `Content` contains a reference to another event, usually encoded as a `nip-19` Bech32 entity in the content field.
- `Source` - external source where content originated (for bridges republishing into nostr)
- `Emoji` - these additionally require acquiring the emoji definition event type in most cases
- `Label`, `Label-Namespace` - labels and label namespaces used in protocol in an indexable form (omitted from main tags due to potential data size for now, also, Nostr has an independent labelling protocol).
- `Content-Warning` - can additionally contain a description of the type of content (possibly, of the `Bio` or just a general NSFW type of content being posted), otherwise the value is empty.
- `Label` and `Label-Namespace` are optional fields that are used to associate the post with categories both abstract and more concrete.

Note we don't currently deal with the `d` tag of events which is used to associate a chain of events as being `parameterized replaceable` events that can be selected and searched for via the `kind:user-id:d-tag`

> [Code location](https://github.com/Hubmakerlabs/hoover/blob/master/pkg/nostr/nostr.go) of implementation.

#### Bluesky

- `Reply-Root-Uri` - a protocol structured URI for the reply root (original post)
- `Reply-Parent-Uri` - a protocol structured URI for the reply parent (for threaded viewing to be a branch under the parent)
- `Embed-Record-Uri` - a reference to an external event other than the reply thread as a URI (this is the bluesky name for the encoded version)
- `Mention-Event-Id` - the bluesky base32 encoding of mentioned event, will appear directly following the `Embed-Record-Id` which is the same thing as the same named object in the Nostr version.
- `Embed-Image` - which appears as a series of numbered items containing subfields with the same numbering prefix so they are easily associated. There is numerous fields that inclued `Ref` (URL), `Mimetype`, `Size`, `Aspect`, and an `Alt` description text, usually these are stored separately on the Bluesky server but we do not fetch or cache them as they are full sized media files, referred to usually by `at://` protocol prefixes.
- `Embed-Record` prefixed fields which contain similar things as `Embed-Image` but usually will refer to video or audio media and can include links to thumbnails, descriptions and titles. 
- `Embed-External` which is like `Embed-Record` but where the URI's are not from atproto `at://`.
- `Entities` - which are references to text within the `Content` field such as full versions of URLs, which are ellipsised in the Content field.
- `Richtext` - more numbered fields which are references to external resources found in the `Content` field, which again are usually ellipsised. *[ Ed: Yes! this makes rendering their events more needlessly complex.]* - these also include hashtags.
- `Language` - which contains a standard 2 letter ISO language code such as `en` or `jp` or `cn` etc.
- `Hashtag` - numbered fields (if more than one) that contain usually a hashtag text

> [Code location](https://github.com/Hubmakerlabs/hoover/blob/master/pkg/bluesky/app.bsky.feed.post.go) of implementation.

#### Farcaster

Farcaster does not use `Reply-Root-Id` so there is a bit more work searching for the complete tree of a discussion thread.

- `Embed-Uri` - references to external media
- `Mention` - references to users mentioned in a post, the following two fields designate the location and this field designates the mentioned `User-Id`. This field is numbered as in `Mention-X` if there is more than one.
- `Mention-X-Start` and `Mention-X-End` marks the positions in the `Content` field where the mentions are found, the number part is absent if there is only one.
- `Reply-Parent-User-Id` - the `User-Id` of the parent event of a reply `Post`
- `Reply-Parent-Uri` - the URI of the reply post, if present.
- `Embed-X-Uri` marks external references to eternal objects.
- `Embed-X-User-Id` and `Embed-X-Event-Id` designates an internal reference to another event, as with most Farcaster event references the `User-Id` is usually provided alongside the `Event-Id` (though probably redundant since the IDs are hashes).

### Kind `Repost`

- `Repost-Event-Id`

  This designates the event that a user is making visible to their followers.

#### Nostr

- `Source` - designates a name of the protocol source, currently mainly this refers to Mastodon via the Mostr bridge relays.
- `Source-Uri` - can encode a protocol specific Uniform Resource Identifier that can be used to search for the original event from the source.
- `Mention` - reposts can tag a user in relation to the repost to specifically notify them.
- `Label`, `Label-Namespace` - labels and label namespaces used in protocol in an indexable form (omitted from main tags due to potential data size for now.

> [Code location](https://github.com/Hubmakerlabs/hoover/blob/master/pkg/nostr/nostr.go) of implementation.

#### Bluesky

- `Repost-Event-Uri` - is the Bluesky protocol URI referring to the `Repost-Event-Id` on protocol.

#### Farcaster

It can be that Farcaster does not have the `Repost-Event-Id` and instead only has a `Repost-Event-Uri`

- `Repost-User-Id` - provides the `User-Id` associated with the `Repost-Event-Id`
- `Repost-Event-Uri` is sometimes present in the data field (todo: this may actually contain the `Event-Id` and `User-Id`)

### Kind `Like`

- `Like-Event-Id`

  The `Event-Id` of the event being responded to.

#### Nostr

The content field can contain a symbol, usually a `+` or `-` or arbitrary Unicode emojis or the names of emojis created in emoji events by the publisher of the Like event, in this form: `:emojiname:`.

In addition, there can also be:

- `Mention` - other user's public keys
- `Source` - designates a name of the protocol source, currently mainly this refers to Mastodon via the Mostr bridge relays.
- `Source-Uri` - can encode a protocol specific Uniform Resource Identifier that can be used to search for the original event from the source.

> [Code location](https://github.com/Hubmakerlabs/hoover/blob/master/pkg/nostr/nostr.go) of implementation.

#### Bluesky

- `Like-Path` - a protocol specific path referring to the `Like-Event-Id`
- `Mention` - refers to the `User-Id` who created the event being liked.

#### Farcaster

Also like the Repost, it can be that the event only has a `Like-Event-Uri` but again todo: probably is the same data and should be decomposed.

- `Like-User-Id` provides the farcaster ID of the user in addition to the common `Like-Event-Id`

### Kind `Follow`

- `Follow-User-Id` - the on-protocol identity of the user being followed.

#### Nostr

Nostr follow events are in fact snapshot events that contain the entire list of followed items. For this reason, being sometimes upwards of 100kb in size, we do not use the tags for them, at this time. There is work in progress to create a CRDT compatible form of these user relations that adds other kinds of connections than follow and mute but it is still in draft.

The data field of the bundle can contain the following two items in lists:

- `Follow-User-Id` - the public key identifiers of other users the publisher subscribes to.
- `Follow-Tag` - a list of hashtags the user subscribes to.

> [Code location](https://github.com/Hubmakerlabs/hoover/blob/master/pkg/nostr/nostr.go) of implementation.

#### Bluesky

No specific protocol tags for this event, except just to note that Bluesky follows contain just one, and we have not implemented the unfollow event, which would be sourced from the delete events streme in Bluesky.

#### Farcaster

Nothing different here, except possibly there are delete events same as Bluesky.

### Kind `Profile`

- `User-Name`

The user-name, usually a short handle containing no spaces, protocols may not have them.

- `Display-Name`

The name that the user wants to be displayed along with their content.

- `Avatar-Image`

A URL referring to an image that should be used as the users avatar image.

- `Banner-Image`

A URL referring to an image that should be shown as a background at the top of their profile page.

#### Nostr

- `Content-Warning`

  Can additionally contain a description of the type of content, otherwise empty.

Additional fields which are placed in the data field of the bundle:

- `Bio` - an arbitrary length text field that can contain nostr protocol references to other users, events, and hash tags
- `Website` - usually a URL for a website the user wants to advertise, often the profile of github or other external social network.
- `Verification` - for Nostr this means a [nip-05](https://github.com/nostr-protocol/nips/blob/master/05.md) `user@example.com`
- `Payment-Address` - currently this is primarily a `user@example.com` lightning network payment address in the LUD16 format.

And a number of fields that refer to items mentioned in the `Bio` field:

- `Mention-Event-Id` - `Event-Id` values that appear in the `Bio`
- `Hashtag` - hashtags.
- `Mention` - tagging another user.
- `Source` - designates a name of the protocol source, currently mainly this refers to Mastodon via the Mostr bridge relays.
- `Source-Uri` - can encode a protocol specific Uniform Resource Identifier that can be used to search for the original event from the source.

> [Code location](https://github.com/Hubmakerlabs/hoover/blob/master/pkg/nostr/nostr.go) of implementation.

#### Bluesky

Bluesky does not have `User-Name` 

- `Avatar-Image` has additional fields in the data for Mimetype and Size.
- `Banner-Image` same as `Avatar-Image`.

#### Farcaster

- Has `Avatar-Image`, `Display-Name`, `Bio`,`Website`, and `User-Name` same as those seen in Nostr above

## Notes

### Verification and Bridging

Currently there is no direct mechanism for reversing the bundling process, some data is omitted and especially with Bluesky the data format is very complex.

Some options for solving this problem include extending the data specification with a full bidirectional transformation matrix, or alternatively, sub-bundling the native raw event data, with the idea that one can be used for most simple on-chain applications where the extra data can be used for more complex data bridging.
