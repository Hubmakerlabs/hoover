# Arweave AO Social Network Data Hoover

#### Aggregator for social network data for use of Arweave developers as part of Arweave AO enabled applications

-----

## Built under contract with Arweave for Hubmaker Labs

Work completed during August and September 2024

#### Authors

David Vennik <me@mleku.dev> - Nostr, Bluesky Aggregation and Arweave upload

Akash Kumar <akashkumar1691@gmail.com> - Farcaster Aggregation and Arweave Smart Contracts

## Introduction

Arweave is a decentralized blockchain project with a central proposition of the idea of persistent, long term storage. The protocol achieves this via a unique monetization method and scaling replication to enable a virtual large scale computer system with parallel processing and progressive convergence to the state of the database across replicas.

Decentralized social network protocols tend to focus on a lower consistency guarantee in exchange for fast data propagation for users, and Arweave presents an opportunity to increase the interconnectivity of social networks across the internet by leveraging its distributed computation and data storage architecture to create a bridge between the networks and the web as a whole, while minimally sacrificing latency in favor of availability.

The ostensible purpose of building an aggregator - as we term it a "data hoover" for pulling together the "firehose" of decentralized social network data is to serve as a backbone for developers of distributed applications that can be made more useful through such data, such as media sharing applications, from videos, to music to books, to market places and corporate helpdesk and community building for users of applications and services that benefit from Arweave's service architecture.

## Bundling Principles

The Hoover's primary goal is to create a largely uniform schema for data that can be used as a baseline for searching multiple social network data feeds with a common basic index of elements that form critical parts in common with all of the social network data formats to enable a single system for indexing and cross-referencing.

There is many types of data that are published on decentralized social networks that are ephemeral, or highly specific to the design of the social network protocol, and quite irrelevant for the purpose of forming a distributed archive of data with a low latency of updating to current events on these networks.

As such, the Hoover attends to only 5 primary event types, which contain common fields and naming conventions, as well as a small amount of necessary extra metadata required for the different protocols due to their differing architectures. But first, a brief discussion of the protocols we are targeting in this initial release:

### The Differences and Similarities of the Protocols We are Supporting

**Nostr** is an extremely minimal design that focuses on publish subscribe architecture and aims to maximize decentralization, and part of the way it achieves this is through a very small specification for database queries and a collaborative specification forum where developers argue all day long about proposals that they often disagree about needing to be standardized. There is a large number of mostly similar, Twitter/X style discussion clients, and a small number of niche clients that are aiming to support decentralized publishing, supporting distributed teams with software and document repositories, and some specializing in instant messaging.

**Bluesky** is a more centralized design, using a federated architecture with Bluesky, Inc. providing a primary aggregation point, through which a tree of federations feed up their user input , that users can tap at any level of the tree that suits them, and thereby protects the decentralization of access, while keeping the network together as a whole as much as possible. Its specification is more strict and complex, and the company that develops the protocol keeps fairly tight control over the specification. There is a large number of client applications for Bluesky, though their structure and format is more limited than Nostr as regards to the potential use cases, due to the federated distribution architecture.

**Farcaster** is a social network that is based on the idea of anchoring data sets on the Ethereum blockchain, while distributing the replicas of events across a mostly voluntary-run or sponsored IPFS based event data replicas. Its protocol design is simpler than Bluesky, but more rigid than Nostr, and as of writing this document, there is only one main application that can be used to access the data of the network.

## Hoover Common Event Data Specification

There is 5 event kinds that are aggregated by the Hoover:

1. **Profile** - User's cryptographic identities, display names, banners and avatars, biography, website, and payment addresses
2. **Post** - Text messages that may potentially be somewhat structured formats - but mostly a very small subset of what would be Markdown. This may include references such as reply parent and original post references, mentions of other users, searchable hashtags and references to external internet resources
3. **Repost** - These are simple reposts that signify other users have made the original post visible to those that follow them, a means to measuring engagement
4. **Like** - Users indicate their agreement or support of another user's post
5. **Follow** - Users designating which other users they are interested in reading posts from, can be used to build a graph that can help with filtering data according to trust and affinity

These form the most essential parts of keeping track of the activity and the users interactions on these protocols, and what is relevant to us is primarily data that is to be used on a read-only basis but retaining enough information to authenticate the data and refer back to its original sources in order to facilitate potential bidirectional bridging, which we are not implementing but are taking care to ensure that the data structure design does not need substantial changes in order to enable.

## Nostr



## Bluesky



## Farcaster

