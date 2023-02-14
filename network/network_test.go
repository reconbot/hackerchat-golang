package network

import (
	"fmt"
	"testing"
)

func TestOpenConnection(t *testing.T) {
	connection, err := NewConnection()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if connection == nil {
		t.Fatalf("No connection created")
	}
}

func TestWrite(t *testing.T) {
	connection, err := NewConnection()
	if err != nil {
		t.Fatalf(err.Error())
	}
	data := "I don't know"
	written, err := connection.Write([]byte(data))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(data) != written {
		t.Fatalf("Wrong number of bytes written, should be %v but got %v", len(data), written)
	}
}

func TestRead(t *testing.T) {
	connection, err := NewConnection()
	if err != nil {
		t.Fatalf(err.Error())
	}
	data := "I don't know"
	done := make(chan bool)
	go func() {
		fmt.Printf("fuck")
		readData, err := connection.Read()
		if err != nil {
			t.Errorf(err.Error())
		}
		if string(readData) != data {
			t.Errorf("Expected %v but got %v", data, string(readData))
		}
		close(done)
	}()

	written, err := connection.Write([]byte(data))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(data) != written {
		t.Fatalf("Wrong number of bytes written, should be %v but got %v", len(data), written)
	}
	<-done
}
