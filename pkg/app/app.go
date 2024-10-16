package app

import (
	"fmt"

	"github.com/kahono0/netfl/pkg/p2p"
	"github.com/kahono0/netfl/pkg/peers"
	"github.com/kahono0/netfl/repo"
	"github.com/kahono0/netfl/utils"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Config struct {
	p2p.P2PConfig
	// dir to get movies from
	Path string

	// port the server will be running on
	SPort int

	// name referring to this peer
	Alias string

	// avatar url
	Avatar string

	// host addr
	HostAddr string
}

type App struct {
	Config    Config
	movieRepo *repo.MovieRepo
	peerStore *peers.PeerStore
	Host      host.Host
	PeerChan  chan peer.AddrInfo
}

func New(cfg Config, avatar, alias string, port int, newStreamHandler func(stream network.Stream)) (*App, error) {
	hostAddr := fmt.Sprintf("http://%s:%d", utils.GetPrivateIP(), port)

	cfg.HostAddr = hostAddr

	cfg.Avatar = hostAddr + avatar
	cfg.StreamHandler = newStreamHandler

	return NewApp(cfg)
}

func NewApp(config Config) (*App, error) {
	peers := peers.NewStore()

	host, peerChan := p2p.Init(config.P2PConfig, peers)
	movies := repo.NewMovieRepo(config.Path, config.HostAddr, false)

	return &App{
		Config:    config,
		movieRepo: movies,
		peerStore: peers,
		Host:      host,
		PeerChan:  peerChan,
	}, nil
}

func (a *App) GetMovieRepo() *repo.MovieRepo {
	return a.movieRepo
}

func (a *App) GetPeerStore() *peers.PeerStore {
	return a.peerStore
}
