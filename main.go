package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
    "html"
    "strings"
    "os"
    "net"
    "strconv"
    "time"
)

const default_listen_address = ":8000"

func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "1.1.1.1:80")
    if err != nil {
        log.Fatalln(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}

func GetPwd() string {
    pwd, err := os.Getwd()
    if err != nil {
        log.Fatalln(err)
    }
    return pwd
}

func IsDockerenvThere() bool {
    _, err := os.Stat("/.dockerenv")
    return err == nil
}

func GetListenAddress(def string) string {
    if addrenv := os.Getenv("LISTEN_ADDRESS"); addrenv != ""{
        return addrenv
    }
    return def
}

func GetResponse(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "HTTP Request info:\n")
    fmt.Fprintf(w, "  method: %q\n", r.Method)
    fmt.Fprintf(w, "  proto: %q\n", r.Proto)
    fmt.Fprintf(w, "  URL path: %q\n", html.EscapeString(r.URL.Path))
    fmt.Fprintf(w, "  remote Addr: %q\n", r.RemoteAddr)
    fmt.Fprint(w, "  headers:\n")
    fmt.Fprintf(w, "    \"Host\": %q\n", html.EscapeString(r.Host))
    for k, v := range r.Header {
        fmt.Fprintf(w,"    %q: [%q]\n", html.EscapeString(k), html.EscapeString(strings.Join(v,", ")))
    }

    fmt.Fprint(w, "\nEnvironment info:\n")
    fmt.Fprint(w, "  environment variables:\n")
    for _, e := range os.Environ() {
        fmt.Fprintf(w, "    %q\n", e)
    }
    fmt.Fprintf(w, "  local outbound IP: %q\n", GetOutboundIP())
    fmt.Fprintf(w, "  current working dir: %q\n", GetPwd())
    fmt.Fprintf(w, "  is /.dockerenv there: %q\n", strconv.FormatBool(IsDockerenvThere()))
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", GetResponse)
    addr := GetListenAddress(default_listen_address)
    log.Println("Listening on " + addr)

    srv := &http.Server{
        Handler:      router,
        Addr:         addr,
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    log.Fatal(srv.ListenAndServe())
}

// vim: tabstop=4 expandtab
