package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nFooFoo:BarBar    \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host, ok := headers.Get("HOST")
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, true, ok)
	host, ok = headers.Get("FooFOo")
	assert.Equal(t, "BarBar", host)
	assert.Equal(t, true, ok)
	assert.Equal(t, len(data), n)
	assert.True(t, done)

	//Test: Valid single header with extra whitespaces
	headers = NewHeaders()
	data = []byte("Host:     localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host, ok = headers.Get("HosT")
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, len(data), n)
	assert.True(t, done)

	//Test: Valid 2 headers
	headers = NewHeaders()
	data = []byte("Host:     localhost:42069       \r\nServer: Google.com    \r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host, ok = headers.Get("Host")
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, true, ok)
	host, ok = headers.Get("SerVer")
	assert.Equal(t, "Google.com", host)
	assert.Equal(t, true, ok)
	assert.Equal(t, len(data), n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid token
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	//Test: 2 or more values per header name
	headers = NewHeaders()
	data = []byte("Host:localhost:42069\r\nHost:Google.com\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host, ok = headers.Get("HosT")
	assert.Equal(t, "localhost:42069,Google.com", host)
	assert.Equal(t, true, ok)
	assert.Equal(t, len(data), n)
	assert.True(t, done)
}
