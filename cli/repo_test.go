package cli

import (
	"fmt"
	"path"
	"testing"
)

func TestRepos(t *testing.T) {
	fakes := []string{
		"fake-centos-install",
		"fake-debian-install",
		"fake-scientificlinux-install",
		"fake-ubuntu-install",
	}
	cliTest(false, false, "params", "create", "-").Stdin(`---
Name: "package-repositories"
Description: "Repositories to use to install packages from"
Schema:
  default:
    - tag: debian-6-install
      os:
        - debian-6
      arch: amd64
      url: https://this.url.is.fake/debian
      distribution: "6"
      installSource: true
      components:
        - main
        - contrib
        - non-free
    - tag: debian-6-security
      os:
        - debian-6
      arch: any
      url: https://this.url.is.secure/debian-security
      distribution: "6"
      securitySource: true
      components:
        - main
        - contrib
        - non-free
    - tag: sweet-debs
      os:
        - debian-6
        - ubuntu-6
      arch: any
      url: https://this.url.is.a.mirror/sweet-debs
      distribution: sweet
      components:
        - on-fire
    - tag: ubuntu-6-install
      os:
        - ubuntu-6
      arch: amd64
      url: https://this.url.is.fake/ubuntu
      distribution: "6"
      installSource: true
      components:
        - main
        - contrib
        - non-free
    - tag: ubuntu-6-security
      os:
        - ubuntu-6
      arch: any
      url: https://this.url.is.secure/ubuntu-security
      distribution: "6"
      securitySource: true
      components:
        - main
        - contrib
        - non-free
    - tag: centos-6-install
      os:
        - centos-6
      arch: x86_64
      url: https://this.url.is.fake/centos/6/os/x86_64
      installSource: true
    - tag: centos-6-security
      os:
        - centos-6
      arch: x86_64
      url: https://this.url.is.secure/centos/6/updates/x86_64
      securitySource: true
    - tag: centos-6-extras
      os:
        - centos-6
      arch: x86_64
      url: https://this.url.is.fake/centos
      distribution: "6"
      components:
        - extras
        - cluster-stuff
        - atomic
    - tag: scientificlinux-6-install
      os:
        - scientificlinux-6
      arch: x86_64
      url: https://this.url.is.fake/scientificlinux/6/x86_64/os
      installSource: true
    - tag: scientificlinux-6-security
      os:
        - scientificlinux-6
      arch: x86_64
      url: https://this.url.is.secure/scientificlinux/6/x86_64/updates
      securitySource: true
    - tag: scientificlinux-6-extras
      os:
        - scientificlinux-6
      arch: x86_64
      url: https://this.url.is.fake/scientificlinux
      distribution: "6"
      components:
        - extras
        - cluster-stuff
        - atomic
    - tag: epel-6
      os:
        - centos-6
        - scientificlinux-6
      arch: x86_64
      url: https://this.url.is.a.mirror/epel/7/$basearch
  type: "array"
  items:
    type: "object"
    required:
      - tag
      - os
      - url
    properties:
      tag:
        type: string
      os:
        type: array
        items:
          type: string
      arch:
        type: string
      url:
        type: string
        format: uri
      packageType:
        type: string
      repoType:
        type: string
      installSource:
        type: boolean
      securitySource:
        type: boolean
      bootloc:
        type: string
      distribution:
        type: string
      components:
        type: array
        items:
          type: string
`).run(t)
	cliTest(false, false, "machines", "create", "-").Stdin(`---
Name: phred
Uuid: c9196b77-deef-4c8e-8130-299b3e3d9a10`).run(t)
	// We have install sources for these repos, so they will be created and available.
	for _, fake := range fakes {
		cliTest(false, false, "bootenvs", "create", "../api/test-data/"+fake+".yml").run(t)
	}
	// All template expansion tests should refer to install sources from package-repositories.
	for _, fake := range fakes {
		cliTest(false, false,
			"machines", "update",
			"c9196b77-deef-4c8e-8130-299b3e3d9a10", "-").
			Stdin(fmt.Sprintf(`{"BootEnv": "%s"}`, fake)).run(t)
		cliTest(false, false, path.Join("machines", "c9196b77-deef-4c8e-8130-299b3e3d9a10", "url")).get(t)
		cliTest(false, false, path.Join("machines", "c9196b77-deef-4c8e-8130-299b3e3d9a10", "install")).get(t)
		cliTest(false, false, path.Join("machines", "c9196b77-deef-4c8e-8130-299b3e3d9a10", "lines")).get(t)
	}
	// This ISO provides fake local kernels, but no install repos.  It is not the most accurate
	// standin for a real install iso, but you get the idea.
	cliTest(false, false, "isos", "upload", "../api/test-data/fake-install.tgz", "as", "fake-install.tgz").run(t)
	// centos and scientificlinux should refer to local repos that are now "present",
	// debian and ubuntu will still refer to upstream.
	for _, fake := range fakes {
		cliTest(false, false,
			"machines", "update",
			"c9196b77-deef-4c8e-8130-299b3e3d9a10", "-").
			Stdin(fmt.Sprintf(`{"BootEnv": "%s"}`, fake)).run(t)
		cliTest(false, false, path.Join("machines", "c9196b77-deef-4c8e-8130-299b3e3d9a10", "url")).get(t)
		cliTest(false, false, path.Join("machines", "c9196b77-deef-4c8e-8130-299b3e3d9a10", "install")).get(t)
		cliTest(false, false, path.Join("machines", "c9196b77-deef-4c8e-8130-299b3e3d9a10", "lines")).get(t)
	}
	cliTest(false, false, "machines", "destroy", "c9196b77-deef-4c8e-8130-299b3e3d9a10").run(t)
	for _, fake := range fakes {
		cliTest(false, false, "bootenvs", "destroy", fake).run(t)
	}
	cliTest(false, false, "params", "destroy", "package-repositories").run(t)
	cliTest(false, false, "isos", "destroy", "fake-install.tgz")
	verifyClean(t)
}
