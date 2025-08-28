package p2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/bsv-blockchain/go-p2p"
	"github.com/bsv-blockchain/teranode-p2p-poc/pkg/parser"
	"github.com/sirupsen/logrus"
)

type Node struct {
	node *p2p.Node
	log  *logrus.Logger

	subsLock sync.Mutex
	subs     map[parser.MessageType][]chan interface{}
}

// Config is a wrapper around p2p.Config with topics
type Config struct {
	P2PConfig p2p.Config
	Topics    []string
}

func New(ctx context.Context, log *logrus.Logger, cfg Config) (*Node, error) {
	n, err := p2p.NewNode(ctx, log, cfg.P2PConfig)
	if err != nil {
		return nil, err
	}
	peer := &Node{
		node: n,
		log:  log,
		subs: make(map[parser.MessageType][]chan interface{}),
	}

	if err := n.Start(ctx, nil, cfg.Topics...); err != nil {
		return nil, err
	}

	// Register generic handler for all topics
	for _, topic := range cfg.Topics {
		t := topic
		if err := n.SetTopicHandler(ctx, t, func(ctx context.Context, data []byte, peerID string) {
			parsedMsg, err := parser.ParseMessage(t, data)
			if err != nil {
				log.Warnf("Failed to parse message on %s: %v", t, err)
				return
			}

			peer.subsLock.Lock()
			for _, c := range peer.subs[parsedMsg.Type] {
				select {
				case c <- parsedMsg.Data:
				default:
				}
			}
			peer.subsLock.Unlock()
		}); err != nil {
			return nil, fmt.Errorf("failed to set handler for %s: %w", t, err)
		}
	}

	return peer, nil
}
