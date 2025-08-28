package p2p

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bsv-blockchain/go-p2p"
	"github.com/bsv-blockchain/go-subtree"
)

func GetSubtree(subtreeMessage p2p.SubtreeMessage) (*subtree.Subtree, error) {
	url := fmt.Sprintf("%s/subtree/%s", subtreeMessage.DataHubURL, subtreeMessage.Hash)
	subtreeBytes, err := fetchBytes(url)
	if err != nil {
		return nil, err
	}

	return subtree.NewSubtreeFromBytes(subtreeBytes)
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
