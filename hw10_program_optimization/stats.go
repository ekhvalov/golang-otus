package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/json-iterator/go"
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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	s := bufio.NewScanner(r)
	i := 0
	for s.Scan() {
		if err = s.Err(); err != nil {
			return
		}
		if err = jsoniter.Unmarshal(s.Bytes(), &result[i]); err != nil {
			return
		}
		i++
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	dotDomain := "." + strings.ToLower(domain)
	for _, user := range u {
		email := strings.ToLower(user.Email)
		if strings.HasSuffix(email, dotDomain) {
			result[strings.SplitN(email, "@", 2)[1]]++
		}
	}
	return result, nil
}
