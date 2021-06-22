# You shall not pass: Mitigating SQL Injection Attacks on Legacy Web
Applications 

This repository contains the source code of SQLBlock for the paper
published in Proceedings of the 2020 on Asia Conference on Computer and
Communications Security

## Citation
Please use the following bibtex for citation:
```
@inproceedings{jahanshahi2020sqlblock,
author = {Jahanshahi, Rasoul and Doup\'{e}, Adam and Egele, Manuel},
title = {You Shall Not Pass: Mitigating SQL Injection Attacks on Legacy Web Applications},
booktitle = {Proceedings of the 15th ACM Asia Conference on Computer and Communications Security},
pages = {445–457},
year = {2020},
url = {https://doi.org/10.1145/3320269.3384760},
}
```

## Resources Links
The following links are for the paper:
[PDF 1](https://megele.io/asiaccs20-sqlblock.pdf)
, [PDF 2](https://dl.acm.org/doi/pdf/10.1145/3320269.3384760)

## Folder Organization
The folder organization is listed below.
```bash
.
|-- Dockerfile 		# dockerfile for building an experimental environment for SQLBlock
|-- runme.sh 		# shell-script for downloading necessary libraries to build docker container
|-- storage
|--    |--sa 		# static analysis of SQLBlock
|--    |--webapp        # web application running in docker
|--    |--sqlblock      # scripts to automate recording and enforcing profiles
|--    |--config        # config files for Nginx, MySQL server, and PHP
```

## Instruction

To facilitate the evaluation of SQLBlock, most of our instructions are based
around Docker containers. We create a docker container which includes a
modified MySQL server, a modified PHP dbi, and our static analysis.

### Build SQLBlock container
`runme.sh` downloads required artifacts for SQLBlock to intercept and record
issued queries to the database, and enforcing a profile.  The Downloaded
artifacts includes the modified MySQL server 5.7, SQLBlock plugin, and PHP
modified dbi. In the next step, `runme.sh` builds the docker container and create a database with the following information.
```bash
databae name: mysqldb
username: admin
password: admin
```

- Before building the docker container, you should copy the web application source-code to the directory of `/storage/webapp`.
### Run the docker container
you can use the following command to run the docker container.
```bash
docker run --name sqlblock-container -d --rm -p 9000:80 -it sqlblock
```
The web application will be accessible on `localhost:9000`.

### Generate a profile for the web application
To generate a profile for the web application use the command below to access a shell from the docker container.
```bash
docker exec -it sqlblock-container bash
```
In the next step, run the following command to load SQLBlock plugin to MySQL server and configure the plugin for recording issued query to the database.
```bash
/home/record_mysql_query
```
Now, you can visit the web application in your browser and perform various operations and SQLBlock will record any issued query to the database.
After finishing the browsing, execute the script `runme` inside the docker. The `runme` script compiles SQLBlock static analysis, analyzes the web application reside under `/var/www/html` directory, and generate a profile for the web application for the issued queries.

### Enforce a profile for the web application
Run the following command in the docker container to enforce the generated profile in the previous step.
```bash
/home/enforce_profile
```

## Contact Us
If you require any further information, send an email to `rasoulj@bu.edu`
