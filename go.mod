module github.com/Hubmakerlabs/hoover

go 1.22.5

require (
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.3.2
	github.com/Hubmakerlabs/replicatr v1.2.17
	github.com/Microsoft/go-winio v0.6.2
	github.com/VictoriaMetrics/fastcache v1.12.2
	github.com/aws/aws-sdk-go-v2 v1.30.3
	github.com/aws/aws-sdk-go-v2/config v1.27.26
	github.com/aws/aws-sdk-go-v2/credentials v1.17.26
	github.com/aws/aws-sdk-go-v2/service/route53 v1.42.3
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/cespare/cp v1.1.1
	github.com/cloudflare/cloudflare-go v0.99.0
	github.com/cockroachdb/errors v1.11.3
	github.com/cockroachdb/pebble v1.1.1
	github.com/consensys/gnark-crypto v0.12.1
	github.com/crate-crypto/go-kzg-4844 v1.0.0
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/deckarep/golang-set/v2 v2.6.0
	github.com/dop251/goja v0.0.0-20240707163329-b1681fb2a2f5
	github.com/ethereum/c-kzg-4844/bindings/go v0.0.0-20230126171313-363c7d7593b4
	github.com/everFinance/gojwk v1.0.0
	github.com/everFinance/ttcrsa v1.1.3
	github.com/fatih/color v1.17.0
	github.com/fjl/gencodec v0.0.0-20230517082657-f9840df7b83e
	github.com/fjl/memsize v0.0.2
	github.com/fsnotify/fsnotify v1.7.0
	github.com/gammazero/deque v0.2.1
	github.com/gammazero/workerpool v1.1.3
	github.com/gballet/go-libpcsclite v0.0.0-20191108122812-4678299bea08
	github.com/getsentry/sentry-go v0.28.1
	github.com/gin-contrib/pprof v1.5.0
	github.com/gin-gonic/gin v1.10.0
	github.com/go-resty/resty/v2 v2.13.1
	github.com/go-stack/stack v1.8.1
	github.com/gofrs/flock v0.12.0
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/golang/protobuf v1.5.4
	github.com/golang/snappy v0.0.4
	github.com/google/gofuzz v1.2.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/graph-gophers/graphql-go v1.5.0
	github.com/hamba/avro v1.8.0
	github.com/hashicorp/go-bexpr v0.1.14
	github.com/holiman/billy v0.0.0-20240322075458-72a4e81ec6da
	github.com/holiman/bloomfilter/v2 v2.0.3
	github.com/holiman/uint256 v1.3.0
	github.com/huin/goupnp v1.3.0
	github.com/iancoleman/strcase v0.3.0
	github.com/inconshreveable/log15 v2.16.0+incompatible
	github.com/influxdata/influxdb-client-go/v2 v2.13.0
	github.com/influxdata/influxdb1-client v0.0.0-20220302092344-a9ab5670611c
	github.com/jackc/pgtype v1.14.3
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jackpal/go-nat-pmp v1.0.2
	github.com/jarcoal/httpmock v1.3.1
	github.com/jedisct1/go-minisign v0.0.0-20230811132847-661be99b8267
	github.com/julienschmidt/httprouter v1.3.0
	github.com/karalabe/usb v0.0.2
	github.com/kylelemons/godebug v1.1.0
	github.com/lestrrat-go/jwx v1.2.29
	github.com/lib/pq v1.10.9
	github.com/linkedin/goavro/v2 v2.13.0
	github.com/mattn/go-colorable v0.1.13
	github.com/mattn/go-isatty v0.0.20
	github.com/minio/sha256-simd v1.0.1
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mleku/btcec/v2 v2.3.2-2
	github.com/mleku/nodl v0.0.2
	github.com/olekukonko/tablewriter v0.0.5
	github.com/panjf2000/ants/v2 v2.10.0
	github.com/peterh/liner v1.2.2
	github.com/prometheus/client_golang v1.19.1
	github.com/protolambda/bls12-381-util v0.1.0
	github.com/robfig/cron v1.2.0
	github.com/rs/cors v1.11.0
	github.com/rubenv/sql-migrate v1.7.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/shopspring/decimal v1.4.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/viper v1.19.0
	github.com/status-im/keycard-go v0.3.2
	github.com/stretchr/testify v1.9.0
	github.com/supranational/blst v0.3.12
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/teivah/onecontext v1.3.0
	github.com/tidwall/gjson v1.17.1
	github.com/tidwall/sjson v1.2.5
	github.com/tyler-smith/go-bip39 v1.1.0
	github.com/urfave/cli/v2 v2.27.2
	go.uber.org/atomic v1.11.0
	go.uber.org/ratelimit v0.3.1
	golang.org/x/crypto v0.25.0
	golang.org/x/exp v0.0.0-20240707233637-46b078467d37
	golang.org/x/sync v0.7.0
	golang.org/x/sys v0.22.0
	golang.org/x/text v0.16.0
	golang.org/x/time v0.5.0
	golang.org/x/tools v0.23.0
	gopkg.in/h2non/gentleman.v2 v2.0.5
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/datatypes v1.2.1
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.11
	nhooyr.io/websocket v1.8.11
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.12.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.9.1 // indirect
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.3 // indirect
	github.com/aws/smithy-go v1.20.3 // indirect
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.13.0 // indirect
	github.com/bytedance/sonic v1.11.9 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/cockroachdb/fifo v0.0.0-20240616162244-4768e80dfb9a // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/dlclark/regexp2 v1.11.2 // indirect
	github.com/gabriel-vasile/mimetype v1.4.4 // indirect
	github.com/garslo/gogen v0.0.0-20230926014519-f497ca02dd4c // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-gorp/gorp/v3 v3.1.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.22.0 // indirect
	github.com/go-sourcemap/sourcemap v2.1.4+incompatible // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/pprof v0.0.0-20240711041743-f6c9dda6c6da // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/influxdata/line-protocol v0.0.0-20210922203350-b1ad95c89adf // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kilic/bls12-381 v0.1.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mitchellh/pointerstructure v1.2.1 // indirect
	github.com/mleku/btcec v1.0.1 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32 // indirect
	github.com/oapi-codegen/runtime v1.1.1 // indirect
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.16.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/templexxx/cpu v0.1.0 // indirect
	github.com/templexxx/xhex v0.0.0-20200614015412-aed53437177b // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/tklauser/numcpus v0.2.2 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/mod v0.19.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/term v0.22.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)
