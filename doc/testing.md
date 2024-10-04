# Setting up a Test Environment

To see the Hoover in action, there are a number of things you need to do:

## Go

You need a working, preferably latest version of Go, any version since 1.22.6 should be fine, as these new versions can automatically install newer versions if the `go.mod` requires it.

Get it here: [https://go.dev/dl](https://go.dev/dl) and follow the instructions given for installation found at the top of that page.

## Arlocal

Arlocal is a simple developer testnet server for Arweave that provides something along similar lines to what is commonly called a "simnet" on many blockchains - you have a tool to mint arbitrary amounts of tokens for a given wallet address, and a tool that bundles a set of transactions into a block on demand.

[arlocal](https://github.com/textury/arlocal) can be installed simply with the instructions provided on the page of that link with nodejs 20+ installed (with npm available):

    npx arlocal

This will set it up to run on the default network listener http://localhost:1984

If you want to expose it to a network port, you can use:

    npx arlocal -- --host

Run this in its own separate terminal instance, and keep it available to watch by using either your GUI's paning interface or if you have a paning capable terminal empulator, split the window.

Note that by default, `arlocal` deletes its database at every run, see the github page for information about how to make it persist the database instead.

## Test Harness

In order to have the network actually run like it is a real network, you need to prompt it to mint tokens and mine blocks. So, you need a test wallet, and one can be found in `cmd/testharnes/keyfile.json` which has the address `27xHJ0MNsBUKFIdOiQ3OlrZdDzSNfBPGnp6YVmWKKxU`

From the root of the repository, you can run the following command, with your Go installation set up to be available from your `PATH`:

    go run ./cmd/testharness/. http://localhost:1984 27xHJ0MNsBUKFIdOiQ3OlrZdDzSNfBPGnp6YVmWKKxU 1000

The test harness starts up and immediately supplements the given wallet balance in the second parameter to the number in the third parameter, and the first parameter sets the address of the `arlocal` instance.

Every second it checks the balance, if it falls too low to send transactions it mints new tokens for the address, and then mines a block.

When you run the test harness, in a separate terminal, you will then see that `arlocal` reports mining new blocks and minting tokens.

## Running the Hoover

In order to run the hoover, you need to set the `WALLET_FILE` environment variable to point to a valid arweave wallet JSON file, which we provide in the repository for test usage:

    WALLET_FILE=cmd/testharness/keyfile.json go run ./cmd/hoover/.

Same as the other commands, this needs to be run from the root of the repository for the paths to be correctly found.

This will then start up the Nostr, Farcaster and Bluesky firehose feeds and bundle up the events in the format described in [data-spec.md](data-spec.md) and publish them to the `arlocal` instance.

## Browser

To actually see some of this data, you can then run the browser web app, which is a simple single page that shows 25 items from the newest on the `arlocal` instance and has the ability to step backwards through the history prior to this and see more events.

> The output is very rudimentary, and only the Nostr event signatures are verified because of the circuitous processes required to perform this on the other two protocols, and we have not implemented this as a result for this proof of concept/tooling project.

This is a simple demo that shows that the data can be accessed on the `arlocal` endpoint, just as it would on the live mainnet, and a more sophisticated browsing application could be built that enables you to unwrap the embeds and links in the `Content` fields of the Hoover bundles and linkify them and enable the ability to browse through the event history as it accumulates.

Run the browser as follows:

    cd browser
    npm run dev

and it will present you with something like this:

```
VITE v5.4.8  ready in 276 ms

➜  Local:   http://localhost:5173/
➜  Network: use --host to expose
➜  press h + enter to show help
```

you should be able to then click on the linkified http link there and it will open the web browser and show you the newest 25 events that arlocal has received.

## Conclusion

With the foregoing instructions, you will be able to see the Hoover in action, and with the help of some Go and Javascript programmers you will be able to build a custom Arweave permaweb app that lets users browse data from these social network protocols.