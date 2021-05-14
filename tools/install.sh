#!/usr/bin/env bash

set -e

# BUMP version on updates
VERSION="v21.05.10-2"

DEFAULT_DRP_VERSION=${DEFAULT_DRP_VERSION:-"stable"}

exit_cleanup() {
  local _x=$1
  shift
  [[ -n "$*" ]] && echo "EXIT MESSAGE: $*"
  exit $_x
}

usage() {
cat <<EOFUSAGE
  USAGE: $0 [ install | upgrade | remove | version ] [ <options: see below> ]

WARNING: 'install' option will OVERWRITE existing installations

OPTIONS:
    --debug=[true|false]    # Enables debug output
    --force=[true|false]    # Forces an overwrite of local install binaries and content
    --upgrade=[true|false]  # Turns on 'force' option to overwrite local binaries/content
    --isolated              # Sets up current directory as install location for drpcli
                            # and dr-provision (makes mess in current directory!)
    --no-content            # Don't add content to the system
    --zip-file=filename.zip # Don't download the dr-provision.zip file, instead use
                            # the referenced zip file (useful for airgap deployments)
                            # NOTE: disables sha256sum checks - do this manually
    --ipaddr=<ip>           # The IP to use for the system identified IP.  The system
                            # will attempt to discover the value if not specified
    --version=<string>      # Version identifier if downloading; stable, tip, or
                            # specific version label, defaults to: $DEFAULT_DRP_VERSION
    --remove-data           # Remove data as well as program pieces
    --skip-run-check        # Skip the process check for 'dr-provision' on new install
                            # only valid in '--isolated' install mode
    --skip-prereqs          # Skip OS dependency checks, for testing 'isolated' mode
    --no-sudo               # Do not use "sudo" prefix on commands (assume you're root)
    --fast-downloader       # (experimental) Use Fast Downloader (uses 'aria2')
    --keep-installer        # In Production mode, do not purge the tmp installer artifacts
    --startup               # Attempt to start the dr-provision service
    --systemd               # Run the systemd enabling commands after installation
    --systemd-services      # Additional services for systemd to wait for before starting DRP.
    --create-self           # DRP will create a machine that represents itself.
                            # Only used with startup/systemd parameters.
    --start-runner          # DRP will start a runner for itself. Implies create self.
                            # Only used with startup/systemd parameters.
    --initial-workflow=<string>
                            # Workflow to assign to the DRP's self machine as install finishes
                            # Only valid with create-self object, only one Workflow may be specified
    --initial-profiles=<string>
                            # Initial profiles to add to the DRP endpoint before starting the workflow,
                            # comma separated list, no spaces, this is only valid with create-self object
    --initial-contents=<string>
                            # Initial content packs to deliver, comma separated with no spaces.
                            # A file, URL, or content-pack name
    --initial-plugins=<string>
                            # Initial plugins to deliver, comma separated list with no spaces.
                            # A file, URL, or content-pack name
    --initial-parameters=<string>
                            # Initial parameters to set on the system.  Simple parameters
                            # as a comma separated list, with no spaces.
    --initial-subnets=<string>
                            # A file or URL containing Subnet definitions in JSON or YAML format,
                            # comma separated, no spaces (NOTE: Subnets can also be injected in
                            # content packs with the '--initial-contents' argument)
    --bootstrap             # Store the install image and the install script in the files bootstrap
    --drp-id=<string>       # String to use as the DRP Identifier (only with --systemd)
    --ha-id=<string>        # String to use as the HA Identifier (only with --systemd)
    --ha-enabled            # Indicates that the system is HA enabled
    --ha-address=<string>   # IP Address to use a VIP for HA system
    --ha-interface=<string> # Interrace to use for HA traffic
    --ha-passive            # Indicates that the system is starting as passive.
    --ha-token=<string>     # The token to use to sync passive to active
    --drp-user=<string>     # DRP user to create after system start (only with --systemd)
    --drp-password=<string> # DRP user password to set after system start (only with --systemd)
    --remove-rocketskates   # Remove the rocketskates user after system start (only with --systemd)
    --local-ui              # Set up DRP to server a local UI
    --system-user           # System user account to create for DRP to run as
    --system-group          # System group name
    --drp-home-dir          # Use with system-user and system-group to set the home directory
                            # for the system-user. This path is where most important drp files live
                            # including the tftp root
    --bin-dir               # Use this as the local of the binaries.  Required for non-root upgrades
    --container             # Force to install as a container, not zipfile
    --container-type=<string>
                            # Container install type, defaults to "$CNT_TYPE"
    --container-name=<string>
                            # Set the "docker run" container name, defaults to "$CNT_NAME"
    --container-restart=<string>
                            # Set the Docker restart option, defaults to "$CNT_RESTART"
                            # options are:  no, on-failure, always, unless-stopped
                            * see: https://docs.docker.com/config/containers/start-containers-automatically/
    --container-volume=<string>
                            # Volume name to use for backing persistent storage, default "$CNT_VOL"
    --container-registry="drp.example.com:5000"
                            # Alternate registry to get container images from, default "$CNT_REGISTRY"
    --container-env="<string> <string> <string>"
                            # Define a space separated list of environment variables to pass to the
                            # container on start (eg "RS_METRICS_PORT=8888 RS_DRP_ID=fred")
                            # see 'dr-provision --help' for complete list of startup variables
    --container-netns="<string>"
                            # Define Network Namespace to start container in. Defaults to "$CNT_NETNS"
                            # If set to empty string (""), then disable setting any network namespace
    --universal             # Load the universal components and bootstrap the system.
                            # This should be first and implies systemd, startup, start-runner, create-self.
                            # Additionally, implies running universal-boostrap and starting discovery.
                            # Subsequent options for '--initial-workflow' will be ignored if this is set.

    version                 # show install.sh script version and exit
    install                 # Sets up an isolated or system 'production' enabled install
    upgrade                 # Sets the installer to upgrade an existing 'dr-provision', for upgrade of
                            # container; kill/rm the DRP container, then upgrade and reattach data volume
    remove                  # Removes the system enabled install.  Requires no other flags
                            # optional: '--remove-data' to wipe all installed data

DEFAULTS:
    |  option:               value:           |  option:               value:
    |  -------------------   ------------     |  ------------------    ------------
    |  remove-rocketskates = false            |  version (*)         = $DEFAULT_DRP_VERSION
    |  isolated            = false            |  nocontent           = false
    |  upgrade             = false            |  force               = false
    |  debug               = false            |  skip-run-check      = false
    |  skip-prereqs        = false            |  systemd             = false
    |  create-self         = false            |  start-runner        = false
    |  drp-id              = unset            |  ha-id               = unset
    |  drp-user            = rocketskates     |  drp-password        = r0cketsk8ts
    |  startup             = false            |  keep-installer      = false
    |  local-ui            = false            |  system-user         = root
    |  system-group        = root             |  drp-home-dir        = /var/lib/dr-provision
    |  bin-dir             = /usr/local/bin   |  container           = false
    |  container-volume    = $CNT_VOL         |  container-registry  = $CNT_REGISTRY
    |  container-type      = $CNT_TYPE           |  container-name      = $CNT_NAME
    |  container-netns     = $CNT_NETNS             |  container-restart   = $CNT_RESTART
    |  bootstrap           = false            |  initial-workflow    = unset
    |  initial-contents    = unset            |  initial-profiles    = unset
    |  initial-plugins     = unset            |  initial-subnets     = unset
    |  ha-enabled          = false            |  ha-address          = unset
    |  ha-passive          = false            |  ha-interface        = unset
    |  ha-token            = unset            |  universal           = unset
    |  systemd-services    = unset

    * version examples: 'tip', 'v4.6.3', 'v4.7.0-beta1.3', or 'stable'

PREREQUISITES:
    NOTE: By default, prerequisite packages will be installed if possible.  You must
          manually install these first on a Mac OS X system. Package names may vary
          depending on your operating system version/distro packaging naming scheme.

    REQUIRED: curl
    OPTIONAL: aria2c (if using experimental "fast downloader")

WARNING: 'install' option will OVERWRITE existing installations

INSTALLER VERSION:  $VERSION
EOFUSAGE
} # end usage()

# control flags
ISOLATED=false
NO_CONTENT=false
DBG=false
UPGRADE=false
REMOVE_DATA=false
SKIP_RUN_CHECK=false
SKIP_DEPENDS=false
FAST_DOWNLOADER=false
SYSTEMD=false
START_RUNNER=false
CREATE_SELF=false
STARTUP=false
REMOVE_RS=false
LOCAL_UI=false
KEEP_INSTALLER=false
CONTAINER=false
BOOTSTRAP=false
INITIAL_WORKFLOW=
INITIAL_PROFILES=
INITIAL_PARAMETERS=
INITIAL_CONTENTS=
INITIAL_PLUGINS=
SYSTEMD_ADDITIONAL_SERVICES=
START_LIMIT_INTERVAL=5
START_LIMIT_BURST=100
UNIVERSAL=false

# download URL locations; overridable via ENV variables
URL_BASE=${URL_BASE:-"https://rebar-catalog.s3-us-west-2.amazonaws.com"}
URL_BASE_CONTENT=${URL_BASE_CONTENT:-"$URL_BASE/drp-community-content"}
DRP_CATALOG=${DRP_CATALOG:-"$URL_BASE/rackn-catalog.json"}

# set some builtin default values
_sudo="sudo"
BIN_DIR=/usr/local/bin
DRP_HOME_DIR=/var/lib/dr-provision
CNT_TYPE=docker
CNT_NAME=drp
CNT_VOL=drp-data
CNT_REGISTRY=index.docker.io
CNT_NETNS=host
CNT_ENV=""
CNT_RESTART="always"
CNT_VOL_REMOVE=true
SYSTEM_USER=root
SYSTEM_GROUP=root
HA_ENABLED=false
HA_ADDRESS=
HA_INTERFACE=
HA_TOKEN=
HA_PASSIVE=false

args=()
while (( $# > 0 )); do
    arg="$1"
    arg_key="${arg%%=*}"
    arg_data="${arg#*=}"
    case $arg_key in
        --help|-h)                  usage; exit_cleanup 0                             ;;
        --debug)                    DBG=true                                          ;;
        --version|--drp-version)    DRP_VERSION=${arg_data}                           ;;
        --isolated)                 ISOLATED=true                                     ;;
        --skip-run-check)           SKIP_RUN_CHECK=true                               ;;
        --skip-dep*|--skip-prereq*) SKIP_DEPENDS=true                                 ;;
        --fast-downloader)          FAST_DOWNLOADER=true                              ;;
        --force)                    force=true                                        ;;
        --remove-data)              REMOVE_DATA=true                                  ;;
        --upgrade)                  UPGRADE=true; force=true
                                    CNT_VOL_REMOVE=false                              ;;
        --nocontent|--no-content)   NO_CONTENT=true                                   ;;
        --no-sudo)                  _sudo=""                                          ;;
        --keep-installer)           KEEP_INSTALLER=true                               ;;
        --startup)                  STARTUP=true; SYSTEMD=true                        ;;
        --systemd)                  SYSTEMD=true                                      ;;
        --systemd-services)         SYSTEMD_ADDITIONAL_SERVICES="${arg_data}"         ;;
        --create-self)              CREATE_SELF=true                                  ;;
        --start-runner)             START_RUNNER=true; CREATE_SELF=true               ;;
        --bootstrap)                BOOTSTRAP=true                                    ;;
        --local-ui)                 LOCAL_UI=true                                     ;;
        --remove-rocketskates)      REMOVE_RS=true                                    ;;
        --initial-workflow)         [[ "$UNIVERSAL" == "false" ]] && INITIAL_WORKFLOW="${arg_data}"
                                                                                      ;;
        --initial-profiles)         INITIAL_PROFILES="${INITIAL_PROFILES}${arg_data}" ;;
        --initial-parameters)       INITIAL_PARAMETERS="${arg_data}"                  ;;
        --initial-contents)         INITIAL_CONTENTS="${INITIAL_CONTENTS}${arg_data}" ;;
        --initial-plugins)          INITIAL_PLUGINS="${INITIAL_PLUGINS}${arg_data}"   ;;
        --initial-subnets)          INITIAL_SUBNETS="${arg_data}"                     ;;
        --drp-user)                 DRP_USER=${arg_data}                              ;;
        --drp-password)             DRP_PASSWORD="${arg_data}"                        ;;
        --drp-id)                   DRP_ID="${arg_data}"                              ;;
        --ha-id)                    HA_ID="${arg_data}"                               ;;
        --ha-enabled)               HA_ENABLED=true                                   ;;
        --ha-address)               HA_ADDRESS="${arg_data}"                          ;;
        --ha-interface)             HA_INTERFACE="${arg_data}"                        ;;
        --ha-passive)               HA_PASSIVE=true                                   ;;
        --ha-token)                 HA_TOKEN="${arg_data}"                            ;;
        --system-user)              SYSTEM_USER="${arg_data}"                         ;;
        --system-group)             SYSTEM_GROUP="${arg_data}"                        ;;
        --drp-home-dir)             DRP_HOME_DIR="${arg_data}"                        ;;
        --container)                CONTAINER=true                                    ;;
        --container-type)           CNT_TYPE="${arg_data}"                            ;;
        --container-name)           CNT_NAME="${arg_data}"                            ;;
        --container-volume)         CNT_VOL="${arg_data}"                             ;;
        --container-restart)        CNT_RESTART="${arg_data}"                         ;;
        --container-registry)       CNT_REGISTRY="${arg_data}"                        ;;
        --container-netns)          CNT_NETNS="${arg_data}"                           ;;
        --container-env)            CNT_ENV="${arg_data}"                             ;;
        --universal)                UNIVERSAL=true; BOOTSTRAP=true;
                                    SYSTEMD=true;
                                    CREATE_SELF=true; START_RUNNER=true;
                                    INITIAL_WORKFLOW="universal-bootstrap";
                                    INITIAL_CONTENTS="universal,";
                                    INITIAL_PROFILES="bootstrap-discovery,"           ;;
        --zip-file)
            ZF=${arg_data}
            ZIP_FILE=$(echo "$(cd $(dirname $ZF) && pwd)/$(basename $ZF)")
            ;;
        --*)
            arg_key="${arg_key#--}"
            arg_key="${arg_key//-/_}"
            # "^^" Paremeter Expansion is a bash v4.x feature; Mac by default is bash 3.x
            #arg_key="${arg_key^^}"
            arg_key=$(echo $arg_key | tr '[:lower:]' '[:upper:]')
            echo "Overriding $arg_key with $arg_data"
            export $arg_key="$arg_data"
            ;;
        *)
            args+=("$arg");;
    esac
    shift
done
set -- "${args[@]}"

if [[ "$HA_ENABLED" == "true" ]] ; then
  if [[ "$HA_TOKEN" == "" ]] ; then
    if [[ "$HA_PASSIVE" == "true" ]] ; then
      usage; exit_cleanup 1 "Passive systems must have a token.  Specify HA_TOKEN."
    fi
    HA_TOKEN="ACTIVE_NO_TOKEN"
  fi
fi

CLI="${BIN_DIR}/drpcli"
CLI_BKUP="${BIN_DIR}/drpcli.drp-installer.backup"
PROVISION="${BIN_DIR}/dr-provision"
DRBUNDLER="${BIN_DIR}/drbundler"
PATH=$PATH:${BIN_DIR}

DRP_VERSION=${DRP_VERSION:-"$DEFAULT_DRP_VERSION"}
DRP_CONTENT_VERSION=stable
if [[ $DRP_VERSION == tip ]]; then
    DRP_CONTENT_VERSION=tip
fi
[[ "$ISOLATED" == "true" ]] && KEEP_INSTALLER=true

[[ $DBG == true ]] && set -x

if [[ $EUID -eq 0 ]]; then
    _sudo=""
else
    if [[ ! -x "$(command -v sudo)" ]]; then
        exit_cleanup 1 "FATAL: Script is not running as root and sudo command is not found. Please be root"
    fi
fi

# verify the container restart directives
case $CNT_RESTART in
  no|on-failure|always|unless-stopped) true ;;
  *) exit_cleanup 1 "FATAL: Container restart directive ('$CNT_RESTART') is not valid. See '--help' for details."
    ;;
esac

# Figure out what Linux distro we are running on.
export OS_TYPE= OS_VER= OS_NAME= OS_FAMILY=

if [[ -f /etc/os-release ]]; then
    . /etc/os-release
    OS_TYPE=${ID,,}
    OS_VER=${VERSION_ID,,}
elif [[ -f /etc/lsb-release ]]; then
    . /etc/lsb-release
    OS_VER=${DISTRIB_RELEASE,,}
    OS_TYPE=${DISTRIB_ID,,}
elif [[ -f /etc/centos-release || -f /etc/fedora-release || -f /etc/redhat-release ]]; then
    for rel in centos-release fedora-release redhat-release; do
        [[ -f /etc/$rel ]] || continue
        OS_TYPE=${rel%%-*}
        OS_VER="$(egrep -o '[0-9.]+' "/etc/$rel")"
        break
    done

    if [[ ! $OS_TYPE ]]; then
        exit_cleanup 1 "FATAL: Cannot determine Linux version we are running on!"
    fi
elif [[ -f /etc/debian_version ]]; then
    OS_TYPE=debian
    OS_VER=$(cat /etc/debian_version)
elif [[ $(uname -s) == Darwin ]] ; then
    OS_TYPE=darwin
    OS_VER=$(sw_vers | grep ProductVersion | awk '{ print $2 }')
fi
OS_NAME="$OS_TYPE-$OS_VER"

case $OS_TYPE in
    amzn|centos|redhat|fedora) OS_FAMILY="rhel"      ;;
    debian|ubuntu|raspbian)        OS_FAMILY="debian"    ;;
    coreos)               OS_FAMILY="container" ;;
    *)                    OS_FAMILY="$OS_TYPE"  ;;
esac

# Wait until DRP is ready, or exit
check_drp_ready() {
    COUNT=0
    while ! drpcli info get > /dev/null 2>&1 ; do
        echo "DRP is not up yet, waiting ($COUNT) ..."
        sleep 2
        # Pre-increment for compatibility with Bash 4.1+
        ((++COUNT))
        if (( $COUNT > 10 )) ; then
            exit_cleanup 1 "FATAL: DRP Failed to start"
        fi
    done
}

# install the EPEL repo if appropriate, and not enabled already
install_epel() {
    if [[ $OS_FAMILY == rhel ]] ; then
        if ( `yum repolist enabled | grep -q "^epel/"` ); then
            echo "EPEL repository installed already."
        else
            if [[ $OS_TYPE != fedora ]] ; then
                $_sudo yum install -y epel-release
            fi
        fi
    fi
}

# set our downloader GET variable appropriately - supports standard
# (curl) downloader or (experimental) aria2c fast downloader
get() {
    [[ -z "$*" ]] && exit_cleanup 1 "FATAL: Internal error, get() expects files to get"

    if [[ "$FAST_DOWNLOADER" == "true" ]]; then
        if which aria2c > /dev/null; then
            GET="aria2c --quiet=true --continue=true --max-concurrent-downloads=10 --max-connection-per-server=16 --max-tries=0"
        else
            exit_cleanup 1 "FATAL: '--fast-downloader' specified, but couldn't find tool ('aria2c')."
        fi
    else
        if which curl > /dev/null; then
            GET="curl -sfL"
        else
            exit_cleanup 1 "FATAL: Unable to find downloader tool ('curl')."
        fi
    fi
    for URL in $*; do
        FILE=${URL##*/}
        echo ">>> Downloading file:  $FILE"
        $GET -o $FILE $URL
    done
}


# setup the system user for drp to run as
setup_system_user() {
    if [[ ${SYSTEM_USER} == "root" ]]; then
        return
    fi
    # 0 or 9 here is fine on Deb or RHEL
    # 0 is success 9 means the account already
    # exists which is expected
    RC=0
    $_sudo groupadd --system ${SYSTEM_GROUP} || RC=$?
    if [[ ${RC} != 0 && ${RC} != 9 ]]; then
        exit_cleanup ${RC} "FATAL: Unable to create system group ${SYSTEM_GROUP}"
    fi
    if [[ ${OS_FAMILY} == "debian" ]]; then
        $_sudo adduser --system --home ${DRP_HOME_DIR} --quiet --group ${SYSTEM_USER}
        return
    else
        RC=0
        if [[ ${OS_FAMILY} == "rhel" ]]; then
            $_sudo adduser --system -d ${DRP_HOME_DIR} --gid ${SYSTEM_GROUP} -m --shell /sbin/nologin ${SYSTEM_USER} || RC=$?
        fi
    fi
    if [[ ${RC} == 0 || ${RC} == 9 ]]; then
        return
    fi
    exit_cleanup ${RC} "FATAL: Unable to create system user ${SYSTEM_USER}"
}

set_ownership_of_drp() {
    # It is possible for the home directory to not exist if
    # a non-root user was specified but already created.
    # Make sure a directory is created so DRP does not hit
    # permissions errors trying to use the home directory.
    if [ ! -d "${DRP_HOME_DIR}" ]; then
        echo "DRP Home directory ${DRP_HOME_DIR} did not exist - creating..."
        $_sudo mkdir -p ${DRP_HOME_DIR}
    fi
    $_sudo chown -R ${SYSTEM_USER}:${SYSTEM_GROUP} ${DRP_HOME_DIR}

    for i in ${PROVISION} ${CLI} ${DRBUNDLER} ${PROVISION}.bak ${CLI}.bak ${DRBUNDLER}.bak ${PROVISION}.new ${CLI}.new ${DRBUNDLER}.new ; do
        [[ -r "$i" ]] && $_sudo chown -R ${SYSTEM_USER}:${SYSTEM_GROUP} $i || true
    done
}

setcap_drp_binary() {
    if [[ ${SYSTEM_USER} != "root" ]]; then
        case ${OS_FAMILY} in
            rhel|debian)
                $_sudo setcap "cap_net_raw,cap_net_bind_service=+ep" ${PROVISION}
            ;;
            *)
                echo "Your OS Family ${OS_FAMILY} does not support setcap" \
                     "and may not be able to bind privileged ports when" \
                     "running as non-root user ${SYSTEM_USER}"
            ;;
        esac
    fi
}

check_bins_darwin() {
  local _bin=$1
     case $_bin in
        aria2c)
            if [[ "$FAST_DOWNLOADER" == "true" ]]; then
                _bin="aria2"
                if ! which $_bin > /dev/null 2>&1; then
                    echo "Must have binary '$_bin' installed."
                    echo "eg:   brew install $_bin"
                    error=1
                fi
            fi
          ;;
        *)
            if ! which $_bin > /dev/null 2>&1; then
                echo "Must have binary '$_bin' installed."
                echo "eg:   brew install $_bin"
                error 1
            fi
         ;;
     esac
} # end check_bins_darwin()

# handle RHEL and Debian types
check_pkgs_linux() {
    # assumes binary and package name are same
    local _pkg=$1
        if ! which $_pkg &>/dev/null; then
            [[ $_pkg == "aria2" && "$FAST_DOWNLOADER" != "true" ]] && return 0
            echo "Missing dependency '$_pkg', attempting to install it... "
            if [[ $OS_FAMILY == rhel ]] ; then
                echo $IN_EPEL | grep -q $_pkg && install_epel
                $_sudo yum install -y $_pkg
            elif [[ $OS_FAMILY == debian ]] ; then
                $_sudo apt-get install -y $_pkg
            fi
        fi
} # end check_pkgs_linux()

ensure_packages() {
    echo "Ensuring required tools are installed"
    case $OS_FAMILY in
        darwin)
            error=0
            BINS="curl aria2c"
            for BIN in $BINS; do
                check_bins_darwin $BIN
            done

            if [[ $error == 1 ]]; then
                echo "After install missing components, restart the terminal to pick"
                echo "up the newly installed commands, and re-run the installer."
                echo
                exit_cleanup 1
            fi
        ;;
        rhel|debian)
            PKGS="curl aria2"
            IN_EPEL="curl aria2"
            for PKG in $PKGS; do
                check_pkgs_linux $PKG
            done
        ;;
        coreos)
            echo "CoreOS does not require any packages to be installed.  DRP will be"
            echo "installed from the Docker Hub registry."
        ;;
        photon)
            if ! which tar > /dev/null 2>&1; then
                echo "Installing packages for Photon Linux..."
                tdnf -y makecache
                tdnf -y install tar
            else
                echo "'tar' already installed on Photon Linux..."
            fi
        ;;
        *)
            exit_cleanup 1 "FATAL: Unsupported OS Family ($OS_FAMILY)."
        ;;
    esac
} # end ensure_packages()

# output a friendly statement on how to download ISOS via fast downloader
show_fast_isos() {
    cat <<FASTMSG
Option '--fast-downloader' requested.  You may download the ISO images using
'aria2c' command to significantly reduce download time of the ISO images.

NOTE: The following genereted scriptlet should download, install, and enable
      the ISO images.  VERIFY SCRIPTLET before running it.

      YOU MUST START 'dr-provision' FIRST! Example commands:

###### BEGIN scriptlet
  export CMD="aria2c --continue=true --max-concurrent-downloads=10 --max-connection-per-server=16 --max-tries=0"
FASTMSG

    for BOOTENV in $*
    do
        echo "  export URL=$(${EP}drpcli bootenvs show $BOOTENV | grep 'IsoUrl' | cut -d '\"' -f 4)"
        echo "  export ISO=$(${EP}drpcli bootenvs show $BOOTENV | grep 'IsoFile' | cut -d '\"' -f 4)"
        echo "  \$CMD -o \$ISO \$URL"
    done
    echo "  # this should move the ISOs to the TFTP directory..."
    echo "  $_sudo mv *.tar *.iso $TFTP_DIR/isos/"
    echo "  $_sudo pkill -HUP dr-provision"
    echo "  echo 'NOTICE:  exploding isos may take up to 5 minutes to complete ... '"
    echo "###### END scriptlet"

    echo
} # end show_fast_isos()

remove_container() {
    case $CNT_TYPE in
        docker)
            $_sudo docker kill $CNT_NAME > /dev/null && echo "Killed docker container '$CNT_NAME'"
            $_sudo docker rm $CNT_NAME > /dev/null && echo "Removed docker container '$CNT_NAME'"
            if [[ "$CNT_VOL_REMOVE" == "true" ]]; then
                $_sudo docker volume rm $CNT_VOL > /dev/null && echo "Removed backing volume named '$CNT_VOL'"
            else
                echo ">>> Not removing container data volume '$CNT_VOL'"
                if [[ "$MODE" == "remove" ]]; then
                    echo "    To remove: '$_sudo docker volume rm $CNT_VOL'"
                    echo "    or use '--remove-data' next time ... "
                fi
            fi
            ;;
        *)  exit_cleanup 1 "Container type '$CNT_TYPE' not supported in installer."
            ;;
    esac
} # end remove_container()

install_container() {
    if [[ -n "$CNT_ENV" ]]; then
        for ENV in $CNT_ENV; do
            OPTS="--env $ENV $OPTS"
        done
        ENV_OPTS=$(echo $OPTS | sed 's/ $//')
    fi
    case $CNT_TYPE in
        docker)
            ! which docker > /dev/null 2>&1 && exit_cleanup 1 "Container install requested but no 'docker' in PATH ($PATH)."
            if [[ "$UPGRADE" == "false" ]]; then
                $_sudo docker volume create $CNT_VOL > /dev/null
                VOL_MNT=$($_sudo docker volume inspect $CNT_VOL | grep Mountpoint | awk -F\" '{ print $4 }')
                echo "Created docker volume named '$CNT_VOL' with mountpoint '$VOL_MNT'"
            else
                if $_sudo docker volume inspect $CNT_VOL > /dev/null 2>&1; then
                    echo "Attempting to reconnect volume '$CNT_VOL'"
                else
                    exit_cleanup 1 "No existing volume '$CNT_VOL' found to reconnect."
                fi
            fi
            if [[ -z "$CNT_NETNS" ]]; then
              echo "WARNING:  No network namespace set - you may have issues with DHCP and TFTP depending on your use case."
            else
              NETNS="--net $CNT_NETNS"
            fi
            CMD="$_sudo docker run $ENV_OPTS --restart "$CNT_RESTART" --volume $CNT_VOL:/provision/drp-data --name \"$CNT_NAME\" -itd $NETNS ${CNT_REGISTRY}/digitalrebar/provision:$DRP_VERSION"
            echo "Starting container with following run command:"
            echo "$CMD"
            eval $CMD

            echo ""
            echo "Digital Rebar Provision container is using backing volume: $CNT_VOL"
            echo "Volume is backed on host filesystem at: $VOL_MNT"
            echo ""
            echo "Docker container run time information:"
            $_sudo docker ps --filter name=$CNT_NAME
            echo ""
            echo ">>>  NOTICE:  If you intend to upgrade this container, record your 'install.sh' options  <<< "
            echo ">>>           for the 'upgrade' command - you must reuse the same options to obtain the  <<< "
            echo ">>>           same installed results in the upgraded container.  You have been warned!   <<< "
            echo ""
            ;;
        *)  exit_cleanup 1 "Container type '$CNT_TYPE' not supported in installer."
            ;;
    esac
} # end install_container()

# main

MODE=$1

if [[ $OS_FAMILY == "container" || $CONTAINER == "true" ]]; then
    RMV_MSG="Remove not supported for container based installer at this time."
    GEN_MSG="Unsupported mode '$MODE'"
    case $MODE in
        install)
            echo "Installing Digital Rebar Provision as a container."
            install_container
            exit_cleanup $?
            ;;
        upgrade)
            echo "Upgrading Digital Rebar Provision as a container."
            UPGRADE=true
            CNT_VOL_REMOVE=false
            remove_container
            install_container
            exit_cleanup $?
            ;;
        remove)  [[ "$REMOVE_DATA" == "true" ]] && CNT_VOL_REMOVE="true" || CNT_VOL_REMOVE="false"
                 remove_container; exit_cleanup $? ;;
        *)       usage; exit_cleanup 1 "$GEN_MSG"  ;;
    esac
fi

arch=$(uname -m)
case $arch in
    x86_64|amd64) arch=amd64   ;;
    aarch64)      arch=arm64   ;;
    armv7l)       arch=arm_v7  ;;
    ppc64le)      arch=ppc64le ;;
    *)            exit_cleanup 1 "FATAL: architecture ('$arch') not supported" ;;
esac

case $(uname -s) in
    Darwin)
        os="darwin"
        binpath="bin/darwin/$arch"
        bindest="${BIN_DIR}"
        tar="command tar"
        # Someday, handle adding all the launchd stuff we will need.
        shasum="command shasum -a 256";;
    Linux)
        os="linux"
        binpath="bin/linux/$arch"
        bindest="${BIN_DIR}"
        tar="command tar"
        if [[ -d /etc/systemd/system ]]; then
            # SystemD
            SYSTEMD=true
            initfile="assets/startup/dr-provision.service"
            initdest="/etc/systemd/system/dr-provision.service"
            starter="$_sudo systemctl daemon-reload && $_sudo systemctl start dr-provision"
            enabler="$_sudo systemctl daemon-reload && $_sudo systemctl enable dr-provision"
        elif [[ -d /etc/init ]]; then
            # Upstart
            initfile="assets/startup/dr-provision.unit"
            initdest="/etc/init/dr-provision.conf"
            starter="$_sudo service dr-provision start"
            enabler="$_sudo service dr-provision enable"
        elif [[ -d /etc/init.d ]]; then
            # SysV
            initfile="assets/startup/dr-provision.sysv"
            initdest="/etc/init.d/dr-provision"
            starter="/etc/init.d/dr-provision start"
            enabler="/etc/init.d/dr-provision enable"
        else
            echo "No idea how to install startup stuff -- not using systemd, upstart, or sysv init"
            exit_cleanup 1
        fi
        shasum="command sha256sum";;
    *)
        # Someday, support installing on Windows.  Service creation could be tricky.
        echo "No idea how to check sha256sums"
        exit_cleanup 1;;
esac

if [[ "$MODE" == "upgrade" ]]
then
    MODE=install
    UPGRADE=true
    force=true
    CNT_VOL_REMOVE=false
fi

case $MODE in
     install)
             # a system under control of another DRP endpoint may have RS_ENDPOINT set, and
             # this will cause various startup check procedures to fail or be wrong
             export RS_ENDPOINT="https://127.0.0.1:8092"

             if [[ "$ISOLATED" == "false" || "$SKIP_RUN_CHECK" == "false" ]]; then
                 if pgrep dr-provision; then
                     echo "'dr-provision' service is running, CAN NOT upgrade ... please stop service first"
                     exit_cleanup 9
                 else
                     echo "'dr-provision' service is not running, beginning install process ... "
                 fi
             else
                 echo "Skipping 'dr-provision' service run check as requested ..."
             fi

             [[ "$SKIP_DEPENDS" == "false" ]] && ensure_packages || echo "Skipping dependency checks as requested ... "

             if [[ "$ISOLATED" == "false" ]]; then
                 TMP_INSTALLER_DIR=$(mktemp -d /tmp/drp.installer.XXXXXX)
                 echo "Using temp directory to extract artifacts to and install from ('$TMP_INSTALLER_DIR')."
                 OLD_PWD=$(pwd)
                 cd $TMP_INSTALLER_DIR
                 TMP_INST=$TMP_INSTALLER_DIR/tools/install.sh
             fi

             # Are we in a build tree
             if [ -e server ] ; then
                 if [ ! -e bin/linux/amd64/drpcli ] ; then
                     echo "It appears that nothing has been built."
                     echo "Please run tools/build.sh and then rerun this command".
                     exit_cleanup 1
                 fi
             else
                 # We aren't a build tree, but are we extracted install yet?
                 # If not, get the requested version.
                 if [[ ! -e sha256sums || $force ]] ; then
                     echo "Installing Version $DRP_VERSION of Digital Rebar Provision"
                     ZIP="dr-provision.zip"
                     SHA="dr-provision.sha256"
                     if [[ -n "$ZIP_FILE" ]]
                     then
                       [[ "$ZIP_FILE" != "dr-provision.zip" ]] && cp "$ZIP_FILE" dr-provision.zip
                       echo "WARNING:  No sha256sum check performed for '--zip-file' mode."
                       echo "          We assume you've already verified your download file."
                     else
                       if [[ ! -e rackn-catalog.json ]] ; then
                         get $DRP_CATALOG
                         FILE=${DRP_CATALOG##*/}
                         mv ${FILE} rackn-catalog.json.gz
                         gunzip rackn-catalog.json.gz
                       fi
                       SOURCE=$(grep "drp-$DRP_VERSION:::" rackn-catalog.json | awk -F\" '{ print $4 }')
                       SOURCE=${SOURCE##*:::}

                       get $SOURCE || true
                       FILE=${SOURCE##*/}
                       # if we couldn't find dr-provision.zip, then try os/arch based zip
                       if [[ ! -e $FILE ]] ; then
                           SOURCE=${SOURCE/.zip/.${os}.${arch}.zip}
                           get $SOURCE
                           FILE=${SOURCE##*/}
                       fi
                       mv $FILE $ZIP

                       # XXX: Put sha back one day
                     fi

                     if ! $tar -xf dr-provision.zip; then
                       # try to fall back to bsdtar
                       echo "Attempting to auto recover from above 'tar' errors."
                       echo "Newer 'install.sh' script being used with older dr-provision.zip file."
                       echo "Attempting to locate and use 'bsdtar' as a fallback..."

                       if which bsdtar > /dev/null; then
                         if ! bsdtar -xf dr-provision.zip; then
                           echo "!!!  FAILED to extract 'dr-provision.zip' with installed 'bsdtar'."
                           echo "!!!  Not cleaning up for forensic analysis."
                           exit 1
                         fi
                       else
                         echo "'bsdtar' not found in path, attempting to install it... "
                         if check_pkgs_linux bsdtar; then
                           if ! bsdtar -xf dr-provision.zip; then
                             echo "!!!  Last ditch attempt FAILED to unpack 'dr-provision.zip' file."
                             echo "!!!  Not attempting cleanup for debug and troubleshooting reasons."
                             exit 1
                           fi
                         else
                           exit_cleanup 1 "FAILED to install 'bsdtar' to satisfy dependencies.  Please install it and retry."
                         fi
                       fi
                     fi
                 fi
                 $shasum -c sha256sums || exit_cleanup 1
             fi

             if [[ $NO_CONTENT == false ]]; then
                 echo "Installing Version $DRP_CONTENT_VERSION of Digital Rebar Provision Community Content"
                 if [[ -n "$ZIP_FILE" ]]; then
                   echo "WARNING: '--zip-file' specified, still trying to download community content..."
                   echo "         (specify '--no-content' to skip download of community content"
                 fi

                 if [[ ! -e rackn-catalog.json ]] ; then
                     get $DRP_CATALOG
                     mv rackn-catalog.json rackn-catalog.json.gz
                     gunzip rackn-catalog.json.gz
                 fi
                 SOURCE=$(grep "drp-community-content-$DRP_CONTENT_VERSION:::" rackn-catalog.json | awk -F\" '{ print $4 }')
                 SOURCE=${SOURCE##*:::}

                 get $SOURCE
                 FILE=${SOURCE##*/}
                 mv $FILE drp-community-content.json
                 # XXX: Add back in sha
             fi

             if [[ $ISOLATED == false ]]; then
                 setup_system_user

                 if [[ $initfile ]]; then
                     if [[ -r $initdest ]]
                     then
                         echo "WARNING ... WARNING ... WARNING"
                         echo "initfile ('$initfile') exists already, not overwriting it"
                         echo "please verify 'dr-provision' startup options are correct"
                         echo "for your environment and the new version .. "
                         echo ""
                         echo "specifically verify: '--file-root=<tftpboot directory>'"
                     else
                         $_sudo sed "s:/usr/local/bin/dr-provision:$PROVISION:g" "$initfile" > "$initdest"
                     fi
                     # output our startup helper messages only if SYSTEMD isn't specified
                     if [[ "$SYSTEMD" == "false" || "$STARTUP" == "false" ]]; then
                        echo
                        echo "######### You can start the DigitalRebar Provision service with:"
                        echo "$starter"
                        echo "######### You can enable the DigitalRebar Provision service with:"
                        echo "$enabler"
                    else
                        echo "######### Will attempt to execute startup procedures ('--startup' specified)"
                        echo "$starter"
                        echo "$enabler"

                    fi
                 fi

                 if [[ ! -e ${DRP_HOME_DIR}/digitalrebar/tftpboot && -e /var/lib/tftpboot ]] ; then
                     echo "MOVING /var/lib/tftpboot to ${DRP_HOME_DIR}/tftpboot location ... "
                     $_sudo mv /var/lib/tftpboot ${DRP_HOME_DIR}
                 fi

                 if [[ $NO_CONTENT == false ]] ; then
                     $_sudo mkdir -p ${DRP_HOME_DIR}/saas-content
                     DEFAULT_CONTENT_FILE="${DRP_HOME_DIR}/saas-content/drp-community-content.json"
                     $_sudo mv drp-community-content.json $DEFAULT_CONTENT_FILE
                 fi

                 # Make sure bindest exists
                 $_sudo mkdir -p "$bindest"

                 # move aside/preserve an existing drpcli - this machine might be under
                 # control of another DRP Endpoint, and this will break the installer (text file busy)
                 if [[ -f "$CLI" ]]; then
                     echo "SAVING '${BIN_DIR}/drpcli' to backup file ($CLI_BKUP)"
                     $_sudo mv "$CLI" "$CLI_BKUP"
                 fi

                 INST="${BIN_DIR}/drp-install.sh"
                 $_sudo cp $TMP_INST $INST && $_sudo chmod 755 $INST
                 echo "Install script saved to '$INST'"
                 echo "(run '$_sudo $INST remove' to uninstall DRP - must be root)"

                 TFTP_DIR="${DRP_HOME_DIR}/tftpboot"
                 $_sudo cp "$binpath"/* "$bindest"

                 setcap_drp_binary

                 set_ownership_of_drp

                 if [[ $SYSTEMD == true ]] ; then
                     if [[ ${SYSTEMD_ADDITIONAL_SERVICES} != "" ]] ; then
                         sed -i "s/^After=network.target\$/After=network.target,${SYSTEMD_ADDITIONAL_SERVICES}/" /etc/systemd/systemd/dr-provision.service
                     fi

                     mkdir -p /etc/systemd/system/dr-provision.service.d
                     cat > /etc/systemd/system/dr-provision.service.d/user.conf <<EOF
[Service]
Restart=always
StartLimitInterval=${START_LIMIT_INTERVAL}
StartLimitBurst=${START_LIMIT_BURST}
User=${SYSTEM_USER}
Group=${SYSTEM_GROUP}
Environment=RS_BASE_ROOT=${DRP_HOME_DIR}
EOF
                     if [[ ${SYSTEM_USER} != "root" ]]; then
                        cat > /etc/systemd/system/dr-provision.service.d/setcap.conf <<EOF
[Service]
PermissionsStartOnly=true
ExecStartPre=-/usr/bin/env setcap "cap_net_raw,cap_net_bind_service=+ep" ${PROVISION}
Environment=RS_EXIT_ON_CHANGE=true
Environment=RS_PLUGIN_COMM_ROOT=pcr
EOF
                     fi
                     if [[ $DRP_ID ]] ; then
                       cat > /etc/systemd/system/dr-provision.service.d/drpid.conf <<EOF
[Service]
Environment=RS_DRP_ID=$DRP_ID
EOF
                     fi
                     if [[ $HA_ID ]] ; then
                       cat > /etc/systemd/system/dr-provision.service.d/haid.conf <<EOF
[Service]
Environment=RS_HA_ID=$HA_ID
EOF
                     fi
                     if [[ "$HA_ENABLED" == "true" ]] ; then
                       cat > /etc/systemd/system/dr-provision.service.d/ha.conf <<EOF
[Service]
Environment=RS_HA_ENABLED=true
Environment=RS_HA_INTERFACE=$HA_INTERFACE
Environment=RS_HA_ADDRESS=$HA_ADDRESS
Environment=RS_HA_PASSIVE=$HA_PASSIVE
Environment=RS_HA_TOKEN=$HA_TOKEN
EOF
                     fi
                     if [[ $IPADDR ]] ; then
                       IPADDR="${IPADDR///*}"
                       cat > /etc/systemd/system/dr-provision.service.d/ipaddr.conf <<EOF
[Service]
Environment=RS_STATIC_IP=$IPADDR
Environment=RS_FORCE_STATIC=true
EOF
                     fi
                     if [[ $LOCAL_UI == true ]] ; then
                       cat > /etc/systemd/system/dr-provision.service.d/local-ui.conf <<EOF
[Service]
Environment=RS_LOCAL_UI=tftpboot/files/ux
Environment=RS_UI_URL=/ux
EOF
                     fi
                     if [[ $CREATE_SELF == true ]] ; then
                       cat > /etc/systemd/system/dr-provision.service.d/create-self.conf <<EOF
[Service]
Environment=RS_CREATE_SELF=$CREATE_SELF
EOF
                     fi
                     if [[ $START_RUNNER == true ]] ; then
                       cat > /etc/systemd/system/dr-provision.service.d/start-runner.conf <<EOF
[Service]
Environment=RS_START_RUNNER=$START_RUNNER
EOF
                     fi

                     eval "$enabler"
                     eval "$starter"

                     # If upgrading, assume DRP user already created
                     if [[ "$UPGRADE" == "true" ]]; then
                         if [[ $DRP_USER ]] ; then
                            export RS_KEY="$DRP_USER:$DRP_PASSWORD"
                         fi
                         check_drp_ready
                     else
                         check_drp_ready
                         if [[ $NO_CONTENT == false ]] ; then
                             drpcli contents upload catalog:task-library-${DRP_CONTENT_VERSION}

                             if [[ "$INITIAL_CONTENTS" != "" ]] ; then
                                 OLD_IFS="${IFS}"
                                 IFS=',' read -ra contents_array <<< "$INITIAL_CONTENTS"
                                 for i in "${contents_array[@]}" ; do
                                     if [[ -f ${OLD_PWD}/$i ]] ; then
                                         drpcli contents upload ${OLD_PWD}/${i}
                                     elif [[ -f "$i" ]] ; then
                                         drpcli contents upload ${i}
                                     elif [[ $i == http* ]] ; then
                                         drpcli contents upload ${i}
                                     else
                                         drpcli catalog item install ${i} --version=${DRP_CONTENT_VERSION} -c $DRP_CATALOG
                                     fi
                                 done
                                 IFS="${OLD_IFS}"
                             fi

                             if [[ "$INITIAL_PLUGINS" != "" ]] ; then
                                 OLD_IFS="${IFS}"
                                 IFS=',' read -ra plugins_array <<< "$INITIAL_PLUGINS"
                                 for i in "${plugins_array[@]}" ; do
                                     if [[ -f ${OLD_PWD}/$i ]] ; then
                                         drpcli plugin_providers upload ${OLD_PWD}/${i}
                                     elif [[ $i == http* ]] ; then
                                         drpcli plugin_providers upload ${i}
                                     else
                                         drpcli catalog item install ${i} --version=${DRP_CONTENT_VERSION} -c $DRP_CATALOG
                                     fi
                                 done
                                 IFS="${OLD_IFS}"
                             fi

                             if [[ "$INITIAL_PROFILES" != "" ]] ; then
                                 if [[ $CREATE_SELF == true ]] ; then
                                   cp $(which drpcli) /tmp/jq
                                   chmod +x /tmp/jq
                                   ID=$(drpcli info get | /tmp/jq .id -r | sed -r 's/:/-/g')
                                   rm /tmp/jq
                                   OLD_IFS="${IFS}"
                                   IFS=',' read -ra profiles_array <<< "$INITIAL_PROFILES"
                                   for i in "${profiles_array[@]}" ; do
                                     drpcli machines addprofile "Name:$ID" "$i" >/dev/null
                                   done
                                   IFS="${OLD_IFS}"
                                 fi
                             fi

                             if [[ "$INITIAL_PARAMETERS" != "" ]] ; then
                                 if [[ $CREATE_SELF == true ]] ; then
                                   cp $(which drpcli) /tmp/jq
                                   chmod +x /tmp/jq
                                   ID=$(drpcli info get | /tmp/jq .id -r | sed -r 's/:/-/g')
                                   rm /tmp/jq
                                   OLD_IFS="${IFS}"
                                   IFS=',' read -ra param_array <<< "$INITIAL_PARAMETERS"
                                   for i in "${param_array[@]}" ; do
                                     IFS="=" read -ra data_array <<< "$i"
                                     drpcli machines set "Name:$ID" param "${data_array[0]}" to "${data_array[1]}" >/dev/null
                                   done
                                   IFS="${OLD_IFS}"
                                 fi
                             fi

                             if [[ "$INITIAL_WORKFLOW" != "" ]] ; then
                                 if [[ $CREATE_SELF == true ]] ; then
                                     cp $(which drpcli) /tmp/jq
                                     chmod +x /tmp/jq
                                     ID=$(drpcli info get | /tmp/jq .id -r | sed -r 's/:/-/g')
                                     rm /tmp/jq
                                     echo "Setting initial workflow to '$INITIAL_WORKFLOW' for Machine '$ID'"
                                     drpcli machines workflow "Name:$ID" "$INITIAL_WORKFLOW" >/dev/null
                                 fi
                             fi

                             if [[ "$INITIAL_SUBNETS" != "" ]] ; then
                                 if [[ $CREATE_SELF == true ]] ; then
                                     OLD_IFS="${IFS}"
                                     IFS=',' read -ra subnet_array <<< "$INITIAL_SUBNETS"
                                     for i in "${subnet_array[@]}" ; do
                                       if [[ $i == http* ]] ; then
                                         echo "Creating subnet from URL '$i'"
                                         drpcli subnets create ${i}
                                       elif [[ -f ${OLD_PWD}/$i ]] ; then
                                         echo "Creating subnet from file '${OLD_PWD}/$i'"
                                         drpcli subnets create ${OLD_PWD}/${i}
                                       elif [[ -f "$i" ]] ; then
                                         echo "Creating subnet from file '${i}'"
                                         drpcli subnets create ${i}
                                       else
                                         echo ">>> WARNING: unable to read subnet file '$i' or '${OLD_PWD}/$i'; no SUBNET created"
                                       fi
                                     done
                                     IFS="${OLD_IFS}"
                                 fi
                             fi
                         fi

                         if [[ $DRP_USER ]] ; then
                             if drpcli users exists $DRP_USER >/dev/null 2>/dev/null ; then
                                 drpcli users update $DRP_USER "{ \"Name\": \"$DRP_USER\", \"Roles\": [ \"superuser\" ] }"
                             else
                                 drpcli users create "{ \"Name\": \"$DRP_USER\", \"Roles\": [ \"superuser\" ] }"
                             fi
                             drpcli users password $DRP_USER "$DRP_PASSWORD"
                             export RS_KEY="$DRP_USER:$DRP_PASSWORD"
                             if [[ $REMOVE_RS == true ]] ; then
                                 drpcli users destroy rocketskates
                             fi
                         fi
                     fi
                     if [[ "$BOOTSTRAP" == "true" ]] ; then
                         drpcli files upload dr-provision.zip as bootstrap/dr-provision.zip
                         drpcli files upload $TMP_INST as bootstrap/install.sh
                     fi
                 else
                     if [[ "$STARTUP" == "true" ]]; then
                         echo "######### Attempting startup of 'dr-provision' ('--startup' specified)"
                         eval "$enabler"
                         eval "$starter"

                         check_drp_ready

                         if [[ "$NO_CONTENT" == "false" ]] ; then
                             drpcli catalog item install task-library --version=${DRP_CONTENT_VERSION} -c $DRP_CATALOG

                             if [[ "$INITIAL_CONTENTS" != "" ]] ; then
                                 IFS=',' read -ra contents_array <<< "$INITIAL_CONTENTS"
                                 for i in "${contents_array[@]}" ; do
                                     if [[ -f ${OLD_PWD}/$i ]] ; then
                                         drpcli contents upload ${OLD_PWD}/${i}
                                     elif [[ $i == http* ]] ; then
                                         drpcli contents upload ${i}
                                     else
                                         drpcli catalog item install ${i} --version=${DRP_CONTENT_VERSION} -c $DRP_CATALOG
                                     fi
                                 done
                             fi

                             if [[ "$INITIAL_PLUGINS" != "" ]] ; then
                                 IFS=',' read -ra plugins_array <<< "$INITIAL_PLUGINS"
                                 for i in "${plugins_array[@]}" ; do
                                     if [[ -f ${OLD_PWD}/$i ]] ; then
                                         drpcli plugin_providers upload ${OLD_PWD}/${i}
                                     elif [[ $i == http* ]] ; then
                                         drpcli plugin_providers upload ${i}
                                     else
                                         drpcli catalog item install ${i} --version=${DRP_CONTENT_VERSION} -c $DRP_CATALOG
                                     fi
                                 done
                             fi

                             if [[ "$INITIAL_PROFILES" != "" ]] ; then
                                 if [[ $CREATE_SELF == true ]] ; then
                                   cp $(which drpcli) /tmp/jq
                                   chmod +x /tmp/jq
                                   ID=$(drpcli info get | /tmp/jq .id -r | sed -r 's/:/-/g')
                                   rm /tmp/jq
                                   IFS=',' read -ra profiles_array <<< "$INITIAL_PROFILES"
                                   for i in "${profiles_array[@]}" ; do
                                     drpcli machines addprofile "Name:$ID" "$i" >/dev/null
                                   done
                                 fi
                             fi

                             if [[ "$INITIAL_PARAMETERS" != "" ]] ; then
                                 if [[ $CREATE_SELF == true ]] ; then
                                   cp $(which drpcli) /tmp/jq
                                   chmod +x /tmp/jq
                                   ID=$(drpcli info get | /tmp/jq .id -r | sed -r 's/:/-/g')
                                   rm /tmp/jq
                                   IFS=',' read -ra param_array <<< "$INITIAL_PARAMETERS"
                                   for i in "${param_array[@]}" ; do
                                     IFS="=" read -ra data_array <<< "$i"
                                     drpcli machines set "Name:$ID" param "${data_array[0]}" to "${data_array[1]}" >/dev/null
                                   done
                                 fi
                             fi

                             if [[ "$INITIAL_WORKFLOW" != "" ]] ; then
                                 if [[ $CREATE_SELF == true ]] ; then
                                     cp $(which drpcli) /tmp/jq
                                     chmod +x /tmp/jq
                                     ID=$(drpcli info get | /tmp/jq .id -r)
                                     rm /tmp/jq
                                     drpcli machines workflow "Name:$ID" "$INITIAL_WORKFLOW" >/dev/null
                                 fi
                             fi
                         fi
                         if [[ "$BOOTSTRAP" == "true" ]] ; then
                             drpcli files upload dr-provision.zip as bootstrap/dr-provision.zip
                             drpcli files upload $TMP_INST as bootstrap/install.sh
                         fi
                     fi
                 fi

                 cd $OLD_PWD
                 if [[ "$KEEP_INSTALLER" == "false" ]]; then
                     rm -rf $TMP_INSTALLER_DIR
                 else
                     echo ""
                     echo "######### Installer artifacts are in '$TMP_INSTALLER_DIR' - to purge:"
                     echo "$_sudo rm -rf $TMP_INSTALLER_DIR"
                 fi

             # do an "isolated" mode install
             else
                 mkdir -p drp-data
                 TFTP_DIR="`pwd`/drp-data/tftpboot"

                 # Make local links for execs
                 rm -f drpcli dr-provision drpjoin
                 ln -s $binpath/drpcli drpcli
                 ln -s $binpath/dr-provision dr-provision
                 if [[ -e $binpath/drpjoin ]] ; then
                     ln -s $binpath/drpjoin drpjoin
                 fi

                 if [[ "$STARTUP" == "false" ]]; then
                     echo
                     echo "********************************************************************************"
                     echo
                     echo "# Run the following commands to start up dr-provision in a local isolated way."
                     echo "# The server will store information and serve files from the drp-data directory."
                     echo
                 else
                     echo
                     echo "********************************************************************************"
                     echo
                     echo "# Will attempt to startup the 'dr-provision' service ... "
                 fi

                 if [[ $IPADDR == "" ]] ; then
                     if [[ $OS_FAMILY == darwin ]]; then
                         ifdefgw=$(netstat -rn -f inet | grep default | awk '{ print $6 }')
                         if [[ $ifdefgw ]] ; then
                                 IPADDR=$(ifconfig en0 | grep 'inet ' | awk '{ print $2 }')
                         else
                                 IPADDR=$(ifconfig -a | grep "inet " | grep broadcast | head -1 | awk '{ print $2 }')
                         fi
                     else
                         gwdev=$(/sbin/ip -o -4 route show default |head -1 |awk '{print $5}')
                         if [[ $gwdev ]]; then
                             # First, advertise the address of the device with the default gateway
                             IPADDR=$(/sbin/ip -o -4 addr show scope global dev "$gwdev" |head -1 |awk '{print $4}')
                         else
                             # Hmmm... we have no access to the Internet.  Pick an address with
                             # global scope and hope for the best.
                             IPADDR=$(/sbin/ip -o -4 addr show scope global |head -1 |awk '{print $4}')
                         fi
                     fi
                 fi

                 if [[ $IPADDR ]] ; then
                     IPADDR="${IPADDR///*}"
                 fi

                 if [[ $OS_FAMILY == darwin ]]; then
                     bcast=$(netstat -rn | grep "255.255.255.255 " | awk '{ print $6 }')
                     if [[ $bcast == "" && $IPADDR ]] ; then
                             echo "# No broadcast route set - this is required for Darwin < 10.9."
                             echo "$_sudo route add 255.255.255.255 $IPADDR"
                             echo "# No broadcast route set - this is required for Darwin > 10.9."
                             echo "$_sudo route -n add -net 255.255.255.255 $IPADDR"
                     fi
                 fi

                 STARTER="$_sudo ./dr-provision --base-root=`pwd`/drp-data > drp.log 2>&1 &"
                 [[ "$STARTUP" == "false" ]] && echo "$STARTER"
                 mkdir -p "`pwd`/drp-data/saas-content"
                 if [[ $NO_CONTENT == false ]] ; then
                     DEFAULT_CONTENT_FILE="`pwd`/drp-data/saas-content/default.json"
                     mv drp-community-content.json $DEFAULT_CONTENT_FILE
                 fi

                 if [[ "$STARTUP" == "true" ]]; then
                     eval $STARTER
                     echo "'dr-provision' running processes:"
                     ps -eo pid,args -o comm  | grep -v grep | grep dr-provision
                     echo

                     if [[ "$BOOTSTRAP" == "true" ]] ; then
                         drpcli files upload dr-provision.zip as bootstrap/dr-provision.zip
                         drpcli files upload tools/install.sh as bootstrap/install.sh
                     fi
                 fi

                 EP="./"
             fi

             echo
             echo "# Once dr-provision is started, setup a base discovery configuration"
             echo "  ${EP}drpcli bootenvs uploadiso sledgehammer"
             echo "  ${EP}drpcli prefs set defaultWorkflow discover-base unknownBootEnv discovery defaultBootEnv sledgehammer defaultStage discover"

             if [[ $NO_CONTENT == true ]] ; then
                 echo "# Add common utilities (sourced from RackN)"
                 echo "  ${EP}drpcli contents upload catalog:task-library-$DRP_CONTENT_VERSION"
             fi
             echo
             echo "# Optionally, locally cache the isos for common community operating systems"
             echo "  ${EP}drpcli bootenvs uploadiso ubuntu-18.04-install"
             echo "  ${EP}drpcli bootenvs uploadiso centos-7-install"
             echo
             [[ "$FAST_DOWNLOADER" == "true" ]] && show_fast_isos "ubuntu-16.04-install" "centos-7-install" "sledgehammer"

         ;;
     remove)
         if [[ $ISOLATED == true ]] ; then
             echo "Remove the directory that the initial isolated install was done in."
             exit_cleanup 0
         fi
         if pgrep dr-provision; then
             echo "'dr-provision' service is running, CAN NOT remove ... please stop service first"
             exit_cleanup 9
         else
             echo "'dr-provision' service is not running, beginning removal process ... "
         fi
         if [[ -f "$CLI_BKUP" ]]
         then
           echo "Restoring original 'drpcli'."
           $_sudo mv "$CLI_BKUP" "$CLI"
           RM_CLI=""
         else
           RM_CLI="$bindest/drpcli"
           echo "No 'drpcli' backup file found ('$CLI_BKUP')."
         fi
         echo "Removing program and service files"
         $_sudo rm -f "$bindest/dr-provision" "$RM_CLI" "$initdest"
         [[ -d /etc/systemd/system/dr-provision.service.d ]] && rm -rf /etc/systemd/system/dr-provision.service.d
         [[ -f ${BIN_DIR}/drp-install.sh ]] && rm -f ${BIN_DIR}/drp-install.sh
         if [[ $REMOVE_DATA == true ]] ; then
             printf "Removing data files and directories ... "
             [[ -d "/usr/share/dr-provision" ]] && RM_DIR="/usr/share/dr-provision " || true
             [[ -d "/etc/dr-provision" ]] && RM_DIR+="/etc/dr-provision " || true
             [[ -d "${DRP_HOME_DIR}" ]] && RM_DIR+="${DRP_HOME_DIR}" || true
             echo "$RM_DIR"
             $_sudo rm -rf $RM_DIR
         fi
         ;;
     version) echo "Installer Version: $VERSION" ;;
     *)
         echo "Unknown action \"$1\". Please use 'install', 'upgrade', or 'remove'";;
esac

exit_cleanup 0
