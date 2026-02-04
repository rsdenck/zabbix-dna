package commands

import (
	"context"
	"fmt"
	"time"

	"zabbix-dna/internal/config"
	"zabbix-dna/internal/observability"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

func newExporterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exporter",
		Short: "Zabbix OTLP Exporter (Grafana Tempo, Beyla, Prometheus)",
		Long: `The Zabbix-DNA Exporter sends Zabbix metrics and traces via OTLP.
Compatible with Grafana Tempo (traces), Grafana Beyla, and Prometheus (via OTLP receiver).`,
	}

	cmd.AddCommand(newMetricsCmd())
	cmd.AddCommand(newTracesCmd())

	return cmd
}

func newMetricsCmd() *cobra.Command {
	var interval string
	var endpoint string

	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Export Zabbix metrics to OTLP",
		Run: func(cmd *cobra.Command, args []string) {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, _ := config.LoadConfig(cfgPath)

			if endpoint == "" && cfg != nil {
				endpoint = cfg.OTLP.Endpoint
			}

			if endpoint == "" {
				handleError(fmt.Errorf("OTLP endpoint not specified"))
				return
			}

			client, err := getZabbixClient(cmd)
			handleError(err)

			headers := []string{"Service", "Endpoint", "Interval", "Status"}
			rows := [][]string{
				{"Metrics Exporter", endpoint, interval, "Starting..."},
			}
			outputResult(cmd, nil, headers, rows)

			engine := observability.NewOTLPEngine(endpoint, "zabbix-dna")
			ctx := context.Background()
			mp, err := engine.InitMetrics(ctx)
			if err != nil {
				handleError(fmt.Errorf("failed to initialize metrics: %v", err))
				return
			}
			defer mp.Shutdown(ctx)

			meter := otel.Meter("zabbix-dna-collector")
			tracer := otel.Tracer("zabbix-dna-collector")
			collector := observability.NewCollector(client, meter, tracer)

			duration, _ := time.ParseDuration(interval)
			ticker := time.NewTicker(duration)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := collector.CollectMetrics(ctx); err != nil {
						handleError(fmt.Errorf("collection error: %v", err))
					}
				case <-ctx.Done():
					return
				}
			}
		},
	}

	cmd.Flags().StringVarP(&interval, "interval", "i", "60s", "Export interval")
	cmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "OTLP endpoint")

	return cmd
}

func newTracesCmd() *cobra.Command {
	var interval string
	var endpoint string

	cmd := &cobra.Command{
		Use:   "traces",
		Short: "Export Zabbix events as OTLP traces",
		Run: func(cmd *cobra.Command, args []string) {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, _ := config.LoadConfig(cfgPath)

			if endpoint == "" && cfg != nil {
				endpoint = cfg.OTLP.Endpoint
			}

			if endpoint == "" {
				handleError(fmt.Errorf("OTLP endpoint not specified"))
				return
			}

			client, err := getZabbixClient(cmd)
			handleError(err)

			headers := []string{"Service", "Endpoint", "Interval", "Status"}
			rows := [][]string{
				{"Traces Exporter", endpoint, interval, "Starting..."},
			}
			outputResult(cmd, nil, headers, rows)

			engine := observability.NewOTLPEngine(endpoint, "zabbix-dna")
			ctx := context.Background()
			tp, err := engine.InitTraces(ctx)
			if err != nil {
				handleError(fmt.Errorf("failed to initialize traces: %v", err))
				return
			}
			defer tp.Shutdown(ctx)

			meter := otel.Meter("zabbix-dna-collector")
			tracer := otel.Tracer("zabbix-dna-collector")
			collector := observability.NewCollector(client, meter, tracer)

			duration, _ := time.ParseDuration(interval)
			ticker := time.NewTicker(duration)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := collector.CollectTraces(ctx); err != nil {
						handleError(fmt.Errorf("collection error: %v", err))
					}
				case <-ctx.Done():
					return
				}
			}
		},
	}

	cmd.Flags().StringVarP(&interval, "interval", "i", "60s", "Export interval")
	cmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "OTLP endpoint")

	return cmd
}
