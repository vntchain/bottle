package core

import (
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/rpc"
)

type server struct {
	ipcEndpoint string       // IPC endpoint to listen at (empty = IPC disabled)
	ipcListener net.Listener // IPC RPC listener socket to serve API requests
	ipcHandler  *rpc.Server  // IPC RPC request handler to process the API requests
	log         log.Logger
	lock        sync.RWMutex
	stop        chan struct{} // Channel to wait for termination notifications
}

func newServer(ipcPath string) *server {
	endpoint := filepath.Join(os.TempDir(), ipcPath)
	return &server{
		ipcEndpoint: endpoint,
		log:         log.New(),
	}
}

// StartIPC initializes and starts the IPC RPC endpoint.
func (s *server) startIPC(apis []rpc.API) error {
	if s.ipcEndpoint == "" {
		return nil // IPC disabled.
	}
	listener, handler, err := rpc.StartIPCEndpoint(s.ipcEndpoint, apis)
	if err != nil {
		return err
	}
	s.ipcListener = listener
	s.ipcHandler = handler
	s.log.Info("IPC endpoint opened", "url", s.ipcEndpoint)
	return nil
}

// stopIPC terminates the IPC RPC endpoint.
func (s *server) stopIPC() {
	if s.ipcListener != nil {
		s.ipcListener.Close()
		s.ipcListener = nil

		s.log.Info("IPC endpoint closed", "endpoint", s.ipcEndpoint)
	}
	if s.ipcHandler != nil {
		s.ipcHandler.Stop()
		s.ipcHandler = nil
	}
}

func (s *server) wait() {
	s.lock.RLock()
	stop := s.stop
	s.lock.RUnlock()
	<-stop
}

// apis returns the collection of RPC descriptors this node offers.
func (s *server) apis() []rpc.API {
	return []rpc.API{
		{
			Namespace: "bottle",
			Version:   "1.0",
			Service:   newBottleAPI(s),
		},
	}
}

type bottleAPI struct {
	server *server
}

func newBottleAPI(s *server) *bottleAPI {
	return &bottleAPI{server: s}
}
