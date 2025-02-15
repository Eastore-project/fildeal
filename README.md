# Fildeal CLI

Fildeal is a command-line interface (CLI) tool for managing Filecoin deals. This tool provides various commands to compare files, generate data segment pieces, split pieces, and initiate deals with miners. It is inspired by the [mkpiece](https://github.com/willscott/mkpiece) tool and [data-segment-library](https://github.com/filecoin-project/go-data-segment).

For easy testing, you can run a 2k lotus-miner setup using these [scripts](https://gist.github.com/lordshashank/fb2fbd53b5520a862bd451e3603b4718), and then use this tool to initiate boost deals with the miner.

## Pre-requisites

- [Go](https://golang.org/doc/install) 1.20 or higher
- [Boost](https://boost.filecoin.io/getting-started)
- For making deals (podsi-deal command):
  - Boost setup with primary boost-wallet funded
  - DataCap allocation if making verified deals

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/yourusername/fildeal.git
   cd fildeal
   ```

2. **Build the CLI:**

   ```sh
   go build -o fildeal ./cmd/fildeal
   ```

3. **Setup environment variables:**

   ```sh
   # For calibration testnet, use https://api.calibration.node.glif.io
   export FULLNODE_API_INFO="your fullnode api info"
   
   # Required only for testnet deals using lighthouse storage
   export LIGHTHOUSE_API_KEY="your lighthouse api key"
   ```

## Commands

### Compare Files (`cmp`)

Compare two files and find the offset of the child file in the parent file.

```sh
fildeal cmp --parent <parentFile> --child <childFile>
# or using short flags
fildeal cmp -p <parentFile> -c <childFile>
```

### Generate Data Segment (`podsi-aggregate`)

Generate a data segment piece from all files in an input folder.

```sh
fildeal podsi-aggregate --input <inputFolder> --output <outputFile>
# or using short flags
fildeal podsi-aggregate -i <inputFolder> -o <outputFile>
```

### Split Piece (`splitpiece`)

Split a podsi-aggregate output file into pieces and save them in the output directory.

```sh
fildeal splitpiece --input <inputFile> --output <outputDir>
# or using short flags
fildeal splitpiece -i <inputFile> -o <outputDir>
```

### Parse Boost Index (`boost-index`)

Parse and index an aggregate file similar to Boost.

```sh
fildeal boost-index <file>
```

### Initiate Deal (`podsi-deal`)

Initiate a deal with a miner using podsi-aggregate for folder aggregation. 

```sh
fildeal podsi-deal --input <inputFolder> --miner <minerID> [options]
```

#### Options:
- `--input, -i`: Input folder containing files to make deal with (required)
- `--miner, -m`: Miner ID to make the deal with (required)
- `--generate-car-path`: Path for generated CAR files (default: "generated_car/")
- `--aggregate-car-path`: Path for aggregate CAR files (default: "aggregate_car_file/")
- `--buffer`: Buffer to use (localhost or lighthouse, default: "localhost")
- `--duration`: Deal duration in epochs (min: 518400 [6 months], max: 1036800 [720 days])
- `--storage-price`: Storage price in attoFIL per epoch per GiB
- `--verified`: Whether the deal is verified
- `--server`: Start a server after initiating the deal
- `--testnet`: Make deal on public testnet

When using `--testnet`, you'll need:
1. Filecoin tokens and datacap from the [faucet](https://faucet.calibnet.chainsafe-fil.io/)
2. - You would need to host the deal CAR file somewhere to serve them to testnet miner. `fildeal` currently supports [lighthouse](https://www.lighthouse.storage/) as the hosting service. You would need to have `LIGHTHOUSE_API_KEY` set in the environment variables. You can get the api key by following [this](https://docs.lighthouse.storage/lighthouse-1/how-to/create-an-api-key).


## Server Mode

When using the `--server` flag with `podsi-deal`, Fildeal starts a local server to serve files to miners in your network. This is particularly useful for local testing with 2k miners (2kb sector size miner)

## Environment Variables

- `FULLNODE_API_INFO`: Your Filecoin node API info
- `LIGHTHOUSE_API_KEY`: Required for testnet deals using Lighthouse storage
- `PORT`: Server port (default: 8000)
- `GENERATE_CAR_PATH`: Path for generated CAR files
- `AGGREGATE_CAR_PATH`: Path for aggregate CAR files
- `LIGHTHOUSE_DOWNLOAD_URL`: Lighthouse download URL
