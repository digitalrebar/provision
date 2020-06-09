#!/bin/bash

rm -rf rel_notes
mkdir -p rel_notes

BR=$(git rev-parse --abbrev-ref HEAD)

PREV=""
git tag | sort -V | while read rev
do
        if [[ $PREV != "" ]] ; then
                git log --name-status --no-merges $PREV..$rev > rel_notes/$rev.txt

                replace=${rev//?/=}
                echo ".. _rs_rel_notes_$rev:" > rel_notes/$rev.rst
                echo "$rev" >> rel_notes/$rev.rst
                echo "$replace" >> rel_notes/$rev.rst
                echo "::" >> rel_notes/$rev.rst
                cat rel_notes/$rev.txt | sed 's/^/  /' >> rel_notes/$rev.rst
                echo "" >> rel_notes/$rev.rst
                rm rel_notes/$rev.txt
        fi
        PREV=$rev

        git log --name-status --no-merges $rev..HEAD > rel_notes/${BR}-pending.txt

        rev=${BR}-pending
        replace=${rev//?/=}
        echo ".. _rs_rel_notes_$rev:" > rel_notes/$rev.rst
        echo "$rev" >> rel_notes/$rev.rst
        echo "$replace" >> rel_notes/$rev.rst
        echo "::" >> rel_notes/$rev.rst
        cat rel_notes/$rev.txt | sed 's/^/  /' >> rel_notes/$rev.rst
        echo "" >> rel_notes/$rev.rst
        rm rel_notes/$rev.txt
done

mkdir -p rebar-catalog/docs/rel-notes/drp
cp rel_notes/* rebar-catalog/docs/rel-notes/drp

aws s3 ls rebar-catalog/docs/rel-notes/drp/ > files.current
ls rebar-catalog/docs/rel-notes/drp >> files.current
sort -r -V files.current >> new-files.current
cp new-files.current rebar-catalog/docs/rel-notes/drp.filelist
rm -f new-files.current files.current

