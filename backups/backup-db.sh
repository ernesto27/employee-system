export $(grep -v '^#' .env | xargs)

CONTAINER_NAME="mysql-employees-system"
DATABASE_NAME=${DB_DATABASE}
BACKUP_PATH="./backups"
TIMESTAMP=$(date +"%F_%T")
BACKUP_FILE="${BACKUP_PATH}/${DATABASE_NAME}_${TIMESTAMP}.sql"
PASSWORD=${DB_PASSWORD}

find ${BACKUP_PATH} -name '*.sql' -mtime +7 -exec rm -f {} \;

docker exec ${CONTAINER_NAME} /usr/bin/mysqldump -u ${DB_USERNAME} --password=${PASSWORD} ${DATABASE_NAME} > ${BACKUP_FILE}

tar -czf "${BACKUP_FILE}.tar.gz" -C ${BACKUP_PATH} $(basename ${BACKUP_FILE})
rm ${BACKUP_FILE}

