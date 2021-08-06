# SQLBlock defense deployed on Drupal

You can download an already setup docker container from the below link
[sqlblock_drupal.tar](https://drive.google.com/drive/folders/1sJslTelkODYPgtIoTWXD_lI-ws3kjLom)

## Running the docker
To run the docker, you should first load the tarball file using the following command.
```bash
docker load -i sqlblock_drupal.tar
```

You can use the following command to run the docker container.
```bash
docker run --name sqlblock-container -d --rm -p 9000:80 -it sqlblock_drupal
```

This docker container includes Drupal 7.0 which is vulnerable to Drupalgeddon SQLi. The docker container is already loaded with the profile of benign browsing in Drupal which includes only simple actions.


## Exploit the vulnerabilities

In the prepared docker container, SQLBlock is defending against the SQli requests. You can turn off SQLBlock using the following command.
```bash
mysql -u root -e "UNINSTALL PLUGIN sqlblock"
```

You can turn SQLBlcok back on, either using the script called `enforce_profile` or you can run the following command.
```bash
mysql -u root -e "INSTALL PLUGIN sqlblock SONAME 'sqlblock.so'"
```

You can exploit the vulnerabilities in Drupal using the python script in the shared google drive from this [link](https://drive.google.com/drive/folders/1sJslTelkODYPgtIoTWXD_lI-ws3kjLom).

# Drupalgeddon vulnerability
You can run the following command to perform the SQLi exploit on auto-suggest plugin in WordPress.
```bash
python ./drupalgeddon.py -t http://localhost:9000 -u admin1 -p admin1
```

The above command will create an admin user with username and password of `admin1`.
