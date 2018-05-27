// http
package http

import (
	"gfs/common"
	"gfs/common/http/cookie"
	"io"
	"net/http"
)

type GFSRequest struct {
	cs cookie.CookieStore
}

func (req *GFSRequest) SetCookieStore(cs cookie.CookieStore) {
	req.cs = cs
}

func (req *GFSRequest) PostObj(url string, body interface{}, result interface{}) error {

	return req.Post(url, nil, common.EncodeToByteBuffer(body), result)
}

func (req *GFSRequest) Post(url string, headers map[string]string, body io.Reader, result interface{}) error {
	httpreq, err1 := http.NewRequest("POST", url, body)
	if err1 != nil {
		return err1
	}
	if headers != nil {
		for k, v := range headers {
			httpreq.Header.Add(k, v)
		}
	}

	if req.cs != nil {
		req.cs.WriteCookie(httpreq)
	}
	resp, err2 := http.DefaultClient.Do(httpreq)
	defer resp.Body.Close()
	if err2 != nil {
		return err2
	}
	if req.cs != nil {
		req.cs.StoreCookie(resp)
	}
	common.DecodeFromReader(result, resp.Body)
	return nil
}
