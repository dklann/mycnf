// mycnf provides access to MySQL configuration files (in .ini format).
// Based on code from https://gist.github.com/nickcarenza/d847ec24455e70a8609b6602ed528133
package mycnf

import (
	"errors"
	"fmt"

	"github.com/go-ini/ini"
)

// ReadMyCnf reads a .my.cnf section "profile", fills in missing values in the passed structure,
// and returns a DSN suitable for use in db.Open(). ReadMyCnf returns an error only if there
// was an actual error attempting to access the file. Not finding the file or the "profile"
// within the file is not an error.
func ReadMyCnf(configFile *string, profile string) (map[string]string, error) {
	dbhost := "localhost" // default MySQL host
	dbport := "3306"      // default MySQL port
	var (
		dbname string
		dbuser string
		dbpass string
	)

	if profile == "" {
		return nil, fmt.Errorf("missing 'profile' name for %s", *configFile)
	}

	// Load the contents of the config file at configFile
	// as well as the two "system" configuration files.
	cfg, _ := ini.LoadSources(
		ini.LoadOptions{
			AllowBooleanKeys: true,
			Insensitive:      true,
			Loose:            true,
		},
		"/etc/mysql/my.cnf",
		"/etc/my.cnf",
		*configFile)

	// Examine .my.cnf (or other configuration file) for named profile and key-values.
	if cfg != nil {
		for _, s := range cfg.Sections() {
			if s.Name() == profile {
				if s.Key("host").Value() != "" {
					dbhost = s.Key("host").Value()
				}
				if s.Key("port").Value() != "" {
					dbport = s.Key("port").Value()
				}
				if s.Key("database").Value() != "" {
					dbname = s.Key("database").Value()
				}
				if s.Key("user").Value() != "" {
					dbuser = s.Key("user").Value()
				}
				if s.Key("password").Value() != "" {
					dbpass = s.Key("password").Value()
				}
				confMap := map[string]string{
					"host":     dbhost,
					"port":     dbport,
					"database": dbname,
					"user":     dbuser,
					"password": dbpass,
				}
				return confMap, nil
			}
		}
	}

	return nil, errors.New("unable to find section '" + profile + "' in " + *configFile)
}
