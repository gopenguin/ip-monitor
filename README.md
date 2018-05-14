# ip-monitor

A simple client application to automatically detect an ip address and report it to a dynamic dns api.

## Usage

```
Monitor all available ip addresses and update the corresponding dyndns provider

Usage:
  ip-monitor [flags]

Flags:
      --config string             config file (default is $HOME/.ip-monitor.yaml)
      --domain string             the domain to update, available via '{.Domain}'
      --expectedResponse string   the expected  response from the server; if empty, it will be ignored
  -h, --help                      help for ip-monitor
      --interval duration         the time between two updates (default 10m0s)
      --publicNet ipNet           the subnet which contains public ipv6 addresses (default 2000::/3)
      --token string              the token required to authenticate; in the url pattern available via '{.Token}'
      --urlTemplate string        the address to call to update the ip address; the ipv6 address is available via '{.IPv6}'
```

## License

This repository is licensed under the MIT license. For more details see [LICENSE](LICENSE)
