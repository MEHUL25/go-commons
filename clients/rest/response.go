package rest

import (
	"net/http"

	"go.nandlabs.io/commons/codec"
	"go.nandlabs.io/commons/errutils"
	"go.nandlabs.io/commons/ioutils"
)

type Response struct {
	raw *http.Response
}

//IsSuccess determines if the response is a success response
func (r Response) IsSuccess() bool {
	return r.raw.StatusCode >= 200 && r.raw.StatusCode <= 204
}

//GetError gets the error with status code and value
func (r Response) GetError() (err error) {
	if !r.IsSuccess() {
		err = errutils.FmtError("Server responded with status code %d and status text %s",
			r.raw.StatusCode, r.raw.Status)
	}
	return
}

//Decode Function decodes the response body to a suitable object. The format of the body is determined by
//Content-Type header in the response
func (r Response) Decode(v interface{}) (err error) {
	var c codec.Codec
	if r.IsSuccess() {
		defer ioutils.CloserFunc(r.raw.Body)
		contentType := r.raw.Header.Get(contentTypeHdr)
		c, err = codec.Get(contentType, codec.DefaultCodecOptions)
		if err == nil {
			err = c.Read(r.raw.Body, v)
		}
	} else {
		err = r.GetError()
	}
	return
}

func (r Response) Raw() *http.Response {
	return r.raw
}
