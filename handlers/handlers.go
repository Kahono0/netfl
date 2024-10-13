package handlers

import (
	"fmt"
	"net/http"

	"github.com/kahono0/netfl/pkg/msgs"
	"github.com/kahono0/netfl/pkg/p2p"
	"github.com/kahono0/netfl/pkg/putils"
	"github.com/kahono0/netfl/utils"
	"github.com/libp2p/go-libp2p/core/host"
)

func ShowPeers(w http.ResponseWriter, r *http.Request) {
	peers := p2p.Peers

	for _, peer := range peers {
		fmt.Fprintf(w, "%s\n", peer.ID)
	}
}

func SendSampleMsg(w http.ResponseWriter, r *http.Request) {
	peerID := r.URL.Query().Get("peer")
	if peerID == "" {
		fmt.Fprint(w, "No peer ID provided")
		return
	}

	peer := p2p.GetPeerByID(peerID)
	if peer == nil {
		fmt.Fprintf(w, "No peer found with ID %s\n", peerID)
		return
	}

	fmt.Printf("Peer: \n%s\n", utils.AsPrettyJson(peer))

	fmt.Fprintf(w, "Sending message to %s\n", peerID)
}

func SendSampleMsgHandler(host host.Host, protocolID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		peerID := r.URL.Query().Get("peer")
		peer := p2p.GetPeerByID(peerID)
		if peer == nil {
			fmt.Fprintf(w, "No peer found with ID %s\n", peerID)
			return
		}

		msg, err := msgs.NewMessage(msgs.Sample, []byte("Hello from server"))
		if err != nil {
			fmt.Fprintf(w, "Error creating message: %s\n", err)
			return
		}

		err = putils.SendMessage(host, *peer, msg, protocolID)
		if err != nil {
			fmt.Fprintf(w, "Error sending message: %s\n", err)
			return
		}

		fmt.Fprintf(w, "Sending message to %s\n", peerID)
	}
}
