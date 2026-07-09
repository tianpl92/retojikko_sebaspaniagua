#Context
Refactor how to save public proposals into poject app.

#Role
Expert developer web frontend services

#Requeriments
1. Use html5 and css for pages
2. Use javascript to enable communication between backend service and front service.
3. Use golang for backend service coding if it's necessary.
3. Check if backend service is running, else, run backend service "public proposals" into local machine to be available with frontend layer. (check with /health endpoint).

#Agents - Subagents
1.  Principal agent to run plans, read documents, generate documentation, read md instructions, use deepseek/deepseek-v4-flash model from openrouter.
2.  Subagent to generate coding html css files and coding scripts, coding testing and run sql scripts use deepseek/deepseek-v4-pro model from openrouter.

#Designer description frontend pages

*   Color palette:
Use the follow color palette to the skeleton front application:
    -primary color: rgb(74, 68, 102) hexadecimal is 4A4466.
    -second color: rgb(110, 173, 188) hexadecimal is 6EADBC.
    -auxiliar colors: rgb(159, 203, 173) hexadecimal is 9FCBAD and rgb(241, 247, 212) hexadecimal is F1F7D4

#Actions

1. Database actions:
    1.1 Delete all data saved into public_call_user_associations and public_calls_proposals tables.
    1.2 Check if public_call_user_associations and public_calls_proposals tables are empty.

2. Backend service actions:
    2.1 Shutdown backend service.
    2.2 Prepare new POST endpoint called "proposal-save". The target is save a public proposal data into public_calls_proposals table.
    2.3 The endpoint keep JWT behavior and cors configurations.
    2.4 Check datos.gov.co documentation, reference about how data must be receives like a struct input and save content into table public_calls_proposals.
    2.5 Check if the proposal information requested was saving yet. For this case update proposals information. Use the params "id del proceso" (or "referencia_del_proceso" is the same),  for the validation exists, this param or column data can't be updated.
    2.6 Update /saved_proposals endpoint (GET), to obtains complete data for proposal saved by user. Apply a join sql script by foreign key from public_call_id column for obtains public proposals data from public_calls_proposals table. Update response struct for returns the list proposals saved with complete information struct.
    2.7 Add unit test and check.
    2.8 Update api contracts documentation.
    2.9 Run service up.


3. Frontend service actions.
    3.1 Dashboard page: When a proposal into dashboard is checked, and press "guardar" action button, update behavior, using next steps:
        3.1.1 First do a request to endpoint /proposal-save (POST) for saving the proposal.(check api contract updated)
        3.1.2 Second, do a request to /saved-proposals (POST) endpoint, for save user id and proposal relation.
        3.1.3 If many proposals are checked, loop the behavior for step 1.1.1 and 1.1.2 for every proposal when button "Guardar..." perform click event.
    3.2 Saved-Proposals page: Change behavior using next steps: 
        3.2.1 When load keeps JWT token validacion and then consumes /saved_proposals (GET).
        3.2.2 Now the endpoint returns all complete information about saved proposals. 
        3.2.3 When show the modal window about proposal saved, show real data , remove mock data information.



#Results expected

1. New endpoint created: "/proposal-save" into backend service.
2. front service can do requests for save proposals and saving user_id and proposal_call_id relation.
3. Front service can shows saved proposals and obtains the complete public proposal data information.

-- 

Do not commit changes yet.

--

Tell me first of all again, just to be sure, what do you plan to do?
