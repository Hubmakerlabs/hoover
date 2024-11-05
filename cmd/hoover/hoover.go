package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/goar/utils"

	"github.com/Hubmakerlabs/hoover/pkg/config"
	"github.com/Hubmakerlabs/hoover/pkg/multi"
	"github.com/Hubmakerlabs/replicatr/pkg/apputil"
	"github.com/Hubmakerlabs/replicatr/pkg/interrupt"
	"github.com/Hubmakerlabs/replicatr/pkg/slog"
	"go-simpler.org/env"
)

type Config struct {
	AppName          string   `env:"APP_NAME" default:"hoover"`
	Root             string   `env:"ROOT_DIR" usage:"root path for all other path configurations (defaults OS user home if empty)"`
	Profile          string   `env:"PROFILE" default:".hoover" usage:"name of directory in root path to store state data and database"`
	WalletFile       string   `env:"WALLET_FILE" default:"keyfile.json" usage:"full path of wallet file to use for uploading to arweave"`
	SpeedFactor      float64  `env:"SPEED_FACTOR" default:"1" usage:"change priority of bundle uploads by this ratio"`
	ArweaveGateways  []string `env:"ARWEAVE_GATEWAYS" usage:""`
	NostrRelays      []string `env:"NOSTR_RELAYS" usage:"nostr relays, comma separated, in standard format 'wss://example.com', ports and insecure ws:// also permissible"`
	FarcasterHubs    []string `env:"FARCASTER_HUBS" usage:"farcaster hub network addresses, comma separated"`
	BlueskyEndpoints []string `env:"BLUESKY_ENDPOINTS" usage:"bluesky endpoints to use, comma separated"`
}

const DefaultEnv = `APP_NAME=hoover
ROOT_DIR=
PROFILE=.hoover
WALLET_FILE=keyfile.json
SPEED_FACTOR=1
ARWEAVE_GATEWAYS=http://localhost:1984
NOSTR_RELAYS=wss://purplepag.es,wss://njump.me,wss://relay.snort.social,wss://relay.damus.io,wss://relay.primal.net,wss://relay.nostr.band,wss://nostr-pub.wellorder.net,wss://relay.nostr.net,wss://nostr.lu.ke,wss://nostr.at,wss://e.nos.lol,wss://nostr.lopp.social,wss://nostr.vulpem.com,wss://relay.nostr.bg,wss://wot.utxo.one,wss://nostrelites.org,wss://wot.nostr.party,wss://wot.sovbit.host,wss://wot.girino.org,wss://relay.lnau.net,wss://wot.siamstr.com,wss://wot.sudocarlos.com,wss://relay.otherstuff.fyi,wss://relay.lexingtonbitcoin.org,wss://wot.azzamo.net,wss://wot.swarmstr.com,wss://zap.watch,wss://satsage.xyz
FARCASTER_HUBS=hub.pinata.cloud,api.hub.wevm.dev,hoyt.farcaster.xyz:2283,lamia.farcaster.xyz:2283,api.farcasthub.com:2283,nemes.farcaster.xyz:2283,hub.farcaster.standardcrypto.vc:2281,hoyt.farcaster.xyz:2281,lamia.farcaster.xyz:2281,api.farcasthub.com:2281,nemes.farcaster.xyz:2281,hub.farcaster.standardcrypto.vc:2282,hoyt.farcaster.xyz:2282,lamia.farcaster.xyz:2282,api.farcasthub.com:2282,nemes.farcaster.xyz:2282,hub.farcaster.standardcrypto.vc:2283,
BLUESKY_ENDPOINTS=wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos
`

// NewConfig creates a new Config struct and loads it from environment variables
// and .env file based on a configured ROOT_DIR/PROFILE directory location which
// can be overridden in the environment to refer to a custom location.
//
// Any values set in the environment will override those in the .env file that
// is loaded to allow one-shot overrides where desired.
func NewConfig() (cfg *Config, err error) {
	cfg = &Config{}
	// load the environment variables so we can catch the root/profile folders if
	// they have been set, and alter where we look for the env file
	if err = env.Load(cfg, &env.Options{SliceSep: ","}); err != nil {
		return
	}
	if cfg.Root == "" {
		var dir string
		if dir, err = os.UserHomeDir(); err != nil {
			return
		}
		cfg.Root = dir
	}
	// if the env file exists where the env root/profile vars say it is, load it, and then
	// override the environment variables that have been set over top so the user can set custom
	// values for one time usage or for purposes for running
	envPath := filepath.Join(filepath.Join(cfg.Root, cfg.Profile), ".env")
	if !apputil.FileExists(envPath) {
		EnsureDir(envPath, 0700)
		// write the default env and load it
		if err = os.WriteFile(envPath, []byte(DefaultEnv), 0600); err != nil {
			return
		}
	}
	var e config.Env
	if e, err = config.GetEnv(envPath); err != nil {
		return
	}
	// load the env file first
	if err = env.Load(cfg, &env.Options{Source: e, SliceSep: ","}); err != nil {
		return
	}
	// load the environment vars again so they can override the .env file (only the
	// vars that have been set will be overwritten).
	if err = env.Load(cfg, &env.Options{SliceSep: ","}); err != nil {
		return
	}
	return
}

// EnsureDir checks a file could be written to a path, creates the directories
// as needed
func EnsureDir(fileName string, perms fs.FileMode) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, perms)
		if merr != nil {
			panic(merr)
		}
	}
}

// FileExists reports whether the named file or directory exists.
func FileExists(filePath string) bool {
	_, e := os.Stat(filePath)
	return e == nil
}

// HelpRequested returns true if any of the common types of help invocation are
// found as the first command line parameter/flag.
func HelpRequested() (help bool) {
	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {
		case "help", "-h", "--h", "-help", "--help", "?":
			help = true
		}
	}
	return
}

// PrintHelp outputs a help text listing the configuration options and default
// values to a provided io.Writer (usually os.Stderr or os.Stdout).
func PrintHelp(cfg *Config, printer io.Writer) (s string) {
	_, _ = fmt.Fprintf(printer,
		"Environment variables that configure %s:\n\n", cfg.AppName)
	env.Usage(cfg, printer, &env.Options{SliceSep: ","})
	_, _ = fmt.Fprintf(printer,
		"\nCLI parameter 'help' also prints this information\n"+
			"\n.env file found at the ROOT_DIR/PROFILE path will be automatically "+
			"loaded for configuration.\nset these two variables for a custom load path,"+
			" this file will be created on first startup.\nenvironment overrides it and "+
			"you can also edit the file to set configuration options\n")
	return
}

// GetWallet loads an arweave wallet JSON file from the provided path for connection with a
// specified arweave gateway endpoint.
func GetWallet(walletFile, endpoint string) (address string, wallet *goar.Wallet, err error) {
	var walletJSON []byte
	if walletJSON, err = os.ReadFile(walletFile); err != nil {
		return
	}
	if wallet, err = goar.NewWallet(walletJSON, endpoint); err != nil {
		return
	}
	address = wallet.Signer.Address
	return
}

// const batchSize = 75000

func main() {
	slog.SetLogLevel(slog.Info)
	var err error
	var cfg *Config
	if cfg, err = NewConfig(); err != nil || HelpRequested() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err)
		}
		PrintHelp(cfg, os.Stderr)
		os.Exit(0)
	}
	if !FileExists(cfg.WalletFile) {
		fmt.Fprintf(os.Stderr, "ERROR: wallet json file `%s` not found\n", cfg.WalletFile)
		PrintHelp(cfg, os.Stderr)
		os.Exit(1)
	}
	var address string
	var gateway string
	var wallet *goar.Wallet
	for _, gateway = range cfg.ArweaveGateways {
		if address, wallet, err = GetWallet(cfg.WalletFile, gateway); err != nil {
			continue
		}
		// successfully loaded wallet for provided gateway, continue
		log.I.F("uploading to arweave gateway %s using wallet address %s", gateway, address)
		break
	}
	c, cancel := context.WithCancel(context.Background())
	batchChan := make(chan *types.BundleItem)
	var itemSigner *goar.ItemSigner
	if itemSigner, err = goar.NewItemSigner(wallet.Signer); chk.E(err) {
		return
	}
	// bundle batcher worker
	go func() {
		for {
			select {
			case <-c.Done():
				return
			case bundle := <-batchChan:
				if bundle == nil {
					continue
				}
				bundle.SignatureType = types.ArweaveSignType
				var item types.BundleItem
				if item, err = itemSigner.CreateAndSignItem([]byte(bundle.Data), bundle.Target, bundle.Anchor, bundle.Tags); chk.E(err) {
					continue
				} else {
					var resp *types.BundlrResp
					if resp, err = utils.SubmitItemToBundlr(item, gateway); chk.E(err) {
						log.E.F("failed to submit item to bundlr: %s", err)
						continue
					}
					log.I.F("successfully submitted item to bundlr, bundler response id: %s", resp.Id)
				}

			}
		}
	}()
	interrupt.AddHandler(cancel)
	var wg sync.WaitGroup

	multi.Firehose(c, cancel, &wg, cfg.NostrRelays, cfg.BlueskyEndpoints, cfg.FarcasterHubs,
		func(bundle *types.BundleItem) (err error) {
			batchChan <- bundle
			return
		})
}
