package metrics

import "github.com/prometheus/client_golang/prometheus"

var AttemptsCreatedTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "attempts_created_total",
		Help: "Total attempts created",
	},
)

var AttemptsSubmittedTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "attempts_submitted_total",
		Help: "Total attempts submitted",
	},
)

var ResultsGeneratedTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "results_generated_total",
		Help: "Total results generated",
	},
)

var FilesUploadedTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "files_uploaded_total",
		Help: "Total files uploaded",
	},
)

var RabbitMQJobsProcessedTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "rabbitmq_jobs_processed_total",
		Help: "Total rabbitmq jobs processed",
	},
)

var RabbitMQJobsFailedTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "rabbitmq_jobs_failed_total",
		Help: "Total rabbitmq jobs failed",
	},
)
