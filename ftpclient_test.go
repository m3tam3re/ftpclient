package ftpclient

import (
	"testing"
)

func TestFtpclient(t *testing.T) {
	f := FtpCon{
		user: "one",
		pass: "1234",
		host: "localhost:21",
	}
	err := f.GetFiletype("/Users/m3tam3re/Tools", "/ftp/one", "", true)
	t.Log(err)
}
