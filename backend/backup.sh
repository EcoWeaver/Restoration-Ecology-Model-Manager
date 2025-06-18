set -a
source .env
set +a

DB_NAME="postgres"
DB_USER=$POSTGRES_USER
CONTAINER_NAME="postgres"
BACKUP_DIR="./backup"
DATE=$(date +%Y-%m-%d_%H-%M)
BACKUP_FILE="backup_$DATE.sql"
BACKUP_PATH="$BACKUP_DIR/$BACKUP_FILE"

echo "Creating Backup to $BACKUP_PATH."

docker exec -e PGPASSWORD=$POSTGRES_PASSWORD $CONTAINER_NAME \
  pg_dump -U $DB_USER -d $DB_NAME \
  -t public.models -t public.metadata \
  --data-only --column-inserts --no-comments --no-owner \
  > "$BACKUP_PATH"

echo "Backup finished."

cd ./backup_preprocessing 

go build -o preprocess_backup main.go
if [ ! -f preprocess_backup ]; then
  echo "Go executable not created."
  exit 1
fi

chmod +x preprocess_backup


../backup_preprocessing/preprocess_backup "../backup/$BACKUP_FILE"

if [ $? -ne 0 ]; then
  echo "Go preprocessor failed."
  exit 1
else
  echo "Go preprocessor finished successfully."
fi

rm preprocess_backup

# Delete some obsolete files. 
if [ -d "$BACKUP_DIR" ]; then
  find "$BACKUP_DIR" -name "*.sql" -type f -mtime +30 -delete
else
  echo "Backup directory $BACKUP_DIR does not exist, skipping cleanup."
fi
