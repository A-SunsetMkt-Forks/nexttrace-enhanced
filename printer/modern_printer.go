package printer

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/xgadget-lab/nexttrace/trace"
)

func RealtimePrinter(res *trace.Result, ttl int) {
	fmt.Printf("%s  ", color.New(color.FgHiYellow, color.Bold).Sprintf("%-2d", ttl+1))

	// 去重
	var latestIP string
	tmpMap := make(map[string][]string)
	for i, v := range res.Hops[ttl] {
		if v.Address == nil && latestIP != "" {
			tmpMap[latestIP] = append(tmpMap[latestIP], fmt.Sprintf("%s ms", "*"))
			continue
		} else if v.Address == nil {
			continue
		}

		if _, exist := tmpMap[v.Address.String()]; !exist {
			tmpMap[v.Address.String()] = append(tmpMap[v.Address.String()], strconv.Itoa(i))
			// 首次进入
			if latestIP == "" {
				for j := 0; j < i; j++ {
					tmpMap[v.Address.String()] = append(tmpMap[v.Address.String()], fmt.Sprintf("%s ms", "*"))
				}
			}
			latestIP = v.Address.String()
		}

		tmpMap[v.Address.String()] = append(tmpMap[v.Address.String()], fmt.Sprintf("%.2f ms", v.RTT.Seconds()*1000))
	}

	if latestIP == "" {
		fmt.Fprintf(color.Output, "%s\n",
			color.New(color.FgWhite, color.Bold).Sprintf("*"),
		)
		return
	}

	var blockDisplay = false
	for ip, v := range tmpMap {
		if blockDisplay {
			fmt.Printf("%4s", "")
		}
		fmt.Fprintf(color.Output, "%s",
			color.New(color.FgWhite, color.Bold).Sprintf("%-15s", ip),
		)

		i, _ := strconv.Atoi(v[0])

		if res.Hops[ttl][i].Geo.Asnumber != "" {
			fmt.Fprintf(color.Output, " %s", color.New(color.FgHiGreen, color.Bold).Sprintf("AS%-6s", res.Hops[ttl][i].Geo.Asnumber))
		} else {
			fmt.Printf(" %-8s", "*")
		}

		if res.Hops[ttl][i].Geo.Country == "" {
			res.Hops[ttl][i].Geo.Country = "LAN Address"
		}

		fmt.Fprintf(color.Output, " %s %s %s %s %s\n    %s   ",
			color.New(color.FgWhite, color.Bold).Sprintf("%s", res.Hops[ttl][i].Geo.Country),
			color.New(color.FgWhite, color.Bold).Sprintf("%s", res.Hops[ttl][i].Geo.Prov),
			color.New(color.FgWhite, color.Bold).Sprintf("%s", res.Hops[ttl][i].Geo.City),
			color.New(color.FgWhite, color.Bold).Sprintf("%s", res.Hops[ttl][i].Geo.District),
			fmt.Sprintf("%-6s", res.Hops[ttl][i].Geo.Owner),
			color.New(color.FgHiBlack, color.Bold).Sprintf("%-22s", res.Hops[ttl][0].Hostname),
		)

		for j := 1; j < len(v); j++ {
			if len(v) == 2 || j == 1 {
				fmt.Fprintf(color.Output, "%s",
					color.New(color.FgHiCyan, color.Bold).Sprintf("%s", v[j]),
				)
			} else {
				fmt.Fprintf(color.Output, " / %s",
					color.New(color.FgHiCyan, color.Bold).Sprintf("%s", v[j]),
				)
			}
		}
		fmt.Println()
		blockDisplay = true
	}
}
