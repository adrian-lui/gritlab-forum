package server

import (
	"errors"
	logger "gritface/log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

var sessionManager SessionManager

// Check for valid session, if not create a new one. Return session user data
func (sm *SessionManager) checkSession(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := sm.isSessionSet(w, r)
	if err != nil {
		// Session is not set, create a new one

		ID := uuid.New().String() // Create a session id

		// Create session cookie
		cookie = &http.Cookie{
			Name:  "session",
			Value: ID,
			Path:  "/",
			// Secure:   true,
			// HttpOnly: true,
			MaxAge: 3600,
		}

		// Send cookie to client
		http.SetCookie(w, cookie)

		// Store cookie in session as uid=0 (unregistered user)
		var thisSessionData = &SessionData{
			UId:  "0",
			Misc: make(map[string]interface{}),
		}
		sm.sessions[ID] = thisSessionData
	}

	return sm.sessions[cookie.Value].UId, nil
}

// This function check if a valid session is alive, returns session ID if alive
func (sm *SessionManager) isSessionSet(w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	c, err := r.Cookie("session")
	if err != nil {
		return nil, err
	}

	_, ok := sessionManager.sessions[c.Value] // Try to get session cookie value which will tell us if a valid session is open

	if ok {
		return c, nil
	} else {
		return nil, errors.New("Cookie value not alive")
	}
}

// function set session UID
func (sm *SessionManager) setSessionUID(uid int, w http.ResponseWriter, r *http.Request) error {

	thisSession, err := sm.isSessionSet(w, r)

	if err != nil {
		// Something wrong with cookie, return error
		return errors.New("Could not retrieve cookie data")
	}

	sm.sessions[thisSession.Value].UId = strconv.Itoa(uid)

	return nil
}

// functio delete session
func (sm *SessionManager) deleteSession(w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	// check if user is logged in
	cookie, err := sm.isSessionSet(w, r)
	if err != nil {
		return nil, err
	}

	// delete the session
	delete(sessionManager.sessions, cookie.Value)

	// remove the cookie
	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	return cookie, nil
}

func (sm *SessionManager) StoreSessionVariable(w http.ResponseWriter, r *http.Request, varName string, varValue interface{}) error {
	thisSession, err := sm.isSessionSet(w, r)
	if err != nil {
		// Something wrong with cookie, return error
		return errors.New("Something is wrong with the session")
	}

	sm.sessions[thisSession.Value].Misc[varName] = varValue

	return nil
}

func (sm *SessionManager) GetSessionVariable(w http.ResponseWriter, r *http.Request, varName string) (interface{}, error) {
	thisSession, err := sm.isSessionSet(w, r)
	if err != nil {
		// Something wrong with cookie, return error
		logger.WTL("Errors while opening session, "+err.Error(), false)
		return nil, errors.New("Something is wrong with the session")
	}

	if len(sm.sessions[thisSession.Value].Misc) < 1 {
		return nil, errors.New("Value not set")
	}

	return sm.sessions[thisSession.Value].Misc[varName], nil
}

func init() {
	sessionManager = SessionManager{
		sessions: make(map[string]*SessionData),
	}
}
