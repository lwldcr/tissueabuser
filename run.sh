#!/bin/bash

CMD=tissueabuser
BINDIR=bin
LOGDIR=logs
LOG=app.log

[ -d ${LOGDIR} ] || mkdir -pv ${LOGDIR}
cmd="./${BINDIR}/${CMD} 2>&1 > ${LOGDIR}/${LOG}"

echo ${cmd}  && eval ${cmd}
