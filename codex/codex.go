package codex

/*
   #cgo LDFLAGS: -L../vendor/nim-codex/build/ -lcodex
   #cgo LDFLAGS: -L../vendor/nim-codex/ -Wl,-rpath,../vendor/nim-codex/

   #include "bridge.h"
   #include <stdlib.h>

   void libcodexNimMain(void);

   static void codex_host_init_once(void){
       static int done;
       if (!__atomic_exchange_n(&done, 1, __ATOMIC_SEQ_CST)) libcodexNimMain();
   }

   // resp must be set != NULL in case interest on retrieving data from the callback
   void callback(int ret, char* msg, size_t len, void* resp);

   static void* cGoCodexNew(const char* configJson, void* resp) {
       void* ret = codex_new(configJson, (CodexCallback) callback, resp);
       return ret;
   }

   static int cGoCodexStart(void* codexCtx, void* resp) {
       return codex_start(codexCtx, (CodexCallback) callback, resp);
   }

   static int cGoCodexStop(void* codexCtx, void* resp) {
       return codex_stop(codexCtx, (CodexCallback) callback, resp);
   }

   static int cGoCodexDestroy(void* codexCtx, void* resp) {
       return codex_destroy(codexCtx, (CodexCallback) callback, resp);
   }

    static int cGoCodexVersion(void* codexCtx, void* resp) {
       return codex_version(codexCtx, (CodexCallback) callback, resp);
   }

   static int cGoCodexRevision(void* codexCtx, void* resp) {
       return codex_revision(codexCtx, (CodexCallback) callback, resp);
   }

   static int cGoCodexRepo(void* codexCtx, void* resp) {
       return codex_repo(codexCtx, (CodexCallback) callback, resp);
   }

   static int cGoCodexSpr(void* codexCtx, void* resp) {
       return codex_spr(codexCtx, (CodexCallback) callback, resp);
   }
*/
import "C"
import (
	"encoding/json"
	"unsafe"
)

type LogFormat string

const (
	LogFormatAuto     LogFormat = "auto"
	LogFormatColors   LogFormat = "colors"
	LogFormatNoColors LogFormat = "nocolors"
	LogFormatJSON     LogFormat = "json"
)

type RepoKind string

const (
	FS      RepoKind = "fs"
	SQLite  RepoKind = "sqlite"
	LevelDb RepoKind = "leveldb"
)

type CodexConfig struct {
	LogFormat                      LogFormat `json:"log-format,omitempty"`
	MetricsEnabled                 bool      `json:"metrics,omitempty"`
	MetricsAddress                 string    `json:"metrics-address,omitempty"`
	DataDir                        string    `json:"data-dir,omitempty"`
	ListenAddrs                    []string  `json:"listen-addrs,omitempty"`
	Nat                            string    `json:"nat,omitempty"`
	DiscoveryPort                  int       `json:"disc-port,omitempty"`
	NetPrivKeyFile                 string    `json:"net-privkey,omitempty"`
	BootstrapNodes                 []byte    `json:"bootstrap-node,omitempty"`
	MaxPeers                       int       `json:"max-peers,omitempty"`
	NumThreads                     int       `json:"num-threads,omitempty"`
	AgentString                    string    `json:"agent-string,omitempty"`
	RepoKind                       RepoKind  `json:"repo-kind,omitempty"`
	StorageQuota                   int       `json:"storage-quota,omitempty"`
	BlockTtl                       int       `json:"block-ttl,omitempty"`
	BlockMaintenanceInterval       int       `json:"block-mi,omitempty"`
	BlockMaintenanceNumberOfBlocks int       `json:"block-mn,omitempty"`
	CacheSize                      int       `json:"cache-size,omitempty"`
	LogFile                        string    `json:"log-file,omitempty"`
}

type CodexNode struct {
	ctx unsafe.Pointer
}

// CodexNew creates a new Codex node with the provided configuration.
// The node is not started automatically; you need to call CodexStart
// to start it.
// It returns a Codex node that can be used to interact
// with the Codex network.
func CodexNew(config CodexConfig) (*CodexNode, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	jsonConfig, err := json.Marshal(config)

	if err != nil {
		return nil, err
	}

	cJsonConfig := C.CString(string(jsonConfig))
	defer C.free(unsafe.Pointer(cJsonConfig))

	ctx := C.cGoCodexNew(cJsonConfig, bridge.resp)

	if _, err := bridge.wait(); err != nil {
		return nil, bridge.err
	}

	return &CodexNode{ctx: ctx}, bridge.err
}

// Start starts the Codex node.
// TODO waits for the node to be fully started,
// by looking into the logs.
func (node CodexNode) Start() error {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoCodexStart(node.ctx, bridge.resp) != C.RET_OK {
		return bridge.callError("cGoCodexStart")
	}

	_, err := bridge.wait()

	return err
}

// StartAsync is the asynchronous version of Start.
func (node CodexNode) StartAsync(onDone func(error)) {
	go func() {
		err := node.Start()
		onDone(err)
	}()
}

// Stop stops the Codex node.
func (node CodexNode) Stop() error {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoCodexStop(node.ctx, bridge.resp) != C.RET_OK {
		return bridge.callError("cGoCodexStop")
	}

	_, err := bridge.wait()
	return err
}

// Destroy destroys the Codex node, freeing all resources.
// The node must be stopped before calling this method.
func (node CodexNode) Destroy() error {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoCodexDestroy(node.ctx, bridge.resp) != C.RET_OK {
		return bridge.callError("cGoCodexDestroy")
	}

	_, err := bridge.wait()
	return err
}

// Version returns the version of the Codex node.
func (node CodexNode) Version() (string, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoCodexVersion(node.ctx, bridge.resp) != C.RET_OK {
		return "", bridge.callError("cGoCodexVersion")
	}

	return bridge.wait()
}

// Revision returns the revision of the Codex node.
func (node CodexNode) Revision() (string, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoCodexRevision(node.ctx, bridge.resp) != C.RET_OK {
		return "", bridge.callError("cGoCodexRevision")
	}

	return bridge.wait()
}

// Repo returns the path of the data dir folder.
func (node CodexNode) Repo() (string, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoCodexRepo(node.ctx, bridge.resp) != C.RET_OK {
		return "", bridge.callError("cGoCodexRepo")
	}

	return bridge.wait()
}
