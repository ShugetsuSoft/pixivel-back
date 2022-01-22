package telemetry

import (
	"github.com/ShugetsuSoft/loki-client-go/lib"
	"log"
	"time"
)

type Label map[string]string

func Log(label Label, s string) {
	log.Printf("[%T] %s", label, s)
	label["type"] = LogType
	err := LoglokiIns.WriteLog(lib.Label(label), s, time.Now())
	if err != nil {
		log.Println(err)
	}
}
