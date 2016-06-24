#!/bin/bash
TIMESTAMP=$(date +%s)
cat >  migrations/mysql/${TIMESTAMP}_$1.sql <<- TEMPLATE
-- +migrate Up

-- +migrate Down

TEMPLATE
