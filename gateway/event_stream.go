package gateway

import (
	"context"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer"
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

	select {
	case <-s.context.Done():
		return s.context.Err()
	case s.events <- e:
	}

	return nil
}

func (s *ChaincodeEventServerStream) Recv(e *peer.ChaincodeEvent) error {
	return s.recv(e)
}

func (s *ChaincodeEventServerStream) RecvMsg(m interface{}) error {
	return s.recv(m.(*peer.ChaincodeEvent))
}

func (s *ChaincodeEventServerStream) recv(ev *peer.ChaincodeEvent) error {
	_ = ev
	if e, ok := <-s.events; ok {
		ev = e
		_ = ev
		return nil
	}
	return ErrEventChannelClosed
}

func (s *ChaincodeEventServerStream) Events() <-chan *peer.ChaincodeEvent {
	return s.events
}

func (s *ChaincodeEventServerStream) Close() {
	s.once.Do(func() {
		s.ready = false
	})
}
