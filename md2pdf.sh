#!/bin/bash

# require
# pandoc: `sudo apt install pandoc`
# wkhtmltopdf:
#    Ubuntu: `sudo apt-get install xvfb libfontconfig wkhtmltopdf`
#    iOS: `brew install wkhtmltopdf`

md_filename=$1
tmp1_md_filename=${md_filename%.*}".tmp1.md"
html_filename=${md_filename%.*}".html"
tmp1_html_filename=${md_filename%.*}".tmp1.html"
tmp2_html_filename=${md_filename%.*}".tmp2.html"
tmp3_html_filename=${md_filename%.*}".tmp3.html"
pdf_filename=${md_filename%.*}".pdf"

if PANDOC=$(which pandoc); then
  awk '{ gsub(/\\\//, "/"); print }' $md_filename > $tmp1_md_filename

  pandoc --from=markdown_github-yaml_metadata_block --standalone \
  --to=html -V -H md.style --output=$tmp1_html_filename $tmp1_md_filename

  pandoc --from=markdown_github-yaml_metadata_block --standalone \
  --to=html -V -H pdf.style --output=$tmp2_html_filename $tmp1_md_filename

  # replace :warninig: image
  awk '{ gsub(/:warning:/, "<span class=\"warn warning\"></span>"); print }' $tmp1_html_filename > $html_filename
  awk '{ gsub(/:warning:/, "<span class=\"warn warning\"></span>"); print }' $tmp2_html_filename > $tmp3_html_filename

  if WKHTMLTOPDF=$(which wkhtmltopdf); then
    wkhtmltopdf --orientation landscape -q -s A4 --enable-internal-links --dpi 300 $tmp3_html_filename $pdf_filename
  else
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')] 'wkhtmltopdf' not found. Generating of html report skipped."
  fi

  rm $tmp1_html_filename
  rm $tmp2_html_filename
  rm $tmp1_md_filename
  rm $tmp3_html_filename

  if [[ -f $html_filename ]]; then
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')] Final .html report is ready at:"
    echo "        '$html_filename'"
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]"
  fi

  if [[ -f $pdf_filename ]]; then
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')] Final .pdf report is ready at:"
    echo "        '$pdf_filename'"
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]"
  fi
else
  echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')] 'pandoc' not found. Generating of html report skipped."
fi
