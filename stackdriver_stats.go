// Copyright 2019, OpenCensus Authors
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
//

package stackdriver

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// The following variables are measures are recorded by ClientHandler:
var (
	timeSeriesCountDistribution = view.Distribution(1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536)
)

var (
	timeSeriesPerRequest     = stats.Int64("contrib.go.opencensus.io/exporter/stackdriver/timeseries_per_request", "Number of timeseries sent per request.", stats.UnitDimensionless)
)

// Predefined views may be registered to collect data for the above measures.
var (
	timeSeriesRequestCountView = &view.View{
		Measure:     timeSeriesPerRequest,
		Name:        "contrib.go.opencensus.io/exporter/stackdriver/create_time_series_request_count",
		Description: "Number of create time series request sent",
		TagKeys:     []tag.Key{},
		Aggregation: view.Count(),
	}

	timeSeriesPerRequestView = &view.View{
		Measure:     timeSeriesPerRequest,
		Name:        "contrib.go.opencensus.io/exporter/stackdriver/time_series_per_request",
		Description: "Number of time series packed in create time series request",
		TagKeys:     []tag.Key{},
		Aggregation: timeSeriesCountDistribution,
	}
)

var stackdriverViews = []*view.View{
	timeSeriesRequestCountView,
	timeSeriesPerRequestView,
}


