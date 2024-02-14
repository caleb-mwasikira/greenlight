package main

import (
	"fmt"
	"net"
	"net/http"
)

func main() {
	host, err := interfaceIPv4Addr("wlp3s0")
	if err != nil {
		panic(err)
	}

	port := 8080
	address := fmt.Sprintf("%v:%d", host, port)

	http.Handle("/", startFileServer(address, "public"))
	http.ListenAndServe(address, nil)
}

/*
* Fetches the IPv4 address of a given interface name
 */
func interfaceIPv4Addr(name string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, ifi := range interfaces {
		if ifi.Name == name {
			addrs, _ := ifi.Addrs()

			for _, addr := range addrs {
				ipv4, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					return "", err
				} else {
					return ipv4.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("interface '%v' could not be found", name)
}

func startFileServer(addr string, dir string) http.Handler {
	msg := fmt.Sprintf("Server running on address: http://%v", addr)
	fmt.Println(msg)

	return http.FileServer(http.Dir(dir))
}
