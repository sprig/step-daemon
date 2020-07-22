package main

import (
	"fmt"
	gio "io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/colinrgodsey/serial"
	"github.com/colinrgodsey/step-daemon/config"
	"github.com/colinrgodsey/step-daemon/io"
	"github.com/colinrgodsey/step-daemon/pipeline"
)

type argMap map[string]string

func (a argMap) has(arg string) bool {
	_, ok := a[arg]
	return ok
}

func handler(head io.Conn, size int, h func(head, tail io.Conn)) (tail io.Conn) {
	head = head.Flip()
	tail = io.NewConn(size, size)

	go h(head, tail)

	return
}

func stepdPipeline(c io.Conn) io.Conn {
	c = handler(c, 8, pipeline.SourceHandler)
	c = handler(c, 8, pipeline.DeltaHandler)
	c = handler(c, 8, pipeline.PhysicsHandler)
	c = handler(c, 8, pipeline.StepHandler)
	c = handler(c, 8, pipeline.DeviceHandler)
	return c
}

func bailArgs() {
	out := []string{
		"Required args (arg=value ...):",
		"",
		"config - Path to config file",
		"device - Path to serial device",
		"baud - Bad rate for serial device",
		"",
	}
	for _, s := range out {
		fmt.Println(s)
	}
	os.Exit(1)
}

func main() {
	args := loadArgs()

	if !args.has("config") || !args.has("device") || !args.has("baud") {
		bailArgs()
	}

	conf, err := config.LoadConfig(args["config"])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c := io.NewConn(8, 8)
	c.Write(conf)
	go io.LinePipe(os.Stdin, os.Stdout, c.Flip())
	c = stepdPipeline(c)
	tailSink(c, args)
}

func closeOnExit(h gio.Closer) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("info: closing device serial")
		h.Close()
		os.Exit(0)
	}()
}

func loadArgs() argMap {
	out := make(argMap)
	for _, arg := range os.Args {
		spl := strings.SplitN(arg, "=", 2)
		if len(spl) == 1 {
			out[spl[0]] = ""
		} else {
			out[spl[0]] = spl[1]
		}
	}
	return out
}

//TODO: need the auto restart loop here
func tailSink(c io.Conn, args argMap) {
	var tail gio.ReadWriteCloser
	var err error

	baud, err := strconv.Atoi(args["baud"])
	if err != nil {
		fmt.Println("Failed to parse baud")
		bailArgs()
	}

	cfg := &serial.Config{Name: args["device"], Baud: baud}
	tail, err = serial.OpenPort(cfg)
	closeOnExit(tail)

	if err != nil {
		panic(fmt.Sprintf("Failed to open dest file %v: %v", os.Args[1], err))
	}

	io.LinePipe(tail, tail, c)
}