
wget http://sourceforge.net/projects/boost/files/boost/1.59.0/boost_1_59_0.tar.gz
tar xzf boost_1_59_0.tar.gz
mv boost_1_59_0 storage/

git clone https://git-seclab.bu.edu/rasoulj/mysql-server
git checkout 5.7
mv mysql-server storage/

mkdir /storage/php
git clone https://git-seclab.bu.edu/rasoulj/php-mysqli /storage/php/php-mysqli
git clone https://git-seclab.bu.edu/rasoulj/pdo-mysql /storage/php/pdo-mysql

docker build -t sqlblock .
