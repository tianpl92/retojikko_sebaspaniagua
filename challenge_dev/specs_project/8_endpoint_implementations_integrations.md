#context
Continue and improve better the backend phase project. The target is coding endpoints using authentication JWT method.

#role
Software engineer developer api services 

#tools requeriments
1. Use golang to coding, with the version local machine.
2. Use connection database from local postgres running.
3. Connecting with database "portal_plan_public_app". (check if exists)
4. Use .env file into project backend for database user and credentials connection.

#Agents - Subagents
1.  Principal agent to run plans, read documents, generate documentation, read md instructions, use deepseek/deepseek-v4-flash model from openrouter.
2.  Subagent to generate go files and coding scripts use deepseek/deepseek-v4-pro model from openrouter.
3.  Use gpt-5.5 model from openrouter only for coding the logic with JWT method authentications, other coding or new sql scripts use deepseek/deepseek-v4-pro.


#actions

1. Build next endpoints into service:
    1.1 Build /healt endpoint (for check service up)
    1.2 Create new endpoint /user-create (POST), for create users. That endpoint let creates users for application, with this specifications:
        a. Users with the same document number identification can't be created again.If existing user with the same document , the endpoint responses this message: "User already exists".
        b. Add validation to analyse if email user already exists. If the email exists , don't create user and show response message: "email already in use".
        c. The params input user there are the same columns into "users" table.
        d. Add validation required if the column database item has not null specification.
    1.3 Build /login endpoint (use for enable token app . Use that token for consume the others endpoint. Only request with token will works and response with status 200 range; And use the next specifications for login endpoint:
        a. Build login service use JWT (JSON Web Token) method.
        b. If there is neccesary aditional libs, please download, install and include into go imports project.
        c. The endpoint needs (required), email and password input params.
        d. The info users must be integrated with "users" table database. 
        e. Prepare token with one hour available session active.
        f. If there is neccesary create a new table to register sessions using user id like foreign key.
    1.4 Prepare and create next endpoints into handlers management, for use with available token generated from /login endpoint. (but Before, this is important:  "/healt" , "/user-create" and "/login" endpoints are the exception, don't needed token JWT to works).
            a. /public-proposals (GET)
            b. /saved_proposals (GET)
            c. /saved-proposals (POST)
            d. /user-info (GET)
            e. /user-modify (POST)
    1.5 Prepare endpoints with token required with the next response message, if the token is not available or is not active ,show: "Credentials invalid".
    1.6 Update database/public_calls_database.sql script with the new tables or new structure implemented.

2. Update plan hermes specification (make a resume) using previos information.

---------

Tell me first of all, just to be sure, what do you plan to do?

  
