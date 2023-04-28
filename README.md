Manual DNS module for Caddy
===========================

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manually manage DNS records.

## Caddy module name

```
dns.providers.manual_dns
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
	"module": "acme",
	"challenges": {
		"dns": {
			"provider": {
				"name": "manual_dns",
				"wait_in_mins": "1",
			}
		}
	}
}
```

or with the Caddyfile:

```
# globally
{
	acme_dns manual_dns ...
}
```

```
# one site
tls {
	dns manual_dns ...
}
```
