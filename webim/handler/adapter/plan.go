package adapter

import "time"

type Plan struct {
	ID               string    `json:"id"`
	CreatedOn        time.Time `json:"created_on"`
	ExpirationTime   time.Time `json:"expiration_time"`
	LoginAgentLimit  int       `json:"login_agent_limit"`
	PayAgentQuantity int       `json:"pay_agent_quantity"`
	PayAmount        int       `json:"pay_amount"`
	Plan             int       `json:"plan"`
	TrialStatus      int       `json:"trial_status"`
	Valid            bool      `json:"valid"`
	VisitorLimit     int       `json:"visitor_limit"`
}
