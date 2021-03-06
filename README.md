Telegram-parser
----
The Telegram parser finds the most popular messages in the user's Telegram feed and provide sorted lists of them via API and UI. 

![](examples/png/all-time.png)

## Prerequisites and launch
1. Install and start database (PostgreSQL)
2. Run RabbitMQ in docker container
3. Install and build tdlib following the [instructions](install-tdlib.txt)
3. Set environment variables:

   - PORT - The port that the parser will listen to.
   
   - RABBITPORT - Port on which rabbitMQ is listened.
   - PGHOST - Host on which postgreSQL is listened.
   - PGPORT - Port on which postgreSQL is listened.
   - PGUSER - PostgreSQL user name.
   - PGPASSWORD - Password to access postgreSQL.
   - TGTELNUMBER - Phone number required to connect to the telegram client.
   - TGAPIID - Application identifier for Telegram API access, which can be obtained at https://my.telegram.org.
   - TGAPIHASH - Application identifier hash for Telegram API access, which can be obtained at https://my.telegram.org.
   
4. Untar modified vendor: 

        tar -xvf vendor.tar
        
5. Build and run app:
        
        go build 
        
        ./telegram-parser -mod=vendor


## UI
UI is available on port {env.PORT}

![](examples/png/today.png)

> Example of using UI


## API
GET `http://localhost:{PORT}/best?period=${period}`

Returns the most popular posts for a period:

Available time periods:

- today
- yesterday
- daybeforeyesterday
- thisweek
- lastweek
- thismonth
- whole (Denotes the entire period from 1970-01-01T00: 00: 00Z to the present)

GET `http://localhost:{PORT}/best/3hour` 

Returns the most popular posts in the last 3 hours
