package main

import (
	"github.com/high-creek-software/weatherlight-workbench/internal"
	"github.com/lmittmann/tint"
	"golang.org/x/exp/slog"
	"os"
	"time"
)
import _ "net/http/pprof"

func main() {
	slog.SetDefault(slog.New(tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	}.NewHandler(os.Stderr)))

	weatherlightWorkbench := internal.NewWeatherlightWorkbench()
	weatherlightWorkbench.Start()
}
