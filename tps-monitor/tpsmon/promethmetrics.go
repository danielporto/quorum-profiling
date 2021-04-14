package tpsmon

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type PrometheusMetricsService struct {
	tpsGauge    prometheus.Gauge
	txnsGauge   prometheus.Gauge
	blocksGauge prometheus.Gauge
	blockNumGauge prometheus.Gauge
	blockTxnGauge prometheus.Gauge
	blockTimeGauge prometheus.Gauge
	blockGasGauge prometheus.Gauge
	port        int
}

func NewPrometheusMetricsService(port int) *PrometheusMetricsService {
	ps := &PrometheusMetricsService{
		tpsGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "Quorum",
			Subsystem: "TransactionProcessing",
			Name:      "TPS",
			Help:      "Transactions processed per second",
		}),
		blocksGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "Quorum",
			Subsystem: "TransactionProcessing",
			Name:      "total_blocks",
			Help:      "total blocks processed",
		}),
		txnsGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "Quorum",
			Subsystem: "TransactionProcessing",
			Name:      "total_transactions",
			Help:      "total transactions processed",
		}),

		blockNumGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "Quorum",
			Subsystem: "TransactionProcessing",
			Name:      "block_ids",
			Help:      "total of block blocks generated (id)",
		}),

		blockTxnGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "Quorum",
			Subsystem: "TransactionProcessing",
			Name:      "block_transactions",
			Help:      "total transactions in the block",
		}),
		blockTimeGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "Quorum",
			Subsystem: "TransactionProcessing",
			Name:      "block_time",
			Help:      "When the block was generated",
		}),
		blockGasGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "Quorum",
			Subsystem: "TransactionProcessing",
			Name:      "block_gas",
			Help:      "Used gas to process this block",
		}),


		port: port,
	}
	return ps
}

func (ps *PrometheusMetricsService) Start() {
	prometheus.MustRegister(ps.tpsGauge)
	prometheus.MustRegister(ps.txnsGauge)
	prometheus.MustRegister(ps.blocksGauge)
	prometheus.MustRegister(ps.blockNumGauge)
	prometheus.MustRegister(ps.blockTxnGauge)
	prometheus.MustRegister(ps.blockTimeGauge)
	prometheus.MustRegister(ps.blockGasGauge)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", ps.port), nil))
}

func (ps *PrometheusMetricsService) publishMetrics(ref time.Time, tps uint64, txnCnt uint64, blkCnt uint64) {
	ps.tpsGauge.SetToCurrentTime()
	ps.tpsGauge.Set(float64(tps))
	ps.txnsGauge.SetToCurrentTime()
	ps.txnsGauge.Set(float64(txnCnt))
	ps.blocksGauge.SetToCurrentTime()
	ps.blocksGauge.Set(float64(blkCnt))

	log.Debug("published metrics to prometheus")
}

func (ps *PrometheusMetricsService) publishBlockMetrics(btime, bnum uint64, btx int, gas uint64) {
	ps.blockTimeGauge.SetToCurrentTime()
	ps.blockTimeGauge.Set(float64(btime))
	ps.blockNumGauge.SetToCurrentTime()
	ps.blockNumGauge.Set(float64(bnum))
	ps.blockTxnGauge.SetToCurrentTime()
	ps.blockTxnGauge.Set(float64(btx))
	ps.blockGasGauge.SetToCurrentTime()
	ps.blockGasGauge.Set(float64(gas))


	log.Debug("published block metrics to prometheus")
}