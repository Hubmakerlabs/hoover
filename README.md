

# Hoover - Decentralized Social Network Data Aggregator
<img src="doc/hoover.PNG" width="30%" />
Hoover is a data aggregator built for Arweave developers to collect and unify event data from decentralized social networks like Nostr, Bluesky, and Farcaster. By consolidating this data on Arweave, Hoover provides a unified feed that developers can leverage for various decentralized applications, helping to bridge the gap between fragmented communities spread across different social media protocols.

## Key Features

- Aggregates data from Nostr, Bluesky, and Farcaster into a common [schema](#architecture-and-data-formats).
- Stores data on Arweave for long-term, decentralized storage.
- Supports five primary event types: `Post`, `Repost`, `Like`, `Follow`, `Profile`. 
- Includes a basic event [browser](doc/browser.md) for viewing and searching indexed data. Additionally, this browser performs in-place signature validations of Nostr and Farcaster posts ensuring their integrity.

## Supported Social Networks

- **Nostr**: A minimalist, decentralized protocol with a publish-subscribe architecture.
- **Bluesky**: A federated protocol developed by Bluesky, Inc., with strict control over the specification.
- **Farcaster**: Built on Ethereum, Farcaster uses a blockchain-based model for social interactions.

## Getting Started

To get started with Hoover, follow these steps:

1. **Install Go** (version 1.22.6 or higher) [Go Download](https://go.dev/dl)
2.  **Generate Key**: Generate new RSA key or load an old one for Arweave uploading
   
      Run the following command from the project root:
      ```bash
      go run ./cmd/keygen/. > keyfile.json
      ```
4. **Run the Hoover**: You can now start Hoover:
   ```bash
   ARWEAVE_GATEWAYS=https://up.arweave.net WALLET_FILE=keyfile.json go run ./cmd/hoover/.
   ```
5. **Start the Browser Interface**: Finally, start the browser to view and interact with the Hoover data:
   ```bash
   cd browser
   npm install
   npm run dev
   ```
   Upon running this command, you should be presented with a clickable http link which allows you to view the Hoover's output.
   > To learn more about the browser, click [here](doc/browser.md)
   

> To run the `hoover` on arlocal for testing purposes, see the [Testing Environment Setup](doc/testing.md)

## Architecture and Data Formats

Hoover works by aggregating data from decentralized social networks and formatting it for storage on Arweave. It supports five primary event types:

- `Post` - Basic text content, potentially containing links and references to other users or events.
- `Repost` - An event re-published by a user for their followers.
- `Like` - A user’s positive reaction to a post or event.
- `Follow` - A public declaration of subscribing to another user.
- `Profile` - User metadata, such as display name, avatar, and other contact information.

For full details on data formats and bundling principles, see the [Data Specification](doc/data-spec.md).

## Project Milestones

Our project milestones cover the development process, goals, and task completion. Check out the detailed milestones in the [Project Milestones Document](doc/milestones.md).

## Credits

This project was developed by [Hubmaker Labs](https://hubmaker.io/), a company specializing in decentralized applications and infrastructure, under contract with [Arweave](https://www.arweave.org), whose support and funding made this work possible.

### Team:
- **David Vennik** (Nostr & Bluesky Aggregation, Arweave Test Harness) - <me@mleku.dev>
- **Akash Kumar** (Farcaster Aggregation, Arweave Integration) - <akashkumar1691@gmail.com>
- **Selami Cetinguney** (Data Browser example app) - <selami.c@sbytes-it.com>

We are grateful to Arweave for enabling the development of Hoover and its role in advancing decentralized, long-term data storage.



