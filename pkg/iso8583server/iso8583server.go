package iso8583server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/jattento/go-iso8583/pkg/mti"
)

type Server struct {
	config       configuration
	handlerRules []handlerRule

	// Connections contains all active connections and may contain non active.
	// Non active connections are not deleted instantly, only if the configured limit is hitted.
	// Clients should feel free to delete non active connections.
	Connections *connectionDB
}

type configuration struct {
	// DeactivatedConnectionCapacity is the amount of deactivated commection persisted
	// in the Connections DB,if the db is at full capacity the connection with the oldest
	// read/write is deleted.
	DeactivatedConnectionCapacity int
	Listener                      net.Listener
	LogInfo                       LogFunc
	LogErr                        LogFunc
	NetRead                       NetReadFunc
	ReadMTI                       ReadMTIFunc
	UnknownHandler                HandlerFunc
	LogOn                         *LogOn
	ConnIdGenerator               ConnectionIdGenerator

	// Timeout is the max time the server waits for new messages before closing the connection
	Timeout time.Duration
}

type LogOn struct {
	ErrorSettingConnectionReadDeadline bool
	ErrorReadMTI                       bool
	ErrorUndefinedHandler              bool
	ErrorReadConnection                bool
	ErrorAcceptIncomingConnection      bool

	ServingConnection bool
}

type handlerRule struct {
	rule    Rule
	handler HandlerFunc
}

type HandlerFunc func(r Response, message []byte)

type ReadMTIFunc = func([]byte) (mti.MTI, error)

type LogFunc = func(v ...interface{})

type NetReadFunc = func(io.Reader) (chan []byte, chan error)

type Rule = func(mti.MTI) bool

type ConnectionIdGenerator = func() (string, error)

func New(options ...Option) (*Server, error) {
	config := configuration{}
	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}

	if config.Listener == nil {
		l, err := net.Listen("tcp", ":8080")
		if err != nil {
			return nil, err
		}

		config.Listener = l
	}

	if config.Timeout == 0 {
		config.Timeout = 3 * time.Minute
	}

	if config.NetRead == nil {
		config.NetRead = DefaultReader{
			SizeChunkLength: 2,
			SizeReader:      func(b []byte) int { return int(binary.BigEndian.Uint16(b)) },
		}.Read
	}

	if config.ReadMTI == nil {
		config.ReadMTI = DefaultReadMTI
	}

	if config.UnknownHandler == nil {
		// nop handler
		config.UnknownHandler = func(r Response, message []byte) { return }
	}

	if config.LogInfo == nil {
		config.LogInfo = log.New(os.Stdout, "iso8583server_info ", log.LstdFlags).Print
	}

	if config.LogErr == nil {
		config.LogErr = log.New(os.Stderr, "iso8583server_error ", log.LstdFlags).Print
	}

	if config.LogOn == nil {
		config.LogOn = &DefaultLogOnConfiguration
	}

	if config.ConnIdGenerator == nil {
		config.ConnIdGenerator = defaultIdGenerator
	}

	return &Server{config: config, Connections: &connectionDB{m: map[string]*Connection{}}}, nil
}

func (server *Server) handleConnection(c *Connection) {
	extendReadDeadline := func() {
		if err := c.SetReadDeadline(time.Now().Add(server.config.Timeout)); err != nil {
			if server.config.LogOn.ErrorSettingConnectionReadDeadline {
				server.config.LogErr(fmt.Sprintf("failed setting read deadline: %v", handleErrorNet(err, c)))
			}
			return
		}
	}

	if server.config.LogOn.ServingConnection {
		server.config.LogInfo(fmt.Sprintf("serving %s\n", c.RemoteAddr().String()))
	}

	extendReadDeadline()

	messageChan, errChan := server.config.NetRead(c)

	for {
		select {
		case message := <-messageChan:
			func() {
				extendReadDeadline()
				c.LastRead = time.Now()

				mtiValue, err := server.config.ReadMTI(message)
				if err != nil {
					if server.config.LogOn.ErrorReadMTI {
						server.config.LogErr("cant read MTI from message: ", err)
					}
					go startHandler(server.config.UnknownHandler, c, message)
					return
				}

				for _, hr := range server.handlerRules {
					if hr.rule(mtiValue) {
						go startHandler(hr.handler, c, message)
						return
					}
				}

				if server.config.LogOn.ErrorUndefinedHandler {
					server.config.LogErr("no handler defined for MTI ", mtiValue, " sending to unknown manager")
				}

				go startHandler(server.config.UnknownHandler, c, message)
			}()
		case err := <-errChan:
			c.Active = false
			server.Connections.deleteOldest(server.config.DeactivatedConnectionCapacity)

			if server.config.LogOn.ErrorReadConnection {
				server.config.LogErr(c.RemoteAddr().String(), ": ", handleErrorNet(err, c))
			}

			return
		}
	}
}

func (server *Server) Close() error {
	return server.config.Listener.Close()
}

func (server *Server) AddBottomPriorityHandler(handler HandlerFunc, rules ...Rule) {
	if handler == nil {
		return
	}

	for _, rule := range rules {
		if rule != nil {
			server.handlerRules = append(server.handlerRules, handlerRule{rule: rule, handler: handler})
		}
	}
}

func (server *Server) AddTopPriorityHandler(handler HandlerFunc, rules ...Rule) {
	if handler == nil {
		return
	}

	for _, rule := range rules {
		if rule != nil {
			server.handlerRules = append([]handlerRule{{rule: rule, handler: handler}}, server.handlerRules...)
		}
	}
}

func (server *Server) Start() {
	for {
		c, err := server.config.Listener.Accept()
		if err != nil {
			if server.config.LogOn.ErrorAcceptIncomingConnection {
				server.config.LogErr(fmt.Sprintf("cant accept incoming connection %v", err))
			}

			continue
		}

		connection, err := server.newConnection(c)
		if err != nil {
			if server.config.LogOn.ErrorAcceptIncomingConnection {
				server.config.LogErr(fmt.Sprintf("cant accept incoming connection, error generating id: %v", err))
			}
			_ = c.Close()
		}
		go server.handleConnection(connection)
	}
}

func DefaultReadMTI(b []byte) (mti.MTI, error) {
	if len(b) < 4 {
		return "", errors.New("message to short to read MTI")
	}

	_, err := strconv.Atoi(string(b[:4]))
	if err != nil {
		return "", fmt.Errorf("first 4 characters arent numbers: %v", err)
	}

	return mti.MTI(b[:4]), nil
}

type DefaultReader struct {
	// SizeChunkLength are the amount of bytes that must be read at message start to determine length
	SizeChunkLength int

	// SizeReader gets as input SizeChunkLength bytes and must return the represented length.
	SizeReader func([]byte) int
}

func (reader DefaultReader) Read(net io.Reader) (chan []byte, chan error) {
	messageChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		var sizeBuffer []byte
		var messageBuffer []byte

		upToFillSizeBuffer := func() int { return reader.SizeChunkLength - len(sizeBuffer) }
		upToFillMessageBuffer := func() int { return reader.SizeReader(sizeBuffer) - len(messageBuffer) }

		for {
			if upToFillSizeBuffer() > 0 {
				sizeChunk := make([]byte, upToFillSizeBuffer())

				n, err := net.Read(sizeChunk)
				if err != nil {
					errChan <- err
					return
				}

				if n > 0 {
					sizeBuffer = append(sizeBuffer, sizeChunk[:n]...)
				}

				continue
			}

			if upToFillMessageBuffer() > 0 {
				messageChunk := make([]byte, upToFillMessageBuffer())

				n, err := net.Read(messageChunk)
				if err != nil {
					errChan <- err
					return
				}

				if n > 0 {
					messageBuffer = append(messageBuffer, messageChunk[:n]...)
				}

				continue
			}

			message := make([]byte, len(messageBuffer))
			copy(message, messageBuffer)

			messageChan <- message

			// Restart all buffer
			sizeBuffer = nil
			messageBuffer = nil
		}
	}()

	return messageChan, errChan
}

func handleErrorNet(err error, c net.Conn) error {
	if clErr := c.Close(); clErr != nil {
		return fmt.Errorf("%v and failed closing connection: %v", err, clErr)
	}
	return err
}

func startHandler(h HandlerFunc, c *Connection, message []byte) {
	r := Response{
		buff:       bytes.NewBuffer(nil),
		connection: c,
	}

	h(r, message)

}

type Response struct {
	buff       *bytes.Buffer
	connection *Connection
}

func (r Response) Write(b []byte) (int, error) {
	return r.buff.Write(b)
}

func (r Response) Close() error {
	_, err := r.connection.Write(r.buff.Bytes())
	if err == nil {
		r.connection.LastWrite = time.Now()
	}
	return err
}
