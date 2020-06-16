
This Makefile builds various distro packages of:

  dr-provision
  drpcli
  service startup files

The Makefile supports several targets for building packages.  Current
package build support is limited to the FPM tool support.

See FPM website(s) for more details:
  https://fpm.readthedocs.io/en/latest/
  https://github.com/jordansissel/fpm


The following primary make targets are supported:

  make                 # builds dr-provision, drpcli, services, stage-outputs
  make clean           # removes intermediary built packages
  make clean-pkgs      # nukes the pkgs/ staged directory
  make clean-all       # returns to factory default
  make dr-provision    # builds dr-provision packages
  make drpcli          # builds seprate drpcli packages
  make services        # builds start unit files for dr-provision
  make stage-outputs   # moves intermediary packages to pkgs/<TYPE>/ dir

If you run an individual target (eg 'dr-provision'), you also will want to
run 'stage-outputs' optionally.

A 'setup-fpm' target exists which attempts to setup the FPM packages and
gems necessary for the resulting builds.  This has only been tested on
CentOS 7 at this point.
