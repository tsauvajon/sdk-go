package websocket

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/gorilla/websocket"
	"github.com/kuzzleio/sdk-go/collection"
	"github.com/kuzzleio/sdk-go/connection"
	"github.com/kuzzleio/sdk-go/event"
	"github.com/kuzzleio/sdk-go/state"
	"github.com/kuzzleio/sdk-go/types"
	"net/url"
	"sync"
	"time"
)

const (
	MAX_EMIT_TIMEOUT = 10
)

type webSocket struct {
	ws      *websocket.Conn
	mu      *sync.Mutex
	queuing bool
	state   int

	listenChan     chan []byte
	channelsResult map[string]chan<- types.KuzzleResponse // todo use sync.Map when go-1.9 will be released
	subscriptions  types.RoomList                         // todo use sync.Map when go-1.9 will be released
	lastUrl        string
	host           string
	wasConnected   bool
	eventListeners map[int]chan<- interface{}

	autoQueue             bool
	autoReconnect         bool
	autoReplay            bool
	autoResubscribe       bool
	queueTTL              time.Duration
	offlineQueue          []types.QueryObject
	offlineQueueLoader    OfflineQueueLoader
	queueFilter           QueueFilter
	queueMaxSize          int
	reconnectionDelay     time.Duration
	replayInterval        time.Duration
	retrying              bool
	stopRetryingToConnect bool
	RequestHistory        map[string]time.Time
}

type QueueFilter interface {
	Filter(interface{}) bool
}

type defaultQueueFilter struct{}

func (qf defaultQueueFilter) Filter(interface{}) bool {
	return true
}

type OfflineQueueLoader interface {
	load() []types.QueryObject
}

func (ws *webSocket) SetQueueFilter(queueFilter QueueFilter) {
	ws.queueFilter = queueFilter
}

func NewWebSocket(host string, options types.Options) connection.Connection {
	var opts types.Options

	if options == nil {
		opts = types.NewOptions()
	} else {
		opts = options
	}
	ws := &webSocket{
		mu:                    &sync.Mutex{},
		queueTTL:              opts.GetQueueTTL(),
		offlineQueue:          make([]types.QueryObject, 0),
		queueMaxSize:          opts.GetQueueMaxSize(),
		channelsResult:        make(map[string]chan<- types.KuzzleResponse),
		subscriptions:         make(types.RoomList),
		eventListeners:        make(map[int]chan<- interface{}),
		RequestHistory:        make(map[string]time.Time),
		autoQueue:             opts.GetAutoQueue(),
		autoReconnect:         opts.GetAutoReconnect(),
		autoReplay:            opts.GetAutoReplay(),
		autoResubscribe:       opts.GetAutoResubscribe(),
		reconnectionDelay:     opts.GetReconnectionDelay(),
		replayInterval:        opts.GetReplayInterval(),
		state:                 state.Ready,
		retrying:              false,
		stopRetryingToConnect: false,
		queueFilter:           &defaultQueueFilter{},
	}
	ws.host = host

	if opts.GetOfflineMode() == types.Auto {
		ws.autoReconnect = true
		ws.autoQueue = true
		ws.autoReplay = true
		ws.autoResubscribe = true
	}
	return ws
}

func (ws *webSocket) Connect() (bool, error) {
	if !ws.isValidState() {
		return false, nil
	}

	addr := flag.String("addr", ws.host, "http service address")
	u := url.URL{Scheme: "ws", Host: *addr}
	resChan := make(chan []byte)

	socket, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		if ws.autoReconnect && !ws.retrying && !ws.stopRetryingToConnect {
			for err != nil {
				ws.retrying = true
				ws.EmitEvent(event.NetworkError, err)
				time.Sleep(ws.reconnectionDelay)
				ws.retrying = false
				socket, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
			}
		}
		ws.EmitEvent(event.NetworkError, err)
		return ws.wasConnected, err
	}

	if ws.lastUrl != ws.host {
		ws.wasConnected = false
		ws.lastUrl = ws.host
	}

	if ws.wasConnected {
		ws.EmitEvent(event.Reconnected, nil)
		ws.state = state.Connected
	} else {
		ws.EmitEvent(event.Connected, nil)
		ws.state = state.Connected
	}

	ws.ws = socket
	ws.stopRetryingToConnect = false

	if ws.autoReplay {
		ws.cleanQueue()
		ws.dequeue()
	}

	ws.RenewSubscriptions()

	go func() {
		for {
			_, message, err := ws.ws.ReadMessage()

			if err != nil {
				close(resChan)
				ws.ws.Close()
				ws.state = state.Offline
				if ws.autoQueue {
					ws.queuing = true
				}
				ws.EmitEvent(event.Disconnected, nil)
				return
			}
			go func() {
				resChan <- message
			}()
		}
	}()

	ws.listenChan = resChan
	go ws.listen()
	return ws.wasConnected, err
}

func (ws *webSocket) Send(query []byte, options types.QueryOptions, responseChannel chan<- types.KuzzleResponse, requestId string) error {
	if ws.state == state.Connected || (options != nil && !options.GetQueuable()) {
		ws.emitRequest(types.QueryObject{
			Query:     query,
			ResChan:   responseChannel,
			RequestId: requestId,
		})
	} else if ws.queuing || (options != nil && options.GetQueuable()) || ws.state == state.Initializing || ws.state == state.Connecting {
		ws.cleanQueue()

		if ws.queueFilter.Filter(query) {
			qo := types.QueryObject{
				Timestamp: time.Now(),
				ResChan:   responseChannel,
				Query:     query,
				RequestId: requestId,
				Options:   options,
			}
			ws.offlineQueue = append(ws.offlineQueue, qo)
			ws.EmitEvent(event.OfflineQueuePush, qo)
		}
	} else {
		ws.discardRequest(responseChannel, query)
	}
	return nil
}

func (ws *webSocket) discardRequest(responseChannel chan<- types.KuzzleResponse, query []byte) {
	if responseChannel != nil {
		responseChannel <- types.KuzzleResponse{Error: types.MessageError{Message: "Unable to execute request: not connected to a Kuzzle server.\nDiscarded request: " + string(query), Status: 400}}
	}
}

// Clean up the queue, ensuring the queryTTL and queryMaxSize properties are respected
func (ws *webSocket) cleanQueue() {
	now := time.Now()
	now = now.Add(-ws.queueTTL * time.Millisecond)

	// Clean queue of timed out query
	if ws.queueTTL > 0 {
		var query types.QueryObject
		for _, query = range ws.offlineQueue {
			if query.Timestamp.Before(now) {
				ws.offlineQueue = ws.offlineQueue[1:]
			} else {
				break
			}
		}
	}

	if ws.queueMaxSize > 0 && len(ws.offlineQueue) > ws.queueMaxSize {
		for len(ws.offlineQueue) > ws.queueMaxSize {
			eventListener := ws.eventListeners[event.OfflineQueuePop]
			if eventListener != nil {
				ws.eventListeners[event.OfflineQueuePop] <- ws.offlineQueue[0]
			}
			ws.offlineQueue = ws.offlineQueue[1:]
		}
	}
}

func (ws *webSocket) RegisterRoom(roomId, id string, room types.IRoom) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	if ws.subscriptions[roomId] == nil {
		ws.subscriptions[roomId] = make(map[string]types.IRoom)
	}
	ws.subscriptions[roomId][id] = room
}

func (ws *webSocket) UnregisterRoom(roomId string) {
	if ws.subscriptions[roomId] != nil {
		ws.mu.Lock()
		defer ws.mu.Unlock()
		delete(ws.subscriptions, roomId)
	}
}

func (ws *webSocket) listen() {
	for {
		var message types.KuzzleResponse
		var r collection.Room

		msg := <-ws.listenChan

		json.Unmarshal(msg, &message)
		m := message

		json.Unmarshal(m.Result, &r)

		if m.RoomId != "" && ws.subscriptions[m.RoomId] != nil {
			var notification types.KuzzleNotification

			json.Unmarshal(m.Result, &notification.Result)
			for _, rooms := range ws.subscriptions[m.RoomId] {
				channel := rooms.GetRealtimeChannel()
				if channel != nil {
					rooms.GetRealtimeChannel() <- notification
				}
			}
		}

		if len(ws.channelsResult) > 0 {
			if ws.channelsResult[m.RequestId] != nil {
				if message.Error.Message == "Token expired" {
					ws.EmitEvent(event.JwtExpired, nil)
				}

				// If this is a response to a query we simply broadcast the response to the corresponding channel
				ws.channelsResult[m.RequestId] <- message
				close(ws.channelsResult[m.RequestId])
				delete(ws.channelsResult, m.RequestId)
			}
		}
	}
}

// Adds a listener to a Kuzzle global event. When an event is fired, listeners are called in the order of their insertion.
func (ws *webSocket) AddListener(event int, channel chan<- interface{}) {
	ws.eventListeners[event] = channel
}

// Removes all listeners, either from all events and close channels
func (ws *webSocket) RemoveAllListeners() {
	for k := range ws.eventListeners {
		if ws.eventListeners[k] != nil {
			close(ws.eventListeners[k])
		}
		delete(ws.eventListeners, k)
	}
}

// Removes a listener from an event.
func (ws *webSocket) RemoveListener(event int) {
	delete(ws.eventListeners, event)
}

// Emit an event to all registered listeners
func (ws *webSocket) EmitEvent(event int, arg interface{}) {
	if ws.eventListeners[event] != nil {
		ws.eventListeners[event] <- arg
	}
}

func (ws *webSocket) StartQueuing() {
	if ws.state == state.Offline && !ws.autoQueue {
		ws.queuing = true
	}
}

func (ws *webSocket) StopQueuing() {
	if ws.state == state.Offline && !ws.autoQueue {
		ws.queuing = false
	}
}

func (ws *webSocket) FlushQueue() {
	ws.offlineQueue = ws.offlineQueue[:cap(ws.offlineQueue)]
}

func (ws *webSocket) ReplayQueue() {
	if ws.state != state.Offline && !ws.autoReplay {
		ws.cleanQueue()
		ws.dequeue()
	}
}

func (ws *webSocket) mergeOfflineQueueWithLoader() error {
	type query struct {
		requestId  string `json:"requestId"`
		controller string `json:"controller"`
		action     string `json:"action"""`
	}

	additionalOfflineQueue := ws.offlineQueueLoader.load()

	for _, additionalQuery := range additionalOfflineQueue {
		for _, offlineQuery := range ws.offlineQueue {
			q := query{}
			json.Unmarshal(additionalQuery.Query, &q)
			if q.requestId != "" || q.action != "" || q.controller != "" {
				offlineQ := query{}
				json.Unmarshal(offlineQuery.Query, &offlineQ)
				if q.requestId != offlineQ.requestId {
					ws.offlineQueue = append(ws.offlineQueue, additionalQuery)
				} else {
					additionalOfflineQueue = additionalOfflineQueue[:1]
				}
			} else {
				return errors.New("Invalid offline queue request. One or more missing properties: requestId, action, controller.")
			}
		}
	}
	return nil
}

func (ws *webSocket) dequeue() error {
	if ws.offlineQueueLoader != nil {
		err := ws.mergeOfflineQueueWithLoader()
		if err != nil {
			return err
		}
	}
	// Example from sdk where we have a good use of _
	if len(ws.offlineQueue) > 0 {
		for _, query := range ws.offlineQueue {
			ws.emitRequest(query)
			ws.offlineQueue = ws.offlineQueue[:1]
			ws.EmitEvent(event.OfflineQueuePop, query)
			time.Sleep(ws.replayInterval * time.Millisecond)
			ws.offlineQueue = ws.offlineQueue[:1]
		}
	} else {
		ws.queuing = false
	}
	return nil
}

func (ws *webSocket) emitRequest(query types.QueryObject) error {
	now := time.Now()
	now = now.Add(-MAX_EMIT_TIMEOUT * time.Second)

	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.channelsResult[query.RequestId] = query.ResChan

	err := ws.ws.WriteMessage(websocket.TextMessage, query.Query)
	if err != nil {
		return err
	}

	// Track requests made to allow Room.subscribeToSelf to work
	ws.RequestHistory[query.RequestId] = time.Now()
	for i, request := range ws.RequestHistory {
		if request.Before(now) {
			delete(ws.RequestHistory, i)
		}
	}

	return nil
}

func (ws *webSocket) Close() error {
	ws.stopRetryingToConnect = true
	ws.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	ws.state = state.Disconnected

	return ws.ws.Close()
}

func (ws webSocket) GetOfflineQueue() *[]types.QueryObject {
	return &ws.offlineQueue
}

func (ws *webSocket) isValidState() bool {
	switch ws.state {
	case state.Initializing, state.Ready, state.Disconnected, state.Error, state.Offline:
		return true
	}
	return false
}

func (ws *webSocket) GetState() *int {
	return &ws.state
}

func (ws webSocket) GetRequestHistory() *map[string]time.Time {
	return &ws.RequestHistory
}

func (ws webSocket) RenewSubscriptions() {
	for _, rooms := range ws.subscriptions {
		for _, room := range rooms {
			room.Renew(room.GetFilters(), room.GetRealtimeChannel(), room.GetResponseChannel())
		}
	}
}

func (ws webSocket) GetRooms() types.RoomList {
	return ws.subscriptions
}