package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/bsv-blockchain/teranode-p2p-poc/pkg/parser"
	"github.com/sirupsen/logrus"

	"github.com/bsv-blockchain/go-p2p"
	"github.com/spf13/viper"

	teranodep2p "teranodep2p/pkg/p2p"
)

func main() {
	ctx := context.Background()
	log := logrus.New()

	// Initialize Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("teranode_p2p")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}

	// Load P2P settings from config
	bootstrapAddresses := viper.GetStringSlice("p2p.bootstrap_addresses")
	statisPeers := viper.GetStringSlice("p2p.static_peers")
	sharedKey := viper.GetString("p2p.shared_key")
	dhtProtocolID := viper.GetString("p2p.dht_protocol_id")
	port := viper.GetInt("p2p.port")
	listenAddresses := viper.GetStringSlice("p2p.listen_addresses")
	advertise := viper.GetBool("p2p.advertise")
	usePrivateDHT := viper.GetBool("p2p.use_private_dht")
	// Get networks from config and generate topics
	var topics []string
	networks := viper.GetStringSlice("networks")
	if len(networks) == 0 {
		// Fallback to old topics config if networks not specified
		topics = viper.GetStringSlice("topics")
		if len(topics) == 0 {
			log.Fatalf("neither 'networks' nor 'topics' configured")
		}
		log.Warn("Using deprecated 'topics' config. Please migrate to 'networks' config.")
	} else {
		topics = parser.GenerateTopics(networks)
		log.Infof("Generated %d topics from %d networks", len(topics), len(networks))
	}

	config := p2p.Config{
		ProcessName:        "teranode-p2p-poc",
		Port:               port,
		ListenAddresses:    listenAddresses,
		Advertise:          advertise,
		UsePrivateDHT:      usePrivateDHT,
		SharedKey:          sharedKey,
		BootstrapAddresses: bootstrapAddresses,
		StaticPeers:        statisPeers,
		DHTProtocolID:      dhtProtocolID,
	}

	node, err := teranodep2p.New(ctx, log, teranodep2p.Config{
		P2PConfig: config,
		Topics:    topics,
	})
	if err != nil {
		panic(err)
	}

	// subscribe and log subtree messages
	subtrees := node.SubscribeSubtree()
	for msg := range subtrees {
		var evt p2p.SubtreeMessage
		var ok bool
		if evt, ok = msg.(p2p.SubtreeMessage); ok {
			log.Infof("Received subtree from %s: %+v", evt.PeerID, evt)
		}

		subtree, err := teranodep2p.GetSubtree(evt)
		if err != nil {
			log.Errorf("Failed to get subtree %s: %v", evt.Hash, err)
			continue
		}

		// print subtree
		log.Infof("Subtree %s data: %v", evt.Hash, subtree)
	}
}
