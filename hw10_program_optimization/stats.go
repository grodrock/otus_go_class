package hw10programoptimization

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return getStats(r, domain)
}

func (u *User) getDomain() string {
	return strings.ToLower(strings.SplitN(u.Email, "@", 2)[1])
}

func getStats(r io.Reader, domain string) (DomainStat, error) {
	ds := make(DomainStat)
	var user User

	br := bufio.NewReader(r)
	var lastStr bool

	for {
		lineB, err := br.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				lastStr = true
			} else {
				return nil, err
			}
		}

		if !bytes.Contains(lineB, []byte(domain)) && !lastStr {
			continue
		}
		if err := user.UnmarshalJSON(lineB); err != nil {
			return ds, err
		}

		if strings.HasSuffix(user.Email, "."+domain) {
			ds[user.getDomain()]++
		}
		if lastStr {
			break
		}
	}
	return ds, nil
}
