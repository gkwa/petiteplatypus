#!/bin/bash

petiteplatypus generate /tmp/trash --verbose
obsidian-cli set-default /tmp/trash
obsidian-cli print-default
