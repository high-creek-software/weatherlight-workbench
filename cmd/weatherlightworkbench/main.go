package main

import (
	"github.com/high-creek-software/weatherlight-workbench/internal"
)
import _ "net/http/pprof"

func main() {
	weatherlightWorkbench := internal.NewWeatherlightWorkbench()
	weatherlightWorkbench.Start()
}
