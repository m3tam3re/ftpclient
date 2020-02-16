// Package ftpclient implements some helper functions for an ftp client
package ftpclient

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/m3tam3re/errors"
	"io/ioutil"
	"path/filepath"
	"time"
)

const path errors.Path = "github.com/m3tam3re/ftpclient/ftpclient.go"

type FtpCon struct {
	con  *ftp.ServerConn
	User string
	Pass string
	Host string
}

// GetFiletype downloads all files of the given ftype from the rfolder on the FTP Server
// to the local lfolder. if del is true files will be deleted after downloading
//
// Function params:
// lfolder = local folder for storing the files
// rfolder = remote folder to download the files from
// ftype = the file extension of the files that should be downloaded / "" for all files
// del = true for deleting files from the remote folder after download, false to keep them
func (fc *FtpCon) GetFiletype(lfolder string, rfolder string, ftype string, del bool) ([]string, error) {
	var op errors.Op = "method: GetFiletype()"
	var downloaded []string

	if fc.con == nil {
		err := fc.Connect()
		if err != nil {
			return nil, errors.E(errors.Internal, path, op, err)
		}
	}
	defer fc.con.Logout()
	if lfolder[len(lfolder)-1:] != "/" {
		lfolder = lfolder + "/"
	}
	if rfolder[len(rfolder)-1:] != "/" {
		rfolder = rfolder + "/"
	}
	if ftype != "" && ftype[0:1] != "." {
		ftype = "." + ftype
	}
	err := fc.con.Login(fc.User, fc.Pass)
	if err != nil {
		return nil, errors.E(path, op, err, "could not login")
	}
	fl, err := fc.con.List(rfolder)
	if err != nil {
		return nil, errors.E(errors.Internal, path, op, err, "could not get file list")
	}
	for _, file := range fl {
		if filepath.Ext(rfolder+file.Name) == ftype || ftype == "" {
			f, err := fc.con.Retr(rfolder + file.Name)
			if err != nil {
				return nil, errors.E(errors.Internal, path, op, err, fmt.Sprintf("could not download %s", rfolder+file.Name))
			}
			buf, err := ioutil.ReadAll(f)
			if err != nil {
				return nil, errors.E(errors.Internal, path, op, err, fmt.Sprintf("could not download %s", rfolder+file.Name))
			}
			err = ioutil.WriteFile(lfolder+file.Name, buf, 0644)
			if err != nil {
				return nil, errors.E(errors.Internal, path, op, err, "could not write to local folder")
			}
			downloaded = append(downloaded, lfolder+file.Name)
			f.Close()
			if del {
				err = fc.con.Delete(rfolder + file.Name)
				if err != nil {
					return nil, errors.E(errors.Internal, path, op, err, "could not delete file")
				}
			}
		}
	}
	return downloaded, err
}

func (fc *FtpCon) Connect() error {
	var op errors.Op = "method: Connect()"
	c, err := ftp.Dial(fc.Host, ftp.DialWithTimeout(time.Second*10))
	if err != nil {
		return errors.E(errors.Internal, path, op, err, "connection error")
	}
	fc.con = c
	return nil
}
