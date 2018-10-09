package deribit

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// RPCRequest is what we send to the remote
type RPCRequest struct {
	Action    string                 `json:"action"`
	ID        uint64                 `json:"id"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	Sig       string                 `json:"sig,omitempty"`
}

// GenerateSig creates the signature required for private endpoints
func (r *RPCRequest) GenerateSig(key, secret string) error {
	nonce := time.Now().UnixNano() / int64(time.Millisecond)
	sigString := fmt.Sprintf("_=%d&_ackey=%s&_acsec=%s&_action=%s", nonce, key, secret, r.Action)

	// Append args if present
	if len(r.Arguments) != 0 {
		var argsString string

		// We have to this to sort by keys
		keys := make([]string, 0, len(r.Arguments))
		for key := range r.Arguments {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := r.Arguments[k]
			var s string

			switch t := v.(type) {
			case []SubscriptionEvent:
				var str = make([]string, len(t))
				for _, j := range t {
					str = append(str, string(j))
				}
				s = strings.Join(str, "")
			case []string:
				s = strings.Join(t, "")
			case bool:
				s = strconv.FormatBool(t)
			case int:
				s = string(t)
			case float64:
				s = fmt.Sprintf("%f", t)
			case string:
				s = t
			default:
				// Absolutely panic here
				panic(fmt.Sprintf("Cannot generate sig string: Unable to handle arg of type %T", t))
			}
			argsString += fmt.Sprintf("&%s=%s", k, s)
		}
		sigString += argsString
	}
	hasher := sha256.New()
	hasher.Write([]byte(sigString))
	sigHash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	r.Sig = fmt.Sprintf("%s.%d.%s", key, nonce, sigHash)
	return nil
}

// RPCResponse is what we receive from the remote
type RPCResponse struct {
	ID            uint64             `json:"id"`
	Success       bool               `json:"success"`
	Error         int                `json:"error"`
	Testnet       bool               `json:"testnet"`
	Message       string             `json:"message"`
	UsIn          uint64             `json:"usIn"`
	UsOut         uint64             `json:"usOut"`
	UsDiff        uint64             `json:"usDiff"`
	Result        interface{}        `json:"result"`
	APIBuild      string             `json:"apiBuild"`
	Notifications []*RPCNotification `json:"notifications"`
}

// RPCCall represents the entire call from request to response
type RPCCall struct {
	Req   RPCRequest
	Res   RPCResponse
	Done  chan bool
	Error error
}

// NewRPCCall returns a new RPCCall initialised with a done channel and request
func NewRPCCall(req RPCRequest) *RPCCall {
	done := make(chan bool)
	return &RPCCall{
		Req:  req,
		Done: done,
	}
}

// makeRequest makes a request over the websocket and waits for a response with a timeout
func (e *Exchange) makeRequest(req RPCRequest) (*RPCResponse, error) {
	e.mutex.Lock()
	id := e.counter
	e.counter++
	req.ID = id
	call := NewRPCCall(req)
	e.pending[id] = call

	//j, _ := json.Marshal(req)
	//fmt.Printf("Req: %s\n", j)

	if err := e.conn.WriteJSON(&req); err != nil {
		delete(e.pending, id)
		e.mutex.Unlock()
		return nil, err
	}
	e.mutex.Unlock()
	select {
	case <-call.Done:
	case <-time.After(2 * time.Second):
		call.Error = fmt.Errorf("Request %d timed out", id)
	}

	if !call.Res.Success {
		call.Error = fmt.Errorf("Request failed with: %s", call.Res.Message)
	}
	if call.Error != nil {
		return nil, call.Error
	}
	return &call.Res, nil
}

// read takes messages off the websocket and deals with them accordingly
func (e *Exchange) read() {
	var err error
Loop:
	for {
		select {
		case <-e.stop:
			fmt.Println("HELLO")
			break Loop
		default:
			var res RPCResponse
			if err := e.conn.ReadJSON(&res); err != nil {
				err = fmt.Errorf("Error reading message: %q", err)
				break Loop
			}
			//j, _ := json.Marshal(res)
			//fmt.Printf("Res: %s\n", j)

			// Notifications do not have an ID field
			if res.Notifications != nil {
				for _, n := range res.Notifications {
					e.mutex.Lock()
					sub := e.subscriptions[n.Message]
					e.mutex.Unlock()
					if sub == nil {
						// Send error to main error channel
						e.errors <- fmt.Errorf("No subscription found for %s", n.Message)
						break Loop
					}
					// Send the notification to the right channel
					sub.Data <- n.Result
				}
			} else {
				e.mutex.Lock()
				call := e.pending[res.ID]
				delete(e.pending, res.ID)
				e.mutex.Unlock()

				if call == nil {
					err = fmt.Errorf("No pending request found for response ID %d", res.ID)
					break Loop
				}
				// Attach the response to the call and close
				call.Res = res
				call.Done <- true
			}
		}
	}
	// On error, terminate all calls by sending on their done channels
	if err != nil {
		e.mutex.Lock()
		for _, call := range e.pending {
			call.Error = err
			call.Done <- true
		}
		e.mutex.Unlock()
	}
}