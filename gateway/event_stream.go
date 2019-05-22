package gateway

import (
	"context"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/protos/peer"
	"google.golang.org/grpc/metadata"
)

type ChaincodeEventServerStream struct {
	context context.Context
	events  chan *peer.ChaincodeEvent
	ready   bool
	opts    []EventOpt
	once    sync.Once
}

func NewChaincodeEventServerStream(ctx context.Context, opts ...EventOpt) (stream *ChaincodeEventServerStream) {
	stream = &ChaincodeEventServerStream{
		context: ctx,
		events:  make(chan *peer.ChaincodeEvent),
		ready:   true,
		opts:    opts,
	}

	go func() {
		<-ctx.Done()
		stream.Close()
	}()

	return stream
}

func (*ChaincodeEventServerStream) SetHeader(metadata.MD) error {
	return nil
}

func (*ChaincodeEventServerStream) SendHeader(metadata.MD) error {
	return nil
}

func (*ChaincodeEventServerStream) SetTrailer(metadata.MD) {
	return
}

func (s *ChaincodeEventServerStream) Context() context.Context {
	return s.context
}

func (s *ChaincodeEventServerStream) SendMsg(m interface{}) (err error) {
	if !s.ready {
		return ErrEventChannelClosed
	}

	e := proto.Clone(m.(*peer.ChaincodeEvent)).(*peer.ChaincodeEvent)
	for _, o := range s.opts {
		if err = o(e); err != nil {
			return err
		}
	}

	s.events <- e
	return nil
}

func (s *ChaincodeEventServerStream) Recv(e *peer.ChaincodeEvent) error {
	return s.RecvMsg(e)
}

func (s *ChaincodeEventServerStream) RecvMsg(m interface{}) error {
	m, ok := <-s.events
	if ok {
		return nil
	}
	return ErrEventChannelClosed
}

func (s *ChaincodeEventServerStream) Events() <-chan *peer.ChaincodeEvent {
	return s.events
}

func (s *ChaincodeEventServerStream) Close() {
	s.once.Do(func() {
		close(s.events)
		s.ready = false
	})
}
