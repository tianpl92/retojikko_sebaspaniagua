#Context
Schema SQL generation using tech instructions from file specs_project/2_database_spec.md

#Role
Database engineer PostgreSQL

#Action

1. Use "quality_check_database_query_skill" skill hermes for this task.
2. path or file name origin database md instructions schema is "specs_project/2_database_spec.md".
3. Use the following models for generation schema process:
    3.1 gpt-5.5 from openrouter
    3.1 deepseek/deepseek-v4-pro from openrouter
    3.3 GLM 4.7 Flash from openrouter
4. Use MiniMax M2.5 model from openrouter like judge and quality score.
5. Update database/compare_rubric/analysis.md file.
6. Compare and replace database/public_calls_database.sql with the best schema generation code.
7. Change model to deepseek/deepseek-v4-flash model by default into hermes.
8. If the score pass, commit and push changes.


#Delivery
1.  The best schema generated script write into database/public_calls_database.sql
2.  If don''t pass, show the message info about it. If pass Sync changes in github repo.




----


Tell me first of all, just to be sure, what do you plan to do?



