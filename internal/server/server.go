package server

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/9seconds/ddoctor/internal/checkers"
	"github.com/9seconds/ddoctor/internal/config"
	"github.com/9seconds/ddoctor/internal/presenter"
)

func Serve(conf *config.Config, ctx context.Context, checks []checkers.Checker) {
	results := make(map[string]*checkers.CheckResult)
	channel := make(chan *checkers.CheckResult, len(checks))
	defer close(channel)

	srv := &http.Server{
		Addr: net.JoinHostPort(conf.Host, strconv.Itoa(int(conf.Port))),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "ddoctor")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		ok := true
		toSerialize := make([]*checkers.CheckResult, 0, len(results))
		for _, v := range results {
			ok = ok && v.Ok
			toSerialize = append(toSerialize, v)
		}

		if ok {
			w.WriteHeader(conf.OkStatus)
		} else {
			w.WriteHeader(conf.NokStatus)
		}

		resp, err := presenter.Serialize(toSerialize, false)
		if err == nil {
			w.Write(resp)
		}
	})

	go func() {
		<-ctx.Done()
		srv.Shutdown(nil)
	}()
	go srv.ListenAndServe()

	doWork := func() {
		for _, v := range checks {
			go v.Run(ctx, channel)
		}
		for i := 0; i < len(checks); i++ {
			res := <-channel
			results[res.Producer] = res
		}
	}

	go doWork()
	timer := time.After(conf.Periodicity.Duration)
	for {
		select {
		case <-timer:
			timer = time.After(conf.Periodicity.Duration)
			go doWork()
		case <-ctx.Done():
			return
		}
	}
}
