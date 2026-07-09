#Task Domain
I want generate schema sql again, but, this specific task has a target, when finish the task, create a new skill called "quality_check_database_query_skill".

#Role
Expert database engineer using PostgreSQL.

#Context
I want to generate three scripts schema generation using diferent models, and then compare the three sql generated to score the quality generation about the next rubric.

  
#Quality Rubric
Category	        Max	Score    What to check
Structure	        30	        Required tables, FK relationships, no hallucinated tables
Naming	            15	        snake_case, id/{table}_id, created_at/updated_at/deleted_at
Integrity	        20	        PK/FK, NOT NULL, UNIQUE, soft delete, CASCADE/RESTRICT
Comments	        15	        Tables and non-obvious decisions documented
Query feasibility	10	        Key queries supported, indexes on FK + hot paths
Spec adherence	    10	        Followed spec, no invented features

#Actions

1.  Read file "database/database_spec.md" for load sql instructions generation.
2.  Generate three scripts sql from the instructions using "gpt-5.4mini", "gpt-5.5" and "openai/gpt-oss-20b:free" models.
3.  Output results file name: 
    3.1 With gpt-5.5 must be named "public_calls_database_gpt55.sql". 
    3.2 With openai/gpt-oss-20b:free must be named "public_calls_database_gpt_oss.sql".
    3.3 With gpt-5.4mini must be named "public_calls_database_gpt54mini.sql"
4. Write and storage files into "database/compare_rubric/" folder.
5. Analyze every sql generated and score them about quality rubric using a fourth model (1 model for judge and the other models to generate content).
6. Sum total score per model-file. total score = (category score + naming score + integrity score + comments score + query feasibility score + spec adherence score)  
7. Generate a matrix 3x1 to compare the total results score per model-file, using this format:

MODEL        | SQL_FILE                 | SCORE
===================================================================
model1-name  | sql_file_generated_name  | total points score number
model2-name  | sql_file_generated_name  | total points score number
model3-name  | sql_file_geneerated_name | total points score number

    7.1 For judge quality rubric use google/gemma-4-31b-it:free  model.

8. Export quality rubric results matrix into "database/compare_rubric/analyse.md" file.

9. Replace content file "database/public_calls_database.sql" with the content from file generated that has max total score > 80 points.

10. I there is a file generation with total score > 80 points, so confirm results (run git add . command) and commit using comment "feat (reto jikko) challenge dev - database schema creation". keep current branch "feature/challenge_dev_public_proposals" and finally run push command (git push -u origin feature/challende_dev_public_proposals). If all files generated has < 80 point total. Dońt do anything, and only shows message "Dont sync changes because the quality check dont pass".

#Target results

* Only generate scripts sql schema generation.
* Don't run query scripts.
* Score the best quality sql schema generation.
* Commit and push changes.


#Creation skill
Create new skill for all this into hermes called "quality_check_database_query_skill". When the skill perform again, please ask:
        -The three model names generation.
        -The fourth model for judge.
        -The path or file name origin database md instructions schema.
        -The branch name for commit and push results.


 -----

Tell me first of all, just to be sure, what do you plan to do?
