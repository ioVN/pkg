// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package database

import (
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
)

func ParseMongoURI(uri string) (dbURI, dbname string) {
	conn, err := connstring.ParseAndValidate(uri)
	if err != nil {
		log.Fatalf("- GetMongoURI error: %s\n", err.Error())
		return
	}
	if conn.UsernameSet && !conn.PasswordSet {
		arr := strings.Split(conn.Original, "@")
		if len(arr) > 1 {
			var pwd string
			print("Enter password of MongoDB: ")
			if _, err := fmt.Scan(&pwd); err != nil {
				log.Fatalf("%s\n", err.Error())
			}
			print("\b\rPassword Fingerprint:      " + strings.Repeat("*", len(pwd)) + "\n") // hide password input
			conn.PasswordSet = true
			conn.Password = pwd
			conn.Original = arr[0] + ":" + pwd + "@" + strings.Join(arr[1:], "@")
			if _, err := connstring.ParseAndValidate(conn.Original); err != nil {
				log.Fatalf("- GetMongoURI error: %s\n", "password invalid format")
				return
			}
		}
	}
	return conn.String(), conn.Database
}
