// +build windows darwin

package process

// #include "process.h"
// #cgo CFLAGS: -std=c99
import "C"
import (
	"github.com/polyverse-security/masche/cresponse"
	"unsafe"
)

type process struct {
	hndl C.process_handle_t
	pid  C.pid_t
}

func (p process) Pid() int {
	return int(p.pid)
}

func (p process) Handle() uintptr {
	return uintptr(p.hndl)
}

func (p process) Close() (harderror error, softerrors []error) {
	resp := C.close_process_handle(p.hndl)
	defer C.response_free(resp)
	return cresponse.GetResponsesErrors(unsafe.Pointer(resp))
}

func openFromPid(pid int) (p Process, harderror error, softerrors []error) {
	var result process

	resp := C.open_process_handle(C.pid_t(pid), &result.hndl)
	harderror, softerrors = cresponse.GetResponsesErrors(unsafe.Pointer(resp))
	C.response_free(resp)

	if harderror == nil {
		result.pid = C.pid_t(pid)
	} else {
		resp = C.close_process_handle(result.hndl)
		C.response_free(resp)
	}

	return result, harderror, softerrors
}
