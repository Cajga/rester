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
    "regexp"
    "strconv"
    "time"
    "sort"
)

const default_listen_address = ":8000"
const version = "0.1.2"

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
    fmt.Fprintf(w, "  UTC time of the request: %q\n", time.Now().UTC().Format(time.RFC3339Nano))
    fmt.Fprint(w, "  headers:\n")
    headers := make([]string, 0, len(r.Header))
    for header := range r.Header {
        headers = append(headers,header)
    }
    headers = append(headers,"Host")
    sort.Strings(headers)
    for _, header := range headers {
        if header == "Host" {
            fmt.Fprintf(w, "    \"Host\": %q\n", html.EscapeString(r.Host))
        } else {
            fmt.Fprintf(w,"    %q: [%q]\n", html.EscapeString(header), html.EscapeString(strings.Join(r.Header[header],", ")))
        }
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
    router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
        match, _ := regexp.MatchString("/.*", r.URL.Path)
        // TODO handle error from regex
        return match
    }).HandlerFunc(GetResponse)
    addr := GetListenAddress(default_listen_address)
    log.Println("Listening on " + addr + " with version " + version)

    srv := &http.Server{
        Handler:      router,
        Addr:         addr,
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    log.Fatal(srv.ListenAndServe())
}

// vim: tabstop=4 expandtab
