package app

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
	"{{shortName}}/config"
	"{{shortName}}/core/database"
	"{{shortName}}/members"
)

func CreateSuperUser() {
	config.LoadConfiguration("config.cfg")
	db, err := database.ConnToDB(config.DataBase.GetURL())
	if nil != err {
		log.Fatalf("Fail connect to db: %s", err)
	}
	err = members.CreateMemberStructure(db)
	if nil != err {
		log.Fatalf("Fail create members table: %s", err)
	} else {
		log.Info("Create members table")
	}

	var email string
	err = survey.AskOne(&survey.Input{
		Message: "Email:",
	}, &email, nil)
	if nil != err {
		log.WithError(err).Fatal("can't get email")
	}
	var name string
	err = survey.AskOne(&survey.Input{
		Message: "Name:",
	}, &name, nil)
	if nil != err {
		log.WithError(err).Fatal("can't get name")
	}
	var password string
	err = survey.AskOne(&survey.Password{
		Message: "Password:",
	}, &password, nil)
	if nil != err {
		log.WithError(err).Fatal("can't get password")
	}

	err = members.CreateMember(db, &members.MemberData{
		Name:        name,
		Email:       email,
		Password:    password,
		IsSuperuser: true,
	})
	if nil == err {
		log.Info("Superuser created: ", name)
	} else {
		log.WithError(err).Fatal("can't create superuser")
	}
}
