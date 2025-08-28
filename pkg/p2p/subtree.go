package p2p

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bsv-blockchain/go-bt/v2/chainhash"
	"github.com/bsv-blockchain/go-p2p"
	"github.com/bsv-blockchain/go-subtree"
)

func GetSubtree(subtreeMessage p2p.SubtreeMessage) (*subtree.Subtree, error) {
	url := fmt.Sprintf("%s/subtree/%s", subtreeMessage.DataHubURL, subtreeMessage.Hash)
	subtreeNodeBytes, err := fetchBytes(url)
	if err != nil {
		return nil, err
	}

	// in the subtree validation, we only use the hashes of the FileTypeSubtreeToCheck, which is what is returned from the peer
	numberOfNodes := len(subtreeNodeBytes) / chainhash.HashSize
	st, err := subtree.NewIncompleteTreeByLeafCount(numberOfNodes)
	if err != nil {
		return nil, fmt.Errorf("failed to create subtree with %d nodes: %w", numberOfNodes, err)
	}

	// Sanity check, subtrees should never be empty
	if numberOfNodes == 0 {
		return nil, errors.New("subtree has zero nodes")
	}

	// Deserialize the subtree nodes from the bytes
	for i := 0; i < numberOfNodes; i++ {
		// Each node is a chainhash.Hash, so we read chainhash.HashSize bytes
		nodeBytes := subtreeNodeBytes[i*chainhash.HashSize : (i+1)*chainhash.HashSize]
		nodeHash, err := chainhash.NewHash(nodeBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to create hash from bytes at index: %d %w", i, err)
		}

		if i == 0 && nodeHash.Equal(subtree.CoinbasePlaceholderHashValue) {
			if err = st.AddCoinbaseNode(); err != nil {
				return nil, fmt.Errorf("failed to add coinbase node to subtree at index: %d %w", i, err)
			}

			continue
		}

		// Add the node to the subtree, we do not know the fee or size yet, so we use 0
		if err = st.AddNode(*nodeHash, 0, 0); err != nil {
			return nil, fmt.Errorf("failed to add node to subtree at index %d %w", i, err)
		}
	}

	return st, nil
}

func fetchBytes(url string) ([]byte, error) {
	fmt.Printf("Fetching subtree from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}
	return data, nil
}
