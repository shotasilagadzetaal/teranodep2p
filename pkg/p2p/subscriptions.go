package p2p

import (
	"github.com/bsv-blockchain/teranode-p2p-poc/pkg/parser"
)

// --- BestBlock ---
func (n *Node) SubscribeBestBlock() <-chan interface{} {
	n.subsLock.Lock()
	defer n.subsLock.Unlock()
	n.subs[parser.TypeBestBlock] = append(n.subs[parser.TypeBestBlock], make(chan interface{}))
	return n.subs[parser.TypeBestBlock][len(n.subs[parser.TypeBestBlock])-1]
}

// --- Block ---
func (n *Node) SubscribeBlock() <-chan interface{} {
	n.subsLock.Lock()
	defer n.subsLock.Unlock()
	n.subs[parser.TypeBlock] = append(n.subs[parser.TypeBlock], make(chan interface{}))
	return n.subs[parser.TypeBlock][len(n.subs[parser.TypeBlock])-1]
}

// --- MiningOn ---
func (n *Node) SubscribeMiningOn() <-chan interface{} {
	n.subsLock.Lock()
	defer n.subsLock.Unlock()
	n.subs[parser.TypeMiningOn] = append(n.subs[parser.TypeMiningOn], make(chan interface{}))
	return n.subs[parser.TypeMiningOn][len(n.subs[parser.TypeMiningOn])-1]
}

// --- Subtree ---
func (n *Node) SubscribeSubtree() <-chan interface{} {
	n.subsLock.Lock()
	defer n.subsLock.Unlock()
	n.subs[parser.TypeSubtree] = append(n.subs[parser.TypeSubtree], make(chan interface{}))
	return n.subs[parser.TypeSubtree][len(n.subs[parser.TypeSubtree])-1]
}

// --- Handshake ---
func (n *Node) SubscribeHandshake() <-chan interface{} {
	n.subsLock.Lock()
	defer n.subsLock.Unlock()
	n.subs[parser.TypeHandshake] = append(n.subs[parser.TypeHandshake], make(chan interface{}))
	return n.subs[parser.TypeHandshake][len(n.subs[parser.TypeHandshake])-1]
}

// --- RejectedTx ---
func (n *Node) SubscribeRejectedTx() <-chan interface{} {
	n.subsLock.Lock()
	defer n.subsLock.Unlock()
	n.subs[parser.TypeRejectedTx] = append(n.subs[parser.TypeRejectedTx], make(chan interface{}))
	return n.subs[parser.TypeRejectedTx][len(n.subs[parser.TypeRejectedTx])-1]
}
