// conf_validator
package gconf

import (
	"fmt"
	"net/url"
)

type ConfValidator func(string) error

var URL_VALIDATOR = ConfValidator(func(s string) error {
	if url, err := url.Parse(s); err == nil {
		if url.Scheme == "" || url.Host == "" {
			return fmt.Errorf("%s is not a valid url,must contains protocol and host", s)
		} else {
			return nil
		}
	} else {
		return err
	}
})
