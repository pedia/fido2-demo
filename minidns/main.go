package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"

	"github.com/miekg/dns"
)

var parent = flag.String("parent", "114.114.114.114", "specify parent(upstream) DNS")

type Stub struct {
	hosts map[string]string
	c     *dns.Client
}

// dig A @127.0.0.1 w3c.com.
// nslookup w3c.com 127.0.0.1

func (stub *Stub) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if len(r.Question) == 0 {
		return
	}

	q := r.Question[0]

	if q.Qtype == dns.TypeA {
		val, ok := stub.hosts[q.Name]
		if ok {
			fmt.Printf("> %d %s %s\n", r.Question[0].Qtype, q.Name, val)

			a := new(dns.Msg)
			a.SetReply(r)
			a.Authoritative = false
			a.Answer = make([]dns.RR, 1)
			a.Answer[0] = &dns.A{
				Hdr: dns.RR_Header{
					Name:   r.Question[0].Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    3600, // 1hr
				},
				A: net.ParseIP(val).To4(),
			}
			w.WriteMsg(a)
			return
		}
	}

	// exchange from parent
	if stub.c == nil {
		stub.c = new(dns.Client)
	}

	in, _, err := stub.c.Exchange(r.Copy(), fmt.Sprintf("%s:53", *parent))
	if err == nil && len(in.Answer) > 0 {
		fmt.Printf("  %s\n", in.Answer[0].String())

		w.WriteMsg(in)
	}
}

func read_hosts(fn string) map[string]string {
	file, err := os.Open(fn)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed %v", fn, err))
	}

	defer file.Close()

	res := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(t, "#") || strings.HasPrefix(t, ";") {
			continue
		}

		arr := strings.FieldsFunc(t, func(r rune) bool {
			return r == '\t' || r == ' '
		})
		if len(arr) == 2 {
			res[fmt.Sprintf("%s.", arr[1])] = arr[0]
		}
	}

	return res
}

func main() {
	flag.Parse()

	stub := new(Stub)

	var fn string
	if len(flag.Args()) > 0 {
		fn = flag.Args()[0]
	} else {
		if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
			fn = "/etc/hosts"
		} else if runtime.GOOS == "windows" {
			fn = "c:\\windows\\system32\\drivers\\etc\\hosts"
		}
	}
	if fn != "" {
		stub.hosts = read_hosts(fn)
	}

	for k, v := range stub.hosts {
		fmt.Printf("%s %s\n", k, v)
	}

	dns.ListenAndServe(":53", "udp", stub)
}
