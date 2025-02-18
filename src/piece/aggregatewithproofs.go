package piece

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	pieceutils "github.com/eastore-project/fildeal/src/piece/utils"
	"github.com/eastore-project/fildeal/src/utils"
)

type ProofData struct {
	Path  string `json:"path"`  // hex encoded path
	Index uint64 `json:"index"` // numeric index
}

type InclusionProofData struct {
	SubtreeProof ProofData `json:"subtree_proof"` // proof of inclusion in data tree
	IndexProof   ProofData `json:"index_proof"`   // proof of inclusion in index tree
}

// AggregateWithProofs aggregates files from inputDir into a single piece file at outputFile
// and stores inclusion proofs in proofDir
func AggregateWithProofs(inputDir, outputFile, proofDir string) error {
	// Ensure input and proof directories exist
	if err := utils.EnsureDirectoriesExist(inputDir, proofDir, filepath.Dir(outputFile)); err != nil {
		return fmt.Errorf("failed to ensure directories exist: %w", err)
	}

	readers, err := utils.GetReaders(inputDir)
	if err != nil {
		return fmt.Errorf("failed to get readers from input folder: %w", err)
	}
	defer func() {
		for _, r := range readers {
			if closer, ok := r.(io.Closer); ok {
				closer.Close()
			}
		}
	}()

	result, err := pieceutils.MakeDataSegmentPieceWithProof(readers)
	if err != nil {
		return fmt.Errorf("failed to make data segment piece with proof: %w", err)
	}

	// Create output file and write aggregated piece
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, result.Reader); err != nil {
		return fmt.Errorf("failed to write to output file: %w", err)
	}

	// Get list of files in input directory to map proofs to filenames
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("failed to read input directory: %w", err)
	}

	// Write proofs for each file
	for i, proof := range result.InclProofs {
		if i >= len(files) {
			break
		}

		proofFileName := filepath.Join(proofDir, fmt.Sprintf("%s.proof.json", files[i].Name()))
		proofFile, err := os.Create(proofFileName)
		if err != nil {
			return fmt.Errorf("failed to create proof file: %w", err)
		}

		// Convert paths to hex with 0x prefix and create proof data
		proofData := InclusionProofData{
			SubtreeProof: ProofData{
				Path:  "0x" + hex.EncodeToString(proof.ProofSubtree.Path[0][:]),
				Index: proof.ProofSubtree.Index,
			},
			IndexProof: ProofData{
				Path:  "0x" + hex.EncodeToString(proof.ProofIndex.Path[0][:]),
				Index: proof.ProofIndex.Index,
			},
		}

		encoder := json.NewEncoder(proofFile)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(map[string]interface{}{
			"filename": files[i].Name(),
			"proof":    proofData,
		}); err != nil {
			proofFile.Close()
			return fmt.Errorf("failed to write proof data: %w", err)
		}
		proofFile.Close()
	}

	return nil
}
