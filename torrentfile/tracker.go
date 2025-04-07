package torrentfile

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/url"
	"strconv"
)

func (t *TorrentFile) buildTrackerUrl(peerId [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}

	base.RawQuery = params.Encode()

	return base.String(), nil
}

type Peer struct {
	IP   net.IP
	Port uint16
}

func Unmarshal(peersBin []byte) ([]Peer, error) {
	const peerSize = 6
	numPeers := len(peersBin) / peerSize
	if len(peersBin)%peerSize != 0 {
		err := fmt.Errorf("Received malformed peers")
		return nil, err
	}

	peers := make([]Peer, numPeers)
	for i := range numPeers {
		offset := i * peerSize
		peers[i].IP = net.IP(peersBin[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset+4 : offset+6])
	}

	return peers, nil
}
