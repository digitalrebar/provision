package provision

const RS_MAJOR_VERSION = "3"
const RS_MINOR_VERSION = "0"
const RS_PATCH_VERSION = "2"
const RS_EXTRA = "-pre-alpha"

var GitHash = "NotSet"
var BuildStamp = "Not Set"

var RS_VERSION = "v" + RS_MAJOR_VERSION + "." + RS_MINOR_VERSION + "." + RS_PATCH_VERSION + RS_EXTRA + "-" + GitHash
