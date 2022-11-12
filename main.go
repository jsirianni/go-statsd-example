package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/cactus/go-statsd-client/v5/statsd"
)

const (
	envHost             = "STATSD_HOST"
	envIncludeResources = "OTEL_INCLUDE_RESOURCES"
	envPodName          = "OTEL_RESOURCE_ATTRIBUTES_POD_NAME"
	envNamespace        = "OTEL_RESOURCE_ATTRIBUTES_NAMESPACE"
	envAppName          = "STATSD_APP_NAME"

	port        = 8125
	flushInt    = time.Second * 5
	resInterval = time.Minute * 2
)

func main() {
	host := os.Getenv(envHost)
	if host == "" {
		fmt.Printf("failed to read %s\n", envHost)
		os.Exit(1)
	}

	appName := os.Getenv(envAppName)
	if appName == "" {
		fmt.Printf("failed to read %s\n", envAppName)
		os.Exit(1)
	}

	config := &statsd.ClientConfig{
		Address:       fmt.Sprintf("%s:%d", host, port),
		Prefix:        appName,
		FlushInterval: flushInt,

		// How often to re-resolve the hostname. Ignored
		// if IP address is used. Should be a high value
		// because the hostname should be a clusterIP service, which
		// will not change frequently.
		ResInterval: resInterval,
	}

	client, err := statsd.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Close()

	includeResources := false
	if x := os.Getenv(envIncludeResources); x != "" {
		b, err := strconv.ParseBool(x)
		if err != nil {
			panic(err.Error())
		}
		includeResources = b
	}

	tags := []statsd.Tag{}

	if includeResources {
		tags = append(tags, statsd.Tag{"k8s.pod.name", os.Getenv(envPodName)})
		tags = append(tags, statsd.Tag{"k8s.namespace.name", os.Getenv(envNamespace)})
	}

	for {
		if err := client.Gauge("request.latency", value(), 1.0, tags...); err != nil {
			fmt.Println(err)
		}

		if err := client.Gauge("request.body.bytes", value(), 1.0, tags...); err != nil {
			fmt.Println(err)
		}

		if err := client.Gauge("request.headers.bytes", value(), 1.0, tags...); err != nil {
			fmt.Println(err)
		}

		if err := client.Gauge("connections.count", value(), 1.0, tags...); err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Millisecond * 10)
	}
}

// random int64 between 0 and 20
func value() int64 {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(20)
	return int64(n)
}
