const { getSSLHubRpcClient } = require('@farcaster/hub-nodejs');
const fs = require('fs');

const hubRpcEndpoint = 'hoyt.farcaster.xyz:2283';
const client = getSSLHubRpcClient(hubRpcEndpoint);
const outputFilePath = 'output.txt';
client.$.waitForReady(Date.now() + 5000, (e) => {
  if (e) {
    console.error(`Failed to connect to ${hubRpcEndpoint}:`, e);
    process.exit(1);
  } else {
    console.log(`Connected to ${hubRpcEndpoint}`);

    client.getCastsByFid({ fid: 5650 }).then((castsResult) => {
          // Access the messages within the value property
      if (castsResult && castsResult.value && Array.isArray(castsResult.value.messages)) {
        const messages = castsResult.value.messages.map((message) => JSON.stringify(message)).join('\n');
        
        // Write the messages to a file
        fs.writeFileSync(outputFilePath, messages, 'utf8');
        console.log(`Output saved to ${outputFilePath}`);
      } else {
        console.error('Unexpected structure of castsResult:', castsResult);
      }

      client.close();
    }).catch((err) => {
      console.error('Error fetching casts:', err);
      client.close();
    });
  }
});
