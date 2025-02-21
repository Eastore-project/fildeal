# Fildeal CLI

Fildeal is a command-line interface (CLI) tool providing general Filecoin utilities. It includes commands for data preparation, making regular and podsi-aggregated deals, comparing files, generating data segment pieces, splitting pieces, and boost indexing.

## Table of Contents
- [Pre-requisites](#pre-requisites)
  - [Required for all commands](#required-for-all-commands)
  - [Required only for deal-making commands](#required-only-for-deal-making-commands-deal-and-podsi-deal-and-boost-index)
- [Installation](#installation)
- [Commands](#commands)
  - [Make Normal Deal](#make-normal-deal-deal)
  - [Make Podsi Deal](#make-podsi-deal-podsi-deal)
      - [Making deal in testnet](#making-deal-in-testnet)
      - [Making deal in localnet](#making-deal-in-localnet)
  - [Data Prep](#data-prep-data-prep)
  - [Generate Data Segment](#generate-data-segment-podsi-aggregate)
  - [Split Piece](#split-piece-splitpiece)
  - [Parse Boost Index](#parse-boost-index-boost-index)
  - [Compare Files](#compare-files-cmp)
- [References](#references)
- [Contributing](#contributing)

## Pre-requisites

### Required for all commands

- [Go](https://golang.org/doc/install) 1.20 or higher

### Required only for deal-making commands (`deal` and `podsi-deal`) and `boost-index`

- [Boost](https://boost.filecoin.io/getting-started)
  1. Install boost from [guidelines](https://boost.filecoin.io/getting-started#building-and-installing). For testnet (calibnet), make sure to build the calibnet version of boost
  2. Run the following commands inside the boost repository:
     ```sh
     sudo make install
     
     # Set FULLNODE_API_INFO globally (choose based on your OS):
     # For calibration testnet, use https://api.calibration.node.glif.io
     
     # For Linux:
     # Add to one of these files based on your shell setup:
     # - ~/.bashrc (for Bash interactive shells)
     # - ~/.bash_profile (for Bash login shells)
     # - ~/.profile (for any shell's login)
     echo 'export FULLNODE_API_INFO="your fullnode api info"' >> ~/.bashrc  # or your preferred config file
     source ~/.bashrc  # or the file you modified
     
     # For macOS:
     # Add to one of these files based on your shell setup:
     # - ~/.zshrc (for Zsh, default shell on modern macOS)
     # - ~/.bash_profile (for Bash if you use it)
     echo 'export FULLNODE_API_INFO="your fullnode api info"' >> ~/.zshrc  # or your preferred config file
     source ~/.zshrc  # or the file you modified
     
     boost init
     boost wallet list
     ```
     
  3. Fund your boost wallet:
     - For testnet: Get FIL and DataCap from the [faucet](https://faucet.calibnet.chainsafe-fil.io/)
     - For mainnet: Fund your primary boost-wallet with FIL
     - DataCap allocation required if making verified deals

Note: All other commands (`cmp`, `podsi-aggregate`, `splitpiece`, `data-prep`) work without boost setup.

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/eastore-project/fildeal.git
   cd fildeal
   ```

2. **Build and install the CLI:**

   ```bash
   go build -o fildeal ./cmd/fildeal
   sudo mv fildeal /usr/local/bin/
   ```

3. **Setup environment variables:**

   ```bash
   # Required if using lighthouse as buffer for deals (recommended)
   export LIGHTHOUSE_API_KEY="your lighthouse api key"

   # Optional: Set custom paths for generated files
   export GENERATE_CAR_PATH="path for generated CAR files" # default: "generated_car/"
   export AGGREGATE_CAR_PATH="path for aggregate CAR files" # default: "aggregate_car_file/"
   export PROOFS_DIR="path for inclusion proofs" # default: "proofs-dir/"
   export PORT="8000"  # Server port (default: 8000)
   ```
   OR 

    Using .env file
   ```sh
   cp .env.example .env    # Copy the example config
      # Edit .env with your values
   source .env           # Load the variables
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

### Make Normal Deal (`deal`)

Initiate a normal deal with a miner without using podsi-aggregation. This is useful when you want to make a simple deal without the aggregation feature.

```bash
fildeal deal --input <inputFolder> --miner <minerID> [options]
```

### Make Podsi Deal (`podsi-deal`)

Initiate a deal with a miner using podsi-aggregate for folder aggregation. This is recommended when you want to aggregate multiple files efficiently.

```bash
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
- `--payload-cid`: Payload CID for the deal (default: "bafkreibtkdcncmofmavpdsar6msrmb2h4d7oetwtwtkz5cv3zsnwoyrrfq")
- `--lighthouse-download-url`: URL for downloading from Lighthouse (default: "https://gateway.lighthouse.storage/ipfs/")
- `--lighthouse-api-key`: API key for Lighthouse storage (required when using lighthouse buffer)

#### Making deal in testnet

For calibration testnet, miner `t017840` supports verified deals LARGER THAN 1MB for free. You can make a verified deal with this miner using:

```bash
fildeal deal --input <inputFolder> --miner t017840 --verified --buffer lighthouse
```
OR
```bash
fildeal podsi-deal --input <inputFolder> --miner t017840 --verified --buffer lighthouse
```

When making deals in testnet, you'll need:

1. Filecoin tokens and datacap from the [faucet](https://faucet.calibnet.chainsafe-fil.io/)
2. You would need to host the deal CAR file somewhere to serve them to testnet miner. `fildeal` currently supports [lighthouse](https://www.lighthouse.storage/) as the hosting service. You would need to have `LIGHTHOUSE_API_KEY` set in the environment variables. You can get the api key by following [this](https://docs.lighthouse.storage/lighthouse-1/how-to/create-an-api-key).

After making a deal, you can check its status using the boost CLI:
```bash
boost deal-status --provider t017840 --deal-uuid <uuid> # deal UUID will be logged in the deal output
```

#### Making deal in localnet

For easy testing, you can run a 2k lotus-miner setup using these [scripts](https://gist.github.com/lordshashank/fb2fbd53b5520a862bd451e3603b4718), and then use this tool to initiate deals with the miner.

For testing with local miners, use the `--server` flag. This starts a local server to serve files to miners from your device:

```bash
fildeal deal --input <inputFolder> --miner <minerID> --server
```

### Data Prep (`data-prep`)

Prepare data for a deal and show deal parameters without actually making the deal. This command is useful for:

- Generating CAR files
- Getting deal parameters (Piece CID, Payload CID, sizes)
- Testing file preparation before making actual deals
- Uploading to Lighthouse storage (if using lighthouse buffer)
- Using the parameters to make deals directly with boost

```bash
fildeal data-prep --input <inputPath> [options]
```


### Generate Data Segment (`podsi-aggregate`)

Generate a data segment piece from all files in an input folder and generate inclusion proofs for each file.

```bash
fildeal podsi-aggregate --input <inputFolder> --output <outputFile> [--proof-dir <proofDir>]
# or using short flags
fildeal podsi-aggregate -i <inputFolder> -o <outputFile>
```

Options:
- `--input, -i`: Input folder containing files to aggregate (required)
- `--output, -o`: Output file path where the aggregated piece will be written (required)
- `--proof-dir`: Directory where inclusion proofs will be stored (default: "proofs-dir/", can be set via PROOFS_DIR env var)

The command will generate:
1. An aggregated piece file at the specified output path
2. A proof file for each input file in the proof directory (if specified)
   - Each proof file will be named `<original_filename>.proof.json`
   - Proofs contain both subtree and index proofs with hex-encoded paths (0x-prefixed)
   - These proofs can be used to verify the inclusion of individual files in the aggregate

### Split Piece (`splitpiece`)

Split a podsi-aggregate output file into pieces and save them in the output directory.

```bash
fildeal splitpiece --input <inputFile> --output <outputDir>
# or using short flags
fildeal splitpiece -i <inputFile> -o <outputDir>
```

### Parse Boost Index (`boost-index`)

Parse and index podsi aggregate similar to Boost.

```bash
fildeal boost-index <file>
```

### Compare Files (`cmp`)

Compare two files and find the offset of the child file in the parent file.

```bash
fildeal cmp --parent <parentFile> --child <childFile>
# or using short flags
fildeal cmp -p <parentFile> -c <childFile>
```

## References

- [mkpiece](https://github.com/willscott/mkpiece)
- [data-segment-library](https://github.com/filecoin-project/go-data-segment)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please make sure to follow the existing coding style.
