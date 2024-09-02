package main

import (
	"fmt"
	"io"
	"os"

	configurations "fildeal/src/config"
	"fildeal/src/deal"
	mkpiece "fildeal/src/mkpiece"
	"fildeal/src/server"
	"fildeal/src/utils"
)

func main() {

    // Detailed usage information
    usage := `Usage: fildeal <command> [arguments]
    Commands:
    cmp <parentFile> <childFile>       Compare two files and find the offset of the child file in the parent file.
    generate <files...>                Generate a data segment piece from the given files and output it to stdout.
    splitpiece <file> <outputDir>      Split the specified file into pieces and save them in the output directory.
    initiate <inputFolder> <miner>     Initiate a deal with the specified input folder and miner.

    Examples:
    fildeal cmp a.car b.car
    fildeal generate a.car b.car c.car > out.dat
    fildeal splitpiece input.car outputDir
    fildeal initiate inputFolder miner [--server]
    `

    // Check for --help flag
    if len(os.Args) < 2 || os.Args[1] == "--help" {
        fmt.Println(usage)
        return
    }

    command := os.Args[1]

    switch command {
    case "cmp":
        if len(os.Args) != 4 {
            fmt.Println("Usage: fildeal cmp <parentFile> <childFile>")
            return
        }
        fileAPath := os.Args[2]
        fileBPath := os.Args[3]
        offset, err := utils.FindOffset(fileAPath, fileBPath)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        fmt.Printf("Child file starts at offset %d in parent\n", offset)

    case "generate":
        if len(os.Args) < 3 {
            fmt.Println("Usage: fildeal generate <files...> > out.dat")
            return
        }
        readers := make([]io.ReadSeeker, 0)
        for _, arg := range os.Args[2:] {
            r, err := os.Open(arg)
            if err != nil {
                panic(err)
            }
            readers = append(readers, r)
        }
        out := mkpiece.MakeDataSegmentPiece(readers)
        io.Copy(os.Stdout, out)

    case "splitpiece":
        if len(os.Args) != 3 {
            fmt.Println("Usage: fildeal splitpiece <file> <outputDir>")
            return
        }
        filePath := os.Args[2]
        outputDir := os.Args[3]
        err := mkpiece.SplitPiece(filePath, outputDir)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

    case "initiate":
        if len(os.Args) < 4 {
            fmt.Println("Usage: fildeal initiate <inputFolder> <miner> [--server]")
            return
        }
        inputFolder := os.Args[2]
        miner := os.Args[3]
        err := deal.MakeDeal(inputFolder, miner)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        if len(os.Args) == 5 && os.Args[4] == "--server" {
            config := configurations.Configurations{Port: configurations.LoadConfigurations().Port} // Example configuration
            handler := server.SetupRouter()
            server.StartServer(config, handler)
        }

    default:
        fmt.Println("Unknown command:", command)
		fmt.Println(usage)
    }
}



