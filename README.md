# ShotSeek

#postgres DB access
To access your Postgres DB instance, run the following command in your terminal

`docker exec -it shotseek-db-1 psql -U admin -d shotseek`

### Generate localhost https certs for dev
```bash
openssl req -new -x509 -keyout .keys/https/localhost.key -out .keys/https/localhost.crt -days 365 -nodes -config .keys/https/openssl.cnf
```

#### Import into the trusted store:
##### Mac:
```bash
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain /Users/m/dev/shotseek/.keys/https/localhost.crt
```
Replace `/Users/m/dev/shotseek/.keys/https/localhost.crt` with the actual path to the certificate

##### Windows
Run the following commands in powershell
```powershell
$certPath = "C:\path\to\your\localhost.crt"
$certStore = "Cert:\LocalMachine\Root"
Import-Certificate -FilePath $certPath -CertStoreLocation $certStore
```

---
---
Ignore the below, here for Mike's reference -- temporarily
#### Cleaning the geonames allCountries.txt
```
grep -P '^\d' allCountries.txt > cleaned_data.txt
```

```
awk -F'\t' '$6 >= -90 && $6 <= 90 && $7 >= -180 && $7 <= 180' cleaned_data.txt > final_data.txt
```

```
awk -F'\t' 'BEGIN {OFS=","} {
    # Handle missing state values
    state_name = ($4 == "" ? "NULL" : "\"" $4 "\"")
    state_code = ($5 == "" ? "NULL" : "\"" $5 "\"")
    latitude = ($6 == "" ? "NULL" : $6)
    longitude = ($7 == "" ? "NULL" : $7)

    # Escape single quotes in fields
    gsub(/'\''/, "''", $3)   # Escape city name
    gsub(/'\''/, "''", $4)   # Escape state name
    gsub(/'\''/, "''", $5)   # Escape state code

    # Create SQL insert statement
    print "INSERT INTO locations (country_code, postal_code, city_name, state_name, state_code, latitude, longitude) VALUES (\"" $1 "\", \"" $2 "\", \"" $3 "\", " state_name ", " state_code ", " latitude ", " longitude ");"
}' final_data.txt > final_sql.sql

```

#### Import into db
```postgres
\i /path/to/final_sql.sql
```


# RabbitMQ
#### IN PROGRESS

To maybe share queues and such
```
docker run --rm -v rabbitmq_rabbitmq_data:/data -v $(pwd):/backup alpine tar czf /backup/rabbitmq_backup.tar.gz -C /data .
```

