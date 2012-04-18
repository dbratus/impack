package web 

import (
	"net/http"
	"mime/multipart"
	"fmt"
	//"impack/loading"
	"archive/zip"
)

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
}

func index(resp http.ResponseWriter, req *http.Request) {
	http.ServeFile(resp, req, "index.html")
}

func upload(resp http.ResponseWriter, req *http.Request) {
	var arch multipart.File

	if f, _, err := req.FormFile("zip"); err == nil {
		arch = f
	} else {
		fmt.Fprintf(resp, "%s", err.Error())
		return
	}
	
	var fileSize int64
	
	if sz, err := arch.Seek(0, 2); err == nil {
		fileSize = sz
	} else {
		fmt.Fprintf(resp, "%s", err.Error())
		return		
	}
	
	if _, err := zip.NewReader(arch, fileSize); err == nil {
	} else {
		fmt.Fprintf(resp, "%s", err.Error())
		return
	}
}