package models

import "time"

// LicenseBundle is a collection of Licenses.
type LicenseBundle struct {
	Contact           string
	ContactEmail      string
	ContactId         string
	Owner             string
	OwnerEmail        string
	OwnerId           string
	Grantor           string
	GrantorEmail      string
	Version           string
	GenerationVersion string
	Endpoints         []string `json:",omitempty"`
	// MaxDrpVersion is the highest major/minor version to allow.
	// For example, v4.3 would mean any v4.3.* release or previous.
	MaxDrpVersion string

	Licenses []License
}

// License is used to authorize acess to extended dr-provision
// features and functionality.
type License struct {
	Name                    string
	Version                 string
	Data                    interface{}
	PurchaseDate            time.Time
	StartDate               time.Time
	SoftExpireDate          time.Time
	HardExpireDate          time.Time
	ContractTerminationDate time.Time
	ShortLicense            string
	LongLicense             string
	Active, Expired         bool
}

func (l *License) Check(ref time.Time) (active, expired bool) {
	active = l.StartDate.Before(ref) && ref.Before(l.HardExpireDate)
	expired = l.SoftExpireDate.Before(ref)
	return
}
