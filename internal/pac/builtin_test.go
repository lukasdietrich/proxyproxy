package pac

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPlainHostName(t *testing.T) {
	assert.False(t, isPlainHostName("www.mozilla.org"))
	assert.True(t, isPlainHostName("www"))
}

func TestDnsDomainIs(t *testing.T) {
	assert.True(t, dnsDomainIs("www.mozilla.org", ".mozilla.org"))
	assert.False(t, dnsDomainIs("www", ".mozilla.org"))
}

func TestLocalHostOrDomainIs(t *testing.T) {
	assert.True(t, localHostOrDomainIs("www.mozilla.org", "www.mozilla.org"))
	assert.True(t, localHostOrDomainIs("www", "www.mozilla.org"))

	assert.False(t, localHostOrDomainIs("www.google.com", "www.mozilla.org"))
	assert.False(t, localHostOrDomainIs("home.mozilla.org", "www.mozilla.org"))
}

func TestIsResolvable(t *testing.T) {
	assert.True(t, isResolvable("www.mozilla.org"))
	assert.False(t, isResolvable("absolute.not.resolvable.mozilla.org"))

	assert.False(t, isResolvable(""))
}

func TestIsInNet(t *testing.T) {
	assert.True(t, isInNet("192.168.178.123", "192.168.178.10", "255.255.255.0"))
	assert.False(t, isInNet("192.168.178.123", "192.168.179.10", "255.255.255.0"))

	assert.False(t, isInNet("", "", ""))
}

func TestDnsResolve(t *testing.T) {
	assert.Equal(t, "127.0.0.1", dnsResolve("localhost"))
}

func TestConvertAddr(t *testing.T) {
	assert.EqualValues(t, 3221226156, convertAddr("192.0.2.172"))
	assert.EqualValues(t, 0, convertAddr(""))
}

func TestMyIpAddress(t *testing.T) {
	assert.Equal(t, "127.0.0.1", myIpAddress())
}

func TestDnsDomainLevels(t *testing.T) {
	assert.Equal(t, 0, dnsDomainLevels("www"))
	assert.Equal(t, 1, dnsDomainLevels("mozilla.org"))
	assert.Equal(t, 2, dnsDomainLevels("www.mozilla.org"))
}

func TestShExpMatch(t *testing.T) {
	assert.True(t, shExpMatch("http://home.netscape.com/people/ari/index.html", "*/ari/*"))
	assert.False(t, shExpMatch("http://home.netscape.com/people/montulli/index.html", "*/ari/*"))
}
