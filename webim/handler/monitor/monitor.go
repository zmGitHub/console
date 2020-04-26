package monitor

import "github.com/prometheus/client_golang/prometheus"

const metricsNamespace = "custm_chat"

var (
	RequestTotalCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "http_request_total",
			Help:      "HTTP requests processed.",
		},
		[]string{"code", "method", "host", "url"},
	)
	RequestDurationInSec = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: metricsNamespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latencies in seconds.",
		},
		[]string{"method", "host", "url"},
	)
	ResponseSizeInBytes = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: metricsNamespace,
			Name:      "http_response_size_bytes",
			Help:      "HTTP response bytes.",
		},
	)

	EnterprisesRegisterCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricsNamespace,
		Subsystem: "enterprise",
		Name:      "enterprises_register_count",
		Help:      "Number of enterprise registered.",
	}, []string{"type"})

	MessagesSentCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricsNamespace,
		Subsystem: "message",
		Name:      "messages_sent_count",
		Help:      "Number of messages sent.",
	}, []string{"type"})

	ConversationsCreationCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricsNamespace,
		Subsystem: "conversation",
		Name:      "conversations_create_count",
		Help:      "Number of conversation created.",
	}, []string{"type"})

	VisitorsComesIn = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricsNamespace,
		Subsystem: "visitor",
		Name:      "visitors_count",
		Help:      "Number of conversation created.",
	}, []string{"type"})
)

func RegisterMetrics() {
	prometheus.MustRegister(RequestTotalCount)
	prometheus.MustRegister(RequestDurationInSec)
	prometheus.MustRegister(ResponseSizeInBytes)
	prometheus.MustRegister(EnterprisesRegisterCount)
	prometheus.MustRegister(MessagesSentCount)
	prometheus.MustRegister(ConversationsCreationCount)
	prometheus.MustRegister(VisitorsComesIn)
}
