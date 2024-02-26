#!/bin/sh
# Dump local dev database structure and required data from specific tables.
#
# Usage: ./mysql-dump.sh
# ------------------------------------------------------------------------------
# shellcheck disable=SC2046 # Ignore SC2046: Quote this to prevent word splitting.

SD=$(dirname "$(readlink -f -- "$0")")
cd "$SD" || exit 1 && cd ".." || exit 1

if [ ! -f ".env" ]; then
  echo "Missing .env file. Please copy .env.example to .env and edit it."
  exit 1
fi

export $(grep -v '^#' .env | sed 's/#.*//g'| xargs)
## Structure only dump
mariadb-dump --no-data --skip-comments \
  -h "${DB_HOST}"                      \
  -u "${DB_USER}"                      \
  -p"${DB_PASSWORD}"                   \
  "${DB_NAME}" |                       \
  sed 's/ AUTO_INCREMENT=[0-9]*//g' >  \
  "./tools/resources/database/structure.sql"

# vim: set ts=2 sw=2 tw=0 et ft=sh:
