package remote

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pyroscope-io/pyroscope/pkg/storage"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Putter interface {
	Put(ctx context.Context, pi *storage.PutInput) error
}

type RemoteWriter struct {
	log *logrus.Logger
}

func NewRemoteWriter(logger *logrus.Logger) *RemoteWriter {
	return &RemoteWriter{
		log: logger,
	}
}

func (r *RemoteWriter) Put(_ context.Context, pi *storage.PutInput) error {
	logrus.Debugf("Putting in remote storage")
	return nil
}

// TODO(eh-am): rename to something else more clear
type StorageOrchestrator struct {
	log     *logrus.Logger
	putters []Putter
}

// TODO(eh-am): move this somewhere else, probably on its own package
// TODO(eh-am): should we use an array of putters?
// maybe this orchestrator should create the putters?
func NewStorageOrchestrator(log *logrus.Logger, putters ...Putter) *StorageOrchestrator {
	// TODO(eh-am): wrap the logger with a svc?
	return &StorageOrchestrator{
		log:     log,
		putters: putters,
	}
}

// Put puts data into the underlying Putters
// By default it uses a parallel strategy
func (so *StorageOrchestrator) Put(ctx context.Context, pi *storage.PutInput) error {
	logrus.Debugf("Putting in storage orchestrator")

	// TODO(eh-am): maybe we should have different strategies here?
	// like writing in parallel, or writing sequentially

	// Write to all putters in parallel
	g, ctx := errgroup.WithContext(ctx)
	for _, p := range so.putters {
		// https://golang.org/doc/faq#closures_and_goroutines
		p := p

		g.Go(func() error {
			return p.Put(ctx, pi)
		})
	}

	return g.Wait()
}

// TODO(eh-am): move to somwhere else, likely next to where
// requestToPutInput is done
func putInputToRequest(pi *storage.PutInput) error {
	// TODO (eh-am): URL
	address := "http://localhost:4040"
	//	req, err := http.NewRequest(http.MethodPost, url)
	//	if err != nil {
	//		return err
	//	}

	u, err := url.Parse(address)
	if err != nil {
		return err
	}

	// Query Params
	// This is basically the same as ingestParamsFromRequest from ingest.go
	q := u.Query()

	// TODO (eh-am): just copied this from ingestParamsFromRequest to refer more easily
	//	[ ] pi.Format = q.Get("format")
	//	[ ] pi.ContentType = r.Header.Get("Content-Type")
	// [x] if qt := q.Get("from"); qt != "" {
	// 	pi.StartTime = attime.Parse(qt)
	// [x] if qt := q.Get("until"); qt != "" {
	// 	pi.EndTime = attime.Parse(qt)
	// [x] if sr := q.Get("sampleRate"); sr != "" {
	// 	sampleRate, err := strconv.Atoi(sr)
	// [X] if sn := q.Get("spyName"); sn != "" {
	// 	pi.SpyName = sn
	// [X] if u := q.Get("units"); u != "" {
	// 	pi.Units = metadata.Units(u)
	// [X] if at := q.Get("aggregationType"); at != "" {
	// 	 pi.AggregationType = metadata.AggregationType(at)

	q.Set("aggregationType", pi.AggregationType.String())
	q.Set("units", pi.Units.String())
	q.Set("spyName", pi.SpyName)
	// TODO(eh-am): since this is a hotpath check how slow using fmt.sprintf is
	q.Set("sampleRate", fmt.Sprint(pi.SampleRate))
	q.Set("until", pi.EndTime.String())
	q.Set("from", pi.StartTime.String())

	// TODO(eh-am): how about format and content type

	// TODO (eh-am): write body
	//  body := &bytes.Buffer{}
	//	writer := multipart.NewWriter(body)

}
