package bytewriter

import (
	"os"
	path "path/filepath"
	"testing"
)

func TestMemoryMappedWriter(t *testing.T) {
	filename := "bytebuffer_memorymappedwriter_test.tmp"
	loc := path.Join(os.TempDir(), filename)

	if _, err := os.Stat(loc); err == nil {
		err = os.Remove(loc)
		if err != nil {
			t.Fatal("Cannot proceed with test as cannot remove spec file")
		}
	}

	w, err := NewMemoryMappedWriter(loc, 10)
	if err != nil {
		t.Fatal("Cannot proceed with test as create writer failed:", err)
	}

	if _, err = os.Stat(loc); err != nil {
		t.Fatalf("No File created at %v despite the Buffer being initialized", loc)
	}

	_, err = w.WriteString("x", 5)
	if err != nil {
		t.Fatal("Cannot Write to MemoryMappedWriter")
	}

	reader, err := os.Open(loc)
	if err != nil {
		t.Fatal("Cannot open memory mapped file")
	}

	data := make([]byte, 10)
	_, err = reader.Read(data)
	if err != nil {
		t.Fatal("Cannot read data from memory mapped file")
	}

	if data[5] != 'x' {
		t.Error("Data Written in buffer not getting reflected in file")
	}
	reader.Close()

	testUnmap(w, loc, t)
}

func testUnmap(w *MemoryMappedWriter, loc string, t *testing.T) {
	var err = w.Unmap(true)
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(loc); err == nil {
		t.Error("Memory Mapped File not getting deleted on Unmap")
	}
}
