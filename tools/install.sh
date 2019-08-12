#!/usr/bin/env bash

set -e

DEFAULT_DRP_VERSION=${DEFAULT_DRP_VERSION:-"stable"}

usage() {
cat <<EOFUSAGE
Usage: $0 [--version=<Version to install>] [--no-content]
          [--isolate] [--ipaddr=<ip>] install | upgrade | remove

Options:
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
    --drp-id=<string>       # String to use as the DRP Identifier (only with --systemd)
    --ha-id=<string>        # String to use as the HA Identifier (only with --systemd)
    --drp-user=<string>     # DRP user to create after system start (only with --systemd)
    --drp-password=<string> # DRP user passowrd to set after system start (only with --systemd)
    --remove-rocketskates   # Remove the rocketskates user after system start (only with --systemd)
    --local-ui              # Set up DRP to server a local UI
    --system-user           # System user account to create for DRP to run as
    --system-group          # System group name

    install                 # Sets up an isolated or system 'production' enabled install.
    upgrade                 # Sets the installer to upgrade an existing 'dr-provision'
    remove                  # Removes the system enabled install.  Requires no other flags

Defaults are:
    option:               value:           option:               value:
    -------------------   ------------     ------------------    ------------
    remove-rocketskates = false            version (*)         = $DEFAULT_DRP_VERSION
    isolated            = false            nocontent           = false
    upgrade             = false            force               = false
    debug               = false            skip-run-check      = false
    skip-prereqs        = false            systemd             = false
    drp-id              = unset            ha-id               = unset
    drp-user            = rocketskates     drp-password        = r0cketsk8ts
    startup             = false            keep-installer      = false
    local-ui            = false            system-user         = root
    system-group        = root

    * version examples: 'tip', 'v3.13.6' or 'stable'

Prerequisites:
    NOTE: By default, prerequisite packages will be installed if possible.  You must
          manually install these first on a Mac OS X system. Package names may vary
          depending on your operating system version/distro packaging naming scheme.

    REQUIRED: 7zip, curl, jq, bsdtar
    OPTIONAL: aria2c (if using experimental "fast downloader")
EOFUSAGE

exit 0
}

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
STARTUP=false
REMOVE_RS=false
LOCAL_UI=false
KEEP_INSTALLER=false
BIN_DIR=/usr/local/bin
DRP_HOME_DIR=/var/lib/dr-provision

_sudo="sudo"
CLI="${BIN_DIR}/drpcli"
CLI_BKUP="${BIN_DIR}/drpcli.drp-installer.backup"
PROVISION="${BIN_DIR}/dr-provision"

# download URL locations; overridable via ENV variables
URL_BASE=${URL_BASE:-"https://rebar-catalog.s3-us-west-2.amazonaws.com"}
URL_BASE_DRP=${URL_BASE_DRP:-"$URL_BASE/drp"}
URL_BASE_CONTENT=${URL_BASE_CONTENT:-"$URL_BASE/drp-community-content"}
DRP_CATALOG=${DRP_CATALOG:-"$URL_BASE/rackn-catalog.json"}
SYSTEM_USER=root
SYSTEM_GROUP=root

args=()
while (( $# > 0 )); do
    arg="$1"
    arg_key="${arg%%=*}"
    arg_data="${arg#*=}"
    case $arg_key in
        --help|-h)
            usage
            exit 0
            ;;
        --debug)
            DBG=true
            ;;
        --version|--drp-version)
            DRP_VERSION=${arg_data}
            ;;
        --zip-file)
            ZF=${arg_data}
            ZIP_FILE=$(echo "$(cd $(dirname $ZF) && pwd)/$(basename $ZF)")
            ;;
        --isolated)
            ISOLATED=true
            ;;
        --skip-run-check)
            SKIP_RUN_CHECK=true
            ;;
        --skip-dep*|--skip-prereq*)
            SKIP_DEPENDS=true
            ;;
        --fast-downloader)
            FAST_DOWNLOADER=true
            ;;
        --force)
            force=true
            ;;
        --remove-data)
            REMOVE_DATA=true
            ;;
        --commit)
            # UNUSED NOW
            ;;
        --upgrade)
            UPGRADE=true
            force=true
            ;;
        --nocontent|--no-content)
            NO_CONTENT=true
            ;;
        --no-sudo)
            _sudo=""
            ;;
        --keep-installer)
            KEEP_INSTALLER=true
            ;;
        --startup)
            STARTUP=true
            SYSTEMD=true
            ;;
        --systemd)
            SYSTEMD=true
            ;;
        --local-ui)
            LOCAL_UI=true
            ;;
        --remove-rocketskates)
            REMOVE_RS=true
            ;;
        --drp-user)
            DRP_USER=${arg_data}
            ;;
        --drp-password)
            DRP_PASSWORD="${arg_data}"
            ;;
        --drp-id)
            DRP_ID="${arg_data}"
            ;;
        --ha-id)
            HA_ID="${arg_data}"
            ;;
        --system-user)
            SYSTEM_USER="${arg_data}"
            ;;
        --system-group)
            SYSTEM_GROUP="${arg_data}"
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

DRP_VERSION=${DRP_VERSION:-"$DEFAULT_DRP_VERSION"}
DRP_CONTENT_VERSION=stable
if [[ $DRP_VERSION == tip ]]; then
    DRP_CONTENT_VERSION=tip
fi
[[ "$ISOLATED" == "true" ]] && KEEP_INSTALLER=true

[[ $DBG == true ]] && set -x


if [[ $EUID -eq 0 ]]; then
   _sudo=""
fi

if [[ -x "$(command -v sudo)" && $_sudo != "" ]]; then
  echo "Script is not running as root and sudo command is not found. Please be root"
  exit 1
fi


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
        echo "Cannot determine Linux version we are running on!"
        exit 1
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
    centos|redhat|fedora) OS_FAMILY="rhel";;
    debian|ubuntu) OS_FAMILY="debian";;
    *) OS_FAMILY=$OS_TYPE;;
esac

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
    if [[ -z "$*" ]]; then
        echo "Internal error, get() expects files to get"
        exit 1
    fi

    if [[ "$FAST_DOWNLOADER" == "true" ]]; then
        if which aria2c > /dev/null; then
            GET="aria2c --quiet=true --continue=true --max-concurrent-downloads=10 --max-connection-per-server=16 --max-tries=0"
        else
            echo "'--fast-downloader' specified, but couldn't find tool ('aria2c')."
            exit 1
        fi
    else
        if which curl > /dev/null; then
            GET="curl -sfL"
        else
            echo "Unable to find downloader tool ('curl')."
            exit 1
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
        echo "Unable to create system group ${SYSTEM_GROUP}"
        exit ${RC}
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
    echo "Unable to create system user ${SYSTEM_USER}"
    exit ${RC}
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

ensure_packages() {
    echo "Ensuring required tools are installed"
    if [[ $OS_FAMILY == darwin ]] ; then
        error=0
        VER=$(tar -h | grep "bsdtar " | awk '{ print $2 }' | awk -F. '{ print $1 }')
        if [[ $VER != 3 ]] ; then
            echo "Please update tar to greater than 3.0.0"
            echo
            echo "E.g: "
            echo "  brew install libarchive --force"
            echo "  brew link libarchive --force"
            echo
            error=1
        fi
        if ! which 7z &>/dev/null; then
            echo "Must have 7z"
            echo "E.g: brew install p7zip"
            echo
            error=1
        fi
        if ! which jq &>/dev/null; then
            echo "Must have jq installed"
            echo "E.g: brew install jq"
            echo
            error=1
        fi
        if ! which curl &>/dev/null; then
            echo "Must have curl installed"
            echo "E.g: brew install curl"
            echo
            error=1
        fi
        if [[ "$FAST_DOWNLOADER" == "true" ]]; then
          if ! which aria2c  &>/dev/null; then
            echo "Install 'aria2' package"
            echo
            echo "E.g: "
            echo "  brew install aria2"
          fi
        fi
        if [[ $error == 1 ]] ; then
            echo "After install missing components, restart the terminal to pick"
            echo "up the newly installed commands."
            echo
            exit 1
        fi
    else
        if ! which bsdtar &>/dev/null; then
            echo "Installing bsdtar"
            if [[ $OS_FAMILY == rhel ]] ; then
                $_sudo yum install -y bsdtar
            elif [[ $OS_FAMILY == debian ]] ; then
                $_sudo apt-get install -y bsdtar
            fi
        fi
        if ! which jq &>/dev/null; then
            echo "Installing jq"
            if [[ $OS_FAMILY == rhel ]] ; then
                install_epel
                $_sudo yum install -y jq
            elif [[ $OS_FAMILY == debian ]] ; then
                $_sudo apt-get install -y jq
            fi
        fi
        if ! which curl &>/dev/null; then
            echo "Installing curl"
            if [[ $OS_FAMILY == rhel ]] ; then
                install_epel
                $_sudo yum install -y curl
            elif [[ $OS_FAMILY == debian ]] ; then
                $_sudo apt-get install -y curl
            fi
        fi
        if ! which 7z &>/dev/null; then
            echo "Installing 7z"
            if [[ $OS_FAMILY == rhel ]] ; then
                install_epel
                $_sudo yum install -y p7zip
            elif [[ $OS_FAMILY == debian ]] ; then
                $_sudo apt-get install -y p7zip-full
            fi
        fi
        if [[ "$FAST_DOWNLOADER" == "true" ]]; then
          if ! which aria2 &>/dev/null; then
            echo "Installing aria2 for 'fast downloader'"
            if [[ $OS_FAMILY == rhel ]] ; then
                install_epel
                $_sudo yum install -y aria2
            elif [[ $OS_FAMILY == debian ]] ; then
                $_sudo apt-get install -y aria2
            fi
          fi
        fi
    fi
}

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
        echo "  export URL=\`${EP}drpcli bootenvs show $BOOTENV | grep 'IsoUrl' | cut -d '\"' -f 4\`"
        echo "  export ISO=\`${EP}drpcli bootenvs show $BOOTENV | grep 'IsoFile' | cut -d '\"' -f 4\`"
        echo "  \$CMD -o \$ISO \$URL"
    done
    echo "  # this should move the ISOs to the TFTP directory..."
    echo "  $_sudo mv *.tar *.iso $TFTP_DIR/isos/"
    echo "  $_sudo pkill -HUP dr-provision"
    echo "  echo 'NOTICE:  exploding isos may take up to 5 minutes to complete ... '"
    echo "###### END scriptlet"

    echo
}

# main
LOCAL_JQ=$(which jq || :)
arch=$(uname -m)
case $arch in
  x86_64|amd64)
                if [[ $LOCAL_JQ == "" ]] ; then
                    get https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
                    LOCAL_JQ=$(pwd)/jq-linux64
                    chmod +x $LOCAL_JQ
                fi
                arch=amd64  ;;
  aarch64)      arch=arm64  ;;
  armv7l)       arch=arm_v7 ;;
  *)            echo "FATAL: architecture ('$arch') not supported"
                exit 1;;
esac

if [[ $LOCAL_JQ == "" ]] ; then
        echo "Must have jq installed to install"
        exit 1
fi

case $(uname -s) in
    Darwin)
        binpath="bin/darwin/$arch"
        bindest="${BIN_DIR}"
        tar="command bsdtar"
        # Someday, handle adding all the launchd stuff we will need.
        shasum="command shasum -a 256";;
    Linux)
        binpath="bin/linux/$arch"
        bindest="${BIN_DIR}"
        tar="command bsdtar"
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
            exit 1
        fi
        shasum="command sha256sum";;
    *)
        # Someday, support installing on Windows.  Service creation could be tricky.
        echo "No idea how to check sha256sums"
        exit 1;;
esac

MODE=$1
if [[ "$MODE" == "upgrade" ]]
then
    MODE=install
    UPGRADE=true
    force=true
fi

case $MODE in
     install)
             if [[ "$ISOLATED" == "false" || "$SKIP_RUN_CHECK" == "false" ]]; then
                 if pgrep dr-provision; then
                     echo "'dr-provision' service is running, CAN NOT upgrade ... please stop service first"
                     exit 9
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
                     exit 1
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
                       if [[ $DRP_VERSION == tip ]] || [[ $DRP_VERSION == stable ]] ; then
                               get $DRP_CATALOG
                               mv rackn-catalog.json rackn-catalog.json.gz
                               gunzip rackn-catalog.json.gz
                               DDV=$($LOCAL_JQ -r ".sections.catalog_items[\"drp-$DRP_VERSION\"].ActualVersion" rackn-catalog.json)
                       else
                           DDV=$DRP_VERSION
                       fi

                       get $URL_BASE_DRP/$DDV.zip
                       mv $DDV.zip $ZIP

                       # XXX: Put sha back one day
                       #get $URL_BASE_DRP/$DDV.sha256
                       #mv $DDV.sha256 $SHA
                       #$shasum -c dr-provision.sha256
                     fi
                     $tar -xf dr-provision.zip
                 fi
                 $shasum -c sha256sums || exit 1
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
                 CC_VERSION=$($LOCAL_JQ -r ".sections.catalog_items[\"drp-community-content-$DRP_CONTENT_VERSION\"].ActualVersion" rackn-catalog.json)

                 CC_JSON=${CC_VERSION}.json
                 get $URL_BASE_CONTENT/$CC_JSON
                 mv $CC_JSON drp-community-content.json
                 # XXX: Add back in sha
             fi

             if [[ $ISOLATED == false ]]; then
                 INST="${BIN_DIR}/drp-install.sh"
                 $_sudo cp $TMP_INST $INST && $_sudo chmod 755 $INST
                 echo "Install script saved to '$INST'"
                 echo "(run '$INST remove' to uninstall DRP)"

                 # move aside/preserve an existing drpcli - this machine might be under
                 # control of another DRP Endpoint, and this will break the installer (text file busy)
                 if [[ -f "$CLI" ]]; then
                     echo "SAVING '${BIN_DIR}/drpcli' to backup file ($CLI_BKUP)"
                     $_sudo mv "$CLI" "$CLI_BKUP"
                 fi

                 setup_system_user

                 TFTP_DIR="${DRP_HOME_DIR}/tftpboot"
                 $_sudo cp "$binpath"/* "$bindest"

                 setcap_drp_binary

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
                         $_sudo cp "$initfile" "$initdest"
                     fi
                     # output our startup helper messages only if SYSTEMD isn't specified
                     if [[ "$SYSTEMD" == "false" || "$STARTUP" == "false" ]]; then
                        echo
                        echo "######### You can start the DigitalRebar Provision service with:"
                        echo "$starter"
                        echo "######### You can enable the DigitalRebar Provision service with:"
                        echo "$enabler"
                    else
                        echo "######### Attempt to execute startup procedures ('--startup' specified)"
                        echo "$starter"
                        echo "$enabler"
                    fi
                 fi

                 # handle the v3.0.X to v3.1.0 directory structure.
                 if [[ ! -e ${DRP_HOME_DIR}/digitalrebar && -e ${DRP_HOME_DIR} ]] ; then
                     DIR=$(mktemp -d)
                     $_sudo mv ${DRP_HOME_DIR} $DIR
                     $_sudo mkdir -p ${DRP_HOME_DIR}
                     $_sudo mv $DIR/* ${DRP_HOME_DIR}/digitalrebar
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

                 set_ownership_of_drp

                 if [[ $SYSTEMD == true ]] ; then
                     mkdir -p /etc/systemd/system/dr-provision.service.d
                     cat > /etc/systemd/system/dr-provision.service.d/user.conf <<EOF
[Service]
User=${SYSTEM_USER}
Group=${SYSTEM_GROUP}
EOF
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
                     if [[ $IPADDR ]] ; then
                       IPADDR="${IPADDR///*}"
                       cat > /etc/systemd/system/dr-provision.service.d/ipaddr.conf <<EOF
[Service]
Environment=RS_STATIC_IP=$IPADDR
Environment=RS_FORCE_STATIC=true
EOF
                     fi
                     if [[ $LOCAL_UI ]] ; then
                       cat > /etc/systemd/system/dr-provision.service.d/local-ui.conf <<EOF
[Service]
Environment=RS_LOCAL_UI=tftpboot/files/ux
Environment=RS_UI_URL=/ux
EOF
                     fi

                     eval "$enabler"
                     eval "$starter"

                     if [[ $NO_CONTENT == false ]] ; then
                         drpcli contents upload catalog:task-library-${DRP_CONTENT_VERSION}
                     fi

                     if [[ $DRP_USER ]] ; then
                         drpcli users create "{ \"Name\": \"$DRP_USER\", \"Roles\": [ \"superuser\" ] }"
                         drpcli users password $DRP_USER "$DRP_PASSWORD"
                         export RS_KEY="$DRP_USER:$DRP_PASSWORD"
                         if [[ $REMOVE_RS == true ]] ; then
                             drpcli users destroy rocketskates
                         fi
                     fi
                 else
                     if [[ "$STARTUP" == "true" ]]; then
                         echo "######### Attempting startup of 'dr-provision' ('--startup' specified)"
                         eval "$enabler"
                         eval "$starter"

                         drpcli info get > /dev/null 2>&1
                         START_CHECK=$?

                         if [[ "$NO_CONTENT" == "false" && "$START_CHECK" == "0" ]] ; then
                             drpcli contents upload catalog:task-library-${DRP_CONTENT_VERSION}
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

#SYG
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
             exit 0
         fi
         if pgrep dr-provision; then
             echo "'dr-provision' service is running, CAN NOT remove ... please stop service first"
             exit 9
         else
             echo "'dr-provision' service is not running, beginning removal process ... "
         fi
         [[ -f "$CLI_BKUP" ]] && ( echo "Restoring original 'drpcli'."; $_sudo mv "$CLI_BKUP" "$CLI"; )
         echo "Removing program and service files"
         $_sudo rm -f "$bindest/dr-provision" "$bindest/drpcli" "$initdest"
         [[ -d /etc/systemd/system/dr-provision.service.d ]] && rm -rf /etc/systemd/system/dr-provision.service.d
         [[ -f ${BIN_DIR}/drp-install.sh ]] && rm -f ${BIN_DIR}/drp-install.sh
         if [[ $REMOVE_DATA == true ]] ; then
             echo "Removing data files"
             $_sudo rm -rf "/usr/share/dr-provision" "/etc/dr-provision" "${DRP_HOME_DIR}"
         fi
         ;;
     *)
         echo "Unknown action \"$1\". Please use 'install', 'upgrade', or 'remove'";;
esac

exit 0
