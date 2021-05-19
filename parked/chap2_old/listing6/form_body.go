package pkgregister

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
)

func createMultiPartMessage(data pkgData) ([]byte, string, error) {
	var b bytes.Buffer
	var err error
	var fw io.Writer
	var d []byte

	mw := multipart.NewWriter(&b)

	fw, err = mw.CreateFormField("name")
	if err != nil {
		return nil, "", err
	}
	fw.Write([]byte(data.Name))

	fw, err = mw.CreateFormField("version")
	if err != nil {
		return nil, "", err
	}
	fw.Write([]byte(data.Version))

	fw, err = mw.CreateFormFile("filedata", data.Filename)
	if err != nil {
		return nil, "", err
	}
	d, err = ioutil.ReadAll(data.Bytes)
	if err != nil {
		return nil, "", err
	}
	fw.Write(d)
	err = mw.Close()
	if err != nil {
		return nil, "", err
	}

	contentType := mw.FormDataContentType()
	return b.Bytes(), contentType, nil
}
