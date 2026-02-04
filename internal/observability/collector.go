package observability

import (
	"context"
	"encoding/json"
	"strconv"

	"zabbix-dna/internal/api"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Collector struct {
	Client *api.ZabbixClient
	Meter  metric.Meter
	Tracer trace.Tracer
}

func NewCollector(client *api.ZabbixClient, meter metric.Meter, tracer trace.Tracer) *Collector {
	return &Collector{
		Client: client,
		Meter:  meter,
		Tracer: tracer,
	}
}

type ZabbixItem struct {
	ItemID    string `json:"itemid"`
	HostID    string `json:"hostid"`
	Name      string `json:"name"`
	Key       string `json:"key_"`
	LastValue string `json:"lastvalue"`
	ValueType string `json:"value_type"`
	Units     string `json:"units"`
	Hosts     []struct {
		Name string `json:"name"`
		Host string `json:"host"`
	} `json:"hosts"`
}

type ZabbixProblem struct {
	EventID   string `json:"eventid"`
	ObjectID  string `json:"objectid"`
	Name      string `json:"name"`
	Severity  string `json:"severity"`
	Clock     string `json:"clock"`
	R_EventID string `json:"r_eventid"`
	Tags      []struct {
		Tag   string `json:"tag"`
		Value string `json:"value"`
	} `json:"tags"`
}

func (c *Collector) CollectMetrics(ctx context.Context) error {
	// Fetch items with numeric values
	params := map[string]interface{}{
		"output":      []string{"itemid", "hostid", "name", "key_", "lastvalue", "value_type", "units"},
		"selectHosts": []string{"name", "host"},
		"filter": map[string]interface{}{
			"value_type": []string{"0", "3"}, // numeric float and numeric unsigned
		},
		"monitored": true,
		"limit":     100, // Limit for demo/safety
	}

	result, err := c.Client.Call("item.get", params)
	if err != nil {
		return err
	}

	var items []ZabbixItem
	if err := json.Unmarshal(result, &items); err != nil {
		return err
	}

	for _, item := range items {
		if item.LastValue == "" {
			continue
		}

		val, err := strconv.ParseFloat(item.LastValue, 64)
		if err != nil {
			continue
		}

		hostName := item.HostID
		if len(item.Hosts) > 0 {
			hostName = item.Hosts[0].Name
		}

		// Create a gauge for each item key
		gauge, err := c.Meter.Float64ObservableGauge(
			item.Key,
			metric.WithDescription(item.Name),
			metric.WithUnit(item.Units),
		)
		if err != nil {
			continue
		}

		_, err = c.Meter.RegisterCallback(func(_ context.Context, obs metric.Observer) error {
			obs.ObserveFloat64(gauge, val, metric.WithAttributes(
				attribute.String("zabbix.itemid", item.ItemID),
				attribute.String("zabbix.hostid", item.HostID),
				attribute.String("zabbix.host", hostName),
			))
			return nil
		}, gauge)
	}

	return nil
}

func (c *Collector) CollectTraces(ctx context.Context) error {
	params := map[string]interface{}{
		"output":     []string{"eventid", "objectid", "name", "severity", "clock", "r_eventid"},
		"selectTags": "extend",
		"recent":     true,
		"sortfield":  []string{"eventid"},
		"sortorder":  "DESC",
		"limit":      50,
	}

	result, err := c.Client.Call("problem.get", params)
	if err != nil {
		return err
	}

	var problems []ZabbixProblem
	if err := json.Unmarshal(result, &problems); err != nil {
		return err
	}

	for _, p := range problems {
		attrs := []attribute.KeyValue{
			attribute.String("zabbix.eventid", p.EventID),
			attribute.String("zabbix.objectid", p.ObjectID),
			attribute.String("zabbix.severity", p.Severity),
			attribute.String("zabbix.clock", p.Clock),
		}

		if p.R_EventID != "" && p.R_EventID != "0" {
			attrs = append(attrs, attribute.String("zabbix.r_eventid", p.R_EventID))
			attrs = append(attrs, attribute.Bool("zabbix.resolved", true))
		} else {
			attrs = append(attrs, attribute.Bool("zabbix.resolved", false))
		}

		// Add tags as attributes
		for _, tag := range p.Tags {
			attrs = append(attrs, attribute.String("zabbix.tag."+tag.Tag, tag.Value))
		}

		// Zabbix 7.0+ might have symptom problems
		if p.Severity == "0" {
			attrs = append(attrs, attribute.String("zabbix.problem_type", "not_classified"))
		}

		_, span := c.Tracer.Start(ctx, p.Name, trace.WithAttributes(attrs...))
		span.End()
	}

	return nil
}


