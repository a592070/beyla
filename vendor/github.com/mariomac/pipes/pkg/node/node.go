// Package node provides functionalities to create nodes and interconnect them.
// A Node is a function container that can be connected via channels to other nodes.
// A node can send data to multiple nodes, and receive data from multiple nodes.
// Deprecated package. Use github.com/mariomac/pipes/pipe package
//
//nolint:unused
package node

import (
	"errors"
	"reflect"

	"github.com/mariomac/pipes/pkg/node/internal/connect"
)

// StartFunc is a function that receives a writable channel as unique argument, and sends
// value to that channel during an indefinite amount of time.
// Deprecated package. Use github.com/mariomac/pipes/pipe package
type StartFunc[OUT any] func(out chan<- OUT)

// MiddleFunc is a function that receives a readable channel as first argument,
// and a writable channel as second argument.
// It must process the inputs from the input channel until it's closed.
// Deprecated package. Use github.com/mariomac/pipes/pipe package
type MiddleFunc[IN, OUT any] func(in <-chan IN, out chan<- OUT)

// TerminalFunc is a function that receives a readable channel as unique argument.
// It must process the inputs from the input channel until it's closed.
// Deprecated package. Use github.com/mariomac/pipes/pipe package
type TerminalFunc[IN any] func(in <-chan IN)

// TODO: OutType and InType methods are candidates for deprecation

// Sender is any node that can send data to another node: node.Start, node.Middle and node.Bypass
// Deprecated package. Use github.com/mariomac/pipes/pipe package
type Sender[OUT any] interface {
	// SendTo connect a sender with a group of receivers
	SendTo(...Receiver[OUT])
	// OutType returns the inner type of the Sender's output channel
	OutType() reflect.Type
}

// Receiver is any node that can receive data from another node: node.Bypass, node.Middle and node.Terminal
// Deprecated package. Use github.com/mariomac/pipes/pipe package
type Receiver[IN any] interface {
	isStarted() bool
	start()
	// joiners will usually return only one joiner instance but in
	// the case of a BypassNode, which might return the joiners of
	// all their destination nodes
	joiners() []*connect.Joiner[IN]
	// InType returns the inner type of the Receiver's input channel
	InType() reflect.Type
}

// SenderReceiver is any node that can both send and receive data: node.Bypass or node.Middle.
// Deprecated package. Use github.com/mariomac/pipes/pipe package
type SenderReceiver[IN, OUT any] interface {
	Receiver[IN]
	Sender[OUT]
}

// Start nodes are the starting points of a graph. This is, all the nodes that bring information
// from outside the graph: e.g. because they generate them or because they acquire them from an
// external source like a Web Service.
// A graph must have at least one Start or StartDemux node.
// An Start node must have at least one output node.
// Deprecated package. Use github.com/mariomac/pipes/pipe package
type Start[OUT any] struct {
	receiverGroup[OUT]
	funs []StartFunc[OUT]
}

// Middle is any intermediate node that receives data from another node, processes/filters it,
// and forwards the data to another node.
// An Middle node must have at least one output node.
// Deprecated package. Use github.com/mariomac/pipes/pipe package
type Middle[IN, OUT any] struct {
	outs    []Receiver[OUT]
	inputs  connect.Joiner[IN]
	started bool
	fun     MiddleFunc[IN, OUT]
	outType reflect.Type
	inType  reflect.Type
}

func (mn *Middle[IN, OUT]) joiners() []*connect.Joiner[IN] {
	return []*connect.Joiner[IN]{&mn.inputs}
}

func (mn *Middle[IN, OUT]) isStarted() bool {
	return mn.started
}

func (mn *Middle[IN, OUT]) SendTo(outputs ...Receiver[OUT]) {
	mn.outs = append(mn.outs, outputs...)
}

func (mn *Middle[IN, OUT]) OutType() reflect.Type {
	return mn.outType
}

func (mn *Middle[IN, OUT]) InType() reflect.Type {
	return mn.inType
}

// Terminal is any node that receives data from another node and does not forward it to another node,
// but can process it and send the results to outside the graph (e.g. memory, storage, web...)
// Deprecated package. Use github.com/mariomac/pipes/pipe package
type Terminal[IN any] struct {
	inputs  connect.Joiner[IN]
	started bool
	fun     TerminalFunc[IN]
	done    chan struct{}
	inType  reflect.Type
}

func (tn *Terminal[IN]) joiners() []*connect.Joiner[IN] {
	if tn == nil {
		return nil
	}
	return []*connect.Joiner[IN]{&tn.inputs}
}

func (tn *Terminal[IN]) isStarted() bool {
	if tn == nil {
		return false
	}
	return tn.started
}

// Done returns a channel that is closed when the Terminal node has ended its processing. This
// is, when all its inputs have been also closed. Waiting for all the Terminal nodes to finish
// allows blocking the execution until all the data in the graph has been processed and all the
// previous stages have ended
// Deprecated package. Use github.com/mariomac/pipes/pipe package
func (tn *Terminal[IN]) Done() <-chan struct{} {
	if tn == nil {
		closed := make(chan struct{})
		close(closed)
		return closed
	}
	return tn.done
}

func (tn *Terminal[IN]) InType() reflect.Type {
	return tn.inType
}

// AsStart wraps a group of StartFunc with the same signature into a Start node.
func AsStart[OUT any](funs ...StartFunc[OUT]) *Start[OUT] {
	var out OUT
	return &Start[OUT]{
		funs: funs,
		receiverGroup: receiverGroup[OUT]{
			outType: reflect.TypeOf(out),
		},
	}
}

// AsMiddle wraps an MiddleFunc into an Middle node.
func AsMiddle[IN, OUT any](fun MiddleFunc[IN, OUT], opts ...Option) *Middle[IN, OUT] {
	var in IN
	var out OUT
	options := getOptions(opts...)
	return &Middle[IN, OUT]{
		inputs:  connect.NewJoiner[IN](options.channelBufferLen),
		fun:     fun,
		inType:  reflect.TypeOf(in),
		outType: reflect.TypeOf(out),
	}
}

// AsTerminal wraps a TerminalFunc into a Terminal node.
func AsTerminal[IN any](fun TerminalFunc[IN], opts ...Option) *Terminal[IN] {
	var i IN
	options := getOptions(opts...)
	return &Terminal[IN]{
		inputs: connect.NewJoiner[IN](options.channelBufferLen),
		fun:    fun,
		done:   make(chan struct{}),
		inType: reflect.TypeOf(i),
	}
}

// Start starts the function wrapped in the Start node. This method should be invoked
// for all the start nodes of the same graph, so the graph can properly start and finish.
// Deprecated package. Use github.com/mariomac/pipes/pipe package
func (sn *Start[OUT]) Start() {
	// a nil start node can be started without no effect on the graph.
	// this allows setting optional nillable start nodes and let start all of them
	// as a group in a more convenient way
	if sn == nil {
		return
	}
	forker, err := sn.receiverGroup.StartReceivers()
	if err != nil {
		panic("Start: " + err.Error())
	}
	for fn := range sn.funs {
		fun := sn.funs[fn]
		go func() {
			fun(forker.AcquireSender())
			forker.ReleaseSender()
		}()
	}
}

func (mn *Middle[IN, OUT]) start() {
	if len(mn.outs) == 0 {
		panic("Middle node should have outputs")
	}
	mn.started = true
	joiners := make([]*connect.Joiner[OUT], 0, len(mn.outs))
	for _, out := range mn.outs {
		joiners = append(joiners, out.joiners()...)
		if !out.isStarted() {
			out.start()
		}
	}
	forker := connect.Fork(joiners...)
	go func() {
		mn.fun(mn.inputs.Receiver(), forker.AcquireSender())
		forker.ReleaseSender()
	}()
}

func (tn *Terminal[IN]) start() {
	if tn == nil {
		return
	}
	tn.started = true
	go func() {
		tn.fun(tn.inputs.Receiver())
		close(tn.done)
	}()
}

func getOptions(opts ...Option) creationOptions {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// receiverGroup connects a sender node with a collection
// of Receiver nodes through a common connect.Forker instance.
// Deprecated package. Use github.com/mariomac/pipes/pipe package
type receiverGroup[OUT any] struct {
	Outs    []Receiver[OUT]
	outType reflect.Type
}

// SendTo connects a group of receivers to the current receiverGroup
func (sn *Start[OUT]) SendTo(outputs ...Receiver[OUT]) {
	// a nil start node can be operated without no effect on the graph.
	// this allows connecting optional nillable start nodes and let start all of them
	// as a group in a more convenient way
	if sn != nil {
		sn.receiverGroup.SendTo(outputs...)
	}
}

func (rg *receiverGroup[OUT]) SendTo(outputs ...Receiver[OUT]) {
	rg.Outs = append(rg.Outs, outputs...)
}

// OutType is the common input type of the receivers
// (output of the receiver group)
func (rg *receiverGroup[OUT]) OutType() reflect.Type {
	return rg.outType
}

// StartReceivers start the receivers and return a connection
// forker to them
func (rg *receiverGroup[OUT]) StartReceivers() (*connect.Forker[OUT], error) {
	if len(rg.Outs) == 0 {
		return nil, errors.New("node should have outputs")
	}
	joiners := make([]*connect.Joiner[OUT], 0, len(rg.Outs))
	for _, out := range rg.Outs {
		joiners = append(joiners, out.joiners()...)
		if !out.isStarted() {
			out.start()
		}
	}
	forker := connect.Fork(joiners...)
	return &forker, nil
}