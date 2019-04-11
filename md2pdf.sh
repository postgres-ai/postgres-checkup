#!/bin/bash

# require  
# pandoc: `sudo apt install pandoc`
# wkhtmltopdf: `sudo apt-get install xvfb libfontconfig wkhtmltopdf`

md_filename=$1
#tmp1_md_filename=${md_filename%.*}".tmp1.md"
#tmp2_md_filename=${md_filename%.*}".tmp2.md"
tmp2_md_filename=$md_filename
html_filename=${md_filename%.*}".html"
tmp_html_filename=${md_filename%.*}".tmp.html"
pdf_filename=${md_filename%.*}".pdf"

#awk '{ gsub(/\|\ _/, "\\_"); print }' $md_filename > $tmp1_md_filename
#awk '{ gsub(/\\\//, "/"); print }' $tmp1_md_filename > $tmp2_md_filename

pandoc --from=markdown_github-yaml_metadata_block --smart --standalone \
--to=html -V -H md.style --output=$tmp_html_filename $tmp2_md_filename

awk '{ gsub(/:warning:/, "<span class=\"warn warning\"></span>"); print }' $tmp_html_filename > $html_filename
rm $tmp_html_filename
#rm $tmp1_md_filename
#rm $tmp2_md_filename

wkhtmltopdf -B 10mm -T 10mm -L 10mm -R 10mm \
-q -s Letter $html_filename $pdf_filename

ls ${md_filename%.*}*
