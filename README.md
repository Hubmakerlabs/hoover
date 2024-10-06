

# Hoover - Decentralized Social Network Data Aggregator
<img src="doc/hoover.PNG" width="20%" />
Hoover is a data aggregator that collects and indexes event data from decentralized social networks, including Nostr, Bluesky, and Farcaster, and stores it on Arweave. This enables Arweave developers to integrate social network data into their distributed applications, enhancing user experiences in areas such as media sharing, community building, and more.

## Key Features

- Aggregates data from Nostr, Bluesky, and Farcaster into a common schema.
- Stores data on Arweave for long-term, decentralized storage.
- Supports five primary event types: `Post`, `Repost`, `Like`, `Follow`, `Profile`.
- Includes a basic event browser for viewing and searching indexed data.

## Supported Social Networks

- **Nostr**: A minimalist, decentralized protocol with a publish-subscribe architecture.
- **Bluesky**: A federated protocol developed by Bluesky, Inc., with strict control over the specification.
- **Farcaster**: Built on Ethereum, Farcaster uses a blockchain-based model for social interactions.

## Getting Started

To get started with Hoover, follow these steps:

1. **Install Go** (version 1.22.6 or higher) [Go Download](https://go.dev/dl)
2. **Set up Arlocal** for local Arweave testing [Arlocal Setup](https://github.com/textury/arlocal)
3. **Run the Test Harness**:
   The test harness will simulate a running Arweave environment.

   Run the following command from the project root:
   ```bash
   go run ./cmd/testharness/. http://localhost:1984 27xHJ0MNsBUKFIdOiQ3OlrZdDzSNfBPGnp6YVmWKKxU 1000
   ```
4. **Run the Hoover**: With the test harness running, you can now start Hoover:
   ```bash
   WALLET_FILE=cmd/testharness/keyfile.json go run ./cmd/hoover/.
   ```
5. **Start the Browser Interface**: Finally, start the browser to view and interact with the Hoover data:
   ```bash
   cd browser
   npm run dev
   ```
   Upon running this command, you should be presented with a clickable http link which allows you to view the Hoover's output.
  
> For a full guide, see the [Testing Environment Setup](doc/testing.md)

## Architecture & Data Formats

Hoover works by aggregating data from decentralized social networks and formatting it for storage on Arweave. It supports five primary event types:

- `Post` - Basic text content, potentially containing links and references to other users or events.
- `Repost` - An event re-published by a user for their followers.
- `Like` - A user’s positive reaction to a post or event.
- `Follow` - A public declaration of subscribing to another user.
- `Profile` - User metadata, such as display name, avatar, and other contact information.

For full details on data formats and bundling principles, see the [Data Specification](doc/data-spec.md).

## Project Milestones

Our project milestones cover the development process, goals, and task completion. Check out the detailed milestones in the [Project Milestones Document](doc/milestones.md).

## Contributors

- **David Vennik** (Nostr & Bluesky Aggregation, Arweave Test Harness) - <me@mleku.dev>
- **Akash Kumar** (Farcaster Aggregation) - <akashkumar1691@gmail.com>
- **Selami** (Data Browser example app)

