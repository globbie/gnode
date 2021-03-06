package session

import (
	"errors"
	"crypto/rsa"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/text/language"
)

type ChatThread struct {
	ThreadId     string
}

type ChatSession struct {
	UserId      string
	ShardId     string
	UserAgent   string
	UserIP      string
	Langs       []language.Tag
	Roles       []string
	Threads     []ChatThread
}

type Claims struct {
	*jwt.StandardClaims
	UserId    string   `json:"uid,required"`
	ShardId   string   `json:"shard,required"`
	UserRoles []string `json:"roles,omitempty"`
}

func New(r *http.Request) (*ChatSession, error) {
	cs := ChatSession{
		UserAgent: r.UserAgent(),
	}
	ip, err := GetSessionIP(r)
	if err == nil {
		cs.UserIP = ip
	}
	cs.Langs, _, _ = language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	
	return &cs, nil
}

func BuildSessionCookie(name string, val string, domain string) (*http.Cookie, error) {
	cookie := http.Cookie{}
	cookie.Name = name
	cookie.Value = val
	cookie.Path = "/"
	if domain != "localhost" {
		cookie.Domain = domain
	}
	cookie.Expires = time.Now().Add(356 * 24 * time.Hour)
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	return &cookie, nil
}

func GetSessionIP(r *http.Request) (string, error) {
    ip := r.Header.Get("X-REAL-IP")
    netIP := net.ParseIP(ip)
    if netIP != nil {
        return ip, nil
    }
    ips := r.Header.Get("X-FORWARDED-FOR")
    splitIps := strings.Split(ips, ",")
    for _, ip := range splitIps {
        netIP := net.ParseIP(ip)
        if netIP != nil {
            return ip, nil
        }
    }
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return "", err
    }
    netIP = net.ParseIP(ip)
    if netIP != nil {
        return ip, nil
    }
    return "", errors.New("No valid IP found")
}

func IssueAccessToken(ses *ChatSession, signKey *rsa.PrivateKey, expiry int) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims = &Claims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(expiry)).Unix(),
		},
		ShardId:   ses.ShardId,
		UserId:    ses.UserId,
		UserRoles: ses.Roles,
	}
	return token.SignedString(signKey)
}

/*
A link to activate your account has been emailed to the address provided.
*/
