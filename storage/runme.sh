#! /bin/bash

# compile static analysis
echo "[*] Compile the static analysis"
cd ~/.go/src/parse-function/
go build
echo "[*] Run the static analysis on /var/www/html \n"
./parse-function /var/www/html/ /home/sa_output

echo "[*] Generate the profile based on /tmp/mysql_record"
/home/createProfile.py /tmp/mysql_record /tmp/profilemysql /home/sa_output
cp /home/sa_output /tmp/classes
