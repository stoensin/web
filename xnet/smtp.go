package xnet

import (
	"crypto/tls"
	//"encoding/base64"
	//"errors"
	//"io"
	//"net"
	//"net/textproto"
	"net/smtp"
	"strings"
)

// Dial returns a new Client connected to an SMTP server at addr.
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil) // 修改为tls.Dial
	if err != nil {
		return nil, err
	}
	host := addr[:strings.Index(addr, ":")]
	return smtp.NewClient(conn, host)
}

func SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	c, err := Dial(addr)
	if err != nil {
		return err
	}

	if err := c.Hello(""); err != nil {
		return err
	}

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err = c.StartTLS(nil); err != nil {
			return err
		}
	}
	/*
		if a != nil && c.ext != nil {
			if _, ok := c.ext["AUTH"]; ok {
				if err = c.Auth(a); err != nil {
					return err
				}
			}
		}
	*/
	// 修改后
	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
