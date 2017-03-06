package bytewriter

import (
	"fmt"
	"os"
	"reflect"
	"unsafe"
	path "path/filepath"
	"syscall"
)

// MemoryMappedWriter is a ByteBuffer that is also mapped into memory
type MemoryMappedWriter struct {
	*ByteWriter
	h syscall.Handle
	f *os.File
	loc  string // location of the memory mapped file
	size int    // size in bytes
}

// NewMemoryMappedWriter will create and return a new instance of a MemoryMappedWriter
func NewMemoryMappedWriter(loc string, size int) (*MemoryMappedWriter, error) {
	if _, err := os.Stat(loc); err == nil {
		err = os.Remove(loc)
		if err != nil {
			return nil, err
		}
	}

	// ensure destination directory exists
	dir := path.Dir(loc)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(loc, syscall.O_CREAT|syscall.O_RDWR|syscall.O_EXCL, 0644)
	if err != nil {
		return nil, err
	}

	l, err := f.Write(make([]byte, size))
	if err != nil {
		return nil, err
	}
	if l < size {
		return nil, fmt.Errorf("Could not initialize %d bytes", size)
	}

	maxSizeHigh := uint32(size >> 32)
	maxSizeLow := uint32(size & 0xFFFFFFFF)
	flProtect := uint32(syscall.PAGE_READWRITE)
	hfile := uintptr(f.Fd())
	h, errno := syscall.CreateFileMapping(syscall.Handle(hfile), nil, flProtect, maxSizeHigh, maxSizeLow, nil)
	if h == 0 {
		return nil, os.NewSyscallError("CreateFileMapping", errno)
	}
	dwDesiredAccess := uint32(syscall.FILE_MAP_WRITE)
	addr, errno := syscall.MapViewOfFile(h, dwDesiredAccess, 0, 0, uintptr(size))
	if addr == 0 {
		return nil, os.NewSyscallError("MapViewOfFile", errno)
	}

	/*
	b, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	*/

	d := []byte{}
	dh := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	dh.Data = addr
	dh.Len = size
	dh.Cap = dh.Len

	return &MemoryMappedWriter{
		NewByteWriterSlice(d),
		h, f,
		loc,
		size,
	}, nil
}

// Unmap will manually delete the memory mapping of a mapped buffer
func (b *MemoryMappedWriter) Unmap(removefile bool) error {

	dh := (*reflect.SliceHeader)(unsafe.Pointer(&(b.buffer)))
	addr := dh.Data
	err := syscall.UnmapViewOfFile(addr)
	if err != nil {
		return err
	}

	syscall.CloseHandle(syscall.Handle(b.h))
	b.f.Close()
	if removefile {
		if err := os.Remove(b.loc); err != nil {
			return err
		}
	}

	return nil
}
