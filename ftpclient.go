// Package modules implements some helper functions for an ftp client
package ftpclient

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"path/filepath"
)

type ftpCon struct {
	*ftp.ServerConn
}

// GetFiles downloads all files of the given ftype from the rfolder on the FTP Server
// to the local lfolder. if del is true files will be deleted after downloading
func (fc *ftpCon) GetFiletype(rfolder string, ftype string, lfolder string, del bool) error {

	if lfolder[len(lfolder)-1:] != "/" {
		rfolder = rfolder + "/"
	}
	if rfolder[len(rfolder)-1:] != "/" {
		rfolder = rfolder + "/"
	}
	fl, err := fc.List(rfolder)
	if err != nil {
		return fmt.Errorf("Could not get file list: %s", err)
	}
	for _, file := range fl {
		if filepath.Ext(rfolder+file.Name) == ftype {
			f, err := fc.Retr(rfolder + file.Name)
			if err != nil {
				return fmt.Errorf("Could not download %s.", rfolder+file.Name)
			}
			buf, err := ioutil.ReadAll(f)
			if err != nil {
				return fmt.Errorf("Error reading file: %s.", rfolder+file.Name)
			}
			err = ioutil.WriteFile(lfolder+file.Name, buf, 0644)
			if err != nil {
				return fmt.Errorf("Could not write to folder: %s", err)
			}
			f.Close()
			if del {
				err = fc.Delete(rfolder + file.Name)
				if err != nil {
					return fmt.Errorf("Error deleting source file %s", err)
				}
			}
		}
	}
	return err
}
