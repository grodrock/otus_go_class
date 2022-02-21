package hw10programoptimization

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
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

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	var user User
	for i, line := range lines {

		if err = user.UnmarshalJSON([]byte(line)); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if !strings.HasSuffix(user.Email, "."+domain) {
			continue
		}
		result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
	}
	return result, nil
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
