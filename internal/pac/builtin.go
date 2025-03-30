package pac

import (
	"encoding/binary"
	"log/slog"
	"net"
	"strings"

	"github.com/dop251/goja"
	"github.com/gobwas/glob"
)

func declareBuiltins(vm *goja.Runtime) {
	vm.Set("isPlainHostName", vm.ToValue(isPlainHostName))
	vm.Set("dnsDomainIs", vm.ToValue(dnsDomainIs))
	vm.Set("localHostOrDomainIs", vm.ToValue(localHostOrDomainIs))
	vm.Set("isResolvable", vm.ToValue(isResolvable))
	vm.Set("isInNet", vm.ToValue(isInNet))

	vm.Set("dnsResolve", vm.ToValue(dnsResolve))
	vm.Set("convert_addr", vm.ToValue(convertAddr))
	vm.Set("myIpAddress", vm.ToValue(myIpAddress))
	vm.Set("dnsDomainLevels", vm.ToValue(dnsDomainLevels))

	vm.Set("shExpMatch", vm.ToValue(shExpMatch))

	// TODO implement time based conditions
	// vm.Set("weekdayRange", vm.ToValue(weekdayRange))
	// vm.Set("dateRange", vm.ToValue(dateRange))
	// vm.Set("timeRange", vm.ToValue(timeRange))

	vm.Set("alert", vm.ToValue(alert))
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

	return network.Contains(addr.IP)
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
	addr := net.ParseIP(ip).To4()
	return binary.BigEndian.Uint32(addr)
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
