package main

import (
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar"
)

func main() {
	arClient := goar.NewClient("https://arweave.net")
	var data []byte
	var err error
	if data, err = arClient.GraphQL(`{
		transactions(
			first: 1,
			sort: HEIGHT_DESC,
			tags: [
				{
					name: "App-Name",
					values: ["com.hubmakerlabs.hoover"]
				}
			]
		)
		{
			pageInfo {
				hasNextPage
			}
			edges {
				cursor
				node {
					id
					signature
					recipient
					owner {
						address
					}
					tags {
						name
						value
					}
				}
			}
		}
	}
`); chk.E(err) {
		log.E.F("failed to query arweave: %s", err)

	}
	log.I.F("arweave query response: %s", string(data))
}
