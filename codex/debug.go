package codex

/*
   #include "bridge.h"

   static int cGoCodexDebug(void* codexCtx, void* resp) {
       return codex_debug(codexCtx, (CodexCallback) callback, resp);
   }
*/
import "C"
import "encoding/json"

type Node struct {
	NodeId  string  `json:"nodeId"`
	PeerId  string  `json:"peerId"`
	Record  string  `json:"record"`
	Address *string `json:"address"`
	Seen    bool    `json:"seen"`
}

type RoutingTable struct {
	LocalNode Node   `json:"localNode"`
	Nodes     []Node `json:"nodes"`
}

type DebugInfo struct {
	ID                string       `json:"id"`    // Peer ID
	Addrs             []string     `json:"addrs"` // Peer info addresses
	Spr               string       `json:"spr"`   // Signed Peer Record
	AnnounceAddresses []string     `json:"announceAddresses"`
	PeersTable        RoutingTable `json:"table"`
}

// Debug retrieves debugging information from the Codex node.
func (node CodexNode) Debug() (DebugInfo, error) {
	var info DebugInfo

	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoCodexDebug(node.ctx, bridge.resp) != C.RET_OK {
		return info, bridge.callError("cGoCodexDebug")
	}

	value, err := bridge.wait()
	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(value), &info)
	return info, err
}
