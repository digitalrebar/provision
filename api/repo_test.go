package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/digitalrebar/provision/models"
)

func testRepoGen(t *testing.T, env *models.BootEnv, items map[string]string) {
	fakeMachine := &models.Machine{Name: "phred", BootEnv: env.Name}
	fakeMachine.Fill()
	rt(t, fmt.Sprintf("Create machine to test %s", env.Name), nil, nil,
		func() (interface{}, error) {
			return nil, session.CreateModel(fakeMachine)
		}, nil)
	if !fakeMachine.Available || fakeMachine.BootEnv != env.Name {
		t.Fatalf("Machine was not set to bootenv %s, it has %s", env.Name, fakeMachine.BootEnv)
	}
	if fakeMachine.OS != env.OS.Name {
		t.Errorf("Machine OS was not automatically set to %s", env.OS.Name)
	}
	for k, v := range items {
		rt(t, fmt.Sprintf("Check to see if %s rendered properly", k), v, nil,
			func() (interface{}, error) {
				urlPath := fmt.Sprintf("http://127.0.0.1:10002/machines/%s/%s", fakeMachine.Key(), k)
				t.Logf("Testing URL: %s", urlPath)
				resp, err := http.Get(urlPath)
				if err != nil {
					t.Errorf("Error fetching reponse: %v", err)
					return nil, err
				}
				defer resp.Body.Close()
				buf := &bytes.Buffer{}
				io.Copy(buf, resp.Body)
				return buf.String(), nil
			}, nil)
	}
	rt(t, "Delete machine", nil, nil,
		func() (interface{}, error) {
			_, err := session.DeleteModel("machines", fakeMachine.Key())
			return nil, err
		}, nil)
}

func TestFakeRepoManagement(t *testing.T) {
	packageRepoParam := mustDecode(&models.Param{}, `
Name: "package-repositories"
Description: "Repositories to use to install packages from"
Schema:
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
      distribution:
        type: string
      components:
        type: array
        items:
          type: string
`).(*models.Param)
	type ans struct {
		uri     string
		repos   []string
		install string
		lines   string
	}
	type fakeInstall struct {
		env         *models.BootEnv
		repos       []string
		local, repo map[string]string
	}
	fakeInstalls := map[string]*fakeInstall{
		"fake-debian-install": {
			repos: []string{"debian-6-install", "debian-6-security", "debian-6"},
			local: map[string]string{
				"url": "http://127.0.0.1:10002/debian-6/install",
				"install": `d-i mirror/protocol string http
d-i mirror/http/hostname string 127.0.0.1:10002
d-i mirror/http/directory string /debian-6/install

`,
			},
			repo: map[string]string{
				"url": "https://this.url.is.fake/debian",
				"install": `d-i mirror/protocol string https
d-i mirror/http/hostname string this.url.is.fake
d-i mirror/http/directory string /debian

d-i apt-setup/security_host string this.url.is.secure
d-i apt-setup/security_path string /debian-security

`,
				"lines": `deb https://this.url.is.fake/debian 6 main contrib non-free
deb https://this.url.is.secure/debian-security 6 main contrib non-free
deb https://this.url.is.a.mirror/sweet-debs sweet on-fire

`,
			},
		},
		"fake-ubuntu-install": {
			repos: []string{"ubuntu-6-install", "ubuntu-6-security", "ubuntu-6"},
			local: map[string]string{
				"url": "http://127.0.0.1:10002/ubuntu-6/install",
				"install": `d-i mirror/protocol string http
d-i mirror/http/hostname string 127.0.0.1:10002
d-i mirror/http/directory string /ubuntu-6/install

`,
			},
			repo: map[string]string{
				"url": "https://this.url.is.fake/ubuntu",
				"install": `d-i mirror/protocol string https
d-i mirror/http/hostname string this.url.is.fake
d-i mirror/http/directory string /ubuntu

d-i apt-setup/security_host string this.url.is.secure
d-i apt-setup/security_path string /ubuntu-security

`,
				"lines": `deb https://this.url.is.a.mirror/sweet-debs sweet on-fire
deb https://this.url.is.fake/ubuntu 6 main contrib non-free
deb https://this.url.is.secure/ubuntu-security 6 main contrib non-free

`,
			},
		},
		"fake-centos-install": {
			repos: []string{"centos-6-install", "centos-6-security", "centos-6"},
			local: map[string]string{
				"url": "http://127.0.0.1:10002/centos-6/install",
				"install": `install
url --url http://127.0.0.1:10002/centos-6/install
repo --name="fake-centos-install" --baseurl=http://127.0.0.1:10002/centos-6/install --cost=100
`,
			},
			repo: map[string]string{
				"url": "https://this.url.is.fake/centos/6/os/x86_64",
				"install": `install
url --url https://this.url.is.fake/centos/6/os/x86_64
repo --name="centos-6-install" --baseurl=https://this.url.is.fake/centos/6/os/x86_64 --cost=100
repo --name="centos-6-security" --baseurl=https://this.url.is.secure/centos/6/updates/x86_64 --cost=100
`,
				"lines": `
[centos-6-install]
name=centos-6 - centos-6-install
baseurl=https://this.url.is.fake/centos/6/os/x86_64
gpgcheck=0

[centos-6-security]
name=centos-6 - centos-6-security
baseurl=https://this.url.is.secure/centos/6/updates/x86_64
gpgcheck=0

[centos-6-extras-extras]
name=centos-6-extras - extras
baseurl=https://this.url.is.fake/centos/6/extras/$basearch
gpgcheck=0

[centos-6-extras-cluster-stuff]
name=centos-6-extras - cluster-stuff
baseurl=https://this.url.is.fake/centos/6/cluster-stuff/$basearch
gpgcheck=0

[centos-6-extras-atomic]
name=centos-6-extras - atomic
baseurl=https://this.url.is.fake/centos/6/atomic/$basearch
gpgcheck=0

[epel-6]
name=centos-6 - epel-6
baseurl=https://this.url.is.a.mirror/epel/7/$basearch
gpgcheck=0

`,
			},
		},
		"fake-scientificlinux-install": {
			repos: []string{"scientificlinux-6-install", "scientificlinux-6-security", "scientificlinux-6"},
			local: map[string]string{
				"url": "http://127.0.0.1:10002/scientificlinux-6/install",
				"install": `install
url --url http://127.0.0.1:10002/scientificlinux-6/install
repo --name="fake-scientificlinux-install" --baseurl=http://127.0.0.1:10002/scientificlinux-6/install --cost=100
`,
			},
			repo: map[string]string{
				"url": "https://this.url.is.fake/scientificlinux/6/x86_64/os",
				"install": `install
url --url https://this.url.is.fake/scientificlinux/6/x86_64/os
repo --name="scientificlinux-6-install" --baseurl=https://this.url.is.fake/scientificlinux/6/x86_64/os --cost=100
repo --name="scientificlinux-6-security" --baseurl=https://this.url.is.secure/scientificlinux/6/x86_64/updates --cost=100
`,
				"lines": `
[scientificlinux-6-install]
name=scientificlinux-6 - scientificlinux-6-install
baseurl=https://this.url.is.fake/scientificlinux/6/x86_64/os
gpgcheck=0

[scientificlinux-6-security]
name=scientificlinux-6 - scientificlinux-6-security
baseurl=https://this.url.is.secure/scientificlinux/6/x86_64/updates
gpgcheck=0

[scientificlinux-6-extras-extras]
name=scientificlinux-6-extras - extras
baseurl=https://this.url.is.fake/scientificlinux/6/$basearch/extras
gpgcheck=0

[scientificlinux-6-extras-cluster-stuff]
name=scientificlinux-6-extras - cluster-stuff
baseurl=https://this.url.is.fake/scientificlinux/6/$basearch/cluster-stuff
gpgcheck=0

[scientificlinux-6-extras-atomic]
name=scientificlinux-6-extras - atomic
baseurl=https://this.url.is.fake/scientificlinux/6/$basearch/atomic
gpgcheck=0

[epel-6]
name=scientificlinux-6 - epel-6
baseurl=https://this.url.is.a.mirror/epel/7/$basearch
gpgcheck=0

`,
			},
		},
	}
	rt(t, "Install fake bootenvs", nil, nil,
		func() (interface{}, error) {
			for name, env := range fakeInstalls {
				tgt := &models.BootEnv{}
				var err error
				envFile := fmt.Sprintf("test-data/%s.yml", name)
				tgt, err = session.InstallBootEnvFromFile(envFile)
				if err != nil {
					return nil, err
				}
				err = session.InstallISOForBootenv(tgt, "test-data/fake-install.tgz", true)
				if err != nil {
					return nil, err
				}
				env.env = tgt
			}
			return nil, nil
		}, nil)
	rt(t, "Verify that the fake bootenvs are available", nil, nil,
		func() (interface{}, error) {
			time.Sleep(time.Second * 5)
			for name, td := range fakeInstalls {
				if err := session.FillModel(td.env, name); err != nil {
					return nil, err
				} else if td.env.HasError() != nil {
					return nil, td.env.HasError()
				}
			}
			return nil, nil
		}, nil)
	t.Logf("Testing repo items when operating purely locally")
	for _, td := range fakeInstalls {
		testRepoGen(t, td.env, td.local)
	}
	rt(t, "Create the package-repositories parameter", nil, nil,
		func() (interface{}, error) {
			return nil, session.CreateModel(packageRepoParam)
		}, nil)
	rt(t, "Create an external repo definition for our fake repo", nil, nil,
		func() (interface{}, error) {
			obj, err := session.GetModel("profiles", "global")
			if err != nil {
				return nil, err
			}
			prof := obj.(*models.Profile)
			ref := []interface{}{}
			if err := DecodeYaml([]byte(`
- tag: debian-6-install
  os:
    - debian-6
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
  url: https://this.url.is.a.mirror/sweet-debs
  distribution: sweet
  components:
    - on-fire
- tag: ubuntu-6-install
  os:
    - ubuntu-6
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
  url: https://this.url.is.fake/centos/6/os/x86_64
  installSource: true
- tag: centos-6-security
  os:
    - centos-6
  url: https://this.url.is.secure/centos/6/updates/x86_64
  securitySource: true
- tag: centos-6-extras
  os:
    - centos-6
  url: https://this.url.is.fake/centos
  distribution: "6"
  components:
    - extras
    - cluster-stuff
    - atomic
- tag: scientificlinux-6-install
  os:
    - scientificlinux-6
  url: https://this.url.is.fake/scientificlinux/6/x86_64/os
  installSource: true
- tag: scientificlinux-6-security
  os:
    - scientificlinux-6
  url: https://this.url.is.secure/scientificlinux/6/x86_64/updates
  securitySource: true
- tag: scientificlinux-6-extras
  os:
    - scientificlinux-6
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
  url: https://this.url.is.a.mirror/epel/7/$basearch
`), &ref); err != nil {
				return nil, err
			}
			prof.Params["package-repositories"] = ref
			if err := session.PutModel(prof); err != nil {
				return nil, err
			}
			return nil, nil
		}, nil)
	t.Logf("Testing repo items with package-repositories set")
	for _, td := range fakeInstalls {
		testRepoGen(t, td.env, td.repo)
	}
}
