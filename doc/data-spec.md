# Arweave Social Media Data and Explorer

The Hoover is a social media protocol data aggregator that pulls a small subset of the most important events in the respective social media protocol data sets, and makes them accessible to applications that can query the Arweave permaweb.

The purpose of this application is to aggregate these disparate social media databases and allow them to be integrated into Arweave applications. It is a possible future use case where the data that can be generated to put into the Hoover data set could also generate events to send out to the respective social networks and create a full bridge between them all.

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

### Kind `Post`

#### Nostr

- `Reply-Root-Id`

  This designates the original post (OP) that this post part of the post thread tree. Older events may not have this field, which makes it more complex to derive a post graph for a given original post.

- `Reply-Parent-Id`

  This designates the immediate previous post that this event is a reply to.

- `Content-Warning` 

  Can additionally contain a description of the type of content, otherwise empty.

In the data field also there can be:

- `Mention`, `Hashtag` - other user's public keys, hashtags
- `Source` - external source where content originated (for bridges republishing into nostr)
- `Emoji` - these additionally require acquiring the emoji definition event type in most cases
- `Label`, `Label-Namespace` - labels and label namespaces used in protocol in an indexable form (omitted from main tags due to potential data size for now, also, Nostr has an independent labelling protocol).

#### Bluesky



#### Farcaster



### Kind `Repost`

#### Nostr

- `Repost-Event-Id`

  This designates the event that a user is making visible to their followers.

In the data field there can also be:

- `Source` - designates a name of the protocol source, currently mainly this refers to Mastodon via the Mostr bridge relays.
- `Source-Uri` - can encode a protocol specific Uniform Resource Identifier that can be used to search for the original event from the source.
- `Mention` - reposts can tag a user in relation to the repost to specifically notify them.
- `Label`, `Label-Namespace` - labels and label namespaces used in protocol in an indexable form (omitted from main tags due to potential data size for now.

#### Bluesky



#### Farcaster



### Kind `Like`

#### Nostr

- `Like-Event-Id`

  The `Event-Id` of the event being responded to.

The content field can contain a symbol, usually a `+` or `-` or arbitrary Unicode emojis or the names of emojis created in emoji events by the publisher of the Like event.

In addition, there can also be:

- `Mention` - other user's public keys
- `Source` - designates a name of the protocol source, currently mainly this refers to Mastodon via the Mostr bridge relays.
- `Source-Uri` - can encode a protocol specific Uniform Resource Identifier that can be used to search for the original event from the source.

#### Bluesky



#### Farcaster



### Kind `Follow`

#### Nostr

Nostr follow events are in fact snapshot events that contain the entire list of followed items. For this reason, being sometimes upwards of 100kb in size, we do not use the tags for them, at this time. There is work in progress to create a CRDT compatible form of these user relations that adds other kinds of connections than follow and mute but it is still in draft.

The data field of the bundle can contain the following two items in lists:

- `Follow-User-Id` - the public key identifiers of other users the publisher subscribes to.
- `Follow-Tag` - a list of hashtags the user subscribes to.

#### Bluesky



#### Farcaster



### Kind `Profile`

#### Nostr

- `User-Name`



- `Display-Name`



- `Avatar-Image`



- `Banner-Image`


- `Content-Warning`

  Can additionally contain a description of the type of content, otherwise empty.

Additional fields which are placed in the data field of the bundle:

- `Bio`
- `Website`
- `Verification`
- `Payment-Address`

And a number of fields that refer to items mentioned in the `Bio` field:

- `Mention-Event-Id` - `Event-Id` values that appear in the `Bio`
- `Hashtag` - hashtags 
- `Mention` - tagging another user.
- `Source` - designates a name of the protocol source, currently mainly this refers to Mastodon via the Mostr bridge relays.
- `Source-Uri` - can encode a protocol specific Uniform Resource Identifier that can be used to search for the original event from the source.

#### Bluesky



#### Farcaster



## Explorer

The explorer is a simple one page browser that loads the latest transactions created by the Hoover found in an arweave gateway and displays them, with the ability to step backwards in the history to earlier events.

Its output is extremely rudimentary, and some of the non-searchable fields, contained in JSON data alongside the content field of the event data, is not displayed.

The purpose of the Explorer is to demonstrate a simple GraphQL based search for Hoover data and shows a checkmark symbol on Nostr events that have valid User-Id/Event-Id/Signature data.

Note that as mentioned previously, there are some fields that can be omitted from the event data and full reversal of the conversion cannot currently be achieved to enable full authentication of the data. This could be the subject of later work... for Bluesky especially it is complex. One possibility is to sub-bundle the raw event data, which would be a very simple change, however, this also means a lot more data per event.