#!/usr/bin/env bash

# unlike SIM projects, each compiler project has a name of its own
PROGRAM=WRECK

# this shifts off ${1:-.} to COMPLOC
source "${COMPGRADING}/comp-lib.sh"
source "${GRADERDIR}/grader-exec.sh"

set -e

# make executable
chmod +x "${graderloc}/normscanu"
# look for required execs
for x in aenorm aenfatott cmptt ; do 
	if ! test -x "${COMPGRADING}/$x" ; then 
		grader_echo "ERROR: '${COMPGRADING}/$x' can be executed."
		exit 1
	fi
done

# create a temporary dir to hold tests, this permits multiple 
# simultaneous grading sessions without stumbling over each other...
GRADERBIO="_wrecktest"
mkdir -p "${GRADERBIO}"
mirror "${graderloc}/basicio" "${GRADERBIO}"
# some students need llre.cfg for run
test -s "llre.cfg" || cp "${graderloc}/llre.cfg" "."


NOTOKEN=_test_empty.scan-u_output_file

basiciotallydir="`grader_mktemp -d basiciotally`"

# RUBRIC 1 ###
# inaccessible scan.lut
tally_missing_datafile "${basiciotallydir}/missinglut" "${comploc}/${PROGRAM}" MISSINGDATAFILE ${NOTOKEN} 

# inaccessible output location
rm -rf "${graderloc}/${GRADERBIO}/_dne"
tally_nooutput_exitnonzero "${basiciotallydir}/nowriteu" "${comploc}/${PROGRAM}" "${GRADERBIO}/_abc.lut" "${GRADERBIO}/_dne/_scan.u" 

# unwritable _tokenid.nfa file
nopermnfa=`sed -n 2p "${GRADERBIO}/_bad.lut" | awk '{print $2}'`.nfa
# clear out
rm -rf "${nopermnfa}"
# unwritable file
touch "${nopermnfa}"
chmod -w "${nopermnfa}"
tally_nooutput_exitnonzero "${basiciotallydir}/nopermnfaf" "${comploc}/${PROGRAM}" "${GRADERBIO}/_bad.lut" "${GRADERBIO}/_scan.u"
# should also fail w/ dirs
rm -f "${nopermnfa}"
mkdir -p "${nopermnfa}"
tally_nooutput_exitnonzero "${basiciotallydir}/nopermnfad" "${comploc}/${PROGRAM}" "${GRADERBIO}/_bad.lut" "${GRADERBIO}/_scan.u"
rm -rf "${nopermnfa}"


# syntax and semantic errors
syntaxtallydir="`grader_mktemp -d syntax`"
semantictallydir="`grader_mktemp -d semantic`"
for tsttype in syntax semantic ; do 
	tdir=${tsttype}tallydir
	case ${tsttype} in syntax ) expes=2 ;; semantic ) expes=3 ;; * ) expes=NaN ;; esac
	for tstlut in "${graderloc}/${tsttype}"*/scan.lut ; do 
		tid=${tstlut##*/*${tsttype}}
		rm -f "${GRADERBIO}/"_scan.* 
		cp "${tstlut}" "${GRADERBIO}/_scan.lut"
		tally_ignore_exitvalue "${!tdir}/${tsttype}${tid%/scan.lut}" -eq ${expes} "${comploc}/${PROGRAM}" "${GRADERBIO}/_scan.lut" "${GRADERBIO}/_scan.u"
	done
done



rm -rf "${GRADERBIO}"


compare_scan_u()
{
	local tst="${1}" ; shift
	"${graderloc}/normscanu"  "${1}" >"_${tst}_student.u"
	"${graderloc}/normscanu"  "${2}" >"_${tst}_expected.u"
	if ! diff "_${tst}_student.u" "_${tst}_expected.u" >/dev/null ; then 
		grader_msg <<EOT
scan.u output file for test set ${tst} in '${graderloc}' FAILED".  The output from your ${PROGRAM} is in ${1};
it has been normalized to _${tst}_students.u and compared to _${tst}_expected.u.

You can see which lines differ using "visual diff" tool, or at the console with

  $ diff -u _${tst}_students.u _${tst}_expected.u

(- lines would be removed from _${tst}_students.u and + lines would be added to 
match _${tst}_expected.u).

IF YOU WANT TO INSPECT THIS FAILURE, CTRL-C now!
EOT
		return 1
	else
		rm -f "_${tst}_student.u" "_${tst}_expected.u"
	fi
	
}

compare_nfa()
{
	local tst="${1}" ; shift
	local tokid="${1}" ; shift
	local stalph=`head -n 1 "${1}" | ( read n l x ; echo "$x"; ) |tee "_comparenfa-sa.lis" | "${COMPGRADING}/aenorm"`	
	local exalph=`head -n 1 "${2}" | ( read n l x ; echo "$x"; ) |tee "_comparenfa-ea.lis" | "${COMPGRADING}/aenorm"`	
	local ret=0
	if ! test "${stalph}" = "${exalph}" ; then
		ret=1
		grader_msg <<EoT
The alphabet in the first line of '${1}' after normalization is
  ${stalph:-empty}
which does not compare favorably to the expected alphabet
  ${exalph}

IF YOU WANT TO INSPECT THIS FAILURE, CTRL-C now!
EoT
	fi
	${COMPGRADING}/aenfatott "${1}" >"_${tst}_${tokid}_student_dfa.tt"
	${COMPGRADING}/aenfatott "${2}" >"_${tst}_${tokid}_expected_dfa.tt"
	if ! ${COMPGRADING}/cmptt "_${tst}_${tokid}_expected_dfa.tt" "_${tst}_${tokid}_student_dfa.tt" 2>"_${tst}_${tokid}_dfacmp.err"; then
		ret=1
		grader_msg <<EoT
The NFA stored for tokenid ${tokid} (${1}) has been converted to an optimized
DFA at _${tst}_${tokid}_student_dfa.tt.  This should be the same as the DFA
in _${tst}_${tokenid}_expected_dfa.tt; but it does not compare favorably.

There *should* be PDFs of your NFA's initial, simplified and minimized
representations stored at  ${1%.nfa}-input.pdf,  ${1%.nfa}-simplified.pdf  and
${1%.nfa}-output.pdf.  The same style PDFs for ${tst} ${tokid} that generate
the grader script results are provided in an aptly named ${tst/test/result}/
directory alongside the grader.sh you are using.  Comparing these PDFs to yours
may provide a hint as to where your NFA construction routine isn't quite up to
snuff. KEEP IN MIND: the output pdfs should match precisely, the input pdfs not
so (since they are NFAs).

The  ${1%.nfa}-simplified  NFA  is the original input NFA but with irrelevant 
states removed.  I've learned it can be easier to see construction logic flaws
in the simplified form --- but the input NFA is what your WRECK is actually 
generating!

The diff(1) differences between the transition table files (.tt) are stored in
_${tst}_${tokid}_dfacmp.err, this may or may not be helpful since it can be a
long way between NFA and optimized DFA.

IF YOU WANT TO INSPECT THIS FAILURE, CTRL-C now!
EoT
	else 
		rm "_${tst}_${tokid}_student_dfa.tt"
		rm "_${tst}_${tokid}_expected_dfa.tt"
		rm "${1%.nfa}"-*.pdf 
		rm -f "${1%.nfa}"-*.gv
	fi
	return $ret
}

test_re_set()
{
	# $1=graderloc test set name
	# $2=tally prefix for scanu gen
	# $3=tally prefix for nfa gen
	local tallyu="${2}"
	local tallynfa="${3}"
	local testdir=_wrecktest
	rm -rf "${testdir}/${1}"
	mkdir -p "${testdir}/${1}"
	mirror "${graderloc}/${1}" "${testdir}/${1}/"
	# some students need llre.cfg for run
	test -s "llre.cfg" || cp "${graderloc}/llre.cfg" "."

	
	local nfaerr=0
	local badu="${tallyu}${1}.bad"
	if ! test_run "${comploc}/${PROGRAM}" "${testdir}/${1}/scan.lut" "_output.u" && test ${GraderRunES} -ne 0 ; then
		grader_echo "ERROR: exit status ${GraderRunES} non-zero."
		touch "${badu}"
	elif ! test -s ./_output.u; then
		grader_echo "ERROR: No ./_output.u generated."
		touch "${badu}"
	elif ! compare_scan_u ${1} ./_output.u "${graderloc}/${1/test/result}/scan.u" ; then
		touch "${badu}"
	else 
		touch "${badu%.bad}.good"
	fi

	for nfa in "${graderloc}/${1/test/result}/"*.nfa ; do 
		prognfa="${nfa##*/}"
		badn="${tallynfa}${1}${prognfa%.nfa}.bad"	
		if ! test -s "${prognfa}" ; then
			grader_echo "ERROR: missing ${prognfa}."
			touch "${badn}"
			nfaerr=1
		elif ! compare_nfa ${1} ${prognfa%.nfa} "${prognfa}" "${nfa}" ; then
			touch "${badn}"
			nfaerr=1
		else 
			touch "${badn%.bad}.good"
		fi
	done


	if test -f "${badu%.bad}.good" -a ${nfaerr} -eq 0 ; then
		grader_echo "RE compilation and scan.u generation for ${1} is GOOD :)"
		rm -rf "${testdir}/"
		rm -f  "_${1}"* _output.u
		for nfa in "${graderloc}/${1/test/result}/"*.nfa ; do 
			rm -f "${nfa##*/}"
		done
	else 
		grader_msg <<EoT
${1} test set failed --- you should have seen messages with details before this.
EoT
		grader_keystroke
	fi
	
	grader_echo ""
}



tallyudir="`grader_mktemp -d scanu`"
tallynfadir="`grader_mktemp -d nfa`"
for tst in `( cd "${graderloc}" && ls -1d test[0-9] )` ; do 
	test_re_set "${tst}" "${tallyudir}/" "${tallynfadir}/" 
done


width=32
show_tallies $width ${basiciotallydir}  "Basic I/O tests"
show_tallies $width ${semantictallydir} "Semantic tests"
show_tallies $width ${syntaxtallydir}   "Syntax tests"
show_tallies $width ${tallyudir}     	"Generates correct scan.u"
show_tallies $width ${tallynfadir} 		"tokenid.nfa files correct"

# always show cwd at the end, so grader is sure the correct results
# are recorded for the correct submission (the upload id is in the path)
grader_echo ""
pwd

