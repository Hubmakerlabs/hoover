# Arweave AO Social Hoover Project Roadmap

## Phase 1: Initial Development & Proof of Concept

### 1. Analyze Data Structures  
*(Completed)*
- [x] Analyze the key message types of Nostr, Bluesky, and Farcaster.
- [x] Design a unified data structure that supports the specifics of each protocol, ensuring data can be stored validly in Arweave and reconstituted for cryptographic verification.

### 2. Data Conversion  
*(Completed)*
- [x] Build data conversion functions for following protocols into the unified data structure:
  - [x] Nostr
  - [x] Bluesky
  - [x] Farcaster
- [x] Test conversion functions for each protocol to ensure accuracy and validity.

### 3. Basic Event Retrieval  
*(Completed)*
- [x] Implement basic event retrieval (hoovering) within a time window for:
  - [x] Nostr
  - [x] Bluesky
  - [x] Farcaster

## Phase 2: Full-Scale Data Aggregation & Storage

### 4. Continuous Event Retrieval (Spider Services)  
*(Completed)*
- [x] Implement spider services that retrieve social network events in realtime:
  - [x] Nostr
  - [x] Bluesky
  - [x] Farcaster

### 5. Arweave Smart Contract Development  
*(In Progress)*
- [ ] Develop a smart contract that can:
  - [ ] Verify cryptographic signatures for the Arweave bundled posts from:
    - [x] Nostr
    - [ ] Bluesky
    - [x] Farcaster


## Phase 3: Future Enhancements

### 6. Data Reconstitution  
*(Pending)*
- [ ] Add the ability to fully restore original social media posts from the Arweave-formatted data, especially for more complex protocols like Bluesky and Farcaster. This enhancement will enable Arweave applications not only to display posts from these protocols but also to create and publish posts. Thanks to the unified data format, a single post may be seamlessly shared across Nostr, Bluesky, and Farcaster simultaneously, allowing for true bidirectional interaction across multiple platforms.



### 7. Optimize Smart Contract Functionality  
*(Pending)*
- [ ] Refine smart contract performance to ensure efficient conversion and signature verification within Arweave’s environment, considering constraints like data size and processing power.
---

### Disclaimer on Incomplete Tasks:

The incomplete tasks in the current phase are due to several technical challenges:

1. **Data Reconstitution**: Protocols like Bluesky and Farcaster use highly structured and encoded formats, such as Bluesky’s base32 encoding and Farcaster’s Ethereum-based anchoring, making reconstituting the original data for signing verification more complex than anticipated.
   
2. **Signature Verification**: While Nostr's signature verification is straightforward, Bluesky and Farcaster require additional external data, such as user profile events or blockchain lookups. These protocols introduce more complex steps for verifying signatures inside the smart contract, which necessitates further research and optimization. That being said, the Farcaster verification has been completed and only Bluesky now remains.

3. **Smart Contract Constraints**: The size and performance limitations of Arweave’s smart contract environment make implementing signature verification for multiple protocols challenging. Ensuring scalability without compromising the system’s security or efficiency requires additional time and architectural work.

These challenges are actively being addressed, with future iterations focused on resolving these complexities.
