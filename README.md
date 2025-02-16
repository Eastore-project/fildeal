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

The CLI provides two main ways to make deals:

1. **Normal Deal (`deal`)**: Use this when you want to make regular filecoin deal without any special aggregation. It's suitable for:
   - Single file deals
   - You don't care about aggregation cryptographic proof
   - Direct deals without podsi-aggregation overhead

2. **Podsi Deal (`podsi-deal`)**: Use this when you need advanced file aggregation. It's recommended for:
   - Multiple files that need efficient aggregation
   - You need proof of Data Segment Inclusion of files in the aggregate

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

Parse and index podsi aggregate similar to Boost.

```sh
fildeal boost-index <file>
```

### Data Prep (`data-prep`)

Prepare data for a deal and show deal parameters without actually making the deal. This command is useful for:
- Generating CAR files
- Getting deal parameters (Piece CID, Payload CID, sizes)
- Testing file preparation before making actual deals
- Uploading to Lighthouse storage (if using lighthouse buffer)
- Using the parameters to make deals directly with boost

```sh
fildeal data-prep --input <inputPath> [options]
```

#### Options:
- `--input, -i`: Input path to prepare for deal (required)
- `--output, -o`: Output directory for the CAR file (default: "aggregate_car_file/")
- `--buffer`: Buffer to use (localhost or lighthouse, default: "localhost")
- `--lighthouse-api-key`: API key for Lighthouse storage (required when using lighthouse buffer)

### Make Normal Deal (`deal`)

Initiate a normal deal with a miner without using podsi-aggregation. This is useful when you want to make a simple deal without the aggregation feature.

```sh
fildeal deal --input <inputFolder> --miner <minerID> [options]
```

### Make Podsi Deal (`podsi-deal`)

Initiate a deal with a miner using podsi-aggregate for folder aggregation. This is recommended when you want to aggregate multiple files efficiently.

```sh
fildeal podsi-deal --input <inputFolder> --miner <minerID> [options]
```

#### Options (common for both deal commands):
- `--input, -i`: Input folder containing files to make deal with (required)
- `--miner, -m`: Miner ID to make the deal with (required)
- `--generate-car-path`: Path for generated CAR files (default: "generated_car/")
- `--aggregate-car-path`: Path for aggregate CAR files (default: "aggregate_car_file/")
- `--buffer`: Buffer to use (localhost or lighthouse, default: "localhost")
- `--duration`: Deal duration in epochs (min: 518400 [6 months], max: 1814400 [3.5 years])
- `--storage-price`: Storage price in attoFIL per epoch per GiB
- `--verified`: Whether the deal is verified (default: true for testnet, false otherwise)
- `--server`: Start a server after initiating the deal
- `--testnet`: Make deal on public testnet
- `--payload-cid`: Payload CID for the deal (default: "bafkreibtkdcncmofmavpdsar6msrmb2h4d7oetwtwtkz5cv3zsnwoyrrfq")
- `--lighthouse-download-url`: URL for downloading from Lighthouse (default: "https://gateway.lighthouse.storage/ipfs/")
- `--lighthouse-api-key`: API key for Lighthouse storage (required when using lighthouse buffer)

When using `--testnet`, you'll need:
1. Filecoin tokens and datacap from the [faucet](https://faucet.calibnet.chainsafe-fil.io/)
2. - You would need to host the deal CAR file somewhere to serve them to testnet miner. `fildeal` currently supports [lighthouse](https://www.lighthouse.storage/) as the hosting service. You would need to have `LIGHTHOUSE_API_KEY` set in the environment variables. You can get the api key by following [this](https://docs.lighthouse.storage/lighthouse-1/how-to/create-an-api-key).

## Usage

For making deals with miners, you have two main options:

1. **Regular deals** using the `deal` command:
   ```sh
   fildeal deal --input <inputFolder> --miner <minerID> [options]
   ```

2. **Podsi deals** using the `podsi-deal` command:
   ```sh
   fildeal podsi-deal --input <inputFolder> --miner <minerID> [options]
   ```

Both commands support the same set of options and can be used with either local or lighthouse storage. For testnet deals, you'll need:
1. Filecoin tokens and datacap from the [faucet](https://faucet.calibnet.chainsafe-fil.io/)
2. You would need to host the deal CAR file somewhere to serve them to testnet miner. `fildeal` currently supports [lighthouse](https://www.lighthouse.storage/) as the hosting service. You would need to have `LIGHTHOUSE_API_KEY` set in the environment variables. You can get the api key by following [this](https://docs.lighthouse.storage/lighthouse-1/how-to/create-an-api-key).

## Server Mode

When using the `--server` flag with `podsi-deal`, Fildeal starts a local server to serve files to miners in your network. This is particularly useful for local testing with 2k miners (2kb sector size miner)

## Environment Variables

- `FULLNODE_API_INFO`: Your Filecoin node API info
- `LIGHTHOUSE_API_KEY`: Required for testnet deals using Lighthouse storage
- `PORT`: Server port (default: 8000)
- `GENERATE_CAR_PATH`: Path for generated CAR files
- `AGGREGATE_CAR_PATH`: Path for aggregate CAR files
- `LIGHTHOUSE_DOWNLOAD_URL`: Lighthouse download URL
