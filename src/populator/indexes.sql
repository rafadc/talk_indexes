CREATE INDEX people_single_index_name_IDX USING BTREE ON indexes.people_single_index (name);

CREATE INDEX people_multi_column_index_name_company_IDX USING BTREE ON indexes.people_multi_column_index (name, company);

CREATE INDEX people_multi_column_index_name_happy_IDX USING BTREE ON indexes.people_multi_column_index (name, happy);
CREATE INDEX people_multi_column_index_happy_name_IDX USING BTREE ON indexes.people_multi_column_index (happy, name);

CREATE INDEX people_multi_column_index_name_date_of_birth_IDX USING BTREE ON indexes.people_multi_column_index (name, date_of_birth);
CREATE INDEX people_multi_column_index_date_of_birth_name_IDX USING BTREE ON indexes.people_multi_column_index (date_of_birth, name);

CREATE FULLTEXT INDEX people_full_text_search_name_IDX ON indexes.people_full_text_search (name) WITH PARSER ngram;
