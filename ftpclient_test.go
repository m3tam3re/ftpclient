package ftpclient

import (
	"testing"
)

func TestFtpclient(t *testing.T) {
	f := FtpCon{
		User: "one",
		Pass: "1234",
		Host: "localhost:21",
	}
	fl, err := f.GetFiletype("/Users/m3tam3re/Tools", "/ftp/one", "", true)
	t.Log(fl, err)
}
