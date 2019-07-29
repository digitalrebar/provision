package v4

var (
	RSPrePart      = "v"
	RSMajorVersion = "4"
	RSMinorVersion = "0"
	RSPatchVersion = "0"
	RSExtra        = "-pre"
	GitHash        = "NotSet"
	BuildStamp     = "Not Set"
	RSVersion      = RSPrePart + RSMajorVersion + "." + RSMinorVersion + "." + RSPatchVersion + RSExtra + "+" + GitHash
)
