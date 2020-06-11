#!/usr/bin/env bash

###
#  Simple helper script to generate a new Knowledge Base document ID, and populate
#  the doc/kb/ directory with a template that contains the standard KB sections.
###

usage() {
cat <<EOF
USAGE:  $0 [ -n ] | [ -d DOCDIR ] | [ -i INDEX ] [ -t TITLE ]
   OR:  $0 -u

WHERE:  -n           DO NOT Attempt to start editor on generated tempalte file
        -d DOCDIR    Specify the document directory for the KB articles
                     (defaults to "doc")
        -k KBDIR     Specify the Knowledge Base subdir in DOCDIR
                     (defaults to "kb")
        -i INDEX     Specify the article Index number
                     (defaults to generate next numerical value found in
                     the DOCDIR/KBDIR/ directory)
        -t TITLE     Provide an initial title to set for the KB article.
                     (none by default)
        -l LABEL     Provide an initial label for the title (requires TITLE)
                     (none by default)
        -u           print this usage statement

NOTES:  * DOCDIR, KBDIR, INDEX, and TITLE can be passed in as command line variables, like:

            # produces path foo/dir/knowledge/
            DOCDIR=foo/dir KBDIR=knowledge INDEX=00007 TITLE="Foo Title" $0

        * Environment variables override command line flags
        * Indexes must be 5 digit numbers

EOF
} # end usage()

main() {
  [[ -z "$INDEX" ]] && get_next_index
  validate_index "$INDEX" || xiterr "Invalid index number format specified ('$INDEX').  Must be 5 digits."
  KB_FILE="$DOCDIR/$KBDIR/kb-$INDEX.rst"
  echo ">>> Knowledge base file set to: '$(basename $KB_FILE)'"
  generate_kb
  (( $EDIT )) && editor

} # end main()

function xiterr() { [[ $1 =~ ^[0-9]+$ ]] && { XIT=$1; shift; } || XIT=1; printf "!!! FATAL: $*\n"; exit $XIT; }

command_line() {
  # "$@" must be passed to us so we can access the flags
  EDIT=1
  while getopts ":nd:k:i:t:l:u" CmdLineOpts
  do
    case $CmdLineOpts in
      n)      EDIT=0             ;;
      d)      DOCDIR=${OPTARG}   ;;
      k)      KBDIR=${OPTARG}    ;;
      i)      INDEX=${OPTARG}    ;;
      t)      TITLE="${OPTARG}"  ;;
      l)      LABEL="${OPTARG}"  ;;
      u)      usage
              exit 0
              ;;
      \?)
              echo "Incorrect usage.  Invalid flag '${OPTARG}'."
              usage
              exit 1
              ;;
        esac
    done

  DSTAMP=$(date)
  DOCDIR=${DOCDIR:-"doc"}
  INDEX=${INDEX:-""}
  KBDIR=${KBDIR:-"kb"}
  TITLE=${TITLE:-""}
  LABEL=${LABEL:-""}

  [[ -n "$TITLE" ]] && echo ">>> Setting document title to : '$TITLE'"
  [[ -n "$LABEL" && -z "$TITLE" ]] && xiterr 1 "LABEL specified, but no required TITLE"
  [[ -n "$TITLE" && -z "$LABEL" ]] && generate_label "$TITLE"
  LABEL=$(echo $LABEL | sed 's/^_//')
  [[ -n "$LABEL" ]] && echo ">>> Setting document label to : '$LABEL'"

} # end command_line()

generate_label(){
  local _label="$*"
  echo "$_label" | sed -e 's/[^a-zA-Z0-9 ]//g' -e "s/ /_/g"
  LABEL="$_label"
}

editor(){
  local _editor
  which vi > /dev/null && _editor=$(which vi)
  which vim > /dev/null && _editor=$(which vim)
  [[ -z "$_editor" ]] && _editor=$VISUAL
  [[ -z "$_editor" ]] && xiterr 1 "can't find an editor to use (vi, vim, \$VISUAL)"
  echo ">>> Attempt to start editor   : '$_editor $KB_FILE'"
  [[ -w "$KB_FILE" ]] && $_editor "$KB_FILE" || xiterr 1 "can't write to KB file ('$KB_FILE')"
}

validate_index(){
  local _idx=$1
  local _len=${#_idx}

  [[ $_len != 5 ]] && xiterr 1 "Invalid index number format specified ('$INDEX').  Must be 5 digits." || true
} # end validate_index()

get_next_index(){
  local _last _num

  if [[ -z "$INDEX" ]]
  then
    # get next index
    _last=$(find $DOCDIR/$KBDIR -name "kb-[0-9][0-9][0-9][0-9][0-9]*\.rst" -exec basename {} \; | sort -n | tail -1)
    _num=$(echo $_last | sed -e 's/^kb-\([0-9][0-9][0-9][0-9][0-9]\).*rst$/\1/')
    _num=$(echo "$_num" | sed 's/^0*//')
    (( _num++ ))
    INDEX=$(printf "%05s" "$_num")
    echo ">>> New index number generated: '$INDEX'"
  else
    echo ">>> Index number already provided, not generating new dynamic index."
  fi
} # end get_netx_index()


generate_kb(){

  [[ -e "$KB_FILE" ]] && xiterr 1 "Knowledge base file already exists, not overwriting ('$KB_FILE')" || true
  echo ">>> Write template KB file to : '$KB_FILE'"
  YEAR=$(date +%Y)
  [[ -n "$TITLE" ]] && TITLE="${TITLE}" || TITLE="TITLE TITLE TITLE"
  TLEN=$(echo $(( $(echo "$TITLE" | wc -c) )))
  (( TLEN-- ))
  TILDE=$(eval printf '~%.0s' {1..${TLEN}})
  # NO space between TITLE and <> brackets for the LABEL
  [[ -n "$TITLE" && -n "$LABEL" ]] && XREF=":ref:\`$TITLE<$LABEL>\`"

  cat <<KB > $KB_FILE
.. Copyright (c) $YEAR RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license

.. REFERENCE kb-00000 for an example and information on how to use this template.
.. If you make EDITS - ensure you update footer release date information.


.. _rs_kb_${INDEX}:

kb-${INDEX}
~~~~~~~~

.. _${LABEL}:

$TITLE
$TILDE

Description
-----------


Solution
--------


Additional Information
----------------------

Additional resources and information related to this Knowledge Base article.


See Also
========


Versions
========


Keywords
========


Revision Information
====================
  ::

    KB Article     :  kb-$INDEX
    initial release:  $DSTAMP
    updated release:  $DSTAMP

KB
}

command_line "$@"
main
