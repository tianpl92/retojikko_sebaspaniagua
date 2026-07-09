#Context
When project bakend service start running, assumes database was creted. So, The backend service will must create schema database if not exists.

#Roles
-Database Enginner software developer. 
-Software developer for Api services. 

#Actions

1. Create a validation script about check if schema database exists into postgres server. 
2. Using Go language and SQL sintax. 
3. If the database does not exists, use "database/public_calls_database.sql" file script to generate the schema.
4. Use .env file to connect server database.
5. The validation schema script, must be executed when running up backend service.
6. Important annotation: Do not drop schema if exists.

#Expected result

1. Script generated for validate if shema exists.
2. Create database schema if the validation do not pass.
3. Do not drop database or create it again if exists. Just skip.
4. The validation is executing when backend service is running up.

-- 

Do not commit changes yet.

--

Tell me first of all again, just to be sure, what do you plan to do?

