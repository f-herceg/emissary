package agent

import (
	"github.com/datawire/dlib/dlog"
	"github.com/emissary-ingress/emissary/v3/pkg/api/agent"
	envoyMetrics "github.com/emissary-ingress/emissary/v3/pkg/api/envoy/service/metrics/v3"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	counterType    = io_prometheus_client.MetricType_COUNTER
	acceptedMetric = &io_prometheus_client.MetricFamily{
		Name: StrToPointer("cluster.apple_prod_443.upstream_rq_total"),
		Type: &counterType,
		Metric: []*io_prometheus_client.Metric{
			{
				Counter: &io_prometheus_client.Counter{
					Value: Float64ToPointer(42),
				},
				TimestampMs: Int64ToPointer(time.Now().Unix() * 1000),
			},
		},
	}
	ignoredMetric = &io_prometheus_client.MetricFamily{
		Name: StrToPointer("cluster.apple_prod_443.metric_to_ignore"),
		Type: &counterType,
		Metric: []*io_prometheus_client.Metric{
			{
				Counter: &io_prometheus_client.Counter{
					Value: Float64ToPointer(42),
				},
				TimestampMs: Int64ToPointer(time.Now().Unix() * 1000),
			},
		},
	}
)

func agentMetricsSetupTest() (*MockClient, *Agent) {
	clientMock := &MockClient{}

	stubbedAgent := &Agent{
		metricsBackoffUntil: time.Time{},
		comm: &RPCComm{
			client: clientMock,
		},
	}

	return clientMock, stubbedAgent
}

func TestMetricsRelayHandler(t *testing.T) {

	t.Run("will relay the metrics", func(t *testing.T) {
		//given
		clientMock, stubbedAgent := agentMetricsSetupTest()
		ctx := dlog.NewTestContext(t, true)

		//when
		stubbedAgent.MetricsRelayHandler(ctx, &envoyMetrics.StreamMetricsMessage{
			Identifier:   nil,
			EnvoyMetrics: []*io_prometheus_client.MetricFamily{ignoredMetric, acceptedMetric},
		})

		//then
		assert.Equal(t, []*agent.StreamMetricsMessage{{
			EnvoyMetrics: []*io_prometheus_client.MetricFamily{acceptedMetric},
		}}, clientMock.SentMetrics)
	})
	t.Run("will not relay the metrics since it is in cool down period.", func(t *testing.T) {
		//given
		clientMock, stubbedAgent := agentMetricsSetupTest()
		ctx := dlog.NewTestContext(t, true)
		stubbedAgent.metricsBackoffUntil = time.Now().Add(defaultMinReportPeriod)

		//when
		stubbedAgent.MetricsRelayHandler(ctx, &envoyMetrics.StreamMetricsMessage{
			Identifier:   nil,
			EnvoyMetrics: []*io_prometheus_client.MetricFamily{acceptedMetric},
		})

		//then
		assert.Equal(t, 0, len(clientMock.SentMetrics))
	})
}
