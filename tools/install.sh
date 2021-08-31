#!/usr/bin/env bash
# Copyright (c) 2021 RackN Inc.

set -e

###
#  Installation, upgrade, and removal script for the Digital Rebar Platform (DRP)
#  service.  Use '--help' for Usage information.  Supports online and offline (eg
#  "airgap" installs if the DRP, content, and plugin pieces are local to the installer
#  script, or accessible via an alternate HTTP/S location.
###

# BUMP version on updates
VERSION="v21.08.30-1"

DEFAULT_DRP_VERSION=${DEFAULT_DRP_VERSION:-"stable"}

exit_cleanup() {
  local _x=$1
  shift
  [[ -n "$*" ]] && echo -e "EXIT MESSAGE: $*"
  exit $_x
}

COLOR_OK=true
set_color() {
# terminal colors
RCol='\e[0m'    # Text Reset

# Regular           Bold                Underline           High Intensity      BoldHigh Intens     Background          High Intensity Backgrounds
Bla='\e[0;30m';     BBla='\e[1;30m';    UBla='\e[4;30m';    IBla='\e[0;90m';    BIBla='\e[1;90m';   On_Bla='\e[40m';    On_IBla='\e[0;100m';
Red='\e[0;31m';     BRed='\e[1;31m';    URed='\e[4;31m';    IRed='\e[0;91m';    BIRed='\e[1;91m';   On_Red='\e[41m';    On_IRed='\e[0;101m';
Gre='\e[0;32m';     BGre='\e[1;32m';    UGre='\e[4;32m';    IGre='\e[0;92m';    BIGre='\e[1;92m';   On_Gre='\e[42m';    On_IGre='\e[0;102m';
Yel='\e[0;33m';     BYel='\e[1;33m';    UYel='\e[4;33m';    IYel='\e[0;93m';    BIYel='\e[1;93m';   On_Yel='\e[43m';    On_IYel='\e[0;103m';
Blu='\e[0;34m';     BBlu='\e[1;34m';    UBlu='\e[4;34m';    IBlu='\e[0;94m';    BIBlu='\e[1;94m';   On_Blu='\e[44m';    On_IBlu='\e[0;104m';
Pur='\e[0;35m';     BPur='\e[1;35m';    UPur='\e[4;35m';    IPur='\e[0;95m';    BIPur='\e[1;95m';   On_Pur='\e[45m';    On_IPur='\e[0;105m';
Cya='\e[0;36m';     BCya='\e[1;36m';    UCya='\e[4;36m';    ICya='\e[0;96m';    BICya='\e[1;96m';   On_Cya='\e[46m';    On_ICya='\e[0;106m';
Whi='\e[0;37m';     BWhi='\e[1;37m';    UWhi='\e[4;37m';    IWhi='\e[0;97m';    BIWhi='\e[1;97m';   On_Whi='\e[47m';    On_IWhi='\e[0;107m';

# palette
CWarn="$BYel"
CFlag="$IBlu"
CDef="$IBla"
CFile="$IBla"
CNote="$Yel"
CInfo="$Cya"
CNotice="$IYel"
COk="$IGre"
CErr="$IRed"
}

c_def() { echo -en "$CDef$@$RCol";}
c_flag() { echo -en "$CFlag$@$RCol";}
c_warn() { echo -en "$CWarn$@$RCol";}
c_file() { echo -en "$CFile$@$RCol";}
c_err() { echo -en "$CErr$@$RCol";}

_drpcli() {
    if [[ $DBG == true ]]; then
        drpcli "$@"
    else
        drpcli "$@" >/dev/null
    fi
}

usage() {
    [[ "$COLOR_OK" == "true" ]] && set_color
echo -e "
  ${ICya}USAGE:${RCol} ${BYel}$0${RCol} $(c_flag "[ install | upgrade | remove | version ] [ <argument flags: see below> ]")

${CWarn}WARNING${RCol}: '$(c_flag "install")' option will OVERWRITE existing installations

${ICya}OPTIONS${RCol}:
    $(c_flag "install")                 - Sets up an isolated or system 'production' enabled install
    $(c_flag "upgrade")                 - Sets the installer to upgrade an existing 'dr-provision', for upgrade of
                              container; kill/rm the DRP container, then upgrade and reattach data volume
    $(c_flag "remove")                  - Removes the system enabled install.  Requires no other flags
                              optional: '$(c_flag "--remove-data")' to wipe all installed data
    $(c_flag "version")                 - Show install.sh script version and exit

${ICya}ARGUMENT FLAGS${RCol}:
    $(c_flag "--debug")=$(c_def "[true|false]")    - Enables debug output
    $(c_flag "--force")=$(c_def "[true|false]")    - Forces an overwrite of local install binaries and content
    $(c_flag "--upgrade")=$(c_def "[true|false]")  - Turns on 'force' option to overwrite local binaries/content
    $(c_flag "--isolated")              - Sets up current directory as install location for drpcli
                              and dr-provision (makes mess in current directory!)
    $(c_flag "--no-content")            - Don't add content to the system
    $(c_flag "--zip-file")=$(c_def "filename.zip") - Don't download the dr-provision.zip file, instead use
                              the referenced zip file (useful for airgap deployments)
                              NOTE: disables sha256sum checks - do this manually
    $(c_flag "--ipaddr")=$(c_def "<ip>")           - The IP to use for the system identified IP.  The system
                              will attempt to discover the value if not specified
    $(c_flag "--version")=$(c_def "<string>")      - Version identifier if downloading; stable, tip, or
                              specific version label, $(c_def "(defaults to: $DEFAULT_DRP_VERSION)")
    $(c_flag "--remove-data")           - Remove data as well as program pieces
    $(c_flag "--skip-run-check")        - Skip the process check for 'dr-provision' on new install
                              only valid in '$(c_flag "--isolated")' install mode
    $(c_flag "--skip-prereqs")          - Skip OS dependency checks, for testing '$(c_flag "--isolated")' mode
    $(c_flag "--no-sudo")               - Do not use \"sudo\" prefix on commands (assume you're root)
    $(c_flag "--no-colors")             - Installer output will have no terminal colors
    $(c_flag "--fast-downloader")       - (experimental) Use Fast Downloader (uses 'aria2')
    $(c_flag "--keep-installer")        - In Production mode, do not purge the tmp installer artifacts
    $(c_flag "--startup")               - Attempt to start the dr-provision service
    $(c_flag "--systemd")               - Run the systemd enabling commands after installation
    $(c_flag "--systemd-services")      - Additional services for systemd to wait for before starting DRP.
    $(c_flag "--create-self")           - DRP will create a machine that represents itself.
                              Only used with startup/systemd parameters.
    $(c_flag "--start-runner")          - DRP will start a runner for itself. Implies create self.
                              Only used with startup/systemd parameters.
    $(c_flag "--initial-workflow")=$(c_def "<string>")
                            - Workflow to assign to the DRP's self machine as install finishes
                              Only valid with create-self object, only one Workflow may be specified
    $(c_flag "--initial-profiles")=$(c_def "<string>")
                            - Initial profiles to add to the DRP endpoint before starting the workflow,
                              comma separated list, no spaces, this is only valid with create-self object
    $(c_flag "--initial-contents")=$(c_def "<string>")
                            - Initial content packs to deliver, comma separated with no spaces.
                              A file, URL, or content-pack name
    $(c_flag "--initial-plugins")=$(c_def "<string>")
                            - Initial plugins to deliver, comma separated list with no spaces.
                              A file, URL, or content-pack name
    $(c_flag "--initial-parameters")=$(c_def "<string>")
                            - Initial parameters to set on the system.  Simple parameters
                              as a comma separated list, with no spaces.
    $(c_flag "--initial-subnets")=$(c_def "<string>")
                            - A file or URL containing Subnet definitions in JSON or YAML format,
                              comma separated, no spaces
                              NOTE: Subnets can also be injected in content packs with the
                                '$(c_flag "--initial-contents")' argument
    $(c_flag "--bootstrap")             - Store the install image and the install script in the files bootstrap
    $(c_flag "--drp-id")=$(c_def "<string>")       - String to use as the DRP Identifier $(c_def "(only with $(c_flag "--systemd")${CDef})")
    $(c_flag "--drp-user")=$(c_def "<string>")     - DRP user to create after system start $(c_def "(only with $(c_flag "--systemd")${CDef})")
    $(c_flag "--drp-password")=$(c_def "<string>") - DRP user password to set after system start $(c_def "(only with $(c_flag "--systemd")${CDef})")
    $(c_flag "--remove-rocketskates")   - Remove the rocketskates user after system start $(c_def "(only with $(c_flag "--systemd")${CDef})")
    $(c_flag "--local-ui")              - Set up DRP to server a local UI
    $(c_flag "--ha-id")=$(c_def "<string>")        - $(c_err "DEPRECATED pre-v4.6 ONLY"): String to use as the HA Identifier $(c_def "(only with $(c_flag "--systemd")${CDef})")
    $(c_flag "--ha-enabled")            - $(c_err "DEPRECATED pre-v4.6 ONLY"): Indicates that the system is HA enabled
    $(c_flag "--ha-address")=$(c_def "<string>")   - $(c_err "DEPRECATED pre-v4.6 ONLY"): IP Address to use a VIP for HA system
    $(c_flag "--ha-interface")=$(c_def "<string>") - $(c_err "DEPRECATED pre-v4.6 ONLY"): Interrace to use for HA traffic
    $(c_flag "--ha-passive")            - $(c_err "DEPRECATED pre-v4.6 ONLY"): Indicates that the system is starting as passive.
    $(c_flag "--ha-token")=$(c_def "<string>")     - $(c_err "DEPRECATED pre-v4.6 ONLY"): The token to use to sync passive to active
    $(c_flag "--system-user")           - System user account to create for DRP to run as
    $(c_flag "--system-group")          - System group name
    $(c_flag "--drp-home-dir")          - Use with system-user and system-group to set the home directory
                              for the system-user. This path is where most important drp files live
                              including the tftp root
    $(c_flag "--bin-dir")               - Use this as the local of the binaries.  Required for non-root upgrades
    $(c_flag "--container")             - Force to install as a container, not zipfile
    $(c_flag "--container-type")=$(c_def "<string>")
                            - Container install type, $(c_def "(defaults to \"$CNT_TYPE\")")
    $(c_flag "--container-name")=$(c_def "<string>")
                            - Set the \"docker run\" container name, $(c_def "(defaults to \"$CNT_NAME\")")
    $(c_flag "--container-restart")=$(c_def "<string>")
                            - Set the Docker restart option, $(c_def "(defaults to \"$CNT_RESTART\")")
                              options are:  no, on-failure, always, unless-stopped
                            * see: https://docs.docker.com/config/containers/start-containers-automatically/
    $(c_flag "--container-volume")=$(c_def "<string>")
                            - Volume name to use for backing persistent storage, $(c_def "(default to \"$CNT_VOL\")")
    $(c_flag "--container-registry")=$(c_def "\"drp.example.com:5000\"")
                            - Alternate registry to get container images from, $(c_def "(default to \"$CNT_REGISTRY\")")
    $(c_flag "--container-env")=$(c_def "\"<string> <string> <string>\"")
                            - Define a space separated list of environment variables to pass to the
                              container on start $(c_def "(eg \"RS_METRICS_PORT=8888 RS_DRP_ID=fred\")")
                              see 'dr-provision --help' for complete list of startup variables
    $(c_flag "--container-netns")=$(c_def "\"<string>\"")
                            - Define Network Namespace to start container in. $(c_def "(defaults to \"$CNT_NETNS\")")
                              If set to empty string (\"\"), then disable setting any network namespace
    $(c_flag "--universal")             - Load the universal components and bootstrap the system.
                              $(c_warn "**This must be first**") all other argument flags must be after this flag;
                              and implies systemd, startup, start-runner, create-self.  Also implies
                              running universal-boostrap and starting discovery.  Subsequent options for
                              '$(c_flag "--initial-workflow")' will be ignored if this is set.

${ICya}DEFAULTS${RCol}:
    |  argument flag:        def. value:      |  argument flag:        def. value:
    |  -------------------   ------------     |  ------------------    ------------
    |  remove-rocketskates = $(c_def "false")            |  version (*)         = $(c_def "$DEFAULT_DRP_VERSION")
    |  isolated            = $(c_def "false")            |  nocontent           = $(c_def "false")
    |  upgrade             = $(c_def "false")            |  force               = $(c_def "false")
    |  debug               = $(c_def "false")            |  skip-run-check      = $(c_def "false")
    |  skip-prereqs        = $(c_def "false")            |  systemd             = $(c_def "false")
    |  create-self         = $(c_def "false")            |  start-runner        = $(c_def "false")
    |  drp-id              = $(c_def "unset")            |  ha-id               = $(c_def "unset")
    |  drp-user            = $(c_def "rocketskates")     |  drp-password        = $(c_def "r0cketsk8ts")
    |  startup             = $(c_def "false")            |  keep-installer      = $(c_def "false")
    |  local-ui            = $(c_def "false")            |  system-user         = $(c_def "root")
    |  system-group        = $(c_def "root")             |  drp-home-dir        = $(c_def "/var/lib/dr-provision")
    |  bin-dir             = $(c_def "/usr/local/bin")   |  container           = $(c_def "false")
    |  container-volume    = $(c_def "$CNT_VOL")         |  container-registry  = $(c_def "$CNT_REGISTRY")
    |  container-type      = $(c_def "$CNT_TYPE")           |  container-name      = $(c_def "$CNT_NAME")
    |  container-netns     = $(c_def "$CNT_NETNS")             |  container-restart   = $(c_def "$CNT_RESTART")
    |  bootstrap           = $(c_def "false")            |  initial-workflow    = $(c_def "unset")
    |  initial-contents    = $(c_def "unset")            |  initial-profiles    = $(c_def "unset")
    |  initial-plugins     = $(c_def "unset")            |  initial-subnets     = $(c_def "unset")
    |  ha-enabled          = $(c_def "false")            |  ha-address          = $(c_def "unset")
    |  ha-passive          = $(c_def "false")            |  ha-interface        = $(c_def "unset")
    |  ha-token            = $(c_def "unset")            |  universal           = $(c_def "unset")
    |  systemd-services    = $(c_def "unset")

    * version examples: '$(c_def "tip")', '$(c_def "v4.6.3")', '$(c_def "v4.7.0-beta1.3")', or '$(c_def "stable")'

${ICya}PREREQUISITES${RCol}:
    ${CNote}NOTE: By default, prerequisite packages will be installed if possible.  You must
          ${CNote}manually install these first on a Mac OS X system. Package names may vary
          ${CNote}depending on your operating system version/distro packaging naming scheme.${RCol}

    ${ICya}REQUIRED${RCol}: curl, tar
    ${ICya}OPTIONAL${RCol}: aria2c (if using experimental "fast downloader")

${CWarn}WARNING${RCol}: '$(c_flag "install")' option will OVERWRITE existing installations, use '$(c_flag "update")'
         flag for inplace upgrades.

${ICya}INSTALLER VERSION${RCol}:  $(c_def "$VERSION")
"
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

###
#  '--universal' must be first argument flag if specified, make it so
###
if [[ "$*" =~ .*--universal.* ]]
then
  set -- "--universal" $(echo "$*" | sed 's/ --universal//g')
fi

args=()
while (( $# > 0 )); do
    arg="$1"
    arg_key="${arg%%=*}"
    arg_data="${arg#*=}"
    case $arg_key in
        --help|-h)                  usage; exit_cleanup 0                             ;;
        --debug)                    DBG=true;                                         ;;
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
        --no-colors)                COLOR_OK=false                                    ;;
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
                                    SYSTEMD=true; STARTUP=true;
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
            echo -e "$PREF_INFO Overriding $arg_key with $arg_data"
            export $arg_key="$arg_data"
            ;;
        *)
            args+=("$arg");;
    esac
    shift
done

[[ "$COLOR_OK" == "true" ]] && set_color
PREF_OK="$COk>>>$RCol"
PREF_ERR="$CErr!!!$RCol"
PREF_INFO="$CInfo###$RCol"
PREF_WARN="${CWarn}Warning$RCol:"

set -- "${args[@]}"

if [[ "$HA_ENABLED" == "true" ]] ; then
  if [[ "$HA_TOKEN" == "" ]] ; then
    if [[ "$HA_PASSIVE" == "true" ]] ; then
      usage; exit_cleanup 1 "Passive systems must have a token. Specify HA_TOKEN."
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
        exit_cleanup 1 "$(c_err "FATAL"): Script is not running as root and sudo command is not found. Please be root"
    fi
fi

# verify the container restart directives
case $CNT_RESTART in
  no|on-failure|always|unless-stopped) true ;;
  *) exit_cleanup 1 "$(c_err "FATAL"): Container restart directive ('$CNT_RESTART') is not valid. See '$(c_flag "--help")${CErr}' for details."
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
        exit_cleanup 1 "$(c_err "FATAL"): Cannot determine Linux version we are running on!"
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
        if (( $COUNT == 0 )); then
            echo -e -n "$PREF_INFO Waiting for dr-provision to start ..."
        else
            echo -n "."
        fi
        sleep 2
        # Pre-increment for compatibility with Bash 4.1+
        ((++COUNT))
        if (( $COUNT > 10 )) ; then
            exit_cleanup 1 "$(c_err "FATAL"): DRP Failed to start"
        fi
    done
    echo
    INFO=$(drpcli info get)
    echo -e "${COk}${EP}Digital Rebar $(jq -r .version <<< "$INFO") started as Endpoint ID $(jq -r .id <<< "$INFO")${RCol}"
    echo
}

# install the EPEL repo if appropriate, and not enabled already
install_epel() {
    if [[ $OS_FAMILY == rhel ]] ; then
        if ( `yum repolist enabled | grep -q "^epel/"` ); then
            echo -e "$PREF_INFO EPEL repository installed already."
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
    [[ -z "$*" ]] && exit_cleanup 1 "$(c_err "FATAL"): Internal error, get() expects files to get"

    if [[ "$FAST_DOWNLOADER" == "true" ]]; then
        if which aria2c > /dev/null; then
            GET="aria2c --quiet=true --continue=true --max-concurrent-downloads=10 --max-connection-per-server=16 --max-tries=0"
        else
            exit_cleanup 1 "$(c_err "FATAL"): '$(c_flag "--fast-downloader")' specified, but couldn't find tool ('aria2c')."
        fi
    else
        if which curl > /dev/null; then
            GET="curl -sfL"
        else
            exit_cleanup 1 "$(c_err "FATAL"): Unable to find downloader tool ('curl')."
        fi
    fi
    for URL in $*; do
        FILE=${URL##*/}
        echo -e "$PREF_OK Downloading file:  ${BWhi}$FILE${RCol}"
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
        exit_cleanup ${RC} "$(c_err "FATAL"): Unable to create system group ${SYSTEM_GROUP}"
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
    exit_cleanup ${RC} "$(c_err "FATAL"): Unable to create system user ${SYSTEM_USER}"
}

set_ownership_of_drp() {
    # It is possible for the home directory to not exist if
    # a non-root user was specified but already created.
    # Make sure a directory is created so DRP does not hit
    # permissions errors trying to use the home directory.
    if [ ! -d "${DRP_HOME_DIR}" ]; then
        echo -e "$PREF_OK Creating DRP Home directory $(c_def "(${DRP_HOME_DIR})")"
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
                echo -e "$PREF_INFO Your OS Family ${OS_FAMILY} does not support setcap" \
                     "  and may not be able to bind privileged ports when" \
                     "  running as non-root user ${SYSTEM_USER}"
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
                    echo -e "$PREF_ERR Must have binary '$_bin' installed."
                    echo "eg:   brew install $_bin"
                    error=1
                fi
            fi
          ;;
        *)
            if ! which $_bin > /dev/null 2>&1; then
                echo -e "$PREF_ERR Must have binary '$_bin' installed."
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
            echo -e "$PREF_INFO Missing dependency '$_pkg', attempting to install it... "
            if [[ $OS_FAMILY == rhel ]] ; then
                echo $IN_EPEL | grep -q $_pkg && install_epel
                $_sudo yum install -y $_pkg
            elif [[ $OS_FAMILY == debian ]] ; then
                $_sudo apt-get install -y $_pkg
            fi
        fi
} # end check_pkgs_linux()

ensure_packages() {
    echo -e "$PREF_OK Ensuring required tools are installed"
    case $OS_FAMILY in
        darwin)
            error=0
            BINS="curl aria2c"
            for BIN in $BINS; do
                check_bins_darwin $BIN
            done

            if [[ $error == 1 ]]; then
                echo -e "$PREF_ERR After install missing components, restart the terminal to pick"
                echo "  up the newly installed commands, and re-run the installer."
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
            echo -e "$PREF_INFO CoreOS does not require any packages to be installed.  DRP will be"
            echo "  installed from the Docker Hub registry."
        ;;
        photon)
            if ! which tar > /dev/null 2>&1; then
                echo -e "$PREF_OK Installing packages for Photon Linux..."
                tdnf -y makecache
                tdnf -y install tar
            else
                echo -e "$PREF_INFO 'tar' already installed on Photon Linux..."
            fi
        ;;
        *)
            exit_cleanup 1 "$(c_err "FATAL"): Unsupported OS Family ($OS_FAMILY)."
        ;;
    esac
} # end ensure_packages()

# output a friendly statement on how to download ISOS via fast downloader
show_fast_isos() {
    echo -e "
$PREF_INFO Option '$(c_flag "--fast-downloader")' requested.  You may download the ISO images using
  'aria2c' command to significantly reduce download time of the ISO images.

${CNote} NOTE: The following genereted scriptlet should download, install, and enable
      ${CNote}the ISO images.  VERIFY SCRIPTLET before running it.${RCol}

      YOU MUST START 'dr-provision' FIRST! Example commands:

$(c_def "###### BEGIN scriptlet")
  export CMD=\"aria2c --continue=true --max-concurrent-downloads=10 --max-connection-per-server=16 --max-tries=0\"
"

    for BOOTENV in $*
    do
        echo "  export URL=$(${EP}drpcli bootenvs show $BOOTENV | grep 'IsoUrl' | cut -d '"' -f 4)"
        echo "  export ISO=$(${EP}drpcli bootenvs show $BOOTENV | grep 'IsoFile' | cut -d '"' -f 4)"
        echo "  \$CMD -o \$ISO \$URL"
    done
    echo "  # this should move the ISOs to the TFTP directory..."
    echo "  $_sudo mv *.tar *.iso $TFTP_DIR/isos/"
    echo "  $_sudo pkill -HUP dr-provision"
    echo "  echo 'NOTICE:  exploding isos may take up to 5 minutes to complete ... '"
    echo -e $(c_def "###### END scriptlet")

    echo
} # end show_fast_isos()

remove_container() {
    case $CNT_TYPE in
        docker)
            $_sudo docker kill $CNT_NAME > /dev/null && echo -e "$PREF_OK Killed docker container '$CNT_NAME'"
            $_sudo docker rm $CNT_NAME > /dev/null && echo -e "$PREF_OK Removed docker container '$CNT_NAME'"
            if [[ "$CNT_VOL_REMOVE" == "true" ]]; then
                $_sudo docker volume rm $CNT_VOL > /dev/null && echo -e "$PREF_OK Removed backing volume named '$CNT_VOL'"
            else
                echo -e "$PREF_OK Not removing container data volume '$CNT_VOL'"
                if [[ "$MODE" == "remove" ]]; then
                    echo "    To remove: '$_sudo docker volume rm $CNT_VOL'"
                    echo "    or use '--remove-data' next time ... "
                fi
            fi
            ;;
        *)  exit_cleanup 1 "$PREF_ERR Container type '$CNT_TYPE' not supported in installer."
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
            ! which docker > /dev/null 2>&1 && exit_cleanup 1 "$PREF_ERR Container install requested but no 'docker' in PATH ($PATH)."
            if [[ "$UPGRADE" == "false" ]]; then
                $_sudo docker volume create $CNT_VOL > /dev/null
                VOL_MNT=$($_sudo docker volume inspect $CNT_VOL | grep Mountpoint | awk -F\" '{ print $4 }')
                echo -e "$PREF_OK Created docker volume named '$CNT_VOL' with mountpoint '$VOL_MNT'"
            else
                if $_sudo docker volume inspect $CNT_VOL > /dev/null 2>&1; then
                    echo -e "$PREF_INFO Attempting to reconnect volume '$CNT_VOL'"
                else
                    exit_cleanup 1 "$PREF_ERR No existing volume '$CNT_VOL' found to reconnect."
                fi
            fi
            if [[ -z "$CNT_NETNS" ]]; then
              echo -e "$PREF_WARN No network namespace set - you may have issues with DHCP and TFTP depending on your use case."
            else
              NETNS="--net $CNT_NETNS"
            fi
            CMD="$_sudo docker run $ENV_OPTS --restart \"$CNT_RESTART\" --volume $CNT_VOL:/provision/drp-data --name \"$CNT_NAME\" -itd $NETNS ${CNT_REGISTRY}/digitalrebar/provision:$DRP_VERSION"
            echo -e "$PREF_INFO Starting container with following run command:"
            echo "$CMD"
            eval $CMD

            echo ""
            echo -e "$PREF_INFO Digital Rebar container is using backing volume: $CNT_VOL"
            echo "Volume is backed on host filesystem at: $VOL_MNT"
            echo ""
            echo -e "$PREF_INFO Docker container run time information:"
            $_sudo docker ps --filter name=$CNT_NAME
            echo ""
            echo -e "${CNotice}>>>  NOTICE:  If you intend to upgrade this container, record your 'install.sh' options  <<< ${RCol}"
            echo -e "${CNotice}>>>           for the 'upgrade' command - you must reuse the same options to obtain the  <<< ${RCol}"
            echo -e "${CNotice}>>>           same installed results in the upgraded container.  You have been warned!   <<< ${RCol}"
            echo ""
            ;;
        *)  exit_cleanup 1 "$PREF_ERR Container type '$CNT_TYPE' not supported in installer."
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
            echo -e "$PREF_OK Installing Digital Rebar as a container."
            install_container
            exit_cleanup $?
            ;;
        upgrade)
            echo -e "$PREF_OK Upgrading Digital Rebar as a container."
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
    *)            exit_cleanup 1 "$(c_err "FATAL"): architecture ('$arch') not supported" ;;
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
            echo -e "$PREF_ERR No idea how to install startup stuff -- not using systemd, upstart, or sysv init"
            exit_cleanup 1
        fi
        shasum="command sha256sum";;
    *)
        # Someday, support installing on Windows.  Service creation could be tricky.
        echo -e "$PREF_ERR No idea how to check sha256sums"
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
                     echo -e "$PREF_ERR 'dr-provision' service is running, CANNOT upgrade ... please stop service or container first"
                     echo -e "  Run ${IYel}$_sudo systemctl stop dr-provision${RCol}"
                     exit_cleanup 9
                 else
                     echo -e "$PREF_INFO 'dr-provision' service is not running, beginning install process ... "
                 fi
             else
                 echo -e "$PREF_INFO Skipping 'dr-provision' service run check as requested ..."
             fi

             [[ "$SKIP_DEPENDS" == "false" ]] && ensure_packages || echo -e "$PREF_INFO Skipping dependency checks as requested ... "

             if [[ "$ISOLATED" == "false" ]]; then
                 TMP_INSTALLER_DIR=$(mktemp -d /tmp/drp.installer.XXXXXX)
                 echo -e "$PREF_INFO Using temp directory to extract artifacts to and install from ('$(c_file "$TMP_INSTALLER_DIR")')."
                 OLD_PWD=$(pwd)
                 cd $TMP_INSTALLER_DIR
                 TMP_INST=$TMP_INSTALLER_DIR/tools/install.sh
             fi

             # Are we in a build tree
             if [ -e server ] ; then
                 if [ ! -e bin/linux/amd64/drpcli ] ; then
                     echo "It appears that nothing has been built."
                     echo -e "Please run $(c_flag "tools/build.sh") and then rerun this command".
                     exit_cleanup 1
                 fi
             else
                 # We aren't a build tree, but are we extracted install yet?
                 # If not, get the requested version.
                 if [[ ! -e sha256sums || $force ]] ; then
                     echo -e "$PREF_OK Installing Version ${IGre}$DRP_VERSION${RCol} of ${ICya}Digital Rebar${RCol} (dr-provision)"
                     if [[ $DRP_ID ]] ; then
                        echo -e "$PREF_INFO Using Endpoint ID ${IGre}$DRP_ID${RCol}."
                     else
                        echo -e "$PREF_INFO Using MAC based Endpoint ID.  Consider using '$(c_flag "--drp-id")' to set during installation."
                     fi
                     ZIP="dr-provision.zip"
                     SHA="dr-provision.sha256"
                     if [[ -n "$ZIP_FILE" ]]
                     then
                       [[ "$ZIP_FILE" != "dr-provision.zip" ]] && cp "$ZIP_FILE" dr-provision.zip
                       echo -e "$PREF_WARN  No sha256sum check performed for '$(c_flag "--zip-file")' mode."
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
                           SOURCE=${SOURCE/.zip/.${arch}.${os}.zip}
                           get $SOURCE
                           FILE=${SOURCE##*/}
                       fi
                       mv $FILE $ZIP

                       # XXX: Put sha back one day
                     fi

                     if ! $tar -xf dr-provision.zip; then
                       # try to fall back to bsdtar
                       echo -e "$PREF_INFO Attempting to auto recover from above 'tar' errors."
                       echo -e "$PREF_INFO Newer 'install.sh' script being used with older dr-provision.zip file."
                       echo -e "$PREF_INFO Attempting to locate and use 'bsdtar' as a fallback..."

                       if which bsdtar > /dev/null; then
                         if ! bsdtar -xf dr-provision.zip; then
                           echo -e "$PREF_ERR FAILED to extract 'dr-provision.zip' with installed 'bsdtar'."
                           echo -e "$PREF_ERR Not cleaning up for forensic analysis."
                           exit 1
                         fi
                       else
                         echo -e "$PREF_OK 'bsdtar' not found in path, attempting to install it... "
                         if check_pkgs_linux bsdtar; then
                           if ! bsdtar -xf dr-provision.zip; then
                             echo -e "$PREF_ERR Last ditch attempt FAILED to unpack 'dr-provision.zip' file."
                             echo -e "$PREF_ERR Not attempting cleanup for debug and troubleshooting reasons."
                             exit 1
                           fi
                         else
                           exit_cleanup 1 "$PREF_ERR FAILED to install 'bsdtar' to satisfy dependencies.  Please install it and retry."
                         fi
                       fi
                     fi
                 fi
                 $shasum -c sha256sums || exit_cleanup 1
             fi

             if [[ $NO_CONTENT == false ]]; then
                 echo -e "$PREF_OK Installing Version ${IGre}$DRP_CONTENT_VERSION${RCol} of ${ICya}Digital Rebar Community Content${RCol}"
                 if [[ -n "$ZIP_FILE" ]]; then
                   echo -e "$PREF_WARN '$(c_flag "--zip-file")' specified, still trying to download community content..."
                   echo "         (specify '$(c_flag "--no-content")' to skip download of community content"
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
                         echo -e "$PREF_WARN"
                         echo -e "${CWarn}  initfile ('$initfile') exists already, not overwriting it${RCol}"
                         echo -e "${CWarn}  please verify 'dr-provision' startup options are correct${RCol}"
                         echo -e "${CWarn}  for your environment and the new version .. ${RCol}"
                         echo
                         echo -e "${CWarn}  specifically verify: '${CFlag}--file-root=$(c_def "<tftpboot directory>${CWarn}'")"
                     else
                         $_sudo sed "s:/usr/local/bin/dr-provision:$PROVISION:g" "$initfile" > "$initdest"
                     fi
                     # output our startup helper messages only if SYSTEMD isn't specified
                     if [[ "$SYSTEMD" == "false" || "$STARTUP" == "false" ]]; then
                        echo
                        echo -e "$PREF_INFO Tips for ${BIYel}systemd${RCol} management of ${ICya}Digital Rebar${RCol} service:"
                        echo -e "  Start:  ${IYel}$starter${RCol}"
                        echo -e "  Enable: ${IYel}$enabler${RCol}"
                        echo -e "  Status: ${IYel}$_sudo systemctl status dr-provision${RCol}"
                        echo -e "  Logs:   ${IYel}$_sudo journalctl -u dr-provision${RCol}"
                    else
                        echo -e "$PREF_INFO Will attempt to execute startup procedures ('$(c_flag "--startup")' specified)"
                        echo -e "  ${IYel}$starter${RCol}"
                        echo -e "  ${IYel}$enabler${RCol}"
                    fi
                 fi

                 if [[ ! -e ${DRP_HOME_DIR}/digitalrebar/tftpboot && -e /var/lib/tftpboot ]] ; then
                     echo -e "$PREF_OK Moving $(c_file "/var/lib/tftpboot") to $(c_file "$DRP_HOME_DIR/tftpboot") location ... "
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
                     echo -e "$PREF_OK Saving '$(c_file "${BIN_DIR}/drpcli")' to backup file ($(c_file "$CLI_BKUP"))"
                     $_sudo mv "$CLI" "$CLI_BKUP"
                 fi

                 INST="${BIN_DIR}/drp-install.sh"
                 $_sudo cp $TMP_INST $INST && $_sudo chmod 755 $INST
                 echo -e  "$PREF_INFO Install script saved to '$(c_file "$INST")'"
                 echo -e "$PREF_INFO You can ${BIYel}uninstall${RCol} DRP with '${IRed}$_sudo $INST remove${RCol}' - must be root)"

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
                             _drpcli contents upload catalog:task-library-${DRP_CONTENT_VERSION}

                             if [[ "$INITIAL_CONTENTS" != "" ]] ; then
                                 OLD_IFS="${IFS}"
                                 IFS=',' read -ra contents_array <<< "$INITIAL_CONTENTS"
                                 for i in "${contents_array[@]}" ; do
                                     echo -e "$PREF_OK Installing content item '$i'"
                                     if [[ -f ${OLD_PWD}/$i ]] ; then
                                         _drpcli contents upload ${OLD_PWD}/${i}
                                     elif [[ -f "$i" ]] ; then
                                         _drpcli contents upload ${i}
                                     elif [[ $i == http* ]] ; then
                                         _drpcli contents upload ${i}
                                     else
                                         _drpcli catalog item install ${i} --version=${DRP_CONTENT_VERSION} -c $DRP_CATALOG
                                     fi
                                 done
                                 IFS="${OLD_IFS}"
                             fi

                             if [[ "$INITIAL_PLUGINS" != "" ]] ; then
                                 OLD_IFS="${IFS}"
                                 IFS=',' read -ra plugins_array <<< "$INITIAL_PLUGINS"
                                 for i in "${plugins_array[@]}" ; do
                                     echo -e "$PREF_OK Installing plugin item '$i'"
                                     if [[ -f ${OLD_PWD}/$i ]] ; then
                                         _drpcli plugin_providers upload ${OLD_PWD}/${i}
                                     elif [[ $i == http* ]] ; then
                                         _drpcli plugin_providers upload ${i}
                                     else
                                         _drpcli catalog item install ${i} --version=${DRP_CONTENT_VERSION} -c $DRP_CATALOG
                                     fi
                                 done
                                 IFS="${OLD_IFS}"
                             fi

                             if [[ "$INITIAL_PROFILES" != "" ]] ; then
                                 if [[ $CREATE_SELF == true ]] ; then
                                     echo -e "$PREF_OK Installing initial profiles..."
                                   cp $(which drpcli) /tmp/jq
                                   chmod +x /tmp/jq
                                   ID=$(drpcli info get | /tmp/jq .id -r | sed -r 's/:/-/g')
                                   rm /tmp/jq
                                   OLD_IFS="${IFS}"
                                   IFS=',' read -ra profiles_array <<< "$INITIAL_PROFILES"
                                   for i in "${profiles_array[@]}" ; do
                                     _drpcli machines addprofile "Name:$ID" "$i"
                                   done
                                   IFS="${OLD_IFS}"
                                 fi
                             fi

                             if [[ "$INITIAL_PARAMETERS" != "" ]] ; then
                                 if [[ $CREATE_SELF == true ]] ; then
                                   echo -e "$PREF_OK Installing initial parameters..."
                                   cp $(which drpcli) /tmp/jq
                                   chmod +x /tmp/jq
                                   ID=$(drpcli info get | /tmp/jq .id -r | sed -r 's/:/-/g')
                                   rm /tmp/jq
                                   OLD_IFS="${IFS}"
                                   IFS=',' read -ra param_array <<< "$INITIAL_PARAMETERS"
                                   for i in "${param_array[@]}" ; do
                                     IFS="=" read -ra data_array <<< "$i"
                                     _drpcli machines set "Name:$ID" param "${data_array[0]}" to "${data_array[1]}"
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
                                     echo -e "$PREF_OK Setting initial workflow to '$INITIAL_WORKFLOW' for Machine '$ID'"
                                     _drpcli machines workflow "Name:$ID" "$INITIAL_WORKFLOW"
                                 fi
                             fi

                             if [[ "$INITIAL_SUBNETS" != "" ]] ; then
                                 if [[ $CREATE_SELF == true ]] ; then
                                     OLD_IFS="${IFS}"
                                     IFS=',' read -ra subnet_array <<< "$INITIAL_SUBNETS"
                                     for i in "${subnet_array[@]}" ; do
                                       if [[ $i == http* ]] ; then
                                         echo -e "$PREF_OK Creating subnet from URL '$i'"
                                         _drpcli subnets create ${i}
                                       elif [[ -f ${OLD_PWD}/$i ]] ; then
                                         echo -e "$PREF_OK Creating subnet from file '${OLD_PWD}/$i'"
                                         _drpcli subnets create ${OLD_PWD}/${i}
                                       elif [[ -f "$i" ]] ; then
                                         echo -e "$PREF_OK Creating subnet from file '${i}'"
                                         _drpcli subnets create ${i}
                                       else
                                         echo -e "$PREF_WARN unable to read subnet file '$i' or '${OLD_PWD}/$i'; no SUBNET created"
                                       fi
                                     done
                                     IFS="${OLD_IFS}"
                                 fi
                             fi
                         fi

                         if [[ $DRP_USER ]] ; then
                             if _drpcli users exists $DRP_USER >/dev/null 2>/dev/null ; then
                                 _drpcli users update $DRP_USER "{ \"Name\": \"$DRP_USER\", \"Roles\": [ \"superuser\" ] }"
                             else
                                 _drpcli users create "{ \"Name\": \"$DRP_USER\", \"Roles\": [ \"superuser\" ] }"
                             fi
                             _drpcli users password $DRP_USER "$DRP_PASSWORD"
                             if _drpcli info check; then
                                 _drpcli users destroy rocketskates
                             else
                                 echo -e "$(c_err "FATAL"): Newly created user/pass not working - bailing out."
                                 exit 1
                             fi
                             export RS_KEY="$DRP_USER:$DRP_PASSWORD"
                             if [[ $REMOVE_RS == true ]] ; then
                                 _drpcli users destroy rocketskates
                             fi
                         fi
                     fi
                     if [[ "$BOOTSTRAP" == "true" ]] ; then
                         _drpcli files upload dr-provision.zip as bootstrap/dr-provision.zip
                         _drpcli files upload $TMP_INST as bootstrap/install.sh
                     fi
                 else
                     if [[ "$STARTUP" == "true" ]]; then
                         echo -e "$PREF_INFO Attempting startup of 'dr-provision' ('$(c_flag "--startup")' specified)"
                         eval "$enabler"
                         eval "$starter"

                         check_drp_ready

                         if [[ "$NO_CONTENT" == "false" ]] ; then
                             drpcli catalog item install task-library --version=${DRP_CONTENT_VERSION} -c $DRP_CATALOG

                             if [[ "$INITIAL_CONTENTS" != "" ]] ; then
                                 IFS=',' read -ra contents_array <<< "$INITIAL_CONTENTS"
                                 for i in "${contents_array[@]}" ; do
                                     if [[ -f ${OLD_PWD}/$i ]] ; then
                                         _drpcli contents upload ${OLD_PWD}/${i}
                                     elif [[ $i == http* ]] ; then
                                         _drpcli contents upload ${i}
                                     else
                                         _drpcli catalog item install ${i} --version=${DRP_CONTENT_VERSION} -c $DRP_CATALOG
                                     fi
                                 done
                             fi

                             if [[ "$INITIAL_PLUGINS" != "" ]] ; then
                                 IFS=',' read -ra plugins_array <<< "$INITIAL_PLUGINS"
                                 for i in "${plugins_array[@]}" ; do
                                     if [[ -f ${OLD_PWD}/$i ]] ; then
                                         _drpcli plugin_providers upload ${OLD_PWD}/${i}
                                     elif [[ $i == http* ]] ; then
                                         _drpcli plugin_providers upload ${i}
                                     else
                                         _drpcli catalog item install ${i} --version=${DRP_CONTENT_VERSION} -c $DRP_CATALOG
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
                                     _drpcli machines addprofile "Name:$ID" "$i"
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
                                     _drpcli machines set "Name:$ID" param "${data_array[0]}" to "${data_array[1]}"
                                   done
                                 fi
                             fi

                             if [[ "$INITIAL_WORKFLOW" != "" ]] ; then
                                 if [[ $CREATE_SELF == true ]] ; then
                                     cp $(which drpcli) /tmp/jq
                                     chmod +x /tmp/jq
                                     ID=$(drpcli info get | /tmp/jq .id -r)
                                     rm /tmp/jq
                                     _drpcli machines workflow "Name:$ID" "$INITIAL_WORKFLOW"
                                 fi
                             fi
                         fi
                         if [[ "$BOOTSTRAP" == "true" ]] ; then
                             _drpcli files upload dr-provision.zip as bootstrap/dr-provision.zip
                             _drpcli files upload $TMP_INST as bootstrap/install.sh
                         fi
                     fi
                 fi

                 cd $OLD_PWD
                 if [[ "$KEEP_INSTALLER" == "false" ]]; then
                     rm -rf $TMP_INSTALLER_DIR
                 else
                     echo ""
                     echo -e "$PREF_INFO Installer artifacts are in '$(c_file "$TMP_INSTALLER_DIR")' - to purge:"
                     echo -e "  ${IYel}$_sudo rm -rf $TMP_INSTALLER_DIR${RCol}"
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
                     echo -e "${CInfo}********************************************************************************"
                     echo
                     echo -e "$PREF_INFO Run the following commands to start up dr-provision in a local isolated way."
                     echo -e "$PREF_INFO The server will store information and serve files from the drp-data directory."
                     echo
                 else
                     echo
                     echo -e "${CInfo}********************************************************************************"
                     echo
                     echo -e "$PREF_INFO Will attempt to startup the 'dr-provision' service ... "
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
                             echo -e "$PREF_INFO No broadcast route set - this is required for Darwin < 10.9."
                             echo -e "  ${IYel}$_sudo route add 255.255.255.255 $IPADDR${RCol}"
                             echo -e "$PREF_INFO No broadcast route set - this is required for Darwin > 10.9."
                             echo -e "  ${IYel}$_sudo route -n add -net 255.255.255.255 $IPADDR${RCol}"
                     fi
                 fi

                 STARTER="$_sudo ./dr-provision --base-root=`pwd`/drp-data > drp.log 2>&1 &"
                 if [[ "$STARTUP" == "false" ]]; then
                   echo -e "$PREF_INFO Digital Rebar can be started with the following command:"
                   echo -e "  ${IYel}$STARTER${RCol}"
                fi
                 mkdir -p "`pwd`/drp-data/saas-content"
                 if [[ $NO_CONTENT == false ]] ; then
                     DEFAULT_CONTENT_FILE="`pwd`/drp-data/saas-content/default.json"
                     mv drp-community-content.json $DEFAULT_CONTENT_FILE
                 fi

                 if [[ "$STARTUP" == "true" ]]; then
                     eval $STARTER
                     echo -e "$PREF_INFO 'dr-provision' running processes:"
                     ps -eo pid,args -o comm  | grep -v grep | grep dr-provision
                     echo

                     if [[ "$BOOTSTRAP" == "true" ]] ; then
                         _drpcli files upload dr-provision.zip as bootstrap/dr-provision.zip
                         _drpcli files upload tools/install.sh as bootstrap/install.sh
                     fi
                 fi

                 EP="./"
             fi

             echo -e "$PREF_INFO With Digital Rebar started, complete the setup using the ${ICya}System Install Wizard${RCol}"
             if [[ $IPADDR ]] ; then
                 echo -e "  open ${IBlu}https://$IPADDR:8092${RCol} (accept self-signed TLS certificate)"
             else
                 echo -e "  open ${IBlu}https://[host ip]:8092${RCol} (accept self-signed TLS certificate)"
             fi
             echo
             echo -e "$PREF_INFO Or, use the CLI to setup for basic system discovery:"
             echo -e "  ${IYel}${EP}drpcli bootenvs uploadiso sledgehammer${RCol}"
             echo -e "  ${IYel}${EP}drpcli prefs set defaultWorkflow discover-base unknownBootEnv discovery defaultBootEnv sledgehammer defaultStage discover${RCol}"
             if [[ $NO_CONTENT == true ]] ; then
                 echo
                 echo -e "$PREF_INFO Add common utilities (sourced from RackN)"
                 echo -e "  ${IYel}${EP}drpcli contents upload catalog:task-library-$DRP_CONTENT_VERSION${RCol}"
             fi
             echo
             echo -e "$PREF_INFO Optionally, upload popular community operating system ISOs"
             echo -e "  ${IYel}${EP}drpcli bootenvs uploadiso ubuntu-20.04-install${RCol}"
             echo -e "  ${IYel}${EP}drpcli bootenvs uploadiso centos-8-install${RCol}"
             echo
             [[ "$FAST_DOWNLOADER" == "true" ]] && show_fast_isos "ubuntu-20.04-install" "centos-8-install" "sledgehammer"

         ;;
     remove)
         if [[ $ISOLATED == true ]] ; then
             echo -e "$PREF_INFO Remove the directory that the initial isolated install was done in."
             exit_cleanup 0
         fi
         if pgrep dr-provision; then
             echo -e "$PREF_ERR 'dr-provision' service is running and CANNOT be removed. Please stop service or container first."
             echo -e "  Run ${IYel}$_sudo systemctl stop dr-provision${RCol}"
             exit_cleanup 9
         else
             echo -e "$PREF_INFO 'dr-provision' service is not running, beginning removal process..."
         fi
         if [[ -f "$CLI_BKUP" ]]
         then
           echo -e "$PREF_INFO Restoring original 'drpcli'."
           $_sudo mv "$CLI_BKUP" "$CLI"
           RM_CLI=""
         else
           RM_CLI="$bindest/drpcli"
           echo -e "$PREF_INFO No 'drpcli' backup file found ('$(c_file "$CLI_BKUP")')."
         fi
         echo -e "$PREF_OK Removing program and service files"
         $_sudo rm -f "$bindest/dr-provision" "$RM_CLI" "$bindest/drbundler" "$bindest/drpjoin" "$bindest/dr-waltool" "$bindest/drpjq" "$initdest"
         # We might not own jq.
         [[ -L "$bindest/jq" ]] && [[ $(readlink "$bindest/jq") = "$bindest/drpcli" ]] && rm -f "$bindest/jq"
         [[ -d /etc/systemd/system/dr-provision.service.d ]] && rm -rf /etc/systemd/system/dr-provision.service.d
         [[ -f ${BIN_DIR}/drp-install.sh ]] && rm -f ${BIN_DIR}/drp-install.sh
         if [[ $REMOVE_DATA == true ]] ; then
             echo -e "$PREF_OK Removing data files and directories ... "
             [[ -d "/usr/share/dr-provision" ]] && RM_DIR="/usr/share/dr-provision " || true
             [[ -d "/etc/dr-provision" ]] && RM_DIR+="/etc/dr-provision " || true
             [[ -d "${DRP_HOME_DIR}" ]] && RM_DIR+="${DRP_HOME_DIR}" || true
             echo -e "$(c_file "$RM_DIR")"
             $_sudo rm -rf $RM_DIR
         fi
         echo -e "${COk}${EP}Digital Rebar has been removed.${RCol}"
         ;;
     version) echo -e "Installer Version: $(c_def "$VERSION")" ;;
     *)
         echo -e "Unknown action \"$1\". Please use '$(c_flag "install")', '$(c_flag "upgrade")', or '$(c_flag "remove")'";;
esac

exit_cleanup 0
