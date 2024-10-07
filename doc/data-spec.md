# Data Specification
This document outlines the structure and fields of social network events aggregated by the Hoover. It describes common data fields applicable to all event types (`Post`, `Repost`, `Like`, `Follow`, `Profile`), as well as event-specific fields for each social network protocol, including Nostr, Bluesky, and Farcaster. The Hoover’s role is to standardize and store this data on Arweave for long-term accessibility and integration with distributed applications.

> For more detailed information, including developers' insights on the data structure and field descriptions, please refer to the [Detailed Data Specification](detailed-data-spec.md) document.

## Common Data Fields
This section lists the core data fields common across all event types (`Post`, `Repost`, `Like`, `Follow`, `Profile`). These fields are present in every event aggregated by the Hoover and provide essential metadata, such as event identifiers, user information, timestamps, and protocol details.

| **Field Name**  | **Explanation of Field**                                                                                                                                                                                                                                 |
| --------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Kind**        | The event type: one of [`Post`](#post), [`Repost`](#repost), [`Like`](#like), [`Follow`](#follow), or [`Profile`](#profile)<br>(Each event type also has event specific sub-headings. Please click on each event type above to view relevant sub-fields) |
| **Protocol**    | The protocol of the event, e.g., `Bluesky`, `Nostr`, or `Farcaster`                                                                                                                                                                                      |
| **Repository**  | For `Nostr`, this is the relay where the event was found; for `Bluesky`, the main Bluesky API endpoint; for `Farcaster`, one of several API access points                                                                                                |
| **Event-Id**    | The protocol-specific event identifier, which is a hash of the canonically structured event data. Bluesky uses base32, Farcaster and Nostr use hexadecimal                                                                                               |
| **User-Id**     | The network identifier for the user. For Nostr, it is the x-only public key of secp256k1 curve; for Bluesky, a hashed version of the NIST-p256 curve public key; for Farcaster, a serial identifier on a blockchain                                      |
| **Unix-Time**   | The time when the event was created, represented as a UNIX timestamp (from 1 January 1970, 00:00). Farcaster uses an epoch offset from 1 January 2021, adjusted to the UNIX epoch for consistency                                                        |
| **Signature**   | The signature on the `Event-Id`, made using the secret counterpart of the public key represented in the `User-Id`. Verification varies by protocol, with Nostr offering the easiest validation                                                           |
| **App-Name**    | The usual reverse-structured naming scheme as used in many app systems, e.g., `com.hubmakerlabs.hoover`                                                                                                                                                  |
| **App-Version** | A semver-style version string, the final version for this project release is `v1.0.0`                                                                                                                                                                    |

## Event Specific Data
This section explains the fields specific to each event type (`Post`, `Repost`, `Like`, `Follow`, `Profile`) across different social network protocols, detailing how data is structured and handled.

### Post
Posts are the primary form of content shared on social networks. They may include replies to other posts, media embeds, mentions of other users, and more. Each protocol handles this event type slightly differently, and the table below outlines the specific fields used across different protocols.

| **Field**            | **Protocol** | **Explanation**                                                                                               |
|----------------------|--------------|---------------------------------------------------------------------------------------------------------------|
| **Reply-Root-Id**     | All          | Designates the original post (OP) that this post is a part of. Older events may not have this field.           |
| **Reply-Parent-Id**   | All          | Designates the immediate previous post that this event is replying to.                                         |
| **Mention**, **Hashtag** | Nostr    | Public keys or hashtags mentioned in the content.                                                             |
| **Mention-Event-Id**  | Nostr        | Reference to another event, usually encoded as a `nip-19` Bech32 entity in the content field.                  |
| **Source**            | Nostr        | External source where content originated (for bridges republishing into Nostr).                                |
| **Emoji**             | Nostr        | Custom emojis, often requiring a separate event to define them.                                               |
| **Label**, **Label-Namespace** | Nostr | Optional fields for associating posts with categories or namespaces.                                           |
| **Content-Warning**   | Nostr        | A warning indicating the post contains sensitive content (e.g., NSFW).                                         |
| **Reply-Root-Uri**    | Bluesky      | URI for the reply root (original post).                                                                       |
| **Reply-Parent-Uri**  | Bluesky      | URI for the parent reply.                                                                                     |
| **Embed-Record-Uri**  | Bluesky      | Reference to an external event as a URI.                                                                      |
| **Embed-Image**       | Bluesky      | Contains information about attached images, such as URL, mimetype, size, and description (Alt text).           |
| **Embed-Record**      | Bluesky      | Similar to Embed-Image but used for media like video or audio, with links to thumbnails and titles.            |
| **Embed-External**    | Bluesky      | Reference to external non-Bluesky URIs.                                                                       |
| **Entities**          | Bluesky      | References to text within the content field, such as full URLs.                                                |
| **Richtext**          | Bluesky      | Numbered fields referring to external resources found in the content field.                                    |
| **Language**          | Bluesky      | 2-letter ISO language code (e.g., `en`, `jp`).                                                                 |
| **Hashtag**           | Bluesky      | Contains hashtag text (numbered fields if more than one).                                                     |
| **Embed-Uri**         | Farcaster    | References to external media objects.                                                                         |
| **Mention**           | Farcaster    | References to mentioned users, with start and end positions in the content field for locating the mentions.    |
| **Reply-Parent-User-Id** | Farcaster | The `User-Id` of the parent event of a reply post.                                                            |
| **Reply-Parent-Uri**  | Farcaster    | URI of the parent reply post, if present.                                                                     |

### Repost
Reposts are used to share or amplify content from other users. This section describes the fields used to identify the original post being reposted and relevant metadata across different protocols.

| **Field**            | **Protocol** | **Explanation**                                                                                               |
|----------------------|--------------|---------------------------------------------------------------------------------------------------------------|
| **Repost-Event-Id**   | All          | Designates the event being reposted by the user.                                                              |
| **Source**            | Nostr        | Name of the protocol source (e.g., Mastodon via the Mostr bridge relays).                                      |
| **Source-Uri**        | Nostr        | URI used to search for the original event from the source.                                                    |
| **Mention**           | Nostr        | Tags another user in relation to the repost event.                                                            |
| **Label**, **Label-Namespace** | Nostr | Optional labels and namespaces used in protocol, though omitted from main tags due to potential data size.    |
| **Repost-Event-Uri**  | Bluesky      | URI referring to the `Repost-Event-Id`.                                                                       |
| **Repost-User-Id**    | Farcaster    | Provides the `User-Id` associated with the `Repost-Event-Id`.                                                 |
| **Repost-Event-Uri**  | Farcaster    | URI for the reposted event, sometimes containing `Event-Id` and `User-Id`.                                    |


### Like
The Like event represents a user’s interaction with another post, usually signaling approval or acknowledgment. This section outlines the metadata associated with liking a post, including the reference to the liked event and optional emojis.

| **Field**            | **Protocol** | **Explanation**                                                                                               |
|----------------------|--------------|---------------------------------------------------------------------------------------------------------------|
| **Like-Event-Id**     | All          | The event identifier of the event being liked by the user.                                                    |
| **Content**           | Nostr        | Can contain a symbol (`+`, `-`, emojis) representing the "like" action, or names of emojis created by the user.|
| **Mention**           | Nostr        | Other public keys (users) mentioned in the event.                                                             |
| **Source**            | Nostr        | Protocol source name (e.g., Mastodon via Mostr bridge relays).                                                |
| **Source-Uri**        | Nostr        | URI used to search for the original event from the source.                                                    |
| **Like-Path**         | Bluesky      | Protocol-specific path referring to the `Like-Event-Id`.                                                      |
| **Mention**           | Bluesky      | Refers to the `User-Id` who created the liked event.                                                          |
| **Like-Event-Uri**    | Farcaster    | URI of the liked event (could be the same as the `Like-Event-Id`).                                            |
| **Like-User-Id**      | Farcaster    | Provides the Farcaster ID of the user who liked the event, in addition to the `Like-Event-Id`.                |

### Follow
Follow events represent a user’s subscription to another user’s activity. This section details the fields that represent this relationship and how it is handled differently across protocols.

| **Field**            | **Protocol** | **Explanation**                                                                                               |
|----------------------|--------------|---------------------------------------------------------------------------------------------------------------|
| **Follow-User-Id**    | All          | The `User-Id` of the user being followed.                                                                     |
| **Follow-User-Id**    | Nostr        | Public key identifiers of users the publisher subscribes to.                                                  |
| **Follow-Tag**        | Nostr        | A list of hashtags the user subscribes to (contained in the data field).                                      |
| **Follow Event**      | Bluesky      | No specific protocol tags except the user being followed. Unfollow events come from the delete stream.        |
| **Follow Event**      | Farcaster    | Same as Bluesky, with delete events for unfollowing users.                                                    |


### Profile
Profile events contain metadata and personal information about users, such as usernames, display names, avatars, and other personal details. This section covers how each protocol structures and manages user profiles.

| **Field**                          | **Protocol**     | **Explanation**                                                                                                                                            |
| ---------------------------------- | ---------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Display-Name**                   | All              | The name the user wishes to be displayed with their profile.                                                                                               |
| **Avatar-Image**                   | All              | URL pointing to the user’s avatar image. For Bluesky, this contains additional fields for image mimetype and size.                                         |
| **User-Name**                      | Nostr, Farcaster | A short handle for the user, usually containing no spaces. (Bluesky does not have `User-Name`.)                                                            |
| **Bio**                            | Nostr, Farcaster | Arbitrary-length text field containing personal details, hashtags, or references to other users and events.                                                |
| **Website**                        | Nostr, Farcaster | URL for a website the user wants to advertise (e.g., GitHub profile).                                                                                      |
| **Banner-Image**                   | Nostr, Bluesky   | URL pointing to the user’s banner image, usually at the top of the profile page. For Bluesky, this contains additional fields for image mimetype and size. |
| **Content-Warning**                | Nostr            | Warning for sensitive content associated with the profile.                                                                                                 |
| **Verification**                   | Nostr            | For Nostr, this means a [nip-05](https://github.com/nostr-protocol/nips/blob/master/05.md) address (e.g., `user@example.com`).                             |
| **Payment-Address**                | Nostr            | A Lightning Network payment address, typically in LUD16 format (e.g., `user@example.com`).                                                                 |
| **Mention-Event-Id**               | Nostr            | Event-Id values appearing in the bio field.                                                                                                                |
| **Hashtag**                        | Nostr            | Hashtags mentioned in the bio.                                                                                                                             |
| **Source**                         | Nostr            | Designates the name of the protocol source (e.g., Mastodon via Mostr bridge relays).                                                                       |
| **Source-Uri**                     | Nostr            | URI to search for the event from the protocol source.                                                                                                      |
| **Avatar-Image**, **Banner-Image** | Bluesky          | Contain additional fields for image mimetype and size.                                                                                                     |

