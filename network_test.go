// Copyright 2026 RelyChan Pte. Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goutils

import (
	"errors"
	"net"
	"testing"
)

func TestParseSubnet(t *testing.T) {
	t.Run("valid CIDR", func(t *testing.T) {
		subnet, err := ParseSubnet("192.168.1.0/24")
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		if subnet == nil {
			t.Fatal("expected non-nil subnet")
		}
	})

	t.Run("IPv4 address without prefix gets /32", func(t *testing.T) {
		subnet, err := ParseSubnet("10.0.0.1")
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		ones, bits := subnet.Mask.Size()
		if ones != 32 || bits != 32 {
			t.Fatalf("expected /32, got /%d/%d", ones, bits)
		}
	})

	t.Run("IPv6 address without prefix gets /128", func(t *testing.T) {
		subnet, err := ParseSubnet("::1")
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		ones, bits := subnet.Mask.Size()
		if ones != 128 || bits != 128 {
			t.Fatalf("expected /128, got /%d/%d", ones, bits)
		}
	})

	t.Run("empty string returns ErrInvalidSubnet", func(t *testing.T) {
		_, err := ParseSubnet("")
		if !errors.Is(err, ErrInvalidSubnet) {
			t.Fatalf("expected ErrInvalidSubnet, got: %v", err)
		}
	})

	t.Run("invalid IP string returns ErrInvalidSubnet", func(t *testing.T) {
		_, err := ParseSubnet("not-an-ip")
		if !errors.Is(err, ErrInvalidSubnet) {
			t.Fatalf("expected ErrInvalidSubnet, got: %v", err)
		}
	})

	t.Run("invalid CIDR notation", func(t *testing.T) {
		_, err := ParseSubnet("999.999.999.999/24")
		if err == nil {
			t.Fatal("expected error for invalid CIDR, got nil")
		}
	})

	t.Run("valid IPv6 CIDR", func(t *testing.T) {
		subnet, err := ParseSubnet("2001:db8::/32")
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		if subnet == nil {
			t.Fatal("expected non-nil subnet")
		}
	})
}

// parseNetCIDR is a helper that wraps net.ParseCIDR for test use.
func parseNetCIDR(s string) (net.IP, *net.IPNet, error) {
	return net.ParseCIDR(s)
}
