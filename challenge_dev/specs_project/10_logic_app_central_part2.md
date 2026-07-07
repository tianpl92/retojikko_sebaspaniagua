#Context
Continue and improve better the backend phase project, now we go to finish the logic for the last endpoints.

#Role
Software engineer developer api services.


#tools requeriments
1. Use golang to coding, use the golang version current machine.
2. Use connection database from local postgres running.
3. Connecting with database "portal_plan_public_app". (check if exists)
4. Use .env file into project backend for database user and credentials connection.

#Agents - Subagents
1.  Principal agent to run plans, read documents, generate documentation, read md instructions, use deepseek/deepseek-v4-flash model from openrouter.
2.  Subagent to generate go files and coding scripts and coding testing use deepseek/deepseek-v4-pro model from openrouter.




#Actions Coding

1. build the logic service for the others endpoints, follow the instructions per endpoint:

    * endpoint "/saved-proposals" (POST): This endpoint has the target to save public proposals by user authenticated. So features for endpoint are:
        -Input params required: user_id (from login authenticated) and id_proposal to save.
        -Register into database the relation for the user and the proposal.
        -The response for success process must be a message property json with value "Guardado satisfactoriamente".
        -If the relation user id and id proposals exists, don't do anything and don't create duplicated row. The response for this case just must be "Guardado satisfactoriamente" too.
        -keep JWT validation.
        -Update documentation into "apicontract.md" file, about this endpoint.

    * endpoint "/saved-proposals" (GET): This endpoint has the target to search and list public proposals saved by user. Features are:
        -Input params required: user_id (from login authenticated).
        -Responses the list of public proposals saved by user.
        -List proposals is an array of proposals object with complete fields from database.
        -keep JWT validation.
        -if the request by user does not have proposals saved by user, returns an empty array, with success response status.
        -Update documentation into "apicontract.md" file, about this endpoint.

    * endpoint "/user-modify" (POST): This endpoint has the target to update registered info users. Features are:
        -Input params are optional: first_name, last_name, gender, phone_number and status. The params  (id) , email , status and password are dont available to modify by this endpoint.
        -If there are not input params (or all are empty), show message "Debe modificarse un campo por lo menos para actualizar".
        -If the input params contains not available params, do not do anything and do not modify data.
        -If the he process finish without errors, response json object message property shows "Informacion actualizada correctamente".
        -If the request only has not available input params, response object message shows "Data not updated".
        -keep JWT validation.
        -Update documentation into "apicontract.md" file, about this endpoint.

    * endpoint "/user-info" (POST): This endpoint has the target to obtain user info authenticated into application.
        -Input params available and required: user_id.
        -If the process is success, returns a json object with the users info.
        -If id users dont exists, show message response "Usuario no registrado".
        -Do not include information into object response about password.
        -keep JWT validation.
        -Update documentation into "apicontract.md" file, about this endpoint.


#Result actions

1. The endpoints "/saved-proposals" (GET and POST) , "/user-modify" and "/user-info", with logic implementated , validations created and connect with database repository.
2. Update api contract documentation.
3. No commit and don't sync changes with current branch yet.

---------

Tell me first of all, just to be sure, what do you plan to do?



