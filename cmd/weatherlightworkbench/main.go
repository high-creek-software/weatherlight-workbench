package main

import (
	"flag"
	"os"
	"runtime/pprof"
	"time"

	"github.com/high-creek-software/weatherlight-workbench/internal"
	"github.com/lmittmann/tint"
	"golang.org/x/exp/slog"

	_ "net/http/pprof"
)

var profileFlag = flag.Bool("p", false, "Sets to use the cpu profiler.")

func main() {
	flag.Parse()

	slog.SetDefault(slog.New(tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	}.NewHandler(os.Stderr)))

	if *profileFlag {
		out, err := os.Create("weatherlight-cpu.pprof")
		if err != nil {
			slog.Error("error creating cpu output file", "error", err)
			os.Exit(2)
		}
		pprof.StartCPUProfile(out)
		defer pprof.StopCPUProfile()
	}

	weatherlightWorkbench := internal.NewWeatherlightWorkbench()
	weatherlightWorkbench.Start()
}
