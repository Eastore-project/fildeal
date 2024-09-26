# Fildeal CLI

Fildeal is a command-line interface (CLI) tool for managing Filecoin deals. This tool provides various commands to compare files, generate data segment pieces, split pieces, and initiate deals with miners. It is inspired by the [mkpiece](https://github.com/willscott/mkpiece) tool and [data-segment-library](https://github.com/filecoin-project/go-data-segment).
One can run a 2k lotus-miner setup easily using the [scripts](https://gist.github.com/lordshashank/fb2fbd53b5520a862bd451e3603b4718). And then use this tool to initiate boost deals with the miner and test various things.

## Pre-requisites

- [Go](https://golang.org/doc/install) 1.20 or higher
- [Boost](https://boost.filecoin.io/getting-started)

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/yourusername/fildeal.git
   cd fildeal
   ```

2. **Build the CLI:**

   ```sh
   go build -o fildeal main.go
   ```

3. Setup the environment variables:

   ```sh
   export FULLNODE_API_INFO="your fullnode api info" # https://api.calibration.node.glif.io for calibration testnet
   export LOTUS_API_INFO="your miner api info" # Required for --testnet deals
   ```

## Known Error

While building the CLI, you might encounter the following error:

```sh
# github.com/ipfs/go-unixfs/hamt
../../../go/pkg/mod/github.com/tech-greedy/go-unixfs@v0.3.2-0.20220430222503-e8f92930674d/hamt/hamt.go:765:19: assignment mismatch: 2 variables but bitfield.NewBitfield returns 1 value
```

This maybe because [tech-greedy/go-unixfs](https://github.com/tech-greedy/generate-car) is not maintained and is not compatible with the latest version of go-unixfs.

To fix this, you can replace the hamt.go file in your packages with the hamt.txt file provided in the repository.

## Usage

The `fildeal` CLI provides the following commands:

- `cmp <parentFile> <childFile>`: Compare two files and find the offset of the child file in the parent file.
- `generate <files...>`: Generate a data segment piece from the given files and output it to stdout.
- `splitpiece <file> <outputDir>`: Split the specified file into pieces and save them in the output directory.
- `initiate <inputFolder> <miner> [--server]`: Initiate a deal with the specified input folder and miner. Optionally, start a server.

### Detailed Usage Information

```sh
Usage: fildeal <command> [arguments]

Commands:
cmp <parentFile> <childFile>       Compare two files and find the offset of the child file in the parent file.
generate <files...>                Generate a data segment piece from the given files and output it to stdout.
splitpiece <file> <outputDir>      Split the specified file into pieces and save them in the output directory.
initiate <inputFolder> <miner>     Initiate a deal with the specified input folder and miner.

Examples:
fildeal cmp a.car b.car
fildeal generate a.car b.car c.car > out.dat
fildeal splitpiece input.car outputDir
fildeal initiate inputFolder miner [--server] [--testnet]
```

### Commands

#### `cmp`

The `cmp` command compares two files and finds the offset of the child file in the parent file.

```sh
fildeal cmp <parentFile> <childFile>
```

#### `generate`

The `generate` command generates a data segment piece from the given files and outputs it to stdout.

```sh
fildeal generate <files...>
```

#### `splitpiece`

The `splitpiece` command splits the specified file into pieces and saves them in the output directory.

```sh
fildeal splitpiece <file> <outputDir>
```

#### `initiate`

The `initiate` command initiates a deal with the specified input folder and miner. Optionally, it starts a server to serve files to local miners. If you want to host file somewhere else, you can put url [here](https://github.com/lordshashank/filecoin-deals/blob/94338a0ac8338dd0a792ac913555aba82577da89/src/deal/utils/initiateDeal.go#L17).
InputFolder should be the location of files you want to initiate the deal with.

```sh
fildeal initiate <inputFolder> <miner> [--server] [--testnet]
```

While initiating a deal with 2k miner, or any other miner, you would have to have a boost wallet with funds to make deal and FULLNODE_API_INFO set in the environment variables.

`--server` Flag

Starts a server to serve files to local miners. This is useful when you want to host the files locally and make them accessible to the miners in your network.

`--testnet` Flag

Indicates that the deal should be initiated on the Filecoin testnet. This is useful for testing purposes, allowing you to simulate deals without using real Filecoin tokens.
Few important things to note for making deals on testnet:

- You would need to have filecoin tokens and datacap to make currently supported verified deals on testnet. Get them from [faucet](https://faucet.calibnet.chainsafe-fil.io/).
- You would need to host the deal CAR file somewhere to serve them to testnet miner. `fildeal` currently supports [lighthouse](https://www.lighthouse.storage/) as the hosting service. You would need to have `LIGHTHOUSE_API_KEY` set in the environment variables. You can get the api key by following [this](https://docs.lighthouse.storage/lighthouse-1/how-to/create-an-api-key).
