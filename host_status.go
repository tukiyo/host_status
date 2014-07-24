package main

import (
    "fmt"
    "net/http"
    "os/exec"
)

func main() {
    http.HandleFunc("/docker", docker)
    http.HandleFunc("/virsh", virsh)
    http.Handle("/", http.FileServer(http.Dir("/")))
    http.ListenAndServe(":9876", nil)
}

func docker(w http.ResponseWriter, r *http.Request) {
    header(w)
    print_hr(w, run("w"))
    print_hr(w, run("free -m"))
    print_hr(w, run("(df -hm | head -n 1) && (df -hm | sort -nr -k2 | head -n 5)"))
    print_hr(w, run("docker ps -a"))
    print_hr(w, run("for i in $(docker ps -a | grep -v '^CONTAINER' | awk '{print $1}'); do echo -n $i': '; docker inspect --format '{{ .NetworkSettings.IPAddress }}' $i;done"))
    print_hr(w, run("docker images"))
    print_hr(w, run("ps aux --sort user,-rss | grep -v '^root'"))
}

func virsh(w http.ResponseWriter, r *http.Request) {
    header(w)
    print_hr(w, run("virsh list --all"))
    print_hr(w, run("virsh net-list"))
    print_hr(w, run("virsh iface-list"))
    print_hr(w, run("virsh nodeinfo"))
    print_hr(w, run("virsh help"))
}

func run(strCommand string) string {
    cmd := exec.Command("/usr/bin/bash", "-c", strCommand)
    d, _ := cmd.Output()
    return string(d)
}

func print_hr(w http.ResponseWriter, stdout string) {
    fmt.Fprintf(w, "<hr>")
    fmt.Fprintf(w, string(stdout))
}

func header(w http.ResponseWriter) {
    w.Header().Add("Content-type", "text/html charset=utf-8")
    fmt.Fprintf(w, "<style>a{padding-right:1em;}</style>")
    fmt.Fprintf(w, "<a href='/docker'>[docker]</a>")
    fmt.Fprintf(w, "<a href='/virsh'>[virsh]</a>")
    fmt.Fprintf(w, "<a href='/var/lib/libvirt/images/'>virsh/images</a>")
    fmt.Fprintf(w, "<a href='/etc/'>/etc/</a>")
    fmt.Fprintf(w, "<a href='/var/log/'>/var/log/</a>")
    fmt.Fprintf(w, "<a href='/var/lib'>/var/lib/</a>")
    fmt.Fprintf(w, "<pre style='font-size:small;'>")
}
