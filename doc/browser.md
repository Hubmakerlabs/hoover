# Hoover Data Browser

## Overview

The **Hoover Data Browser** is a simple, single-page application that displays the latest social media events aggregated by the Hoover from decentralized networks such as Nostr, Bluesky, and Farcaster. These events are stored on Arweave, and the browser allows users to view and step backward through the history of events, demonstrating how the Hoover system aggregates and manages decentralized social media data.

The primary purpose of the Hoover Data Browser is to showcase how data from multiple decentralized networks is collected and stored on Arweave, as well as provide basic verification for the posts from each protocol. This verification helps ensure the integrity of the event data and confirms that posts have not been tampered with.

## Features

### 1. Display Latest Events

The browser fetches and displays the 25 most recent social media events stored on an Arweave gateway. It provides users the ability to navigate backward through earlier events, giving insight into how data flows from the decentralized networks into Arweave.

### 2. Signature Verification

The Hoover Data Browser performs basic signature verification for supported protocols. Verification status is represented visually within the interface (e.g., a green checkmark icon in the top right-hand corner of validated posts). <img src="doc/verified.png" width="5%" />


The browser supports signature verification for the following protocols:

- **Nostr (Signature Type 3)**:  
  Verified using the `schnorr` signature scheme on the secp256k1 curve.
  
- **Farcaster (Signature Types 1 and 2)**:  
  - Type 1: Verified using the `ed25519` curve.
  - Type 2: Verified using Ethereum smart contract signatures through an Ethereum provider.

- **Bluesky (Signature Type 0)**:  
  Currently, no signature verification is implemented for Bluesky events. Signature Type 0 indicates that verification is not performed in the browser at this time.

### 3. Output

As a proof-of-concept, the browser’s output is minimal and focuses on the core functionality of retrieving and verifying event data. Some non-essential fields from the event data are not displayed. For an explanation of each field found in the outputted posts, please see the  [Data Specification document](data-spec.md)

### 4. Running the Browser
Please refer to the [Getting Started](../README.md#getting-started) section of the main README for instructions on running the browser.

### Conclusion
The Hoover Data Browser is a simple yet effective tool to demonstrate how decentralized social media data can be aggregated, stored on Arweave, and verified using the Hoover system. While the current version focuses on basic event display and verification, it lays the groundwork for more advanced features, such as Bluesky verification and improved browsing capabilities.


