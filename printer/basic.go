package printer

import (
	"fmt"
	"net"

	"github.com/fatih/color"
)

var version = "v0.0.0.alpha"
var buildDate = ""
var commitID = ""

func Version() {
	fmt.Fprintf(color.Output, "%s %s %s %s\n",
		color.New(color.FgWhite, color.Bold).Sprintf("%s", "NextTrace Enhanced"),
		color.New(color.FgHiBlack, color.Bold).Sprintf("%s", version),
		color.New(color.FgHiBlack, color.Bold).Sprintf("%s", buildDate),
		color.New(color.FgHiBlack, color.Bold).Sprintf("%s", commitID),
	)
}

func CopyRight() {
	fmt.Println("OwO Network Vincent (vincent.moe) & Tso (tsosunchia@gmail.com) & zhshch (xzhsh.ch) & Leo (leo.moe)")
}

func PrintTraceRouteNav(ip net.IP, domain string, dataOrigin string) {
	fmt.Println("IP Geo Data Provider: " + dataOrigin)

	if ip.String() == domain {
		fmt.Printf("traceroute to %s, 30 hops max, 32 byte packets\n", ip.String())
	} else {
		fmt.Printf("traceroute to %s (%s), 30 hops max, 32 byte packets\n", ip.String(), domain)
	}
}
