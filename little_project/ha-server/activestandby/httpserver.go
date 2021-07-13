package activestandby

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"k8s.io/klog/v2"
)

func StartHelperHTTPServer()  {
	router := mux.NewRouter()
	router.HandleFunc("/readyz", electionHandler).Methods("GET")
	port := 10000 + rand.Intn(20)
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	klog.Infof("helper server as: %s", addr)
	server := &http.Server{
		Addr: addr,
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		klog.Errorf("err: ", err)
	}
}

func StartMainHTTPServer() {
	router := mux.NewRouter()
	router.HandleFunc("/hello", helloHandler).Methods("GET")
	port := 10010 + rand.Intn(20)
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	klog.Infof("main server as: %s", addr)
	server := &http.Server{
		Addr: addr,
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		klog.Errorf("err: ", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	hostnameOverride, err := os.Hostname()
	if err != nil {
		hostnameOverride = "node"
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("I am " + hostnameOverride)); err != nil {
		klog.Errorf("failed to write http response, err: %v", err)
	}
	return
}

func electionHandler(w http.ResponseWriter, r *http.Request) {
	// if checker is nill not get
	checker := Config.Checker
	if checker == nil {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("server is ready with no leaderelection")); err != nil {
			klog.Errorf("failed to write http response, err: %v", err)
		}
		return
	}
	if checker.Check(r) != nil {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("server is not ready")); err != nil {
			klog.Errorf("failed to write http response, err: %v", err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("server is ready")); err != nil {
			klog.Errorf("failed to write http response, err: %v", err)
		}
	}
}