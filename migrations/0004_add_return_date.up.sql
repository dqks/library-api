ALTER TABLE loans 
ADD COLUMN return_date DATE NULL CHECK (return_date >= loan_date);