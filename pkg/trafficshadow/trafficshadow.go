package trafficshadow

import (
	"context"

	"github.com/pyroscope-io/pyroscope/pkg/parser"
	"github.com/sirupsen/logrus"
)

type Storage interface {
	Put(context.Context, *parser.PutInput) error
}

type TrafficShadower struct {
	targetAddr string
	log        *logrus.Logger
	s          Storage
}

func New(s Storage, log *logrus.Logger, targetAddr string) *TrafficShadower {
	return &TrafficShadower{
		targetAddr: targetAddr,
		log:        log,
		s:          s,
	}
}

// Put shadows traffic to the target server (via HTTP)
// while relaying to the original storage
func (ts *TrafficShadower) Put(ctx context.Context, pi *parser.PutInput) error {
	// TODO: do in in parallel

	err := ts.s.Put(ctx, pi)
	if err != nil {
		return err
	}

	err = ts.relay(ctx, pi)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TrafficShadower) relay(ctx context.Context, pi *parser.PutInput) error {
	logrus.Debugf("Relaying to %s", ts.targetAddr)

	// TODO look at ingestParamsFromRequest in pkg/server/ingest.go
	// and do the reverse effect (putInput to request)
	return nil
}
