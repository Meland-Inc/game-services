module github.com/Meland-Inc/game-services

go 1.17

require (
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/aliyun/aliyun-oss-go-sdk v2.2.5+incompatible // indirect
	github.com/aws/aws-sdk-go v1.44.91 // indirect
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/crackcomm/go-gitignore v0.0.0-20170627025303-887ab5e44cc3 // indirect
	github.com/dapr/dapr v1.8.0 // indirect
	github.com/dapr/go-sdk v1.5.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/go-kit/kit v0.12.0 // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gookit/color v1.3.6 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/ipfs/go-cid v0.0.7 // indirect
	github.com/ipfs/go-ipfs-api v0.3.0 // indirect
	github.com/ipfs/go-ipfs-files v0.0.9 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/libp2p/go-buffer-pool v0.0.2 // indirect
	github.com/libp2p/go-flow-metrics v0.0.3 // indirect
	github.com/libp2p/go-libp2p-core v0.6.1 // indirect
	github.com/libp2p/go-openssl v0.0.7 // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/minio/sha256-simd v0.1.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.0.3 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-multiaddr v0.3.0 // indirect
	github.com/multiformats/go-multibase v0.0.3 // indirect
	github.com/multiformats/go-multihash v0.0.14 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/qiniu/api.v7/v7 v7.8.2 // indirect
	github.com/rjeczalik/notify v0.9.1 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/spacemonkeygo/spacelog v0.0.0-20180420211403-2296661a0572 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	github.com/tklauser/go-sysconf v0.3.6 // indirect
	github.com/tklauser/numcpus v0.2.2 // indirect
	github.com/whyrusleeping/tar-utils v0.0.0-20180509141711-8c6c8ba81d5c // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/crypto v0.0.0-20220511200225-c6db032c6c88 // indirect
	golang.org/x/net v0.0.0-20220621193019-9d032be2e588 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11 // indirect
	google.golang.org/genproto v0.0.0-20220622171453-ea41d75dfa0f // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.3.6 // indirect
	gorm.io/gorm v1.23.8 // indirect
)

require (
	// 使用本地go代码仓库方式: https://zhuanlan.zhihu.com/p/109828249
	game-message-core v0.0.0
	github.com/ethereum/go-ethereum v1.10.26
)

replace game-message-core => ./src/game-message-core/messageGo
