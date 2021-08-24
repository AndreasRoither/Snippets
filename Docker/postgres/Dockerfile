
FROM postgres:13.4

# Set maintainer label: maintainer=[YOUR-EMAIL]
LABEL maintainer='andreas.roither'

# add the 'postgres' admin role
USER postgres

# Copy sql files to the docker entrypoint, it automatically connects and starts psql
# so we don't have to write anything specific in the sql file, just the sql commands
ADD createTables.sql /docker-entrypoint-initdb.d/

# expose Postgres port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Set the default command to run when starting the container
CMD ["postgres"]