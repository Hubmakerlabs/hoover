package farcaster

import (
	"encoding/json"

	farcaster "github.com/ertan/go-farcaster/pkg"
	"github.com/spf13/viper"
)

func prettyPrint(st interface{}) {
	stJson, err := json.Marshal(st)
	if err != nil {
		panic(err)
	}
	println(string(stJson))
}

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	apiUrl := viper.Get("FARCASTER_API_URL").(string)
	mnemonic := viper.Get("FARCASTER_MNEMONIC").(string)
	providerWs := viper.Get("ETHEREUM_PROVIDER_WS").(string)
	farcaster := farcaster.NewFarcasterClient(apiUrl, mnemonic, providerWs)
	println("Farcaster client created")

	// Recent casts
	casts, _, err := farcaster.Casts.GetRecentCasts(10)
	if err != nil {
		panic(err)
	}
	println("Recent casts fetched")
	prettyPrint(casts)
	// Casts by user
	casts, _, err = farcaster.Casts.GetCastsByFname("ertan", 0, "")
	if err != nil {
		panic(err)
	}
	prettyPrint(&casts[0])
	// Cast by hash
	// cast, err := farcaster.Casts.GetCastByHash("new hash format here")
	// if err != nil {
	// 	panic(err)
	// }
	// prettyPrint(&cast)

	// Uncomment for mutating examples
	// // Publish cast
	// cast, err = farcaster.Casts.PublishCast("Testing client")
	// if err != nil {
	// 	panic(err)
	// }
	// prettyPrint(&cast)
	// // Publish a reply cast
	// threadHash := cast.ThreadHash
	// fid, err := farcaster.Registry.GetFidByFname("ertan")
	// if err != nil {
	// 	panic(err)
	// }
	// cast, err = farcaster.Casts.PublishReplyCast("Testing thread", fid, threadHash)
	// if err != nil {
	// 	panic(err)
	// }
	// prettyPrint(&cast)
	// // Clean up published casts
	// err = farcaster.Casts.DeleteCast(cast.Hash)
	// if err != nil {
	// 	panic(err)
	// }
	// println("Cast deleted:", cast.Hash)
	// err = farcaster.Casts.DeleteCast(threadHash)
	// if err != nil {
	// 	panic(err)
	// }
	// println("Cast deleted:", threadHash)
}
