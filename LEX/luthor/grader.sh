#!/usr/bin/env bash

# unlike SIM projects, each compiler project has a name of its own
PROGRAM=LUTHOR

# this shifts off ${1:-.} to COMPLOC
source "${COMPGRADING}/comp-lib.sh"
source "${GRADERDIR}/grader-exec.sh"

set -e

# FIXME:  routines should rm -f _nooutput* before execution,
#         check that _nooutput* are all empty (we permit that they may be created)

# make executable
chmod +x "${graderloc}/normoutput"

# create a temporary dir to hold basicio tests in GRADERLOC, this permits multiple 
# simultaneous grading sessions without stumbling over each other...
GRADERBIO="_luthortest"
mkdir -p "${GRADERBIO}"
mirror "${graderloc}/basicio" "${GRADERBIO}"
# -i.bak compatible between bsd and linux --- but 2021 spring students complained of 
# permissions errors with this operation (from bash -x logging) --- so I'm just 
# squashing the errors into /dev/null 
sed 2>/dev/null -i.bak -e "s#_basicio#${GRADERBIO}#" "${GRADERBIO}/_abc.u" "${GRADERBIO}/_badtt.u"

NOTOKEN=_test_empty.token_output_file
# RUBRIC 1 ###
# inaccessible scan.u, program.src
test_missing_datafile "${comploc}/${PROGRAM}" MISSINGDATAFILE "${GRADERBIO}/_abc.src" ${NOTOKEN} 
test_missing_datafile "${comploc}/${PROGRAM}" "${GRADERBIO}/_abc.u" MISSINGDATAFILE ${NOTOKEN} 

# inaccessible output location
rm -rf "${graderloc}/${GRADERBIO}/_dne"
test_nooutput_exitnonzero "${comploc}/${PROGRAM}" "${GRADERBIO}/_abc.u" "${GRADERBIO}/_abc.src" "${GRADERBIO}/_dne/_output.tok" 

# inaccessible transition table entry in scanner definition entry
test_nooutput_exitnonzero "${comploc}/${PROGRAM}" "${GRADERBIO}/_badtt.u" "${GRADERBIO}/_abc.src" ${NOTOKEN} 

# RUBRIC 2 ###
# empty scanner definition file
test_nooutput_exitnonzero "${comploc}/${PROGRAM}" EMPTYDATAFILE "${GRADERBIO}/_abc.src" ${NOTOKEN} 

# RUBRIC 3 ###
# empty source files should exit 0 and truncate a pre-existing output file
test_nooutput_exitzero "${comploc}/${PROGRAM}" "${GRADERBIO}/_abc.u" EMPTYDATAFILE ${NOTOKEN} 

# RUBRIC 4 ###
# invalid character in source file
test_ignore_exitnonzero "${comploc}/${PROGRAM}" "${GRADERBIO}/_abc.u" "${GRADERBIO}/_xyz.src" "_output.tok" 

rm -rf "${GRADERBIO}"

test_token_source()
{
	# $1=graderloc tokenset dir
	# $2=tally prefix for .tokgood .tokbad
	# $3=graderloc tokenset dir source  (.src)
	local tally="${2}"
	local src="${3}"
	# tokens are either named after the program, or they are called tokens.dat
	local tok="${3/.src/.tok}"
	test -s "${graderloc}/${1}/${tok}" || tok=tokens.dat
	
	local testdir=_luthortest
	mkdir -p "${testdir}/${1}"
	mirror "${graderloc}/${1}" "${testdir}/${1}/"
	rm -f "${testdir}/${1}/${tok}"   # lest it be looked up in a side-channel
	# fix up path to .tt files
	sed 2>/dev/null -i.bak -e '/\.tt/s#^#'"${testdir}/#" "${testdir}/${1}/scan.u"
	
	rm -f ./_output.tok
	if test_run "${comploc}/${PROGRAM}" "${testdir}/${1}/scan.u" "${testdir}/${1}/${3}" ./_output.tok ; then
		touch "${tally}.bad"
	fi
	if ! test -s ./_output.tok ; then
		grader_echo "ERROR: No ./_output.tok generated."
		touch "${tally}.bad"
	fi
	# no point in continuing
	test -f "${tally}.bad" && return 0

	"${graderloc}/normoutput" < ./_output.tok  >./_students.tok
	"${graderloc}/normoutput" < "${graderloc}/${1}/${tok}" >./_expected.tok
	if ! diff _students.tok _expected.tok >/dev/null ; then 
		touch "${tally}.bad"
		grader_msg <<EOT
Tokenization for scanner definition ${1}/scan.u and source ${1}/${src}
in '${graderloc}' FAILED.  The output from your ${PROGRAM} is in _output.tok.
_output.tok was normalized to _students.tok and compared to _expected.tok.

You can see which lines differ using "visual diff" tool, or at the console with

  $ diff -u _students.tok _expected.tok

(- lines would be removed from _students.tok and + lines would be added to 
match _expected.tok.)

IF YOU WANT TO INSPECT THIS FAILURE, CTRL-C now!
EOT
		grader_keystroke
	else 
		touch "${tally}.good"
		grader_echo "Tokenization for scanner definition ${1}/scan.u and source ${1}/${src} in '${graderloc}' is GOOD :)"
		rm -rf "${testdir}"
	fi
}

test_token_set()
{
	# $1=graderloc tokenset dir
	# $2=tally prefix for the tokenset
	for src in `( cd "${graderloc}/${1}" && ls -1d *.src )` ; do 
		test_token_source "${1}" "${2}/${1/\//-}.${src/.src/}" "${src}"
	done
}


# make tally dirs for rubric line items, test the appropriate tokensets
for cat in tied disjoint overnl complex ; do 
	eval "${cat}tallydir=`grader_mktemp -d ${cat}tally`"
	tdv=${cat}tallydir
	for ts in `( cd "${graderloc}" && ls -1d ${cat}/? )` ; do 
		test_token_set ${ts} ${!tdv}
	done
done


show_tallies()
{
	local tallydir="${1}"
	shift
	local descr="${@}" 
	# tallies
	local total=`ls -1 "${tallydir}/"*.* 2>/dev/null|wc -l | tr -d '[:space:]'`
	local good=`ls -1 "${tallydir}/"*.good 2>/dev/null|wc -l | tr -d '[:space:]'`

#   alamode bc incompatible, need to use simple awk
#	pcttt=0
#	test $totaltt -gt 0 && pcttt=`bc -e "scale=5 ; ( ${goodtt} / ${totaltt} ) * 100;" -e quit`
#	pctm=0
#	test $totalm -gt 0 && pctm=`bc -e "scale=5 ; ( ${goodm} / ${totalm} ) * 100;" -e quit`
#
	grader_msg << EOT 
${descr}: ${good} out of ${total} correct.
EOT
}

show_tallies ${tiedtallydir} "Tied matches over one line"
show_tallies ${disjointtallydir} "Disjoint matches over one line"
show_tallies ${overnltallydir} "Matches over multiple lines"
show_tallies ${complextallydir} "Complex matches"

# always show cwd at the end, so grader is sure the correct results
# are recorded for the correct submission (the upload id is in the path)
pwd

