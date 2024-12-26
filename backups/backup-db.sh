export $(grep -v '^#' .env | xargs)

CONTAINER_NAME="mysql-employees-system"
DATABASE_NAME=${DB_DATABASE}
BACKUP_PATH="./backups"
TIMESTAMP=$(date +"%F_%T")
BACKUP_FILE="${BACKUP_PATH}/${DATABASE_NAME}_${TIMESTAMP}.sql"
PASSWORD=${DB_PASSWORD}

# Remove backups older than 7 days
find ${BACKUP_PATH} -name '*.sql' -mtime +7 -exec rm -f {} \;

# Create backup sql on folder backups
docker exec ${CONTAINER_NAME} /usr/bin/mysqldump -u ${DB_USERNAME} --password=${PASSWORD} ${DATABASE_NAME} > ${BACKUP_FILE}

# Compress backup file
tar -czf "${BACKUP_FILE}.tar.gz" -C ${BACKUP_PATH} $(basename ${BACKUP_FILE})

# Remove backup file
rm ${BACKUP_FILE}

