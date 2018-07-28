FROM mysql:latest
COPY setup.sh /mysql/setup.sh
# COPY data-dump.sql /mysql/setup.sql
COPY data-dump.sql /docker-entrypoint-initdb.d/
# ADD data-dump.sql /docker-entrypoint-initdb.d
# RUN /bin/bash -c "/usr/bin/mysqld_safe --skip-grant-tables &" && \
#   sleep 5 && \
#   # mysql -u root -e "CREATE DATABASE fibear" && \
#   mysql -u root fibear < /mysql/setup.sql
# RUN ["chmod", "+x", "/mysql/setup.sh"]
# RUN /mysql/setup.sh