// Copyright 2023 The frp Authors
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

package net

import (
	"context"
	"net"
	"time"
)

var customDNSAddress string

func SetDefaultDNSAddress(dnsAddress string) {
	if _, _, err := net.SplitHostPort(dnsAddress); err != nil {
		dnsAddress = net.JoinHostPort(dnsAddress, "53")
	}
	customDNSAddress = dnsAddress
	createNewResolver()
}

// ClearDNSCache clears the DNS cache by creating a fresh resolver instance.
// This should be called before reconnection attempts to ensure fresh DNS lookups.
func ClearDNSCache() {
	createNewResolver()
}

func createNewResolver() {
	dnsAddr := customDNSAddress
	if dnsAddr == "" {
		// Use system default
		net.DefaultResolver = &net.Resolver{
			PreferGo: true,
		}
		return
	}
	
	// Use custom DNS address
	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
			var d net.Dialer
			dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return d.DialContext(dialCtx, network, dnsAddr)
		},
	}
}
