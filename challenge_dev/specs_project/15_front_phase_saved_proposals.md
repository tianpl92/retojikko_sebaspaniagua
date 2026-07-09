#Context
Continue to project plan about coding the frontend layer. Now we are integrate the service for saving proposals.

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

#Designer description

*   Color palette:
Use the follow color palette to the skeleton front application:
    -primary color: rgb(74, 68, 102) hexadecimal is 4A4466.
    -second color: rgb(110, 173, 188) hexadecimal is 6EADBC.
    -auxiliar colors: rgb(159, 203, 173) hexadecimal is 9FCBAD and rgb(241, 247, 212) hexadecimal is F1F7D4

#Actions


1. Integrate dashboard page with /saved-proposals (POST) endpoint, The logic add new features are:
    1.1 Add a floating button, show caption "Guardar...". the button position must be to right side page position. If the page scrolls, the button keeps position.
    1.2 Enable click event to "Guardar..." button. The target is integrate endpoint /saved-proposals (POST).
    1.3 Consumes  /saved-proposals (POST) using the current id user login, and the id proposals selected by checkbox.
    1.4 By default the "Guardar..." button is disabled.
    1.5 When one proposal  or more are checked the button turns to available.
    1.6 Use color palette suggestions to switch disable and available button color.
    1.7 If many proposals are checked, prepare a loop for sending the request save for every proposal.
    1.8 Use animation modal while the proposals are saving. Close it when finished.

4. Run up again the frontend service.

#Results expected

1. The application can save public proposals selected by current user.

-- 

Do not commit changes yet.

--

Tell me first of all again, just to be sure, what do you plan to do?
