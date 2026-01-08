package cap

import (
	"encoding/xml"
	"time"

	"github.com/google/uuid"
)

// CAPParams represents parameters for building a CAP alert
type CAPParams struct {
	Sender      string
	Event       string
	Urgency     string
	Severity    string
	Certainty   string
	Area        string
	Instruction string
}

// capAlert represents the CAP 1.2 alert structure
type capAlert struct {
	XMLName    xml.Name `xml:"alert"`
	Xmlns      string   `xml:"xmlns,attr,omitempty"`
	Identifier string   `xml:"identifier"`
	Sender     string   `xml:"sender"`
	Sent       string   `xml:"sent"`
	Status     string   `xml:"status"`
	MsgType    string   `xml:"msgType"`
	Scope      string   `xml:"scope"`
	Info       capInfo  `xml:"info"`
}

type capInfo struct {
	Category    string  `xml:"category"`
	Event       string  `xml:"event"`
	Urgency     string  `xml:"urgency"`
	Severity    string  `xml:"severity"`
	Certainty   string  `xml:"certainty"`
	Instruction string  `xml:"instruction"`
	Area        capArea `xml:"area"`
}

type capArea struct {
	AreaDesc string `xml:"areaDesc"`
}

// BuildCAPXML builds a CAP 1.2 compliant XML alert
func BuildCAPXML(params CAPParams) string {
	alert := capAlert{
		Xmlns:      "urn:oasis:names:tc:emergency:cap:1.2",
		Identifier: uuid.NewString(),
		Sender:     params.Sender,
		Sent:       time.Now().UTC().Format(time.RFC3339),
		Status:     "Actual",
		MsgType:    "Alert",
		Scope:      "Public",
		Info: capInfo{
			Category:    "Safety",
			Event:       params.Event,
			Urgency:     params.Urgency,
			Severity:    params.Severity,
			Certainty:   params.Certainty,
			Instruction: params.Instruction,
			Area:        capArea{AreaDesc: params.Area},
		},
	}

	output, err := xml.MarshalIndent(alert, "", "  ")
	if err != nil {
		return ""
	}

	return xml.Header + string(output)
}

// ValidateUrgency checks if urgency value is valid CAP urgency
func ValidateUrgency(urgency string) bool {
	valid := []string{"Immediate", "Expected", "Future", "Past", "Unknown"}
	for _, v := range valid {
		if v == urgency {
			return true
		}
	}
	return false
}

// ValidateSeverity checks if severity value is valid CAP severity
func ValidateSeverity(severity string) bool {
	valid := []string{"Extreme", "Severe", "Moderate", "Minor", "Unknown"}
	for _, v := range valid {
		if v == severity {
			return true
		}
	}
	return false
}

// ValidateCertainty checks if certainty value is valid CAP certainty
func ValidateCertainty(certainty string) bool {
	valid := []string{"Observed", "Likely", "Possible", "Unlikely", "Unknown"}
	for _, v := range valid {
		if v == certainty {
			return true
		}
	}
	return false
}
