#!/bin/bash
set -o errexit

## --------------- game config DB settings ----------------
export GAME_CONFIG_DB_HOST=127.0.0.1
export GAME_CONFIG_DB_USER=root
export GAME_CONFIG_DB_PASS=123456
export GAME_CONFIG_DB_PORT=3306
export GAME_CONFIG_DB_DATABASE=meland_cnf_dev

## --------------- game data DB settings ----------------
export GAME_DB_HOST=127.0.0.1
export GAME_DB_PORT=3306
export GAME_DB_USER=root
export GAME_DB_PASS=123456
export GAME_DB_DATABASE=meland_game_data

## -------------- game TOKEN KEY ----------------
export JWT_SECRET=token key