package investigations

import "context"

type Evidence struct {
	ID         string         `json:"id"`
	RecordType string         `json:"record_type"`
	RecordID   string         `json:"record_id"`
	Summary    string         `json:"summary"`
	Data       map[string]any `json:"data"`
}

type Claim struct {
	Text        string   `json:"text"`
	EvidenceIDs []string `json:"evidence_ids"`
}

type Report struct {
	Finding               string     `json:"finding"`
	Summary               string     `json:"summary"`
	Explanation           Claim      `json:"explanation"`
	ConfirmedFacts        []Claim    `json:"confirmed_facts"`
	Hypotheses            []Claim    `json:"hypotheses"`
	Confidence            float64    `json:"confidence"`
	ConfidenceExplanation string     `json:"confidence_explanation"`
	Recommendation        string     `json:"recommendation"`
	RequiresHumanApproval bool       `json:"requires_human_approval"`
	UnresolvedQuestions   []string   `json:"unresolved_questions"`
	Evidence              []Evidence `json:"evidence"`
	AnalysisSource        string     `json:"analysis_source"`
	ModelStatus           string     `json:"model_status"`
	ModelMessage          string     `json:"model_message,omitempty"`
}

type Snapshot struct {
	Order    OrderRecord     `json:"order"`
	Attempts []AttemptRecord `json:"attempts"`
	Events   []EventRecord   `json:"events"`
}

type OrderRecord struct {
	ID                string `json:"id"`
	MerchantReference string `json:"merchant_reference"`
	AmountMinor       int64  `json:"amount_minor"`
	Currency          string `json:"currency"`
	Description       string `json:"description"`
	Status            string `json:"status"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type AttemptRecord struct {
	ID              string  `json:"id"`
	PayPalOrderID   *string `json:"paypal_order_id"`
	PayPalCaptureID *string `json:"paypal_capture_id"`
	AmountMinor     int64   `json:"amount_minor"`
	Currency        string  `json:"currency"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type EventRecord struct {
	ID               string  `json:"id"`
	ProviderEventID  string  `json:"provider_event_id"`
	EventType        string  `json:"event_type"`
	PayPalOrderID    *string `json:"paypal_order_id"`
	PayPalCaptureID  *string `json:"paypal_capture_id"`
	ProcessingStatus string  `json:"processing_status"`
	ReceivedAt       string  `json:"received_at"`
}

type Store interface {
	Load(context.Context, string) (Snapshot, error)
}

type Model interface {
	Analyze(context.Context, string, Snapshot, []Evidence) (Report, error)
}
