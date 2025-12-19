#!/bin/bash

set -e

BACKUP_DIR="/var/backups/dna-relations"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/dna_backup_$TIMESTAMP.sql.gz"

mkdir -p $BACKUP_DIR

source /opt/dna-relations/.env

docker exec dna-postgres pg_dump -U $POSTGRES_USER $POSTGRES_DB | gzip > $BACKUP_FILE

find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete

echo "Backup created: $BACKUP_FILE"

if [ -n "$BACKUP_S3_BUCKET" ]; then
    aws s3 cp $BACKUP_FILE s3://$BACKUP_S3_BUCKET/dna-relations/
    echo "Uploaded to S3"
fi
