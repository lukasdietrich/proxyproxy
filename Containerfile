from docker.io/library/golang:1.24-alpine as build
	workdir /build

	copy internal ./internal
	copy cmd ./cmd
	copy go* .

	run ls -l
	run go build -v ./cmd/proxyproxy

from docker.io/library/alpine
	workdir /app

	copy --from=build /build/proxyproxy .
	copy LICENSE .

	label org.opencontainers.image.authors="Lukas Dietrich <lukas@lukasdietrich.com>"
	label org.opencontainers.image.source="https://github.com/lukasdietrich/proxyproxy"

	env PROXYPROXY_AUTOCONFIGURE_ENABLED="true"
	env PROXYPROXY_AUTOCONFIGURE_ROOT="/auto-configure-root"

	expose 8080/tcp
	volume /auto-configure-root
	
	cmd [ "/app/proxyproxy" ]
