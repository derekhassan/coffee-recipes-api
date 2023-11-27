# Coffee Recipe API

Simple API for creating and viewing specialty coffee recipes.

## Getting Started

### Setting up the Database

This application uses MySQL for storage. The recommended approach is running the database in a Docker container. To pull and run the MySQL image, run the following command with the following variables:
- `MYSQL_ROOT_PASSWORD`: The password to the root user
- `MYSQL_DATABASE`: The name of the database
- `MYSQL_USER`: The name of the new user to create
- `MYSQL_PASSWORD`: The password of the new user to create

```
docker run --name coffee_recipes -p 3306:3306  \
    -e MYSQL_ROOT_PASSWORD=<your-password> \
    -e MYSQL_DATABASE=coffee_recipes \
    -e MYSQL_USER=<your-username> \
    -e MYSQL_PASSWORD=<your-password> \
    -d mysql:latest
```

This will run a new container with the name `coffee_recipes`.

### Initializing the Database

1. Copy the `init.sql` file to the container you created:
```sh
$ docker cp init.sql coffee_recipes:/init.sql
```

2. Connect to the container:
```sh
$ docker exec -it coffee_recipes bash
```

3. Log into the `mysql` CLI
```sh
$ mysql -u root -p
```

4. Select the database that was created:
```sh
mysql> USE coffee_recipes;
```

5. Run the script:
```sh
mysql> source init.sql
```

### Creating .env File

Copy and rename the `.env.example` file into a file named `.env` and fill in the variables based on what you used in the database creation setup.