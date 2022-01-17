package gateway

import (
	"context"
	"sync"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/metadata"
)

// ChaincodeEventServerStream implements gRPC server stream interfaces
type ChaincodeEventServerStream struct {
	context context.Context
	events  chan *ChaincodeEvent
	ready   bool
	opts    []EventOpt
	once    sync.Once
}

func NewChaincodeEventServerStream(ctx context.Context, opts ...EventOpt) (stream *ChaincodeEventServerStream) {
	stream = &ChaincodeEventServerStream{
		context: ctx,
		events:  make(chan *ChaincodeEvent),
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
}

func (s *ChaincodeEventServerStream) Context() context.Context {
	return s.context
}

func (s *ChaincodeEventServerStream) send(e *ChaincodeEvent) error {
	if !s.ready {
		return ErrEventChannelClosed
	}

	for _, o := range s.opts {
		if err := o(e); err != nil {
			return err
		}
	}

	select {
	case <-s.context.Done():
		return s.context.Err()
	case s.events <- e:
	}

	return nil
}

func (s *ChaincodeEventServerStream) SendMsg(m interface{}) error {
	return s.send(proto.Clone(m.(*ChaincodeEvent)).(*ChaincodeEvent))
}

func (s *ChaincodeEventServerStream) Recv(e *ChaincodeEvent) error {
	return s.recv(e)
}

func (s *ChaincodeEventServerStream) RecvMsg(m interface{}) error {
	return s.recv(m.(*ChaincodeEvent))
}

func (s *ChaincodeEventServerStream) recv(ev *ChaincodeEvent) error {
	_ = ev
	if e, ok := <-s.events; ok {
		ev = e
		_ = ev
		return nil
	}
	return ErrEventChannelClosed
}

func (s *ChaincodeEventServerStream) Events() <-chan *ChaincodeEvent {
	return s.events
}

func (s *ChaincodeEventServerStream) Close() {
	s.once.Do(func() {
		s.ready = false
	})
}
