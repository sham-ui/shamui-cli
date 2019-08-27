package members

import (
	"database/sql"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"{{ shortName }}/config"
	"{{ shortName }}/models"
)

func uniqueEmail(email string) bool {
	models.ConnToDB(config.DataBase.GetURL())
	var existingEmail string
	row := models.Db.QueryRow("SELECT email FROM members WHERE email = $1", email)
	err := row.Scan(&existingEmail)
	if err != nil {
		log.Error(err)
	}
	if err == sql.ErrNoRows {
		log.Info("Email does not exist.")
		return true
	}

	log.Info("Email exists in the store: ", existingEmail)
	if len(existingEmail) > 0 {
		return false
	}
	return true
}

// emailAvailable returns true if the email is not taken
func emailAvailable(email string) bool {
	if uniqueEmail(email) == true {
		return true
	}
	return false
}

// passwordsMatch confirms that both passwords are matching.  This helps the user avoid typing the incorrect password
func passwordsMatch(pw1 string, pw2 string) bool {
	if pw1 == pw2 {
		return true
	}
	return false
}
