// Copyright 2017, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command stackdriver is an example program that collects data for
// video size. Collected data is exported to
// Stackdriver Monitoring.
package main

import (
	"context"
	"go.opencensus.io/metric/metricdata"
	"log"
	"os"
	"sync"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/monitoredresource"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// Create measures. The program will record measures for the size of
// processed videos and the nubmer of videos marked as spam.
var videoSize = stats.Int64("my.org/measure/video_size", "size of processed videos", stats.UnitBytes)

func main() {
	ctx := context.Background()

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// Collected view data will be reported to Stackdriver Monitoring API
	// via the Stackdriver exporter.
	//
	// In order to use the Stackdriver exporter, enable Stackdriver Monitoring API
	// at https://console.cloud.google.com/apis/dashboard.
	//
	// Once API is enabled, you can use Google Application Default Credentials
	// to setup the authorization.
	// See https://developers.google.com/identity/protocols/application-default-credentials
	// for more details.
	se, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID:         os.Getenv("GOOGLE_CLOUD_PROJECT"), // Google Cloud Console project ID for stackdriver.
		MonitoredResource: monitoredresource.Autodetect(),
		BundleDelayThreshold: 2 * time.Second,
		BundleCountThreshold: 10000,
		TraceSpansBufferMaxBytes: 20 * 1024 * 1024,
	})
	se.StartMetricsExporter()
	defer se.StopMetricsExporter()

	if err != nil {
		log.Fatal(err)
	}

	// Create view to see the processed video size cumulatively.
	// Subscribe will allow view data to be exported.
	// Once no longer need, you can unsubscribe from the view.
	if err := view.Register(&view.View{
		Name:        "my.org/views/video_size_cum",
		Description: "processed video size over time",
		Measure:     videoSize,
		Aggregation: view.Distribution(1<<16, 1<<32),
	}); err != nil {
		log.Fatalf("Cannot subscribe to the view: %v", err)
	}

	trace.RegisterExporter(se)
	log.Println("Start sending")
	wg := sync.WaitGroup{}
	for i := 0; i< 500000; i++ {
		wg.Add(1)
		go func() {
			processVideo(ctx, i)
			wg.Done()
		}()
		time.Sleep(80 * time.Microsecond)
	}
	wg.Wait()
	log.Println("Done sending")

	// Wait for a duration longer than reporting duration to ensure the stats
	// library reports the collected data.
	log.Println("Wait longer than the reporting duration...")
	time.Sleep(10 * time.Second)
}

func getSpanCtxAttachment(ctx context.Context) metricdata.Attachments {
	attachments := map[string]interface{}{}
	span := trace.FromContext(ctx)
	if span == nil {
		return attachments
	}
	spanCtx := span.SpanContext()
	if spanCtx.IsSampled() {
		attachments[metricdata.AttachmentKeySpanContext] = spanCtx
	}
	return attachments
}

func processVideo(ctx context.Context, attempt int) {
	ctx, span := trace.StartSpan(ctx, "example.com/ProcessVideo")
	//span.AddAttributes(trace.StringAttribute("attempt", "01234567890123456789"))
	//span.AddAttributes(trace.StringAttribute("attempt1", "01234567890123456789"))
	//span.AddAttributes(trace.StringAttribute("attempt2", "01234567890123456789"))
	//span.AddAttributes(trace.StringAttribute("attempt3", "01234567890123456789"))
	//span.AddAttributes(trace.StringAttribute("attempt4", "01234567890123456789"))
	defer span.End()
	// Do some processing and record stats.
	stats.RecordWithOptions(ctx, stats.WithMeasurements(videoSize.M(25648)), stats.WithAttachments(getSpanCtxAttachment(ctx)))
}
