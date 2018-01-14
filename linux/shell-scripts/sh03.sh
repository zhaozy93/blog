#!/bin/bash



read -p 'pls input your filename:    ' filename
postfilename=${filename:-"filename"}
echo $postfilename
date1=$(date +%Y%m%d)
touch "$postfilename-----$date1"
exit 0
