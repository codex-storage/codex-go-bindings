package codex

/*
   #include "bridge.h"
   #include <stdlib.h>

   static int cGoCodexConnect(void* codexCtx, char* peerId, const char** peerAddresses, uintptr_t peerAddressesSize,  void* resp) {
       return codex_connect(codexCtx, peerId, peerAddresses, peerAddressesSize, (CodexCallback) callback, resp);
   }
*/
import "C"
import (
	"log"
	"unsafe"
)

func (node CodexNode) Connect(peerId string, peerAddresses []string) error {
	bridge := newBridgeCtx()
	defer bridge.free()

	var cPeerId = C.CString(peerId)
	defer C.free(unsafe.Pointer(cPeerId))

	if len(peerAddresses) > 0 {
		var cAddresses = make([]*C.char, len(peerAddresses))
		for i, addr := range peerAddresses {
			cAddresses[i] = C.CString(addr)
			defer C.free(unsafe.Pointer(cAddresses[i]))
		}

		log.Println("peerAddresses", cAddresses)

		if C.cGoCodexConnect(node.ctx, cPeerId, &cAddresses[0], C.uintptr_t(len(peerAddresses)), bridge.resp) != C.RET_OK {
			return bridge.callError("cGoCodexConnect")
		}
	} else {
		if C.cGoCodexConnect(node.ctx, cPeerId, nil, 0, bridge.resp) != C.RET_OK {
			return bridge.callError("cGoCodexConnect")
		}
	}

	_, err := bridge.wait()
	return err
}
