package iso8583server_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/jattento/go-iso8583/pkg/iso8583server"
	"github.com/jattento/go-iso8583/pkg/mti"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

var testDefaultLogConfig = iso8583server.LogOn{
	ErrorAcceptIncomingConnection:      false,
	ErrorReadMTI:                       false,
	ErrorUndefinedHandler:              false,
	ErrorReadConnection:                false,
	ErrorSettingConnectionReadDeadline: false,

	ServingConnection: false,
}

func TestSimpleHandler(t *testing.T) {
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		t.Fatal(err)
	}

	sv, err := iso8583server.New(
		iso8583server.OptionLogConfiguration(&testDefaultLogConfig),
		iso8583server.OptionListener(l),
	)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	sv.AddTopPriorityHandler(func(r iso8583server.Response, message []byte) {
		if _, err := r.Write(message); err != nil {
			t.Fatal(err)
		}
		if err := r.Close(); err != nil {
			t.Fatal(err)
		}
	}, func(mtiValue mti.MTI) bool {
		return mtiValue.Version() == mti.Version8583To1993
	})

	go sv.Start()
	defer func() {
		assert.Nil(t, sv.Close())
	}()

	cn, err := net.Dial("tcp", "localhost:8081")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	defer func() {
		assert.Nil(t, cn.Close())
	}()

	responseBuffer := bytes.NewBuffer(nil)
	go func() {
		for {
			b := make([]byte, 1)
			n, _ := cn.Read(b)
			responseBuffer.Write(b[:n])
		}
	}()

	bLen := make([]byte, 2)
	binary.BigEndian.PutUint16(bLen, 10)
	_, err = cn.Write(bLen)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	req := []byte("1000111111")
	_, err = cn.Write(req)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	if !waitFor(func() bool { return string(responseBuffer.Bytes()) == string(req) }, time.Minute) {
		t.Fatal("test timeout waiting for server reply")
	}
}

func TestSimpleConcurrent(t *testing.T) {
	concurrentConnection := 100
	amountOfMessagePerConnection := 100

	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		t.Fatal(err)
	}

	sv, err := iso8583server.New(
		iso8583server.OptionLogConfiguration(&testDefaultLogConfig),
		iso8583server.OptionListener(l),
	)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	sv.AddTopPriorityHandler(func(r iso8583server.Response, message []byte) {
		if _, err := r.Write(message); err != nil {
			t.Fatal(err)
		}
		if err := r.Close(); err != nil {
			t.Fatal(err)
		}
	}, func(mtiValue mti.MTI) bool {
		return mtiValue.Version() == mti.Version8583To1993
	})

	go sv.Start()
	defer func() {
		assert.Nil(t, sv.Close())
	}()

	var connections []net.Conn
	for n := 0; n < concurrentConnection; n++ {
		cn, err := net.Dial("tcp", "localhost:8082")
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		defer func() {
			assert.Nil(t, cn.Close())
		}()
		connections = append(connections, cn)
	}

	var rspBuffers []*bytes.Buffer
	for n := 0; n < concurrentConnection; n++ {
		buff := bytes.NewBuffer(nil)

		go func(cn int) {
			io.Copy(buff, connections[cn])
		}(n)

		rspBuffers = append(rspBuffers, buff)
	}

	// Here connection are write concurrently
	for n := 0; n < concurrentConnection; n++ {
		go func(m int) {
			for n := 0; n < amountOfMessagePerConnection; n++ {
				message := "1000" + strconv.Itoa(m)
				bLen := make([]byte, 2)
				binary.BigEndian.PutUint16(bLen, uint16(len(message)))

				_, err = connections[m].Write(bLen)
				if !assert.Nil(t, err) {
					t.FailNow()
				}

				_, err = connections[m].Write([]byte(message))
				if !assert.Nil(t, err) {
					t.FailNow()
				}
			}
		}(n)
	}

	for n := 0; n < concurrentConnection; n++ {
		if !waitFor(func() bool {
			return string(rspBuffers[n].Bytes()) ==
				strings.Repeat("1000"+strconv.Itoa(n), amountOfMessagePerConnection)
		}, time.Minute) {
			t.Fatal("test timeout waiting for server reply")
		}
	}
}

func waitFor(conditional func() bool, timeout time.Duration) bool {
	successChan := make(chan struct{})
	go func() {
		for !conditional() {
			// wait.
		}
		close(successChan)
	}()

	timeoutChan := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		close(timeoutChan)
	}()

	select {
	case <-successChan:
		return true
	case <-timeoutChan:
		return false
	}
}

func TestSimpleConnections(t *testing.T) {
	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		t.Fatal(err)
	}

	var cName int

	sv, err := iso8583server.New(
		iso8583server.OptionLogConfiguration(&testDefaultLogConfig),
		iso8583server.OptionListener(l),
		iso8583server.OptionDeactivatedConnectionsCapacity(1),
		iso8583server.OptionConnIdGenerator(func() (string, error) {
			defer func() {cName++}()
			return strconv.Itoa(cName),nil
		}),
	)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	sv.AddTopPriorityHandler(func(r iso8583server.Response, message []byte) {
		if _, err := r.Write(message); err != nil {
			t.Fatal(err)
		}
		if err := r.Close(); err != nil {
			t.Fatal(err)
		}
	}, func(mtiValue mti.MTI) bool {
		return mtiValue.Version() == mti.Version8583To1993
	})

	go sv.Start()
	defer func() {
		assert.Nil(t, sv.Close())
	}()

	cn, err := net.Dial("tcp", "localhost:8082")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	defer func() {
		assert.Nil(t, cn.Close())
	}()

	keepConnection(cn)

	time.Sleep(time.Second)
	fmt.Println(sv.Connections.GetAll())
	if len(sv.Connections.GetAll()) != 1 {
		fffdgsfmt.Println(sv.Connections.GetAll())
		t.Fatal("connection isnt up")
	}
}

func keepConnection(c net.Conn) {
	go func() {
		for {
			_, _ = io.Copy(io.MultiWriter(), c)
		}
	}()
	go func() {
		bLen := make([]byte, 2)
		binary.BigEndian.PutUint16(bLen, 10)
		_, _ = c.Write(bLen)
		req := []byte("1000111111")
		_, _ = c.Write(req)
		time.Sleep(time.Second/2)
	}()
}
