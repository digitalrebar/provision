#!/usr/bin/env bash

# Check to make sure ruby is installed
if ! [[ -f $(which ruby) ]] || ! [[ -f $(which gem) ]]; then
  echo "Ruby must be installed for the docs indexer to function"
  exit 1
fi

# If json built docs don't exist, generate them
if ! [[ -d _build/json ]]; then
  make json
fi

# Install nokogiri gem (for html parsing)
gem list | grep nokogiri --s || gem install nokogiri

# Run a ruby script to:
#  - extract doc contents as text
#  - pair the contents with the page name
#  - return a single json blob with all of the contents
ruby <<EOF
require 'nokogiri'
require 'json'

# Find all json doc files
doc_files = Dir["_build/json/doc/**/*.fjson"]

indexes = {}

doc_files.each do |f|
  # Read each file as json
  blob = JSON.parse(open(f, "r").read)
  page = blob["current_page_name"]
  title_raw = Nokogiri::HTML.parse(blob["title"]).text
  title = title_raw.match(/^((\d+\.)+ )(.*)$/)[3]
  # Extract only the text from the html
  text = Nokogiri::HTML.parse(blob["body"]).text

  # Store in the indexes object
  indexes[page] = {title: title, text: text}
end

# Print out the indexes
puts indexes.to_json
EOF