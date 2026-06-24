#!/bin/sh
set -e
./app migrate
exec ./app serve
