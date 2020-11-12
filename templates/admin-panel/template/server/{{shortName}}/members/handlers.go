//Package members is for member CRUD
package members

import (
	"encoding/json"
	"github.com/antonlindstrom/pgstore"
	_ "github.com/lib/pq" // github.com/lib/pq
	log "github.com/sirupsen/logrus"
	"net/http"
	"{{ shortName }}/config"
	"{{ shortName }}/models"
	"{{ shortName }}/sessions"
	"strconv"
)

type memberOutput struct {
	Status string
	Errors []string
}

type resDetails struct {
	Status   string
	Messages []string
}

type member struct {
	NewName, NewEmail1, NewEmail2 string
}

// SignupMember creates a single member
func SignupMember(w http.ResponseWriter, r *http.Request) {
	var memberValid = true
	var m models.NewMember
	var signupErrs []string
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		log.WithError(err).Warn("Error decoding new member")
		log.Warn("Bad New Member data provided: ", r.Body)
		signupErrs = append(signupErrs, "Error processing new member data.")
	}

	if m.Name == "" {
		signupErrs = append(signupErrs, "Name must not be empty.")
		// json.NewEncoder(w).Encode("Name must not be empty.")
		memberValid = false
	}

	if emailAvailable(m.Email) != true {
		signupErrs = append(signupErrs, "Email is already in use.")
		// json.NewEncoder(w).Encode("Email is already in use.")
		memberValid = false
	}

	if passwordsMatch(m.Password, m.Password2) != true {
		signupErrs = append(signupErrs, "Passwords do not match.")
		// json.NewEncoder(w).Encode("Passwords do not match.")
		memberValid = false
	}

	if memberValid == true {
		msg := memberOutput{
			Status: "Member Created",
			Errors: signupErrs,
		}
		err := models.CreateMember(&m)
		if err != nil {
			log.Error(err)
		}
		json.NewEncoder(w).Encode(msg)
		log.Info("User Created", m.Email, m.Name)
	} else {
		log.Info("Error creating member.")
		w.WriteHeader(400)
		msg := memberOutput{
			Status: "Member Not Created",
			Errors: signupErrs,
		}
		json.NewEncoder(w).Encode(msg)
	}
	log.Info("User data supplied:", m)
}

// UpdateMemberEmail allows the user to update member information and returns an error or the newly made member name
func UpdateMemberEmail(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	if session == nil {
		msg := resDetails{
			Status:   "Expired session or cookie",
			Messages: []string{"Session Expired.  Log out and log back in."},
		}
		json.NewEncoder(w).Encode(msg)
		return
	}

	var msg resDetails

	var memberUpdate member
	err := json.NewDecoder(r.Body).Decode(&memberUpdate)
	if err != nil {
		log.WithError(err).Warn("Error decoding body")
	}
	// Check for bad email length
	if len(memberUpdate.NewEmail1) <= 0 || len(memberUpdate.NewEmail2) <= 0 {
		msg := resDetails{
			Status:   "Bad Email",
			Messages: append(msg.Messages, "Email must have more than 0 characters."),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
		return
	}
	log.Info("New Email: ", memberUpdate.NewEmail1)
	if len(memberUpdate.NewEmail1) >= 1 || len(memberUpdate.NewEmail2) >= 1 {
		if memberUpdate.NewEmail1 != memberUpdate.NewEmail2 {
			msg := resDetails{
				Status:   "Bad Email",
				Messages: append(msg.Messages, "Emails don't match."),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(msg)
			return
		} else if models.UpdateMemberEmail(session.ID, memberUpdate.NewEmail1) == true {
			msg.Messages = append(msg.Messages, memberUpdate.NewEmail1)
			msg.Status = "OK"
			store, err := pgstore.NewPGStore(config.DataBase.GetURL(), []byte(config.Session.Secret))
			if err != nil {
				log.Println(err)
			}
			defer store.Close()
			session, err := store.Get(r, "{{ shortName }}-session")
			session.Values["email"] = memberUpdate.NewEmail1
			if err = session.Save(r, w); err != nil {
				log.Printf("Error saving session: %v", err)
			}
			json.NewEncoder(w).Encode(msg)
		} else {
			msg := resDetails{
				Status:   "Fail update email",
				Messages: append(msg.Messages, "Fail update email"),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(msg)
			return
		}
	}

	log.Info("Member's ID: ", session.ID)
	log.Info(msg)
}

// UpdateMemberName will update the existing member name for authorized sessions
func UpdateMemberName(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	if session == nil {
		msg := resDetails{
			Status:   "Expired session or cookie",
			Messages: []string{"Session Expired.  Log out and log back in."},
		}
		json.NewEncoder(w).Encode(msg)
		return
	}

	var msg resDetails

	var memberUpdate member
	err := json.NewDecoder(r.Body).Decode(&memberUpdate)
	if err != nil {
		log.WithError(err).Warn("Error decoding body")
	}
	// Check for bad name length
	if len(memberUpdate.NewName) < 1 {
		msg := resDetails{
			Status:   "Bad Name",
			Messages: append(msg.Messages, "Name must have more than 0 characters."),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
		return
	}
	log.Println("New Name: ", memberUpdate.NewName)
	if len(memberUpdate.NewName) >= 1 {
		if models.UpdateMemberName(session.ID, memberUpdate.NewName) == true {
			msg.Messages = append(msg.Messages, memberUpdate.NewName)
			msg.Status = "OK"
			store, err := pgstore.NewPGStore(config.DataBase.GetURL(), []byte(config.Session.Secret))
			if err != nil {
				log.Println(err)
			}
			defer store.Close()
			session, err := store.Get(r, "{{ shortName }}-session")
			session.Values["name"] = memberUpdate.NewName
			if err = session.Save(r, w); err != nil {
				log.Printf("Error saving session: %v", err)
			}
			json.NewEncoder(w).Encode(msg)
		} else {
			msg := resDetails{
				Status:   "Fail update name",
				Messages: append(msg.Messages, "Server error"),
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(msg)
			return
		}
	}

	log.Info("Member's ID: ", session.ID)
	log.Info(msg)
}

func Members(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	if session == nil {
		msg := resDetails{
			Status:   "Expired session or cookie",
			Messages: []string{"Session Expired.  Log out and log back in."},
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	if !session.IsSuperuser {
		msg := resDetails{
			Status:   "Only superuser can get member list",
			Messages: []string{"Only superuser can get member list"},
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(msg)
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if nil != err {
		offset = 0
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if nil != err {
		limit = 20
	}
	res := map[string]interface{}{
		"meta": map[string]int{
			"offset": offset,
			"limit":  limit,
			"total":  models.GetMembersCount(),
		},
		"members": models.GetMembers(offset, limit),
	}
	json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusOK)
}
