#!/usr/bin/env bash

###
#  Manage the RackN install.sh installer script to S3 bucket
#  upload.  See '-h' for HELP option output for usage.
###

#export AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID:-"key"}
#export AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY:-"secret"}
export AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION:-"us-west-2"}
FILE=${INSTALLER_FILE:-"tools/install.sh"}
DESTS=${INSTALLER_DESTS:-"test"}
BKT=${INSTALLER_BUCKET:-"get.rebar.digital"}

usage() {
#  cat  << EOT
echo -e "
${Gre}USAGE${RCol}:   ${IBlu}$0${RCol} -h
   ${Gre}OR${RCol}:   ${IBlu}$0${RCol} -d dest [ -k key -p password -f file -b bucket -r region ]

${Gre}WHERE${RCol}:   -h              print this help statement

         -d dest1,dest2  (required) Destination object name(s), multiple with
                         comma separated list, no spaces
         -k key          AWS Access Key ID to use (default: '$AWS_ACCESS_KEY_ID')
         -p password     AWS Secret for Key ID to use (default: '$AWS_SECRET_ACCESS_KEY')

         -f file         source file to upload (default: '$FILE')
         -b bucket       S3 bucket to upload to (default: '$BKT')
         -r region       AWS S3 region to upload to (default: '$AWS_DEFAULT_REGION')

${Gre}MORE${RCol}:    The following shell environment varialbes can be set:

         AWS_ACCESS_KEY_ID  AWS_SECRET_ACCESS_KEY  AWS_DEFAULT_REGION
         INSTALLER_FILE     INSTALLER_DESTS         INSTALLER_BUCKET

${Red}WARNING${RCol}: AWS Key and Secret MUST BE SET, either via environment varialbes
         command line options, or in the $HOME/.aws/credentials file.

"
}

command_line() {
  while getopts ":hf:d:b:k:p:r:" opt
  do
    case "${opt}" in
      h)  usage; exit 0                 ;;
      f)  FILE=$OPTARG                  ;;
      d)  DESTS=$OPTARG                  ;;
      b)  BKT=$OPTARG                   ;;
      k)  AWS_ACCESS_KEY_ID==$OPTARG    ;;
      p)  AWS_SECRET_ACCESS_KEY=$OPTARG ;;
      r)  AWS_DEFAULT_REGION=$OPTARG    ;;
      \?) echo
          echo "Invalid usage option flag: $OPTARG"
          usage
          exit 1
          ;;
    esac
  done
}

main() {
  PRE="${Cya}>>>${RCol}"

  OLD_IFS="$IFS"
  IFS=","
  echo ""

  for DEST in $DESTS
  do
    OBJ="s3://$BKT/$DEST"
    OBJ_OUT="${UYel}${OBJ}${RCol}"
    DEST_OUT="${UYel}${DEST}${RCol}"

    echo -e "${PRE} ${Blu}MODIFY${RCol}  ${UWhi}S3 object ${DEST_OUT}"
    MSG="${PRE} ${Blu}COPY${RCol}    local file '${FILE}' to object '${OBJ_OUT}'"
    print_msg "$MSG"
    aws --quiet s3 cp $FILE $OBJ
    success

    MSG="${PRE} ${Blu}SET${RCol}     public read policy on object '${OBJ_OUT}'"
    print_msg "$MSG"
    aws s3api put-object-acl --bucket $BKT --key $DEST --acl public-read
    success

    MSG="${PRE} ${Blu}TEST${RCol}    download object '${OBJ_OUT}'"
    print_msg "$MSG"
    curl -fs $BKT/$DEST > /dev/null
    success

    echo ""
  done

  IFS=$OLD_IFS
}

success() {
  echo -e " ${IGre}Success${RCol}"
}

print_msg() {
  local _nws
  local _msg="$*"
  # _len includes the msg plus control characters - 4 sets of on/off sequences
  local _len="120"
  local _cnt=$(printf "$_msg" | wc -c)
  (( _nws = _len - _cnt ))
  WS=$(printf "%${_nws}s" " ")
  echo -en "${_msg}${WS}${Rcol}"
}

set_colors() {
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

}

set_colors
command_line $@
main
