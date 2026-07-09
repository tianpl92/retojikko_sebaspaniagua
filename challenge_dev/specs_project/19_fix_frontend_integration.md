#Context
Refactor frontend api about proposals saved functionality and how obtains saved proposals.

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


1. Frontend service actions.
    1.1 Dashboard page: When a proposal into dashboard is checked, and press "guardar" action button, update behavior, using next steps:
        1.1.1 First do a request to endpoint /proposal-save (POST) for saving the proposal.(check api contract updated)
        1.1.2 Second, do a request to /saved-proposals (POST) endpoint, for save user id and proposal relation.
        1.1.3 If many proposals are checked, loop the behavior for step 1.1.1 and 1.1.2 for every proposal when button "Guardar..." perform click event.
    1.2 Saved-Proposals page: Change behavior using next steps: 
        1.2.1 When load keeps JWT token validacion and then consumes /saved_proposals (GET). (check api contract updated)
        1.2.2 Now the endpoint returns all complete information about saved proposals. 
        1.2.3 When show the modal window about proposal saved, show all real data from response.

2. Run frontend service server



#Results expected

1. front service can do requests for save proposals and saving user_id and proposal_call_id relation.
2. Front service can shows saved proposals and obtains the complete public proposal data information.

-- 

Do not commit changes yet.

--

Tell me first of all again, just to be sure, what do you plan to do?
