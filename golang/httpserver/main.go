package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"httpserver/metrics"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

//辛苦老师辛苦老师
func main() {

	flag.Set("v", "4")
	glog.V(2).Info("starting http server")
	metrics.Register()
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healThz)
	mux.HandleFunc("/", rootHandler)
	mux.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	timer := metrics.NewTimer()
	defer timer.ObserverTotal()
	rand.Seed(start.UnixNano())
	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
	glog.Infof("hello go.")
	io.WriteString(w, "hello go: "+strconv.Itoa(int(time.Since(start).Milliseconds()))+" ms")
}
func healThz(w http.ResponseWriter, r *http.Request) {
	//request header 添加到response header中
	for k, v := range r.Header {
		for _, v2 := range v {
			w.Header().Set(k, v2)
		}
	}
	w.Header().Add("VERSION", os.Getenv("VERSION"))
	//设置返回状态码
	w.WriteHeader(200)
	//客戶端ip和状态码输出
	glog.Infof("client ip: %s status code:%d", r.Host, 200)

	io.WriteString(w, "200")
}
