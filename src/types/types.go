package types

import "github.com/tech-greedy/go-generate-car/util"

type Result struct {
	Ipld      *util.FsNode
	DataCid   string
	PieceCid  string
	PieceSize uint64
	CidMap    map[string]util.CidMapValue
	CarSize  uint64
}

type CarParams struct {
	Input       string
	PieceSize   uint64
	OutDir      string
	Parent      string
	TmpDir      string
	Single      bool
}

type DealFlags struct {
	Testnet bool
	Server bool
}