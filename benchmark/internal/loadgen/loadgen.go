package loadgen

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/pyroscope-io/pyroscope/benchmark/internal/config"
	"github.com/pyroscope-io/pyroscope/pkg/agent/upstream"
	"github.com/pyroscope-io/pyroscope/pkg/agent/upstream/remote"
	"github.com/pyroscope-io/pyroscope/pkg/structs/transporttrie"
	"github.com/sirupsen/logrus"
)

// how many retries to check the pyroscope server is up
const MaxReadinessRetries = 10

var CapacityExceeded = errors.New("capacity exceeded")

type Fixtures [][]*transporttrie.Trie

type LoadGen struct {
	Config    *config.LoadGen
	Rand      *rand.Rand
	SymbolBuf []byte

	runProgressMetric prometheus.Gauge
	uploadErrors      prometheus.Counter
	successfulUploads prometheus.Counter
	pusher            GatewayPusher
}

type GatewayPusher interface {
	Push() error
}
type NoopGatewayPusher struct{}

func (NoopGatewayPusher) Push() error {
	return nil
}

func Cli(cfg *config.LoadGen) error {
	r := rand.New(rand.NewSource(int64(cfg.RandSeed)))

	var pusher GatewayPusher
	if cfg.PushgatewayAddress == "" {
		logrus.Debug("no pushgateway configured")
		pusher = NoopGatewayPusher{}
	} else {
		logrus.Debug("will push metrics to ", cfg.PushgatewayAddress)
		pusher = push.New(cfg.PushgatewayAddress, cfg.ServerAddress).Gatherer(prometheus.DefaultGatherer)
	}

	l := &LoadGen{
		Config:    cfg,
		Rand:      r,
		SymbolBuf: make([]byte, cfg.ProfileSymbolLength),

		runProgressMetric: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "pyroscope",
			Subsystem: "benchmark",
			Name:      "progress",
			Help:      "",
		}),
		uploadErrors: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "pyroscope",
			Subsystem: "benchmark",
			Name:      "upload_errors",
			Help:      "",
		}),
		successfulUploads: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "pyroscope",
			Subsystem: "benchmark",
			Name:      "successful_uploads",
			Help:      "",
		}),
		pusher: pusher,
	}

	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "pyroscope",
			Subsystem: "benchmark",
			Name:      "requests_total",
			Help:      "",
		},
		func() float64 { return float64(cfg.Apps * cfg.Requests * cfg.Clients) },
	)

	if cfg.CapacityBenchmark {
		minClients := cfg.Clients
		//logrus.Infof("minClients: %d", minClients)
		for l.Run(cfg) != CapacityExceeded {
			cfg.Clients *= 2
		}
		maxClients := cfg.Clients
		//logrus.Infof("maxClients: %d", maxClients)
		logrus.Infof("client capacity: %d", sort.Search(maxClients-minClients, func(n int) bool {
			cfg.Clients = minClients + n
			capacityExceeded := 0
			for i := 0; i < cfg.Trials; i++ {
				if err := l.Run(cfg); err == CapacityExceeded {
					capacityExceeded++
				} else if err != nil {
					panic(err)
				}
			}
			return capacityExceeded*2 > cfg.Trials
		}))
		return nil
	}

	return l.Run(cfg)
}

func (l *LoadGen) Run(cfg *config.LoadGen) error {
	logrus.Info("checking server is available...")
	err := waitUntilEndpointReady(cfg.ServerAddress)
	if err != nil {
		return err
	}

	logrus.Info("generating fixtures")
	fixtures := l.generateFixtures()
	logrus.Debug("done generating fixtures.")

	logrus.Info("starting sending requests")
	logrus.Infof("cfg %+v\n", cfg)
	wg := sync.WaitGroup{}
	wg.Add(l.Config.Apps * l.Config.Clients)
	appNameBuf := make([]byte, 25)

	// Memory-inefficient but obviously correct histogram.
	var latencyHistogram []time.Duration
	var errorCounter *uint64
	if cfg.CapacityBenchmark {
		latencyHistogram = make([]time.Duration, cfg.Apps*cfg.Clients*cfg.Requests)
		errorCounter = new(uint64)
	}

	for i := 0; i < l.Config.Apps; i++ {
		// generate a random app name
		l.Rand.Read(appNameBuf)
		appName := hex.EncodeToString(appNameBuf)
		for j := 0; j < l.Config.Clients; j++ {
			offset := i*cfg.Clients*cfg.Requests + j*cfg.Requests
			go l.startClientThread(appName, &wg, fixtures[i], latencyHistogram[offset:offset+cfg.Requests], errorCounter)
		}
	}
	wg.Wait()

	logrus.Debug("done sending requests")

	// Latency quantiles and error rate. If latency's too high or error rate
	// is non-zero, we've exceeded this instance's capacity.
	if cfg.CapacityBenchmark {
		sort.Slice(latencyHistogram, func(i, j int) bool {
			return latencyHistogram[i] < latencyHistogram[j]
		})
		logrus.Infof("50th percentile latency: %v", latencyHistogram[len(latencyHistogram)*50/100])
		logrus.Infof("90th percentile latency: %v", latencyHistogram[len(latencyHistogram)*90/100])
		logrus.Infof("99th percentile latency: %v", latencyHistogram[len(latencyHistogram)*99/100])
		logrus.Infof("errors: %v", atomic.LoadUint64(errorCounter))

		if latencyHistogram[len(latencyHistogram)*99/100] > 7500*time.Millisecond { // 7.5 seconds but keep the math in integers
			return CapacityExceeded
		}
		if atomic.LoadUint64(errorCounter) > 0 {
			return CapacityExceeded
		}
	}

	return nil
}

func (l *LoadGen) generateFixtures() Fixtures {
	var f Fixtures

	for i := 0; i < l.Config.Apps; i++ {
		f = append(f, []*transporttrie.Trie{})

		randomGen := rand.New(rand.NewSource(int64(l.Config.RandSeed + i)))
		p := l.generateProfile(randomGen)
		for j := 0; j < l.Config.Fixtures; j++ {
			f[i] = append(f[i], p)
		}
	}

	return f
}

func (l *LoadGen) startClientThread(appName string, wg *sync.WaitGroup, appFixtures []*transporttrie.Trie, latencyHistogram []time.Duration, errorCounter *uint64) {
	rc := remote.RemoteConfig{
		UpstreamThreads:        1,
		UpstreamAddress:        l.Config.ServerAddress,
		UpstreamRequestTimeout: 10 * time.Second,
	}

	r, err := remote.New(rc, logrus.New())
	if err != nil {
		panic(err)
	}

	requestsCount := l.Config.Requests

	threadStartTime := time.Now().Truncate(10 * time.Second)
	threadStartTime = threadStartTime.Add(time.Duration(-1*requestsCount) * (10 * time.Second))

	st := threadStartTime

	for i := 0; i < requestsCount; i++ {
		t := appFixtures[i%len(appFixtures)]

		st = st.Add(10 * time.Second)
		et := st.Add(10 * time.Second)
		t0 := time.Now()
		err := r.UploadSync(&upstream.UploadJob{
			Name:            appName + "{}",
			StartTime:       st,
			EndTime:         et,
			SpyName:         "gospy",
			SampleRate:      100,
			Units:           "samples",
			AggregationType: "sum",
			Trie:            t,
		})
		if latencyHistogram != nil {
			latencyHistogram[i] = time.Since(t0)
		}
		if err != nil {
			l.uploadErrors.Add(1)
			if errorCounter != nil {
				atomic.AddUint64(errorCounter, 1)
			}
			time.Sleep(time.Second)
		} else {
			l.successfulUploads.Add(1)
		}

		err = l.pusher.Push()
		if err != nil {
			logrus.Error(err)
		}
	}

	wg.Done()
}

func (l *LoadGen) generateProfile(randomGen *rand.Rand) *transporttrie.Trie {
	t := transporttrie.New()

	for w := 0; w < l.Config.ProfileWidth; w++ {
		symbol := []byte("root")
		for d := 0; d < 2+l.Rand.Intn(l.Config.ProfileDepth); d++ {
			randomGen.Read(l.SymbolBuf)
			symbol = append(symbol, byte(';'))
			symbol = append(symbol, []byte(hex.EncodeToString(l.SymbolBuf))...)
			if l.Rand.Intn(100) <= 20 {
				t.Insert(symbol, uint64(l.Rand.Intn(100)), true)
			}
		}

		t.Insert(symbol, uint64(l.Rand.Intn(100)), true)
	}
	return t
}

// TODO(eh-am) exponential backoff and whatnot
func waitUntilEndpointReady(url string) error {
	client := http.Client{Timeout: 10 * time.Second}
	retries := 0

	for {
		_, err := client.Get(url)

		// all good?
		if err == nil {
			return nil
		}
		if retries >= MaxReadinessRetries {
			break
		}

		time.Sleep(time.Second)
		retries++
	}

	return fmt.Errorf("maximum retries exceeded ('%d') waiting for server ('%s') to respond", retries, url)
}
