package sapase

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	// DRIVER: tds
	"github.com/thda/tds"

	"github.com/xo/usql/drivers"
)

var endRE = regexp.MustCompile(`;?\s*$`)

func init() {
	drivers.Register("tds", drivers.Driver{
		AllowMultilineComments:  true,
		RequirePreviousPassword: true,
		LexerName:               "tsql",
		Version: func(db drivers.DB) (string, error) {
			var ver string
			err := db.QueryRow(`SELECT @@version`).Scan(&ver)
			if err != nil {
				return "", err
			}
			return ver, nil
		},
		ChangePassword: func(db drivers.DB, user, newpw, oldpw string) error {
			if user != "" {
				return errors.New("Cannot change password for another user")
			}
			_, err := db.Exec(`exec sp_password '` + oldpw + `', '` + newpw + `'`)
			return err
		},
		Err: func(err error) (string, string) {
			if e, ok := err.(tds.SybError); ok {
				return strconv.Itoa(int(e.MsgNumber)), e.Message
			}

			msg := err.Error()
			if i := strings.LastIndex(msg, "tds:"); i != -1 {
				msg = msg[i:]
			}

			return "", msg
		},
		IsPasswordErr: func(err error) bool {
			return strings.Contains(err.Error(), "Login failed")
		},
		Process: func(prefix string, sqlstr string) (string, string, bool, error) {
			sqlstr = endRE.ReplaceAllString(sqlstr, "")
			typ, q := drivers.QueryExecType(prefix, sqlstr)
			return typ, sqlstr, q, nil
		},
	})
}
