# Indexes

Notes on my small talk about indexes

## Demo toolkit

A playground to play with indexes is included in the src folder. Just enter there and run

```
docker-compose up
```

And you will have a pre-populated db with the tables used in this talk. You can connect to mySQL at 3306 port username as password ```indexes```

The first time you start it, it will take several minutes to create all data. We will have a reasonably big database to play with.

If that is too much for your machine you can always edit docker-compose.yml to set an adequate number of records and workers specific to your environment.