package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/lack/waybar-nmvpn/pkg/nmvpn"

	waybar "github.com/lack/gowaybarplug"
)

func loop(interval time.Duration) {
	wb := waybar.NewUpdater()

	for true {
		names := []string{}
		status := waybar.Status{
			Text: "",
		}

		vpns, err := nmvpn.GetVPNs()
		if err == nil && len(vpns) > 0 {
			active := false
			for _, v := range vpns {
				if v.Active {
					active = true
					break
				}
			}
			if active {
				status.Alt = "connected"
				status.Class = []string{"connected"}
			} else {
				break
			}

			status.Tooltip = ""
			for _, v := range vpns {
				if !v.Active {
					continue
				}

				names = append(names, v.Name)

				if status.Tooltip != "" {
					status.Tooltip += "\n"
				}
				state := "down"
				if v.Active {
					state = "up"
				}
				status.Tooltip += fmt.Sprintf("%s: %s", v.Name, state)
			}
		} else if err != nil {
			status.Alt = "error"
			status.Class = []string{"error"}
			status.Tooltip = fmt.Sprintf("Error fetching VPN information: %v", err)
		} else {
			status.Alt = "none"
			status.Class = []string{"unconfigured"}
			status.Tooltip = "No VPN connections configured"
		}

		status.Text = strings.Join(names, " | ")

		wb.Status <- &status
		time.Sleep(interval)
	}
}

func main() {
	loop(5 * time.Second)
}
