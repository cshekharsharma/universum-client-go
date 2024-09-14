package universum

import (
	"testing"
	"time"
)

func TestOptionsInit(t *testing.T) {
	testCases := []struct {
		name     string
		input    Options
		expected Options
	}{
		{
			name: "All default values",
			input: Options{
				HostAddr:        "",
				ClientName:      "",
				DialTimeout:     0,
				ReadTimeout:     0,
				WriteTimeout:    0,
				MaxRetries:      0,
				RetryBackoff:    0,
				ConnPoolsize:    0,
				ConnMaxLifetime: 0,
			},
			expected: Options{
				HostAddr:        DefaultHostAddr,
				ClientName:      DefaultClientName,
				DialTimeout:     DefaultDialTimeout,
				ReadTimeout:     DefaultReadTimeout,
				WriteTimeout:    DefaultWriteTimeout,
				MaxRetries:      DefaultMaxRetries,
				RetryBackoff:    DefaultRetryBackoff,
				ConnPoolsize:    DefaultConnPoolsize,
				ConnMaxLifetime: DefaultConnMaxLifetime,
			},
		},
		{
			name: "Exceeding max values",
			input: Options{
				DialTimeout:     10 * time.Second,
				ReadTimeout:     10 * time.Second,
				WriteTimeout:    10 * time.Second,
				MaxRetries:      100,
				RetryBackoff:    1 * time.Second,
				ConnPoolsize:    70000,
				ConnMaxLifetime: 40 * time.Minute,
			},
			expected: Options{
				HostAddr:        DefaultHostAddr,
				ClientName:      DefaultClientName,
				DialTimeout:     MaxDialTimeout,
				ReadTimeout:     MaxReadTimeout,
				WriteTimeout:    MaxWriteTimeout,
				MaxRetries:      AllowedMaxRetries,
				RetryBackoff:    MaxRetryBackoff,
				ConnPoolsize:    MaxConnPoolsize,
				ConnMaxLifetime: MaxConnMaxLifetime,
			},
		},
		{
			name: "Valid values within limits",
			input: Options{
				HostAddr:        "customhost:12345",
				ClientName:      "CustomClient",
				DialTimeout:     2 * time.Second,
				ReadTimeout:     2 * time.Second,
				WriteTimeout:    2 * time.Second,
				MaxRetries:      5,
				RetryBackoff:    100 * time.Millisecond,
				ConnPoolsize:    5000,
				ConnMaxLifetime: 20 * time.Minute,
			},
			expected: Options{
				HostAddr:        "customhost:12345",
				ClientName:      "CustomClient",
				DialTimeout:     2 * time.Second,
				ReadTimeout:     2 * time.Second,
				WriteTimeout:    2 * time.Second,
				MaxRetries:      5,
				RetryBackoff:    100 * time.Millisecond,
				ConnPoolsize:    5000,
				ConnMaxLifetime: 20 * time.Minute,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.input.Init()

			if tc.input.HostAddr != tc.expected.HostAddr {
				t.Errorf("Expected HostAddr %s, got %s", tc.expected.HostAddr, tc.input.HostAddr)
			}
			if tc.input.ClientName != tc.expected.ClientName {
				t.Errorf("Expected ClientName %s, got %s", tc.expected.ClientName, tc.input.ClientName)
			}
			if tc.input.DialTimeout != tc.expected.DialTimeout {
				t.Errorf("Expected DialTimeout %s, got %s", tc.expected.DialTimeout, tc.input.DialTimeout)
			}
			if tc.input.ReadTimeout != tc.expected.ReadTimeout {
				t.Errorf("Expected ReadTimeout %s, got %s", tc.expected.ReadTimeout, tc.input.ReadTimeout)
			}
			if tc.input.WriteTimeout != tc.expected.WriteTimeout {
				t.Errorf("Expected WriteTimeout %s, got %s", tc.expected.WriteTimeout, tc.input.WriteTimeout)
			}
			if tc.input.MaxRetries != tc.expected.MaxRetries {
				t.Errorf("Expected MaxRetries %d, got %d", tc.expected.MaxRetries, tc.input.MaxRetries)
			}
			if tc.input.RetryBackoff != tc.expected.RetryBackoff {
				t.Errorf("Expected RetryBackoff %s, got %s", tc.expected.RetryBackoff, tc.input.RetryBackoff)
			}
			if tc.input.ConnPoolsize != tc.expected.ConnPoolsize {
				t.Errorf("Expected ConnPoolsize %d, got %d", tc.expected.ConnPoolsize, tc.input.ConnPoolsize)
			}
			if tc.input.ConnMaxLifetime != tc.expected.ConnMaxLifetime {
				t.Errorf("Expected ConnMaxLifetime %s, got %s", tc.expected.ConnMaxLifetime, tc.input.ConnMaxLifetime)
			}
		})
	}
}
