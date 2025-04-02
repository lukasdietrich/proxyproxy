# proxyproxy

An http proxy to proxy another http proxy 🤪

## Usecase

Imagine you worked for a company that only allowed internet traffic through an http proxy and
rolled it out via pac[^1] file. Now imagine that most programs under linux do not support pac.

Some of us don't have to imagine that.

This application acts as a proxy between applications and the corporate http proxy.
This way, applications can use environment variables[^2] to route their traffic through proxyproxy,
which in turn understands the pac file and routes the traffic to the desired upstream proxy.

```mermaid
sequenceDiagram
  autonumber

  actor Y as You
  create participant C as curl

  participant P as proxyproxy


  Y ->>  C: HTTP GET https://example.org/my-secret-file.txt
  activate C
  C ->>  P: HTTP CONNECT //example.org:443

  create participant F as corporate.pac
  P ->>  F: Should I use another proxy to connect?
  activate F
  note over P,F: The corporate.pac file is a JavaScript that evaluates an url to <br> either "PROXY host:port" or "DIRECT" for a direct connection.
  F -->> P: Connect via corporate proxy.
  deactivate F

  participant U as corporate proxy
  P ->>  U: HTTP CONNECT //example.org:443
  participant T as example.org
  U ->>  T: TCP Connect to example.org:443.
  activate T
  note over C,T: Data is copied between all hops, but the encryption happens end-to-end between curl and example.org unless plain http is used.
  T <<-->>  C: Tunnel encrypted data
  deactivate T

  deactivate C
  destroy C
  C ->>  Y: my-secret-file.txt

```

[^1]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file
[^2]: https://everything.curl.dev/transfers/conn/proxies.html?highlight=https_proxy#proxy-environment-variables
