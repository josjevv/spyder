package spyder

import (
	"github.com/bulletind/spyder/log"

	"gopkg.in/mgo.v2"
	"strings"
	"time"
	"crypto/tls"
	"net"
)

func GetSession(connString string) *mgo.Session {
	log.Debug("Attempting connecting to:", connString)
	// quick hack to allow SSL based connections, may be removed in future when parseURL supports it
	// see also: https://github.com/go-mgo/mgo/issues/84
	const SSL_SUFFIX = "&ssl=true"
	useSsl := false

	if strings.HasSuffix(connString, SSL_SUFFIX) {
		connString = strings.TrimSuffix(connString, SSL_SUFFIX)
		useSsl = true
	}

	dialInfo, err := mgo.ParseURL(connString)
	if err != nil {
		panic(err)
	}

	dialInfo.Timeout = 10 * time.Second

	if useSsl {
		config := tls.Config{}
		config.InsecureSkipVerify = true

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &config)
		}
	}

	// get a mgo session
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}

	log.Info("spyder connected to db:", connString)
	session.SetMode(mgo.Monotonic, true)
	return session
}

type Progress struct {
	Path  string
	Ident []byte
}
