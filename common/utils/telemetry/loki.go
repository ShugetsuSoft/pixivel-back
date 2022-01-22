package telemetry

import (
	loki "github.com/ShugetsuSoft/loki-client-go"
)

var LoglokiIns = &loki.LokiClient{}
var LogType = ""

func RunLoki(uri string, logType string) chan error {
	LoglokiIns = loki.NewLokiClient(uri)
	LogType = logType
	return LoglokiIns.RunPush()
}
