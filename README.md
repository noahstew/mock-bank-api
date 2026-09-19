# Starting:

make postgres

make createdb

make migrateup

# DB Migration:

migrate create -ext sql -dir db/migration -seq <migration_name>
