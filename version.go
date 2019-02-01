package provision

// All of these fields are injected as part of the
// build.sh script.

// RSPrePart is 'v' place holder
var RSPrePart = "v"

// RSMajorVersion is the first number of SemVer
var RSMajorVersion = "3"

// RSMinorVersion is the second number of SemVer
var RSMinorVersion = "10"

// RSPatchVersion is the third number of SemVer
var RSPatchVersion = "1000"

// RSExtra is a mostly free form field that must
// start with a dash, but an be anything else.
var RSExtra = "-pre-alpha"

// GitHash is the injected GIT Hash for the current commit
// for this built program.
var GitHash = "NotSet"

// BuildStamp is the time the build occurred.
var BuildStamp = "Not Set"

// RSVersion is the aggregated SemVer String
var RSVersion = RSPrePart + RSMajorVersion + "." + RSMinorVersion + "." + RSPatchVersion + RSExtra + "-" + GitHash
