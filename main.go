package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/miekg/dns"
	"log"
	"strings"
)

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	requestId := uuid.New().String()
	log.Printf("%s source %s ", requestId, w.RemoteAddr())

	switch r.Opcode {
	case dns.OpcodeQuery:
		for _, q := range m.Question {
			switch q.Qtype {
			case dns.TypeTXT:
				log.Printf("%s source %s txt %s ", requestId, w.RemoteAddr(), q.Name)
				// local dns
				localDnsIP := strings.Split(w.RemoteAddr().String(), ":")[0]
				localDnsRR, err := dns.NewRR(fmt.Sprintf("%s 0 IN TXT localdns=%s", q.Name, localDnsIP))
				if err != nil {
					log.Printf("%s source %s txt %s your.local.dns[%s]", requestId, w.RemoteAddr(), q.Name, err.Error())
					return
				}
				m.Answer = append(m.Answer, localDnsRR)

				// subnet handler
				subnet := ""
				for _, v := range r.Extra {
					opt, ok := v.(*dns.OPT)
					if ok {
						for _, o := range opt.Option {
							switch o.(type) {
							case *dns.EDNS0_SUBNET:
								subnet = o.String()
								subnetRR, err := dns.NewRR(fmt.Sprintf("%s 0 IN TXT subnet=%s", q.Name, subnet))
								if err != nil {
									log.Printf("%s source %s txt %s your.subnet.dns[%s]", requestId, w.RemoteAddr(), q.Name, err.Error())
								}
								m.Answer = append(m.Answer, subnetRR)
							}
						}
					}
				}

				// request id
				requestIdRR, _ := dns.NewRR(fmt.Sprintf("%s 0 IN TXT request_id=%s", q.Name, requestId))
				m.Answer = append(m.Answer, requestIdRR)

				log.Printf("%s source %s txt %s your.local.dns[%s] your.subnet.dns[%s] ", requestId, w.RemoteAddr(), q.Name, localDnsIP, subnet)

			case dns.TypeA:
				log.Printf("%s source %s a %s ", requestId, w.RemoteAddr(), q.Name)
				// local dns
				localDnsIP := strings.Split(w.RemoteAddr().String(), ":")[0]
				localDnsRR, err := dns.NewRR(fmt.Sprintf("%s 0 IN a %s", q.Name, localDnsIP))
				if err != nil {
					log.Printf("%s source %s a %s your.local.dns[%s]", requestId, w.RemoteAddr(), q.Name, err.Error())
					return
				}
				m.Answer = append(m.Answer, localDnsRR)

				// subnet handler
				subnet := ""
				for _, v := range r.Extra {
					opt, ok := v.(*dns.OPT)
					if ok {
						for _, o := range opt.Option {
							switch o.(type) {
							case *dns.EDNS0_SUBNET:
								subnet = o.String()
								subnetRR, err := dns.NewRR(fmt.Sprintf("%s 0 IN a %s", q.Name, strings.Split(subnet, "/")[0]))
								if err != nil {
									log.Printf("%s source %s txt %s your.subnet.dns[%s]", requestId, w.RemoteAddr(), q.Name, err.Error())
								}
								m.Answer = append(m.Answer, subnetRR)
							}
						}
					}
				}

				log.Printf("%s source %s a %s your.local.dns[%s] your.subnet.dns[%s] ", requestId, w.RemoteAddr(), q.Name, localDnsIP, subnet)
			}

		}
	}

	_ = w.WriteMsg(m)
}

func main() {
	addr := flag.String("addr", ":53", "listen addr")
	proto := flag.String("proto", "udp", "listen protocol")
	flag.Parse()

	dns.HandleFunc(".", handleDnsRequest)

	server := dns.Server{Addr: *addr, Net: *proto}
	log.Printf("start dns server listen %s(%s)", *addr, *proto)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
