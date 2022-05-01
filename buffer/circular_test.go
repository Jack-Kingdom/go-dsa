package buffer

import "testing"

func TestCircularBuffer_Write(t *testing.T) {
	buf := Get(20)
	defer Put(buf)

	circularBuffer := NewCircularBuffer(buf)

	content := []byte("Hello World")

	n, err := circularBuffer.Write(content)
	if err!=nil {
		t.Errorf("Error writing to circularBuffer: %s", err)
	}
	if n!=len(content) {
		t.Errorf("Expected to write %d bytes, but wrote %d", len(content), n)
	}
	if circularBuffer.Len()!=len(content) {
		t.Errorf("Expected circularBuffer length to be %d, but was %d", len(content), circularBuffer.Len())
	}

	readBuffer := Get(20)
	n, err = circularBuffer.Read(readBuffer)
	if err !=nil {
		t.Errorf("Error reading from circularBuffer: %s", err)
	}
	if n!=len(content) {
		t.Errorf("Expected to read %d bytes, but read %d", len(content), n)
	}
	if string(readBuffer[:n])!=string(content) {
		t.Errorf("Expected to read %s, but read %s", string(content), string(readBuffer))
	}

	if circularBuffer.Len()!=0 {
		t.Errorf("Expected circularBuffer length to be 0, but was %d", circularBuffer.Len())
	}
}