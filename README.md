## calorie_counter_bot
Telegram Bot for counting calories of your lunchs.

## Installing

For working you need to have mysql-server on your localhost and create `.env` file in root of project with 3 variables:

```sh
TOKEN= // YOUR TELEGRAM-API-TOKEN
DBUSER= // YOUR MYSQL USER
DBPASSWORD= // YOUR MYSQL PASSWORD
```
Run next command from project path to install DB:

```sh
mysql -uroot -h127.0.0.1 -p callorie_counter_bot < callorie_counter_bot.sql
```
## Starting

Service can be start with run Makefile instruction, which run docker-container.</br>

```sh
make run build-image
make run start-container
```

## How it works

The logic is built around the `waiting` parameter, which initially has the value `"no-waiting"`.</br>
Having received a ask from the user, which the bot knows, it will change the `waiting` value to the one necessary for the corresponding ask and start the required processing.</br>
After complete ask or the user violates the execution conditions (for example, writes digits instead of the expected text), the `waiting` value returns to `"no-waiting"`
