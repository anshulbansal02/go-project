package user

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Secret string `json:"secret,omitempty"`
}

func (m User) Public() User {
	m.Secret = ""
	return m
}
