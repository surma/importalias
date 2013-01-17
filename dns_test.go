package main

import (
	"reflect"
	"testing"

	"github.com/miekg/dns"
)

func TestRecursiveDNS(t *testing.T) {
	table := map[string]dns.RR{
		dns.Fqdn("_importalias.go.surmair.de"): &dns.TXT{
			Hdr: dns.RR_Header{
				Name:     "", // Will be set automatically
				Rrtype:   dns.TypeTXT,
				Class:    dns.ClassINET,
				Ttl:      60,
				Rdlength: 33,
			},
			Txt: []string{
				"d41f99784fd644a7b3990f86741b5eb5",
			},
		},
	}

	for domain, expected := range table {
		expected.Header().Name = domain
		rr, err := RecursiveDNS(domain, expected.Header().Rrtype)
		if err != nil {
			t.Fatalf("Failed to resolve DNS: %s", err)
		}
		if len(rr) != 1 {
			t.Fatalf("Empty response")
		}
		if txt, ok := rr[0].(*dns.TXT); !ok || !reflect.DeepEqual(txt, expected) {
			t.Fatalf("Unexpected response. Wanted %#v, got %#v", expected, txt)
		}
	}
}
