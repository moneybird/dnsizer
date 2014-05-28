// dnsizer is a small DNS server which can be used to resolv a development domain to a single IP adress
// This can be used, for example, which a virtual machine where your code lives.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
)

var (
	addr   = flag.String("addr", ":53", "Address to listen on")
	ip     = flag.String("ip", "127.0.0.1", "IP to resolve DNS requests to")
	domain = flag.String("domain", "dev", "Domain to resolve for")
)

func handleResponse(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	rr, err := dns.NewRR(fmt.Sprintf("%s IN A %s", r.Question[0].Name, *ip))

	if err != nil {
		log.Fatalf("Error setting up reply %s", err)
	}

	m.Answer = append(m.Answer, rr)

	w.WriteMsg(m)
}

func main() {

	flag.Parse()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	dns.HandleFunc(fmt.Sprintf("%s.", *domain), handleResponse)

	go func() {
		srv := &dns.Server{Addr: *addr, Net: "udp"}
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to setup UDP listener %s\n", err)
		}
	}()

	go func() {
		srv := &dns.Server{Addr: *addr, Net: "tcp"}
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to setup TCP listener %s\n", err)
		}
	}()

	for {
		select {
		case s := <-sig:
			log.Fatalf("Signal (%s) received, stopping\n", s)
		}
	}
}
