package main

/*
#include <libusb20.h>
#cgo LDFLAGS: -lusb
*/
import "C"
import (
	"errors"
	"fmt"
)

type Error int

const (
	// Success
	Success = C.LIBUSB20_SUCCESS

	// Input/output error
	ErrIO Error = C.LIBUSB20_ERROR_IO
	// Invalid parameter
	ErrInvalidParam Error = C.LIBUSB20_ERROR_INVALID_PARAM
	// Access denied (insufficient permissions)
	ErrAccess Error = C.LIBUSB20_ERROR_ACCESS
	// No such device (it may have been disconnected)
	ErrNoDevice Error = C.LIBUSB20_ERROR_NO_DEVICE
	// Entity not found
	ErrNotFound Error = C.LIBUSB20_ERROR_NOT_FOUND
	// Resource busy
	ErrBusy Error = C.LIBUSB20_ERROR_BUSY
	// Operation timed out
	ErrTimeout Error = C.LIBUSB20_ERROR_TIMEOUT
	// Overflow
	ErrOverflow Error = C.LIBUSB20_ERROR_OVERFLOW
	// Pipe error
	ErrPipe Error = C.LIBUSB20_ERROR_PIPE
	// System call interrupted (perhaps due to signal)
	ErrInterrupted Error = C.LIBUSB20_ERROR_INTERRUPTED
	// Insufficient memory
	ErrNoMem Error = C.LIBUSB20_ERROR_NO_MEM
	// Operation not supported or unimplemented on this platform
	ErrNotSupported Error = C.LIBUSB20_ERROR_NOT_SUPPORTED
	// Other error
	ErrOther Error = C.LIBUSB20_ERROR_OTHER
)

const (
	DevModeHost   = C.LIBUSB20_MODE_HOST
	DevModeDevice = C.LIBUSB20_MODE_DEVICE
)

const (
	DevSpeedUnknown  = C.LIBUSB20_SPEED_UNKNOWN
	DevSpeedLow      = C.LIBUSB20_SPEED_LOW
	DevSpeedFull     = C.LIBUSB20_SPEED_FULL
	DevSpeedHigh     = C.LIBUSB20_SPEED_HIGH
	DevSpeedVariable = C.LIBUSB20_SPEED_VARIABLE
	DevSpeedSuper    = C.LIBUSB20_SPEED_SUPER
)

const (
	DevPowerOff     = C.LIBUSB20_POWER_OFF
	DevPowerOn      = C.LIBUSB20_POWER_ON
	DevPowerSave    = C.LIBUSB20_POWER_SAVE
	DevPowerSuspend = C.LIBUSB20_POWER_SUSPEND
	DevPowerResume  = C.LIBUSB20_POWER_RESUME
)

func (e Error) Error() string {
	return fmt.Sprintf("%s (%d): %s",
		C.GoString(C.libusb20_error_name(C.int(e))),
		int(e),
		C.GoString(C.libusb20_strerror(C.int(e))))
}

type Device struct {
	d *C.struct_libusb20_device
}

func (d *Device) BackendName() string {
	return C.GoString(C.libusb20_dev_get_backend_name(d.d))
}

func (d *Device) Description() string {
	return C.GoString(C.libusb20_dev_get_desc(d.d))
}

func (d *Device) String() string {
	return d.Description()
}

func (d *Device) Open(transferMax uint16) error {
	if res := C.libusb20_dev_open(d.d, C.ushort(transferMax)); res != Success {
		return Error(res)
	}
	return nil
}

func (d *Device) Close() error {
	if res := C.libusb20_dev_close(d.d); res != Success {
		return Error(res)
	}
	return nil
}

type Backend struct {
	b     *C.struct_libusb20_backend
	currd *C.struct_libusb20_device
}

var ErrNoBackend = errors.New("no backend found")

func NewBackend() (*Backend, error) {
	b := C.libusb20_be_alloc_default()
	if b == nil {
		return nil, ErrNoBackend
	}
	return &Backend{b: b}, nil
}

func (b *Backend) Free() {
	C.libusb20_be_free(b.b)
}

func (b *Backend) Scan() bool {
	b.currd = C.libusb20_be_device_foreach(b.b, b.currd)
	return b.currd != nil
}

func (b *Backend) Device() *Device {
	return &Device{d: b.currd}
}
