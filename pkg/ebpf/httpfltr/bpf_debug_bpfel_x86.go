// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64
// +build 386 amd64

package httpfltr

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type bpf_debugHttpConnectionInfoT struct {
	S_h    uint64
	S_l    uint64
	D_h    uint64
	D_l    uint64
	S_port uint16
	D_port uint16
	Flags  uint32
}

type bpf_debugSockArgsT struct {
	Addr       uint64
	AcceptTime uint64
}

// loadBpf_debug returns the embedded CollectionSpec for bpf_debug.
func loadBpf_debug() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_Bpf_debugBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load bpf_debug: %w", err)
	}

	return spec, err
}

// loadBpf_debugObjects loads bpf_debug and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*bpf_debugObjects
//	*bpf_debugPrograms
//	*bpf_debugMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadBpf_debugObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadBpf_debug()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// bpf_debugSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpf_debugSpecs struct {
	bpf_debugProgramSpecs
	bpf_debugMapSpecs
}

// bpf_debugSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpf_debugProgramSpecs struct {
	KprobeTcpConnect    *ebpf.ProgramSpec `ebpf:"kprobe_tcp_connect"`
	KretprobeSockAlloc  *ebpf.ProgramSpec `ebpf:"kretprobe_sock_alloc"`
	KretprobeSysAccept4 *ebpf.ProgramSpec `ebpf:"kretprobe_sys_accept4"`
	KretprobeSysConnect *ebpf.ProgramSpec `ebpf:"kretprobe_sys_connect"`
}

// bpf_debugMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpf_debugMapSpecs struct {
	ActiveAcceptArgs    *ebpf.MapSpec `ebpf:"active_accept_args"`
	ActiveConnectArgs   *ebpf.MapSpec `ebpf:"active_connect_args"`
	Events              *ebpf.MapSpec `ebpf:"events"`
	FilteredConnections *ebpf.MapSpec `ebpf:"filtered_connections"`
}

// bpf_debugObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadBpf_debugObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpf_debugObjects struct {
	bpf_debugPrograms
	bpf_debugMaps
}

func (o *bpf_debugObjects) Close() error {
	return _Bpf_debugClose(
		&o.bpf_debugPrograms,
		&o.bpf_debugMaps,
	)
}

// bpf_debugMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadBpf_debugObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpf_debugMaps struct {
	ActiveAcceptArgs    *ebpf.Map `ebpf:"active_accept_args"`
	ActiveConnectArgs   *ebpf.Map `ebpf:"active_connect_args"`
	Events              *ebpf.Map `ebpf:"events"`
	FilteredConnections *ebpf.Map `ebpf:"filtered_connections"`
}

func (m *bpf_debugMaps) Close() error {
	return _Bpf_debugClose(
		m.ActiveAcceptArgs,
		m.ActiveConnectArgs,
		m.Events,
		m.FilteredConnections,
	)
}

// bpf_debugPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadBpf_debugObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpf_debugPrograms struct {
	KprobeTcpConnect    *ebpf.Program `ebpf:"kprobe_tcp_connect"`
	KretprobeSockAlloc  *ebpf.Program `ebpf:"kretprobe_sock_alloc"`
	KretprobeSysAccept4 *ebpf.Program `ebpf:"kretprobe_sys_accept4"`
	KretprobeSysConnect *ebpf.Program `ebpf:"kretprobe_sys_connect"`
}

func (p *bpf_debugPrograms) Close() error {
	return _Bpf_debugClose(
		p.KprobeTcpConnect,
		p.KretprobeSockAlloc,
		p.KretprobeSysAccept4,
		p.KretprobeSysConnect,
	)
}

func _Bpf_debugClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed bpf_debug_bpfel_x86.o
var _Bpf_debugBytes []byte
