package restapi

import (
	"context"
	"fmt"
	"gitlab.com/devskiller-tasks/messaging-app-golang/fastsmsing"
	"gitlab.com/devskiller-tasks/messaging-app-golang/smsproxy"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type SmsApp struct {
	serveMux        *http.ServeMux
	server          *http.Server
	smsProxyService smsproxy.SmsProxy
}

func NewServer(port int) SmsApp {
	serveMux := *http.NewServeMux()
	server := http.Server{Handler: &serveMux, Addr: fmt.Sprintf(":%d", port)}
	smsProxyService := smsproxy.ProdSmsProxy(fastsmsing.NewInMemoryClient(), smsproxy.MinimumInBatchOption(2))
	return SmsApp{serveMux: &serveMux, server: &server, smsProxyService: smsProxyService}
}

func (app *SmsApp) BindEndpoints() {
	app.serveMux.HandleFunc("/", app.routingHandler)
}

func (app *SmsApp) Run() error {
	idleChan := make(chan struct{})

	go func(app *SmsApp, idleChan chan struct{}) {
		signChan := make(chan os.Signal, 1)
		signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
		sig := <-signChan
		log.Println("shutdown:", sig)

		app.Stop(5 * time.Second)

		// Actual shutdown trigger.
		close(idleChan)
	}(app, idleChan)

	app.smsProxyService.Start()
	return app.server.ListenAndServe()
}

func (app *SmsApp) Stop(t time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()
	app.smsProxyService.Stop()
	err := app.server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}
}

func (app *SmsApp) routingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.RequestURI() == "/sms" || r.URL.RequestURI() == "/sms/" {
		if r.Method == http.MethodPost {
			sendSmsHandler(app.smsProxyService).ServeHTTP(w, r)
			return
		}
	}

	if strings.HasPrefix(r.URL.RequestURI(), "/sms/") {
		rightUriPart := strings.TrimPrefix(r.URL.RequestURI(), "/sms/")
		if !strings.Contains(rightUriPart, "/") && r.Method == http.MethodGet && len(rightUriPart) > 0 {
			getSmsStatusHandler(app.smsProxyService).ServeHTTP(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte(`{"message": "not found"}`))
	if err != nil {
		log.Println(err.Error())
	}
}
