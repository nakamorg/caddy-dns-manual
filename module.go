package template

import (
	"context"
	"fmt"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/libdns"
	"go.uber.org/zap"
)

// Provider lets Caddy read and manipulate DNS records hosted by this DNS provider.
type Provider struct {
	WaitInMins string `json:"wait_in_mins,omitempty"`
}

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.manual_dns",
		New: func() caddy.Module { return new(Provider) },
	}
}

// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	p.WaitInMins = caddy.NewReplacer().ReplaceAll(p.WaitInMins, "1")
	return nil
}

// AppendRecords doesn't do anything and simply returns the records that were asked to be added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	caddy.Log().Named("manual-dns").Info("please append following dns records manually", zap.Reflect("records", records))
	p.wait()
	return records, nil
}

// DeleteRecords doesn't do anything and simply returns the records that were asked to be deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	caddy.Log().Named("manual-dns").Info("please delete following dns records manually", zap.Reflect("records", records))
	p.wait()
	return records, nil
}

func (p *Provider) wait() error {
	minutesToWait, err := time.ParseDuration(fmt.Sprintf("%sm", p.WaitInMins))
	if err != nil {
		caddy.Log().Named("manual-dns").Error("waiting for records", zap.Error(err))
		return err
	}
	caddy.Log().Named("manual-dns").Info("waiting for records", zap.Duration("time", minutesToWait))
	time.Sleep(minutesToWait)
	caddy.Log().Named("manual-dns").Info("wait finished")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	providername [<wait_in_mins>] {
//	    wait_in_mins <wait_in_mins>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.WaitInMins = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "wait_in_mins":
				if d.NextArg() {
					p.WaitInMins = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.WaitInMins == "" {
		p.WaitInMins = "1"
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
