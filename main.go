package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"

	"github.com/muhammadaser/cash_deposit/accounts"
	cf "github.com/muhammadaser/cash_deposit/config"

	"github.com/go-kit/kit/log"
	"github.com/go-pg/pg"
	"github.com/oklog/oklog/pkg/group"
)

func main() {
	// initial configuration
	err := cf.InitConfig(cf.Env)
	if err != nil {
		panic(err)
	}

	// Define our flags. Your service probably won't need to bind listeners for
	// *all* supported transports, or support both Zipkin and LightStep, and so
	// on, but we do it here for demonstration purposes.
	fs := flag.NewFlagSet("cms", flag.ExitOnError)
	var (
		// debugAddr = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
		httpAddr = fs.String("http-addr", cf.Config.Address, "HTTP listen address")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(cf.LogOutput)
		logger = log.With(logger, "ts", log.DefaultTimestamp)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// postgresql connection
	pgDB := pg.Connect(&pg.Options{
		Addr:     cf.Config.Pg.Addr + ":" + cf.Config.Pg.Port, // 172.31.1.178
		User:     cf.Config.Pg.Username,                       // program
		Password: cf.Config.Pg.Password,                       // jatis123
		Database: cf.Config.Pg.Database,
	})
	var n int
	_, err = pgDB.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		panic(err)
	}

	httpLogger := log.With(logger, "component", "http")

	logger = log.With(logger, "p", "accounts")
	var (
		AccountsService     = accounts.New(logger, pgDB)
		AccountsEndpoints   = accounts.NewEndpoint(AccountsService, logger)
		AccountsHTTPHandler = accounts.NewHTTPHandler(AccountsEndpoints, httpLogger)
	)

	logger = log.With(logger, "p", "deposits")
	var (
		DepositsService     = accounts.New(logger, pgDB)
		DepositsEndpoints   = accounts.NewEndpoint(DepositsService, logger)
		DepositsHTTPHandler = accounts.NewHTTPHandler(DepositsEndpoints, httpLogger)
	)

	mux := http.NewServeMux()

	mux.Handle("/cash-deposit/v1/accounts", AccountsHTTPHandler)
	mux.Handle("/cash-deposit/v1/deposits", DepositsHTTPHandler)

	httpHandler := accessControl(mux)

	var g group.Group
	// {
	// 	// The debug listener mounts the http.DefaultServeMux, and serves up
	// 	// stuff like the Prometheus metrics route, the Go debug and profiling
	// 	// routes, and so on.
	// 	debugListener, err := net.Listen("tcp", *debugAddr)
	// 	if err != nil {
	// 		logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
	// 		os.Exit(1)
	// 	}
	// 	g.Add(func() error {
	// 		logger.Log("transport", "debug/HTTP", "addr", *debugAddr)
	// 		return http.Serve(debugListener, http.DefaultServeMux)
	// 	}, func(error) {
	// 		debugListener.Close()
	// 	})
	// }
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", *httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST,PUT,DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
