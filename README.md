# What is an `Indexer`?

Indexer - the App, is a web application written in `Go`, It indexes Ethereum's Consensus Layer (Beacon Chain) and stores it in postgres which can be visualised by visiting http://localhost:8080.
The App is designed to index top 5 slots only - for obvious reasons(storage). 

## How to use the `Indexer`?

It can be run independently connecting to any Beacon node & Postgres instance or by using `Docker`.

To run independently, you need to provide/override database connection credentials in `.env` under `/cmd` and url to Beacon node instance.

Following is a sample `.env` file

```.env
CLIENT_URL=https://{your-endpoint-name}.quiknode.pro/{your-token}/
POSTGRES_HOST=indexer_db:5432
POSTGRES_NAME=db
POSTGRES_USER=db_user
POSTGRES_PASSWORD=db_user_password
POSTGRES_DISABLE_TLS=true
```

###  

To run the app using `Docker` just type

```sh
  make up
```

View the indexed epochs/slots/blocks at http://0.0.0.0:7080

And to tear down

```sh
  make down
```


## Why `PostgresSQL`?

- PostgreSQL ensures data integrity and provides support for ACID (Atomicity, Consistency, Isolation, Durability) properties, making it suitable for handling critical and consistent data.
- PostgreSQL has been around for a long time and has a strong reputation for stability and reliability.
- PostgreSQL supports a wide range of data types, including JSON, JSONB which can be useful when dealing with complex Ethereum data structures. It also offers powerful querying capabilities, such as indexing, which can enhance performance on larger dataset.
- PostgreSQL is capable of handling large datasets and can scale to accommodate future growth if needed - when we plan on not just storing recent 5 epochs.
- PostgreSQL has very good Community Support and Tooling.