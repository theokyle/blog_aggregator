# Gator

Gator is a command line tool for aggregating RSS feeds and viewing the posts.

## Installation

You will need the latest version of the [Go toolchain](https://golang.org/dl/) installed as well as a local postgres database.
You can then install with the command:

```bash
go install ...
```

## Config

Create a '.gatorconfig.json' file in your home directory with the following structure:

```json
{
    "db_url": "postgres://username:@localhost:5432/database?sslmode=?disable"
}
```

Replace the values with your database connection string

## Using the Aggregator

You can create a user with the command:

```bash
gator register <name>
```

You can add a feed with the command:

```bash
gator addfeed <url>
```

After you have added a feed to your username, you can start the aggregator with the command:

```bash
gator agg <duration>
```

The duration will specify how frequently the aggregator will update feeds (such as 30s).

Finally, you can browse feeds with the 'browse' command:

```bash
gator browse <limit>
```

This command will specify a limit of how many posts you would like to view.

## Other commands

- "gator login <name>" - login as a user that already exists. The register command automatically logs in the registered user.
- "gator users" - List all users
- "gator feeds" - List all programs
- "gator follow <url>" - Current user will follow a feed that already exists in the database
- "gator unfollow <url>" - Current user will unfollow a feed that already exists in the database

