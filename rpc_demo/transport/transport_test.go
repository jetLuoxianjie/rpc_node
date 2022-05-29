package transport

import (
	"net"
	"sync"
	"testing"
	"time"
)

func TestTransport_ReadWrite(t *testing.T) {
	addr := "localhost:3212"
	dataToSend := "Hello World"
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {

		defer wg.Done()
		l, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		defer l.Close()
		conn, _ := l.Accept()
		time.Sleep(100 * time.Millisecond)
		s := NewTransport(conn)
		err = s.Send([]byte(dataToSend))
		t.Log("listen and accept")
		if err != nil {
			t.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		tp := NewTransport(conn)
		data, err := tp.Read()
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != dataToSend {
			t.FailNow()
		}
	}()
	wg.Wait()
}
