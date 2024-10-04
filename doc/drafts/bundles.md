# Social Network Hoover Data Specification

## Anatomy of Events

1. **Protocol**
2. **Kind of event** (necessary events: posts, reposts, likes, follow, block, profile)
   standardized
3. **Event data** (standardized: should be same across all protocols)
4. **Verification**/original source info (will be different for each protocol)

## Common

1. `App` - `com.hubmakerlabs.hoover`
2. `Protocol` - nostr/bluesky/farcaster
3. `Event-ID` - protocol specific
4. `User-ID` - protocol specific
5. `Timestamp` - standardized to unix timestamp
6. `Kind` - Kind identifiers are from the list below - `Post`/`Repost`/`Like`/`Follow`/`Block`/
   `Profile`
7. `Path` - needed for bluesky (from the original bluesky repo for event)
8. `Signature` - format different for each protocol, should verify to Event-ID

## Event Data

Reference tags for media and URLs are recorded in the same order they are found in the original
protocol.

### 1. Posts

Kind: `Post`

1. Post content is stored in the bundle `Data` field
2. `Richtext` - these need to be in nested bundles due to size
3. `Image` - standard formatted URL with protocol, address and path.
4. `Media` - same as Image-#
5. `Links` - standard proto://address/path style, one per numbered entry
6. `Embed` - identifiers of nested bundle related to URLs from the previous media references -
   the number represents the position in the nested bundle list, the URL will correspond to the
   reference tag from the previous items, and the data itself stored in the nested bundle Data
   field (optional)

#### Tags

Posts often come with tags aside from those listed above, specifically, references to protocol
events and protocol users, and hashtags.

1. `Mention` - reference to another user on the protocol, can be more than one
2. `Event-Reference` - protocol reference to another event, usually a post or repost.
3. `Source` - when the original post came from another protocol (in nostr this is `#proxy`) with
   source type and URL, for example: `activitypub,http://example.com/post.html`
4. `Hashtag` - when the post contains a hash tag reference added by the user to help others find
   the post
5. `Emoji` - custom emojis that will be found in the post text, formatted as "name,URL" for
   example: `pikachu,https://example.com/pikachu.gif`
6. `Content-Warning` - if the usuer has manually flagged the post as NSFW etc (value may have
   more info)

### 2. Reposts

Kind: `Repost`

1. `Repost-Event-ID` - protocol native reference to another post that is appearing again on a
   user’s feed

### Likes

Kind: `Like`

1. `Like-Event-ID` - protocol native reference to a post from a user

### Follow

Kind: `Follow`

1. `Add` - boolean (true/false)
2. `Follow-User-ID` - User is now/now not following this user

### Profile

Kind: `Profile`

1. `User-Name` - protocol based “user name” separate from display
2. `Display-Name` - protocol based field that should be prominently showed for a user
3. `Avatar-Image` - URL for the image of a user (can be in-protocol embedded in a nested bundle)
4. `Banner-Image` - URL for the image that would be showed in the background of a profile page
5. `Bio` - (2kb max, or if larger, embedded with `embed` and the nested bundle is labeled `Bio`)
6. `Website` - URL related to the user
7. `Verification` - external verification information (nip-05 for nostr)
8. `Payment-Address` - (lightning emoji and LN address eg 🗲mleku@getalby.com for lightning)

## Embedded Media

Bluesky in particular has the facility to store embedded media files. The standard way to
encapsulate such files in Arweave is to use nested bundles, bundles that are wrapped up in
the outer event bundle. These bundles have a separate data structure that is common to HTTP
mimetypes.

This facility will be implemented after the textual data as the uploading of the much larger
blobs of binary data may dissuade those deploying a Hoover to actually upload these files.

1. `Reference` - The numbered `Embed-###` field name of the encapsulating bundle.
2. `Mime-Type` - standard annotations as used with web browsers, eg:
    1. `audio/aac`
    2. `video/mp4`
    3. `image/png`
    4. ...
3. `Origin` - this will generally be `bluesky` and will indicate that the URL can be found via a
   `bsky` protocol query. For general web service, this will be the protocol name `http` `https`
   `magnet` `arweave` etc.
4. `URI` - the identifier that the protocol that stores the data will recognize how to fetch the
   data.
5. `Size` - the size of the embedded data in bytes
6. `Dimensions` - for images and videos, the pixels Horizontal by Vertical `HxV`
7. `Duration` - For video and audio, the playback time in seconds

The `Data` field of the bundle will contain the binary file itself (not implemented in initial
version)
