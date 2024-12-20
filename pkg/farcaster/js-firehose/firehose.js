const { getSSLHubRpcClient, HubEventType } = require('@farcaster/hub-nodejs');
const fs = require('fs');

const hubRpcEndpoints = [
  'hoyt.farcaster.xyz:2283',
  'hub-grpc.pinata.cloud',
  'nemes.farcaster.xyz:2283'
];
let currentEndpointIndex = 0;

// Update the file path to output.jsonl
const outputFilePath = 'output.jsonl';

async function connectToHub() {
  const hubRpcEndpoint = hubRpcEndpoints[currentEndpointIndex];
  const client = getSSLHubRpcClient(hubRpcEndpoint);

  client.$.waitForReady(Date.now() + 5000, async (e) => {
    if (e) {
      console.error(`Failed to connect to ${hubRpcEndpoint}:`, e);

      // Retry connection with the next endpoint if available
      currentEndpointIndex++;
      if (currentEndpointIndex < hubRpcEndpoints.length) {
        console.log(`Retrying with next endpoint: ${hubRpcEndpoints[currentEndpointIndex]}`);
        await connectToHub(); // Retry with the next endpoint
      } else {
        console.error('All connection attempts failed.');
        process.exit(1);
      }
    } else {
      console.log(`Connected to ${hubRpcEndpoint}`);

      const subscribeResult = await client.subscribe({
        eventTypes: [HubEventType.MERGE_MESSAGE],
      });

      if (subscribeResult.isOk()) {
        const stream = subscribeResult.value;
        const fileStream = fs.createWriteStream(outputFilePath, { flags: 'a' });

        for await (const event of stream) {
          // Write each event as a new line in the JSONL file
          const eventJson = JSON.stringify(event) + '\n';
          fileStream.write(eventJson);

          // Print the cast to the console
          console.log('New cast added:', event);
        }

        fileStream.end();
      } else {
        console.error('Failed to subscribe:', subscribeResult.error);
      }

      client.close();
    }
  });
}

// Start the connection process
connectToHub();
