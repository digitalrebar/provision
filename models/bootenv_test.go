package models

import (
	"testing"
)

func nfv(name, family, version string) OsInfo {
	return OsInfo{Name: name, Family: family, Version: version}
}

func teq(t *testing.T, param, a, b string) {
	t.Helper()
	if a == b {
		t.Logf("%s matched: `%s` == `%s`", param, a, b)
	} else {
		t.Errorf("ERROR: %s did not match: `%s` != `%s`", param, a, b)
	}
}

func teqv(t *testing.T, tgt OsInfo, famv string) {
	t.Helper()
	if tgt.VersionEq(famv) {
		t.Logf("%s version matched: `%s` == `%s`", tgt.Name, tgt.FamilyVersion(), famv)
	} else {
		t.Errorf("ERROR: %s version did not match: `%s` != `%s`", tgt.Name, tgt.FamilyVersion(), famv)
	}
}

func nfveq(t *testing.T, tgt OsInfo, family, version, famtype, famveq string) {
	t.Helper()
	teq(t, "FamilyName", tgt.FamilyName(), family)
	teq(t, "FamilyVersion", tgt.FamilyVersion(), version)
	teq(t, "FamilyType", tgt.FamilyType(), famtype)
	teqv(t, tgt, famveq)
}

func TestBootEnvOSStuff(t *testing.T) {
	cent6 := nfv("centos-6", "", "")
	rh7 := nfv("redhat-7.4.1705", "", "")
	ubuntu18 := nfv("ubuntu-18.04", "", "")
	deb9 := nfv("debian-9.4", "", "")
	cent63 := nfv("centos-6.3.1611", "", "")
	centweird := nfv("cewnt", "centos", "6.3.1611")
	nfveq(t, cent6, "centos", "6", "rhel", "6")
	nfveq(t, cent63, "centos", "6.3.1611", "rhel", "6")
	nfveq(t, cent63, "centos", "6.3.1611", "rhel", "6.3")
	nfveq(t, cent63, "centos", "6.3.1611", "rhel", "6.3.1611")
	nfveq(t, centweird, "centos", "6.3.1611", "rhel", "6")
	nfveq(t, centweird, "centos", "6.3.1611", "rhel", "6.3")
	nfveq(t, centweird, "centos", "6.3.1611", "rhel", "6.3.1611")
	nfveq(t, rh7, "redhat", "7.4.1705", "rhel", "7")
	nfveq(t, rh7, "redhat", "7.4.1705", "rhel", "7.4")
	nfveq(t, rh7, "redhat", "7.4.1705", "rhel", "7.4.1705")
	nfveq(t, ubuntu18, "ubuntu", "18.04", "debian", "18")
	nfveq(t, ubuntu18, "ubuntu", "18.04", "debian", "18.04")
	nfveq(t, deb9, "debian", "9.4", "debian", "9")
	nfveq(t, deb9, "debian", "9.4", "debian", "9.4")
}
