package main

import (
	"fmt"

	"github.com/miekg/dns"
)

func RecursiveDNS(domain string, ctype uint16) ([]dns.RR, error) {
	return recursiveDNS(domain, ctype, []string{"198.41.0.4:53"})
}

func recursiveDNS(domain string, ctype uint16, servers []string) ([]dns.RR, error) {
	for _, server := range servers {
		m := new(dns.Msg)
		m.SetQuestion(domain, ctype)
		c := new(dns.Client)
		in, _, err := c.Exchange(m, server)
		if err != nil {
			continue
		}
		if len(in.Answer) > 0 {
			return in.Answer, nil
		}
		authorities := make([]string, 0, len(in.Ns))
		for _, authority := range in.Ns {
			ns, ok := authority.(*dns.NS)
			if !ok {
				continue
			}
			for _, authority_a := range in.Extra {
				a, ok := authority_a.(*dns.A)
				if !ok {
					continue
				}
				if ns.Ns == a.Hdr.Name {
					authorities = append(authorities, a.A.String()+":53")
				}
			}
			if len(in.Extra) <= 0 {
				ips, err := RecursiveDNS(ns.Ns, dns.TypeA)
				if err != nil {
					continue
				}
				for _, ip := range ips {
					a, ok := ip.(*dns.A)
					if !ok {
						continue
					}
					authorities = append(authorities, a.A.String()+":53")
				}
			}
		}
		if len(authorities) > 0 {
			return recursiveDNS(domain, ctype, authorities)
		}
	}
	return nil, fmt.Errorf("No authority gave an answer")
}
