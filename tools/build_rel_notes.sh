#!/bin/bash

subdir="drp"

rm -rf rel_notes
mkdir -p rel_notes

BR=$(git rev-parse --abbrev-ref HEAD)

PREV=""
git tag | sort -V | while read rev
do
        if [[ $PREV != "" ]] ; then
                git log --name-status --no-merges $PREV..$rev > rel_notes/$rev.txt

                revname=${rev}
                revname=$(echo $rev | sed 's/\.0$/.0-release/')
                replace=${rev//?/=}
                echo ".. _rs_rel_notes_$rev:" > rel_notes/$revname.rst
                echo "" >> rel_notes/$revname.rst
                echo "$rev" >> rel_notes/$revname.rst
                echo "$replace" >> rel_notes/$revname.rst
                echo "" >> rel_notes/$revname.rst
                echo "  ::" >> rel_notes/$revname.rst
                echo "" >> rel_notes/$revname.rst
                cat rel_notes/$rev.txt | sed 's/^/    /' >> rel_notes/$revname.rst
                echo "" >> rel_notes/$revname.rst
                echo "    End of Note" >> rel_notes/$revname.rst
                echo "" >> rel_notes/$revname.rst
                rm rel_notes/$rev.txt
        fi
        PREV=$rev

        git log --name-status --no-merges $rev..HEAD > rel_notes/${BR}-pending.txt

        rev=${BR}-pending
        revname=${rev}
        revname=$(echo $rev | sed 's/\.0$/.0-release/')
        replace=${rev//?/=}
        echo ".. _rs_rel_notes_$rev:" > rel_notes/$revname.rst
        echo "" >> rel_notes/$revname.rst
        echo "$rev" >> rel_notes/$revname.rst
        echo "$replace" >> rel_notes/$revname.rst
        echo "" >> rel_notes/$revname.rst
        echo "  ::" >> rel_notes/$revname.rst
        echo "    End of Note" >> rel_notes/$revname.rst
        echo "" >> rel_notes/$revname.rst
        cat rel_notes/$rev.txt | sed 's/^/    /' >> rel_notes/$revname.rst
        echo "" >> rel_notes/$revname.rst
        echo "    End of Note" >> rel_notes/$revname.rst
        echo "" >> rel_notes/$revname.rst
        rm rel_notes/$rev.txt
done

mkdir -p rebar-catalog/docs/rel-notes/$subdir
cp rel_notes/* rebar-catalog/docs/rel-notes/$subdir

aws s3 ls rebar-catalog/docs/rel-notes/$subdir/ | awk '{ print $4 }' > files.current
ls rebar-catalog/docs/rel-notes/$subdir >> files.current
sort -r -V files.current >> new-files.current
cat new-files.current | while read f
do
    echo "https://rebar-catalog.s3-us-west-2.amazonaws.com/docs/rel-notes/$subdir/$f"
done > new-files2.current
cp new-files2.current rebar-catalog/docs/rel-notes/$subdir.filelist
rm -f new-files.current new-files2.current files.current
