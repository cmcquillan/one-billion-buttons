#!/bin/bash
export TF_VAR_DigitalOceanToken=$DO_TOKEN
export TF_VAR_CloudflareToken=$CO_TOKEN

# Set up Digital OCean
doctl auth init --access-token $DO_TOKEN

doctl serverless install

# Install Claude
npm install -g @anthropic-ai/claude-code@latest

