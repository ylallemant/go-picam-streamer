#!/bin/sh

set -e

# if [ $# -eq 0 ]; then
#     echo "ERROR: Need to specify the install repository"
#     exit 1
# fi

# eg. release-lab/whatchanged
owner="ylallemant"
repo="go-picam-streamer"
exe_name="picam-streamer"
githubUrl=""
githubApiUrl=""

relative_executable_folder=".local/bin"
absolute_executable_folder="${HOME}/$relative_executable_folder" # Eventually, the executable file will be placed here

# make sure PATH is properly set, in some pipeline it may not be the case
PATH=${absolute_executable_folder}:$PATH

separator="-"

get_arch() {
    # darwin/amd64: Darwin axetroydeMacBook-Air.local 20.5.0 Darwin Kernel Version 20.5.0: Sat May  8 05:10:33 PDT 2021; root:xnu-7195.121.3~9/RELEASE_X86_64 x86_64
    # linux/amd64: Linux test-ubuntu1804 5.4.0-42-generic #46~18.04.1-Ubuntu SMP Fri Jul 10 07:21:24 UTC 2020 x86_64 x86_64 x86_64 GNU/Linux
    a=$(uname -m)
    case ${a} in
        "x86_64" | "amd64" )
            echo "amd64"
        ;;
        "i386" | "i486" | "i586")
            echo "386"
        ;;
        "aarch64" | "arm64" | "arm")
            echo "arm64"
        ;;
        # "mips64el")
        #     echo "mips64el"
        # ;;
        # "mips64")
        #     echo "mips64"
        # ;;
        # "mips")
        #     echo "mips"
        # ;;
        *)
            echo ${NIL}
        ;;
    esac
}

get_os(){
    # darwin: Darwin
    echo $(uname -s | awk '{print tolower($0)}')
}

# parse flag
for i in "$@"; do
    case $i in
        -v=*|--version=*)
            version="${i#*=}"
            shift # past argument=value
        ;;
        *)
            # unknown option
        ;;
    esac
done

if [ -z "$exe_name" ]; then
    exe_name=$repo
    echo "INFO: file name is not specified, use '$repo'"
    echo "INFO: if you want to specify the name of the executable, set flag --exe=name"
fi

if [ -z "$githubUrl" ]; then
    githubUrl="https://github.com"
fi
if [ -z "$githubApiUrl" ]; then
    githubApiUrl="https://api.github.com"
fi


echo "prepare installation of binary $exe_name from repository github.com/$owner/$repo"


if [ -z "$version" ]; then
  echo "no version has been requested, retrieving latest version from the repository"
  version=$(curl -s https://api.github.com/repos/$owner/$repo/releases/latest | grep -m1 -Eo "$exe_name-[^/]+-linux-amd64.tar.gz" | grep -Eo "([0-9]+\.[0-9]+\.[0-9]+)")
else
  echo "version $version has been requested"
fi

if command -v $exe_name >/dev/null; then
  local_version=$(command $exe_name version --semver)
  echo "$exe_name is already installed in version $local_version"
  if [ $local_version = $version ]; then
    echo "local installation has already the wanted version, nothing to do"
    exit 0
  else
    echo "local $local_version and wanted $version diverge, start installation"
  fi
fi

downloadFolder="${TMPDIR:-/tmp}"
mkdir -p ${downloadFolder} # make sure download folder exists
os=$(get_os)
arch=$(get_arch)
file_name="${exe_name}${separator}${version}${separator}${os}${separator}${arch}.tar.gz" # the file name should be download
downloaded_file="${downloadFolder}/${file_name}" # the file path should be download

mkdir -p $absolute_executable_folder

# if version is empty
if [ -z "$version" ]; then
    asset_path=$(
        command curl -L \
            -H "Accept: application/vnd.github+json" \
            -H "X-GitHub-Api-Version: 2022-11-28" \
            ${githubApiUrl}/repos/${owner}/${repo}/releases |
        command grep -o "/${owner}/${repo}/releases/download/.*/${file_name}" |
        command head -n 1
    )
    if [[ ! "$asset_path" ]]; then
        echo "ERROR: unable to find a release asset called ${file_name}"
        exit 1
    fi
    asset_uri="${githubUrl}${asset_path}"
else
    asset_uri="${githubUrl}/${owner}/${repo}/releases/download/${version}/${file_name}"
fi

echo "[1/3] Download ${asset_uri} to ${downloadFolder}"
rm -f ${downloaded_file}
curl --fail --location --output "${downloaded_file}" "${asset_uri}"

echo "[2/3] Install ${exe_name} to the ${absolute_executable_folder}"
tar -xz -f ${downloaded_file} -C ${absolute_executable_folder}
exe_path=${absolute_executable_folder}/${exe_name}
chmod +x ${exe_path}
echo "      ${exe_name} was installed successfully to ${exe_path}"

echo "[3/3] Set environment variables"
shell_profile_file=".profile"

case $SHELL in
*/zsh)
  shell_profile_file=".zprofile"
  ;;
esac

echo "      assuming shell profile file at $HOME/$shell_profile_file"

if [ -z "$(grep "/$relative_executable_folder" "$HOME/$shell_profile_file")" ]; then
    echo "      add the ${absolute_executable_folder} directory to your \$HOME/$shell_profile_file"
    echo "
# set PATH so it includes user's private bin if it exists
if [ -d \"${absolute_executable_folder@Q}\" ] ; then
    PATH=\"${absolute_executable_folder@Q}:\$PATH\"
fi

" >> $HOME/$shell_profile_file

    export PATH=${absolute_executable_folder}:$PATH
    echo "Run '$exe_name --help' to get started"
else
    echo "Run '$exe_name --help' to get started"
fi

exit 0
