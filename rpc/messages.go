package rpc

// Request represents an incoming JSON-RPC request
type Request struct {
	ID     string `json:"id"`
	Async  bool   `json:"async,omitempty"`
	Method string `json:"method,omitempty"`
	Params []any  `json:"params,omitempty"`
}

// Response represents an outgoing JSON-RPC response
type Response struct {
	ID     string `json:"id"`
	Error  *Error `json:"error,omitempty"`
	Result any    `json:"result,omitempty"`
}

// Error represents a JSON-RPC error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

func (r *Error) Error() string {
	return r.Message
}
