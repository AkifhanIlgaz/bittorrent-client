package client

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/AkifhanIlgaz/bittorrent-client/bitfield"
	"github.com/AkifhanIlgaz/bittorrent-client/handshake"
	"github.com/AkifhanIlgaz/bittorrent-client/peers"
)

// A client is a TCP connection with a peer
type Client struct {
	Connection net.Conn
	Choked     bool
	BitField   bitfield.BitField
	peer       peers.Peer
	infoHash   [20]byte
	peerId     [20]byte
}

// New connects with a peer, completes a handshake, and receives a handshake
// returns an err if any of those fail.
func New(peer peers.Peer, peerId, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Minute)
	if err != nil {
		return nil, err
	}

	
}

func completeHandshake(conn net.Conn, peerId, infoHash [20]byte) (*handshake.Handshake, error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	req := handshake.New(infoHash, peerId)
	_, err := conn.Write(req.Serialize())
	if err != nil {
		return nil, err
	}

	res, err := handshake.Read(conn)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(res.InfoHash[:], infoHash[:]) {
		return nil, fmt.Errorf("Expected infohash %x but got %x", res.InfoHash, infoHash)
	}

	return res, err
}

func recvBitfield(conn net.Conn) (bitfield.BitField, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	// TODO: İmplement message package
	return bitfield.BitField{}, nil
}
