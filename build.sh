#!/bin/bash
#
# Copyright (C) 2018 Roberto Mier Escandon <rmescandon@gmail.com>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License version 3 as
# published by the Free Software Foundation.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

set -e

default_plugins_src_path="./makers/plugins"
default_plugins_dst_path="./_plugins"
plugins_src=$1
plugins_dst=$2
skip_plugins=0

show_help() {
    exec cat <<EOF
Usage: build.sh <plugins_src> <plugins_dst>

Builds CRUDer project along with the built-in plugins

optional arguments:
  --help                Show this help message and exit
  --skip_plugins        Skip built-in plugins building

Positional arguments:
  <plugins_src>         Path to the folder with the plugins source files (default: $default_plugins_src_path)
  <plugins_dst          Path to the folder with the plugins library .so files (default: $default_plugins_dst_path)
EOF
}

build_plugins() {
    if [ -z "$1" ]; then
        plugins_src="$default_plugins_src_path"
    fi

    if [ -z "$2" ]; then
        plugins_dst="$default_plugins_dst_path"
    fi

    [ -d $plugins_dst ] || mkdir -p 755 $plugins_dst

    echo "Building plugins..."
    for plugin in $(find "$plugins_src" -type f -name "*.go" | grep -v ".*_test.go"); do
        filename=$(basename "$plugin")
        filename="${filename%.*}"
        mod="$filename.so"
        go build -buildmode=plugin -o "$plugins_dst/$mod" "$plugin"
        echo "$mod"
    done

    echo "All plugins built."
    echo "Find generated shared modules at $plugins_dst"
}

while [ -n "$1" ]; do
	case "$1" in
        -h)
            ;&
		--help)
			show_help
			exit
			;;
        -s)
            ;&
        --skip_plugins)
			skip_plugins=1
			shift
			;;
		*)
			echo "Unknown command: $1"
			exit 1
			;;
	esac
done

if [ "$skip_plugins" -eq 0 ]; then
    build_plugins
fi
    
go install ./cmd/...
echo "Project built."
echo "All done."
