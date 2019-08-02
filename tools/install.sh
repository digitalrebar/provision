#!/usr/bin/env bash

set -e

DEFAULT_DRP_VERSION=${DEFAULT_DRP_VERSION:-"stable"}

usage() {
cat <<EOFUSAGE
Usage: $0 [--version=<Version to install>] [--no-content] [--commit=<githash>]
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
    --commit=<string>       # github commit file to wait for; unset assumes the files
                            # are in place
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
    local-ui            = false

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
_sudo="sudo"
CLI="/usr/local/bin/drpcli"
CLI_BKUP="/usr/local/bin/drpcli.drp-installer.backup"

# download URL locations; overridable via ENV variables
URL_BASE=${URL_BASE:-"https://github.com/digitalrebar/"}
URL_BASE_DRP=${URL_BASE_DRP:-"$URL_BASE/provision/releases/download"}
URL_BASE_CONTENT=${URL_BASE_CONTENT:-"$URL_BASE/provision-content/releases/download"}

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
            COMMIT=${arg_data}
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
[[ "$ISOLATED" == "true" ]] && KEEP_INSTALLER=true

[[ $DBG == true ]] && set -x

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
arch=$(uname -m)
case $arch in
  x86_64|amd64) arch=amd64  ;;
  aarch64)      arch=arm64  ;;
  armv7l)       arch=arm_v7 ;;
  *)            echo "FATAL: architecture ('$arch') not supported"
                exit 1;;
esac

case $(uname -s) in
    Darwin)
        binpath="bin/darwin/$arch"
        bindest="/usr/local/bin"
        tar="command bsdtar"
        # Someday, handle adding all the launchd stuff we will need.
        shasum="command shasum -a 256";;
    Linux)
        binpath="bin/linux/$arch"
        bindest="/usr/local/bin"
        tar="command bsdtar"
        if [[ -d /etc/systemd/system ]]; then
            # SystemD
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

if [[ $COMMIT != "" ]] ; then
    set +e
    DRP_CMT=dr-provision-hash.$COMMIT
    while ! get $URL_BASE_DRP/$DRP_VERSION/$DRP_CMT ; do
            echo "Waiting for dr-provision-hash.$COMMIT"
            sleep 60
    done
    set -e
fi

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
                       get $URL_BASE_DRP/$DRP_VERSION/$ZIP $URL_BASE_DRP/$DRP_VERSION/$SHA
                       $shasum -c dr-provision.sha256
                     fi
                     $tar -xf dr-provision.zip
                 fi
                 $shasum -c sha256sums || exit 1
             fi

             if [[ $NO_CONTENT == false ]]; then
                 DRP_CONTENT_VERSION=stable
                 if [[ $DRP_VERSION == tip ]]; then
                     DRP_CONTENT_VERSION=tip
                 fi
                 echo "Installing Version $DRP_CONTENT_VERSION of Digital Rebar Provision Community Content"
                 if [[ -n "$ZIP_FILE" ]]; then
                   echo "WARNING: '--zip-file' specified, still trying to download community content..."
                   echo "         (specify '--no-content' to skip download of community content"
                 fi
                 CC_YML=drp-community-content.yaml
                 CC_SHA=drp-community-content.sha256
                 get $URL_BASE_CONTENT/$DRP_CONTENT_VERSION/$CC_YML $URL_BASE_CONTENT/$DRP_CONTENT_VERSION/$CC_SHA
                 $shasum -c $CC_SHA
             fi

             if [[ $ISOLATED == false ]]; then
                 INST="/usr/local/bin/drp-install.sh"
                 $_sudo cp $TMP_INST $INST && $_sudo chmod 755 $INST
                 echo "Install script saved to '$INST'"
                 echo "(run '$INST remove' to uninstall DRP)"

                 # move aside/preserve an existing drpcli - this machine might be under
                 # control of another DRP Endpoint, and this will break the installer (text file busy)
                 if [[ -f "$CLI" ]]; then
                     echo "SAVING '/usr/local/bin/drpcli' to backup file ($CLI_BKUP)"
                     $_sudo mv "$CLI" "$CLI_BKUP"
                 fi

                 TFTP_DIR="/var/lib/dr-provision/tftpboot"
                 $_sudo cp "$binpath"/* "$bindest"
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
                 if [[ ! -e /var/lib/dr-provision/digitalrebar && -e /var/lib/dr-provision ]] ; then
                     DIR=$(mktemp -d)
                     $_sudo mv /var/lib/dr-provision $DIR
                     $_sudo mkdir -p /var/lib/dr-provision
                     $_sudo mv $DIR/* /var/lib/dr-provision/digitalrebar
                 fi

                 if [[ ! -e /var/lib/dr-provision/digitalrebar/tftpboot && -e /var/lib/tftpboot ]] ; then
                     echo "MOVING /var/lib/tftpboot to /var/lib/dr-provision/tftpboot location ... "
                     $_sudo mv /var/lib/tftpboot /var/lib/dr-provision/
                 fi

                 $_sudo mkdir -p /usr/share/dr-provision
                 if [[ $NO_CONTENT == false ]] ; then
                     DEFAULT_CONTENT_FILE="/usr/share/dr-provision/default.yaml"
                     $_sudo mv drp-community-content.yaml $DEFAULT_CONTENT_FILE
                 fi

                 if [[ $SYSTEMD == true ]] ; then
                     mkdir -p /etc/systemd/system/dr-provision.service.d
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
                 rm -f drpcli dr-provision drbundler drpjoin
                 ln -s $binpath/drpcli drpcli
                 ln -s $binpath/dr-provision dr-provision
                 if [[ -e $binpath/drbundler ]] ; then
                     ln -s $binpath/drbundler drbundler
                 fi
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
                 STARTER="$_sudo ./dr-provision --base-root=`pwd`/drp-data --local-content=\"\" --default-content=\"\" > drp.log 2>&1 &"
                 [[ "$STARTUP" == "false" ]] && echo "$STARTER"
                 mkdir -p "`pwd`/drp-data/saas-content"
                 if [[ $NO_CONTENT == false ]] ; then
                     DEFAULT_CONTENT_FILE="`pwd`/drp-data/saas-content/default.yaml"
                     mv drp-community-content.yaml $DEFAULT_CONTENT_FILE
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
                 echo "  ${EP}drpcli contents upload https://api.rackn.io/catalog/content/task-library"
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
         [[ -f /usr/local/bin/drp-install.sh ]] && rm -f /usr/local/bin/drp-install.sh
         if [[ $REMOVE_DATA == true ]] ; then
             echo "Removing data files"
             $_sudo rm -rf "/usr/share/dr-provision" "/etc/dr-provision" "/var/lib/dr-provision"
         fi
         ;;
     *)
         echo "Unknown action \"$1\". Please use 'install', 'upgrade', or 'remove'";;
esac

exit 0
