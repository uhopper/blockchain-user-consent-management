package model

type Consent struct {
	ID                string `json:"ID"`
	UserConsent       bool   `json:"consent"`
	PrivacyPolicyHash string `json:"privacyPolicyHash"`
	LastUpdate        int64  `json:"lastUpdate"`
}

type ConsentWritable struct {
	ID                string `json:"ID"`
	UserConsent       bool   `json:"consent"`
	PrivacyPolicyHash string `json:"privacyPolicyHash"`
}
