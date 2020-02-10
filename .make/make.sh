#!/bin/sh
set -eu


# LOG_SEPARATOR used to friendly output blocks separation.
LOG_SEPARATOR='================================================================================'


# make_install_by_package_manager will try to install the `make` program using a well-known package manager.
make_install_by_package_manager()
{
	if [ -x "$(command -v apk)" ]
	then
		# Alpine based environments.
		apk add make
	elif [ -x "$(command -v apt-get)" ]
	then
		# Ubuntu and debian based environments.
		apt-get -qq update
		apt-get -qq install -y make
	elif [ -x "$(command -v yum)" ]
	then
		# Redhat based environments.
		yum -q install -y make
	else
		echo '`make` cannot be installed in current environment.'
		exit 1
	fi
}


# Phase 1. We need `GNUmake`.
# The only target cannot be resolved by make - GNUmake itself.
# Main idea - we dont know where `make` session was invoked.
# Maybe there is not GNUmake installed. Detect system, package manager and install it.
if ! [ -x "$(command -v make)" ]
then
	# We will detect package managers only.
	# If no package manager installed we would compile gnumake.
	# But it is too far...
	echo ${LOG_SEPARATOR}
	echo 'Autoinstall data:'
	make_install_by_package_manager
fi


# Phase 2. Print information about installed `GNUmake`.
echo ${LOG_SEPARATOR}
make --version


# Get required variables for running make session from caller.
readonly DIFF
readonly DIFF_MOD


# Phase 3. Doing main job.
# Make installed, version is compatible, just forward request to his binary.
echo ${LOG_SEPARATOR}
make DIFF_MOD='STUPID' DIFF='' $@
