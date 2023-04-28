package template

import (
	"fmt"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/libdns"
	"go.uber.org/zap"
)

// Provider lets Caddy read and manipulate DNS records hosted by this DNS provider.
type Provider struct {
	WaitInMins string
}

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.manual_dns",
		New: func() caddy.Module { return &Provider{new(Provider)} },
	}
}

// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	p.Provider.WaitInMins = caddy.NewReplacer().ReplaceAll(p.Provider.WaitInMins, "1")
	return nil
}


// AppendRecords doesn't do anything and simply returns the records that were asked to be added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	caddy.Log().Named("manual-dns").Info("appending dns records", zap.Info(records))
	p.wait()
	return records, nil
}

// DeleteRecords doesn't do anything and simply returns the records that were asked to be deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	caddy.Log().Named("manual-dns").Info("deleting dns records", zap.Info(records))
	p.wait()
	return records, nil
}

func (p *Provider) wait() error {
	minutesToWait, err := time.ParseDuration(fmt.Sprintf("%sm", p.WaitInMins))
	if err != nil {
		caddy.Log().Named("manual-dns").Error("waiting for records", zap.Error(err))
		return err
	}
	caddy.Log().Named("manual-dns").Info("waiting for records", zap.Info(minutesToWait))
	time.Sleep(minutesToWait)
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// providername [<wait_in_mins>] {
//     wait_in_mins <wait_in_mins>
// }
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	waitInMins 
	for d.Next() {
		if d.NextArg() {
			p.Provider.WaitInMins = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "wait_in_mins":
				if d.NextArg() {
					p.Provider.WaitInMins = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.WaitInMins == "" {
		p.Provider.WaitInMins = "1"
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)