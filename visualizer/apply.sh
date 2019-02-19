#!/bin/bash

bq --project_id "slavayssiere-sandbox" --location=EU mk test_bq
bq --project_id "slavayssiere-sandbox" --location=EU mk --external_table_definition=./simple_definition.json test_bq.ms

# cbt -instance "test-instance" read "test-table"

