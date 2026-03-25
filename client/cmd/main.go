//go:build windows

package main

import "rc/pkg/client"

func main() {
	client.Connect()
}
