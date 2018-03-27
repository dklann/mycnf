// Package mycnf provides access to MySQL configuration files (in .ini format).
// Based on code from https://gist.github.com/nickcarenza/d847ec24455e70a8609b6602ed528133
package mycnf

import (
	"errors"
	"fmt"
)

import (
	"github.com/go-ini/ini"
)

// MyConfer grants access to the myconf type.
type MyConfer interface {
	ReadMyCnf(profile string) (string, error)
}

// Mycnf contains the pieces needed to access a MySQL database.
type Mycnf struct {
	DbHost string
	DbPort string
	DbUser string
	DbPass string
	DbName string
}

// NewMyCnf returns a new, empty Mycnf structure.
func NewMyCnf() *Mycnf {
	return &Mycnf{
		DbHost: "",
		DbPort: "",
		DbUser: "",
		DbPass: "",
		DbName: "",
	}
}

// ReadMyCnf reads a .my.cnf section "profile", fills in missing values in the passed structure,
// and returns a DSN suitable for use in db.Open(). ReadMyCnf returns an error only if there
// was an actual error attempting to access the file. Not finding the file or the "profile"
// within the file is not an error.
func (c *Mycnf) ReadMyCnf(configFile *string, profile string) (string, error) {
	dbhost := "localhost" // default MySQL host
	dbport := "3306"      // default MySQL port
	dbname := ""
	dbuser := ""
	dbpass := ""

	if profile == "" {
		return "", errors.New("missing 'profile' name for .my.cnf")
	}

	// Load the contents of the config file at configFile.
	cfg, _ := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true, Insensitive: true}, *configFile)

	// Examine .my.cnf (or other configuration file) for named profile and key-values.
	if cfg != nil {
		for _, s := range cfg.Sections() {
			if s.Name() == profile {
				// Prefer setting "host" in this order: passed as arg, from .my.cnf, default (localhost)
				if c.DbHost != "" {
					dbhost = c.DbHost
				} else if s.Key("host").Value() != "" {
					dbhost = s.Key("host").Value()
				}
				// Prefer setting "port" in the same manner as "host".
				if c.DbPort != "" {
					dbport = c.DbPort
				} else if s.Key("port").Value() != "" {
					dbport = s.Key("port").Value()
				}
				// Nonexistent database name is an error.
				if c.DbName != "" {
					dbname = c.DbName
				} else if s.Key("dbname").Value() != "" {
					dbname = s.Key("dbname").Value()
				} else {
					return "", errors.New("missing database name, cannot continue")
				}
				// Nonexistent user name is an error.
				if c.DbUser != "" {
					dbuser = c.DbUser
				} else if s.Key("user").Value() != "" {
					dbuser = s.Key("user").Value()
				} else {
					return "", errors.New("missing database user name, cannot continue")
				}
				// Nonexistent password is an error.
				if c.DbPass != "" {
					dbpass = c.DbPass
				} else if s.Key("password").Value() != "" {
					dbpass = s.Key("password").Value()
				} else {
					return "", errors.New("missing database password, cannot continue")
				}
				return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbuser, dbpass, dbhost, dbport, dbname), nil
			}
		}
	}

	// Did not find a profile in .my.cnf, so use the values passed.
	// Prefer setting "host" in this order: passed as arg, from .my.cnf, default (localhost).
	if c.DbHost != "" {
		dbhost = c.DbHost
	}
	// Prefer setting "port" in the same order as "host".
	if c.DbPort != "" {
		dbport = c.DbPort
	}
	// Nonexistent database name is an error.
	if c.DbName == "" {
		return "", errors.New("missing database name, cannot continue")
	}
	dbname = c.DbName
	// Nonexistent user name is an error.
	if c.DbUser == "" {
		return "", errors.New("missing database user name, cannot continue")
	}
	dbuser = c.DbUser
	// Nonexistent password is an error.
	if c.DbPass == "" {
		return "", errors.New("missing database password, cannot continue")
	}
	dbpass = c.DbPass

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbuser, dbpass, dbhost, dbport, dbname), nil
}
