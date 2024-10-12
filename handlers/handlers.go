package handlers

import (
	"fmt"
	"net/http"

	"github.com/kahono0/netfl/p2p"
)

func ShowPeers(w http.ResponseWriter, r *http.Request) {
	peers := p2p.Peers

	for _, peer := range peers {
		fmt.Fprintf(w, "%s\n", peer.ID)
	}
}

func SendSampleMsg(w http.ResponseWriter, r *http.Request) {
	peerID := r.URL.Query().Get("peer")

	fmt.Printf("Sending message to %s\n", peerID)

	fmt.Fprintf(w, "Sending message to %s\n", peerID)
}
