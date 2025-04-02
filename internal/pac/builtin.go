package pac

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/dop251/goja"
	"github.com/gobwas/glob"
)

func declareBuiltins(vm *goja.Runtime) {
	for name, fn := range map[string]any{
		"isPlainHostName":     isPlainHostName,
		"dnsDomainIs":         dnsDomainIs,
		"localHostOrDomainIs": localHostOrDomainIs,
		"isResolvable":        isResolvable,
		"isInNet":             isInNet,
		"dnsResolve":          dnsResolve,
		"convert_addr":        convertAddr,
		"myIpAddress":         myIpAddress,
		"dnsDomainLevels":     dnsDomainLevels,
		"shExpMatch":          shExpMatch,
		"alert":               alert,
	} {
		declareFunction(vm, name, fn)
	}

	// TODO implement time based conditions
	//
	//  - weekdayRange
	//  - dateRange
	//  - timeRange
}

func declareFunction(vm *goja.Runtime, name string, fn any) error {
	log := slog.With(slog.String("builtin", name))

	callable, ok := goja.AssertFunction(vm.ToValue(fn))
	if !ok {
		return fmt.Errorf("provided function %s is not a function", name)
	}

	vm.GlobalObject().Set(name, func(call goja.FunctionCall) (returnValue goja.Value) {
		var err error

		defer func() {
			if r := recover(); r != nil {
				returnValue = goja.Undefined()

				log.Warn("panic while calling function",
					slog.Any("call", call),
					slog.Any("panic", r),
				)
			} else {
				log.Debug("called function",
					slog.Any("call", call),
					slog.Any("return", returnValue),
					slog.Any("err", err),
				)
			}
		}()

		returnValue, err = callable(call.This, call.Arguments...)
		if err != nil {
			returnValue = goja.Undefined()
		}

		return
	})

	return nil
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file#isplainhostname
func isPlainHostName(host string) bool {
	return !strings.ContainsRune(host, '.')
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file#dnsdomainis
func dnsDomainIs(host, domain string) bool {
	return strings.HasPrefix(domain, ".") && strings.HasSuffix(host, domain)
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file#localhostordomainis
func localHostOrDomainIs(host, hostDom string) bool {
	if isPlainHostName(host) {
		return strings.HasPrefix(hostDom, host+".")
	}

	return host == hostDom
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file#isresolvable
func isResolvable(host string) bool {
	return dnsResolve(host) != ""
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file#isinnet
func isInNet(host, pattern, mask string) bool {
	addr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return false
	}

	network := net.IPNet{
		IP:   net.ParseIP(pattern),
		Mask: net.IPMask(net.ParseIP(mask)),
	}

	return network.IP != nil && network.Mask != nil && network.Contains(addr.IP)
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file#dnsresolve
func dnsResolve(host string) string {
	ip, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return ""
	}

	return ip.String()
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file#convert_addr
func convertAddr(ip string) uint32 {
	if addr := net.ParseIP(ip).To4(); addr != nil {
		return binary.BigEndian.Uint32(addr)
	}

	return 0
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file#myipaddress
func myIpAddress() string {
	return "127.0.0.1"
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file#dnsdomainlevels
func dnsDomainLevels(host string) int {
	return strings.Count(host, ".")
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file#shexpmatch
func shExpMatch(str, shExp string) bool {
	matcher, err := glob.Compile(shExp)
	if err != nil {
		return false
	}

	return matcher.Match(str)
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file#alert
func alert(message string) {
	slog.Info("pac alert", slog.String("message", message))
}
