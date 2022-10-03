#!/bin/bash
set -o errexit

## --------------- meland table config DB settings ----------------
export MELAND_CONFIG_DB_HOST=127.0.0.1
export MELAND_CONFIG_DB_USER=root
export MELAND_CONFIG_DB_PASS=123456
export MELAND_CONFIG_DB_PORT=3306
export MELAND_CONFIG_DB_DATABASE=meland_cnf_dev

## --------------- meland game data DB settings ----------------
export MELAND_GAME_DB_HOST=127.0.0.1
export MELAND_GAME_DB_USER=root
export MELAND_GAME_DB_PASS=123456
export MELAND_GAME_DB_PORT=3306
export MELAND_GAME_DB_DATABASE=meland_game_data

## -------------- meland TOKEN KEY ----------------
export token key