# Arweave AO Social Hoover Project

## Milestones

### Milestone 1:

#### Technical Specification Document  & Proof of Concept

(4 weeks)

##### Analyze the data structures of Nostr, Bluesky and Farcaster

- Devise a data structure layout that takes into account for all of the specific details of each protocol’s message types, enabling all protocol data to be stored validly in Arweave while maintaining the ability to reconstitute the original for cryptographic verification

- Build bidirectional data conversion functions (for validation of conversion) and test them

- Build basic event retrieval functions for rudimentary hoovering of posts from their native sources (Nostr, Bluesky and Farcaster) within a time window

#### Milestone 2:

##### Full-Scale Data Hoover

(4 weeks)

- Spider services that pull the events from the social network data sources at reasonable frequency (5-10 seconds interval) for Nostr, Bluesky and Farcaster

- An ao smart contract that given an Arweave-formatted social media post, that can convert it to its canonical (for signing) form (Nostr, Bluesky, and Farcaster) and verify the signatures of the native data within the smart contract
