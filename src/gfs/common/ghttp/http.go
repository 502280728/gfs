// http
package ghttp

import (
	"gfs/common"
	"gfs/common/ghttp/cookie"
	"io"
	"net/http"
)

type GFSRequest struct {
	CS cookie.CookieStore
}

func (req *GFSRequest) SetCookieStore(cs cookie.CookieStore) {
	req.CS = cs
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

	if req.CS != nil {
		req.CS.WriteCookie(httpreq)
	}
	resp, err2 := http.DefaultClient.Do(httpreq)
	if err2 != nil {
		return err2
	}
	defer resp.Body.Close()
	if req.CS != nil {
		req.CS.StoreCookie(resp)
	}
	common.DecodeFromReader(result, resp.Body)
	return nil
}
