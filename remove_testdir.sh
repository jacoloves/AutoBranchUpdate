#!/bin/sh

DATE=`date '+%Y%m%d'`

DELREPO="${HOME}/released/AutoBranchUpdate/test/${DATE}"

if [ -e $DELREPO ]; then
    echo $DELREPO "remove"
    rm -rf $DELREPO
fi

if [ ! -e $DELREPO ]; then
    echo $DELREPO "remove ok!!"
fi

DELREPO="${HOME}/released/AutoBranchUpdate/test/log/${DATE}"

if [ -e $DELREPO ]; then
    echo $DELREPO "remove"
    rm -rf $DELREPO
fi

if [ ! -e $DELREPO ]; then
    echo $DELREPO "remove ok!!"
fi
