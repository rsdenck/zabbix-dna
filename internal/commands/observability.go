package commands

import (
	"context"
	"fmt"
	"time"

	"zabbix-dna/internal/api"
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
				fmt.Println("Error: OTLP endpoint not specified")
				return
			}

			client := api.NewClient(cfg.Zabbix.URL, cfg.Zabbix.Token, cfg.Zabbix.Timeout)
			if cfg.Zabbix.Token == "" && cfg.Zabbix.User != "" {
				err := client.Login(cfg.Zabbix.User, cfg.Zabbix.Password)
				if err != nil {
					fmt.Printf("Login failed: %v\n", err)
					return
				}
			}

			fmt.Printf("Starting metrics export to %s...\n", endpoint)

			engine := observability.NewOTLPEngine(endpoint, "zabbix-dna")
			ctx := context.Background()
			mp, err := engine.InitMetrics(ctx)
			if err != nil {
				fmt.Printf("Failed to initialize metrics: %v\n", err)
				return
			}
			defer mp.Shutdown(ctx)

			meter := otel.Meter("zabbix-dna-collector")
			tracer := otel.Tracer("zabbix-dna-collector")
			collector := observability.NewCollector(client, meter, tracer)

			duration, _ := time.ParseDuration(interval)
			ticker := time.NewTicker(duration)
			defer ticker.Stop()

			fmt.Printf("Metrics engine active (interval: %s). Press Ctrl+C to stop.\n", interval)

			for {
				select {
				case <-ticker.C:
					if err := collector.CollectMetrics(ctx); err != nil {
						fmt.Printf("Collection error: %v\n", err)
					} else {
						fmt.Println("Metrics collected and exported.")
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
				fmt.Println("Error: OTLP endpoint not specified")
				return
			}

			client := api.NewClient(cfg.Zabbix.URL, cfg.Zabbix.Token, cfg.Zabbix.Timeout)
			if cfg.Zabbix.Token == "" && cfg.Zabbix.User != "" {
				err := client.Login(cfg.Zabbix.User, cfg.Zabbix.Password)
				if err != nil {
					fmt.Printf("Login failed: %v\n", err)
					return
				}
			}

			fmt.Printf("Starting traces export to %s...\n", endpoint)

			engine := observability.NewOTLPEngine(endpoint, "zabbix-dna")
			ctx := context.Background()
			tp, err := engine.InitTraces(ctx)
			if err != nil {
				fmt.Printf("Failed to initialize traces: %v\n", err)
				return
			}
			defer tp.Shutdown(ctx)

			meter := otel.Meter("zabbix-dna-collector")
			tracer := otel.Tracer("zabbix-dna-collector")
			collector := observability.NewCollector(client, meter, tracer)

			duration, _ := time.ParseDuration(interval)
			ticker := time.NewTicker(duration)
			defer ticker.Stop()

			fmt.Printf("Traces engine active (interval: %s). Press Ctrl+C to stop.\n", interval)

			for {
				select {
				case <-ticker.C:
					if err := collector.CollectTraces(ctx); err != nil {
						fmt.Printf("Collection error: %v\n", err)
					} else {
						fmt.Println("Traces collected and exported.")
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


