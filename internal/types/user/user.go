// Package user contains all the definitions for the user.
package user

// User contains the data for a user. This can be used for authentication
// with third party services, as it contains the necessary fields.
type User struct {
	Email string `json:"email"`

	// Including Email, these are the fields returned when querying the Google
	// API.
	VerifiedEmail bool   `json:"email_verified"`
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
}
