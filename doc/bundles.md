# Social Network Hoover Data Specification

## Anatomy of Events

1. **Protocol**
2. **Kind of event** (necessary events: posts, reposts, likes, follow, block, profile) standardized
3. **Event data** (standardized: should be same across all protocols)
4. **Verification**/original source info (will be different for each protocol)

## Common

1. `Protocol` - nostr/bluesky/farcaster

2. `Event-ID` - protocol specific

3. `User-ID` - protocol specific

4. `Timestamp` - standardized to unix timestamp

5. `Path` - needed for bluesky (from the original bluesky repo for event)

6. `Signature` - format different for each protocol, should verify to Event-ID

## Event Data

Where the name has a ### put in place a zero prefixed 3 character number

### 1. Posts
1. `Post-Text` - (this should probably be as raw binary in the Data field, because of size)
2. `Richtext-###` - these need to be in nested bundles due to size
3. `Image-###` - initially just a URL, which can be protocol specific for bluesky - can be embedded for bluesky
4. `Media-###` - same as Image-#
5. `Links-###` - standard proto://address/path style, one per numbered entry
6. `Embed-###` - identifiers of nested bundle related to URLs from the previous media references - the number represents the position in the nested bundle list, the URL will correspond to the reference tag from the previous items, and the data itself stored in the nested bundle Data field (optional)

### 2. Reposts

1. `Repost-Event-ID` - protocol native reference to another post that is appearing again on a user’s feed

### Likes

1. `Likes-Event-ID` - protocol native reference to a post from a user

### Follow

1. `Add-Remove` - boolean (true/false)
2. `Follow-User-ID` - User is now/now not following this user

### Block

1. `Add-Remove` - boolean (true/false)
2. `Block-User-ID` - User is now/now not following this user

### Profile

1. `User-Name` - protocol based “user name” separate from display
2. `Display-Name` - protocol based field that should be prominently showed for a user
3. `Avatar-Image` - URL for the image of a user (can be in-protocol embedded in a nested bundle)
4. `Banner-Image` - URL for the image that would be showed in the background of a profile page
5. `Bio` - (2kb max)
6. `Website` - URL related to the user
7. `Verification` - external verification information (nip-05 for nostr)
8. `Payment-Address` - (lightning emoji and LN address eg 🗲mleku@getalby.com for lightning)
