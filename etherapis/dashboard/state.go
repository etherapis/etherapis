// Contains the dashboard state assembler and differential pusher.

package dashboard

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/etherapis/etherapis/etherapis"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gopkg.in/inconshreveable/log15.v2"
)

// state contains all the system state that the dashboard requires in order to
// display any data that a user might need. It is maintained throughout the life
// cycle of an active dashboard connection.
type state struct {
	// Accounts contains all the known information about owned accounts
	Accounts map[string]etherapis.Account

	// Ethereum contains all the information pertaining to the Ethereum network
	Ethereum struct {
		Self    *p2p.NodeInfo        // Network connectivity information about the current node
		Peers   []*p2p.PeerInfo      // Network connectivity information about remote peers
		Head    *types.Header        // Current head of the local Ethereum blockchain
		Syncing *downloader.Progress // Sync progress as reported by the Ethereum downloader
	}
}

// Apply injects a state diff into the current state. If the diff specifies a
// non-existing path, it will hard panic!
func (s *state) Apply(diff stateDiff) {
	log15.Debug("Patching dashboard state tree", "path", strings.Join(diff.Path, "/"))

	// Drill down in the state struct to find the node to update
	current := reflect.ValueOf(s)
	for i, child := range diff.Path {
		switch current.Kind() {
		case reflect.Ptr:
			current = current.Elem().FieldByName(strings.ToUpper(child[:1]) + child[1:])
		case reflect.Map:
			if i == len(diff.Path)-1 {
				current.SetMapIndex(reflect.ValueOf(child), reflect.ValueOf(diff.Node))
				return
			}
			current = current.MapIndex(reflect.ValueOf(child))
		default:
			current = current.FieldByName(strings.ToUpper(child[:1]) + child[1:])
		}
	}
	// Update the child, or crash miserably
	current.Set(reflect.ValueOf(diff.Node))
}

// stateUpdate is a push message sent to the dashboard to update one or more nodes
// of the state tree.
type stateUpdate struct {
	Id    uint32      // Id of the state diff, mainly used for tracking and debugging
	Time  string      // Time instance of the state diff, mainly used for tracking and debugging
	Diffs []stateDiff // State updates to apply
}

// stateDiff is a modification to the current state. It contains the path within
// the state tree that needs to be updated and also the new node which overwrites
// the old one in the tree.
type stateDiff struct {
	Path []string    // Branch of the state tree to update
	Node interface{} // New state node to overwrite the old with
}

// stateServer is an HTTP multiplexer that serves up state data, able to handle
// simultaneously both RESTful requests against various branches of the state tree
// as well as WebSocket based push notifications against the root of the tree.
type stateServer struct {
	eapis *etherapis.EtherAPIs

	state *state                     // Current state of the UI that should be displayed
	conns map[uint32]*websocket.Conn // Active connections to push state updates to
	lock  sync.Mutex                 // Lock protecting the state and connection pools

	sessions uint32 // Number of sessions handled, used for contextual logging
	diffs    uint32 // Number of state diffs sent, used for for logging and client side checks
}

// newStateServer creates an etherapis state endpoint to serve RESTful requests,
// and WebSocket push notifications, returning the HTTP route multiplexer to
// embed in the main handler.
func newStateServer(base string, eapis *etherapis.EtherAPIs) *mux.Router {
	// Create a state server to expose various internals
	server := &stateServer{
		eapis: eapis,
		state: new(state),
		conns: make(map[uint32]*websocket.Conn),
	}
	server.start()

	// Register all the API handler endpoints
	router := mux.NewRouter()
	router.HandleFunc(base, server.State)
	return router
}

// start creates the initial dashboard state based on the current system data
// and  starts the update loop.
func (server *stateServer) start() {
	node := server.eapis.Geth().Stack()
	eth := server.eapis.Ethereum()

	// Assemble the initial state
	accounts, _ := server.eapis.Accounts()

	server.state.Accounts = make(map[string]etherapis.Account)
	for _, account := range accounts {
		server.state.Accounts[account.Hex()] = server.eapis.GetAccount(account)
	}
	server.state.Ethereum.Self = node.Server().NodeInfo()
	server.state.Ethereum.Head = eth.BlockChain().CurrentBlock().Header()
	server.state.Ethereum.Peers = node.Server().PeersInfo()

	origin, current, height, pulled, known := eth.Downloader().Progress()
	server.state.Ethereum.Syncing = &downloader.Progress{origin, current, height, pulled, known}

	// Register any event listeners
	go server.loop()
}

// loop is the heart of the state server, which processes the various Ethereum
// events from the underlying systems and converts them into state diffs that it
// sends out afterwards to all active state connections.
func (server *stateServer) loop() {
	// Register for various system events and start a poller too for non-event updates
	node := server.eapis.Geth().Stack()
	eth := server.eapis.Ethereum()

	gethEvents := node.EventMux().Subscribe(core.ChainHeadEvent{}, core.TxPreEvent{}, core.TxPostEvent{}, etherapis.NewAccountEvent{}, etherapis.DroppedAccountEvent{})
	gethPoller := time.NewTicker(time.Second)

	// Quick hack helper method to check for account updates
	updateAccounts := func(update *stateUpdate) {
		addresses, _ := server.eapis.Accounts()
		for _, address := range addresses {
			previous := server.state.Accounts[address.Hex()]
			current := server.eapis.GetAccount(address)

			if previous.CurrentBalance != current.CurrentBalance || previous.PendingBalance != current.PendingBalance {
				update.Diffs = append(update.Diffs, stateDiff{Path: []string{"accounts", address.Hex()}, Node: current})
			}
		}
	}

	for {
		// Prepare the state update diff list for population
		update := &stateUpdate{
			Id:   atomic.AddUint32(&server.diffs, 1),
			Time: fmt.Sprintf("%v", time.Now()),
		}
		logger := log15.New("diff", update.Id)

		// Wait for a state update and update all the paths
		select {
		case event, ok := <-gethEvents.Chan():
			// Don't do anything stupid if the subscription goes down
			if !ok {
				logger.Info("Event sub down, terminating state server...")
				return
			}
			// Figure out what to do based on the event type
			switch event := event.Data.(type) {
			case core.ChainHeadEvent:
				// New head arrived, update the current node and head block infos
				head := event.Block
				logger.Debug("New chain head detected", "number", head.Number(), "hash", fmt.Sprintf("%x", head.Hash().Bytes()))

				update.Diffs = append(update.Diffs, []stateDiff{
					{Path: []string{"ethereum", "self"}, Node: node.Server().NodeInfo()},
					{Path: []string{"ethereum", "head"}, Node: head.Header()},
				}...)
				updateAccounts(update)

			case core.TxPreEvent, core.TxPostEvent:
				// A transaction was initiated or completed, check if we need to update accounts
				logger.Debug("Transaction received/processed, checking accounts")
				updateAccounts(update)

			case etherapis.NewAccountEvent:
				// A new account was created, push it out to the dashboard
				logger.Debug("New account created/imported", "address", event.Address.Hex())

				update.Diffs = append(update.Diffs, []stateDiff{
					{Path: []string{"accounts", event.Address.Hex()}, Node: server.eapis.GetAccount(event.Address)},
				}...)

			case etherapis.DroppedAccountEvent:
				// An existing account was deleted, remove it from any dahsboards
				logger.Debug("Existing account deleted", "address", event.Address.Hex())

				update.Diffs = append(update.Diffs, []stateDiff{
					{Path: []string{"accounts", event.Address.Hex()}, Node: nil},
				}...)

			default:
				// Something strange arrived, issue a warning
				logger.Warn("Unexpected event in state server", "type", fmt.Sprintf("%T", event))
				continue
			}

		case <-gethPoller.C:
			// Periodically collect various infos that do not have events associated.
			//
			// Important: this branch doubles as a heartbeat mechanism so the dashboard
			// can detect if the backend went offline. Even if all data pushes convert
			// over to events, this branch should still push empty updates!
			origin, current, height, pulled, known := eth.Downloader().Progress()

			update.Diffs = append(update.Diffs, []stateDiff{
				{Path: []string{"ethereum", "peers"}, Node: node.Server().PeersInfo()},
				{Path: []string{"ethereum", "syncing"}, Node: &downloader.Progress{origin, current, height, pulled, known}},
			}...)
		}
		// Apply the diff locally to ensure everything's valid
		server.lock.Lock()
		for _, diff := range update.Diffs {
			server.state.Apply(diff)
		}
		update.Time = fmt.Sprintf("%v", time.Now())

		// Encode the diffs and lowercase the fields
		message, err := jsonMarshalLowercase(update)
		if err != nil {
			logger.Crit("Failed to serialize state update", "error", err)
		}
		for session, conn := range server.conns {
			logger.Debug("Sending new state diff to dashboard", "session", session)
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logger.Warn("Failed to send state update, dropping", "session", session, "error", err)
				delete(server.conns, session)
			}
		}
		server.lock.Unlock()
	}
}

// State opens up a push notification based state differential protocol. After
// a client connects to this endpoint, it will constantly receive state updates
// whenever a change occurs. The endpoint does not accept any inbound data apart
// from ping/pong messages to ensure the dashboard is still active.
func (server *stateServer) State(w http.ResponseWriter, r *http.Request) {
	// Define the WebSocket protocol details and open the persistent connection
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  65536,
		WriteBufferSize: 65536,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log15.Error("Failed to upgrade connection to web sockets", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Create a contextual logger for this connection
	session := atomic.AddUint32(&server.sessions, 1)
	logger := log15.New("session", session)

	logger.Debug("New dashboard connection")
	defer logger.Debug("Dashboard connection closed")

	// Lock the state for updates and stream out the current one
	server.lock.Lock()

	message, err := jsonMarshalLowercase(server.state)
	if err != nil {
		logger.Crit("Failed to serialize initial state", "error", err)
		server.lock.Unlock()
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		// Tear down the connection if initial send fails
		logger.Warn("Failed to send initial state", "error", err)
		server.lock.Unlock()
		return
	}
	// Otherwise register the connection for update and ensure it's cleaned up
	logger.Debug("Initial state sent, subscribing to diffs")
	server.conns[session] = conn
	defer func() {
		server.lock.Lock()
		defer server.lock.Unlock()

		delete(server.conns, session)
	}()
	// Updates can proceed to all sessions
	server.lock.Unlock()

	// Keep reading messages until the client closes the connection
	for {
		// This is a push only endpoint, dump whatever the client sends
		if _, _, err := conn.NextReader(); err != nil {
			return
		}
	}
}
