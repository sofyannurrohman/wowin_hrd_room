ALTER TABLE answers ADD CONSTRAINT unique_answer_per_question UNIQUE (participant_id, question_id);
