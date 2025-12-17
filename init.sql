CREATE TABLE professors(
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(100) NOT NULL,
    lastname VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accountants(
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(100) NOT NULL,
    lastname VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

--verified
CREATE TABLE courses(
    id SERIAL PRIMARY KEY,
    course_ID VARCHAR(20) NOT NULL,
    course_name VARCHAR(150) NOT NULL,
    professor_ID INTEGER NOT NULL,
    course_program_ID INTEGER NOT NULL,
    course_program VARCHAR(30) NOT NULL,
    sec VARCHAR(20) NOT NULL,
    semester_ID INTEGER NOT NULL,
    semester VARCHAR(10) NOT NULL,
    class_day_ID INTEGER NOT NULL,
    class_day VARCHAR(10) NOT NULL,
    class_start TIME NOT NULL,
    class_end TIME NOT NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_date TIMESTAMP
);

--verified
CREATE TABLE ta_job_posting(
    id SERIAL PRIMARY KEY,
    professor_ID INTEGER NOT NULL,
    task VARCHAR(200) NOT NULL,
    status_ID INTEGER NOT NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_date TIMESTAMP
);

-- CREATE TABLE ta_application(

-- )

--lookup table

--verified
CREATE TABLE class_days(
    class_day_ID  SERIAL PRIMARY KEY,
    class_day_value VARCHAR(20) NOT NULL
);

--verified
CREATE TABLE course_programs(
    course_program_ID  SERIAL PRIMARY KEY,
    course_program_value VARCHAR(20) NOT NULL
);

--verified
CREATE TABLE status(
    status_ID  SERIAL PRIMARY KEY,
    status_value VARCHAR(20) NOT NULL
);

CREATE TABLE semester(
    semester_ID  SERIAL PRIMARY KEY,
    semester_value VARCHAR(10) NOT NULL,
    start_date DATE,
    end_date DATE
);

-- insert constant values

INSERT INTO professors (firstname, lastname) VALUES
('Amnach', 'Khawne'),
('Charoen', 'Vongchumyen'),
('Kietikul', 'Jearanaitanakij'),
('Orachat', 'Chitsobhuk'),
('Rathachai', 'Chawuthai'),
('Sakchai', 'Thipchaksurat'),
('Surin', 'Kittitornkun'),
('Aranya', 'Walairacht'),
('Chompoonuch', 'Tengcharoen'),
('Chutimet', 'Srinilta'),
('Pakorn', 'Watanachaturaporn'),
('Phongsak', 'Keeratiwintakorn'),
('Thanunchai', 'Threepak'),
('Akkradach', 'Watcharapupong'),
('Bundit', 'Pasaya'),
('Kiatnarong', 'Tongprasert'),
('Sorayut', 'Glomglome'),
('Thana', 'Hongsuwan'),
('Jirayu', 'Petchhan'),
('Parinya', 'Ekparinya'),
('Kanut', 'Tangtisanon'),
('Watjanapong', 'Kasemsiri'),
('Jirasak', 'Sittigorn'),
('Pithawat', 'Kitmongkolchai');

-- semester date
WITH all_start_dates AS (
    -- Generate the starting dates for all semesters (always the 1st of the month)
    SELECT (date_trunc('year', '2025-01-01'::date) + (n || ' months')::interval)::date AS start_date FROM generate_series(6, 6 + 5*12, 12) AS t(n) -- Start Semester 1 (July 1st)
    UNION ALL
    SELECT (date_trunc('year', '2025-01-01'::date) + (n || ' months')::interval)::date AS start_date FROM generate_series(10, 10 + 5*12, 12) AS t(n) -- Start Semester 2 (Nov 1st)
    UNION ALL
    SELECT (date_trunc('year', '2025-01-01'::date) + (n || ' months')::interval)::date AS start_date FROM generate_series(3, 3 + 5*12, 12) AS t(n) -- Start Semester 3 (Apr 1st)
),
semester_calc AS (
    SELECT
        asd.start_date,
        EXTRACT(MONTH FROM asd.start_date) AS start_month,
        EXTRACT(YEAR FROM asd.start_date) + 543 AS base_be_year,
        
        -- 1. Determine Semester Number
        CASE
            WHEN EXTRACT(MONTH FROM asd.start_date) BETWEEN 7 AND 10 THEN 1
            WHEN EXTRACT(MONTH FROM asd.start_date) IN (11, 12, 1, 2, 3) THEN 2
            WHEN EXTRACT(MONTH FROM asd.start_date) BETWEEN 4 AND 6 THEN 3
            ELSE 0
        END AS semester_num,
        
        -- 2. Calculate End Date (Now ending on the last day of the respective end month)
        CASE
            -- Semester 1 (Jul-Oct): Ends Oct 31st (4 months after July 1st)
            WHEN EXTRACT(MONTH FROM asd.start_date) BETWEEN 7 AND 10 THEN (asd.start_date + '4 months'::interval - '1 day'::interval)::date 
            
            -- Semester 2 (Nov-Mar): Ends Mar 31st (5 months after Nov 1st) or (3 months after Jan 1st)
            WHEN EXTRACT(MONTH FROM asd.start_date) BETWEEN 11 AND 12 THEN (asd.start_date + '5 months'::interval - '1 day'::interval)::date
            WHEN EXTRACT(MONTH FROM asd.start_date) BETWEEN 1 AND 3 THEN (asd.start_date + '3 months'::interval - '1 day'::interval)::date
            
            -- Semester 3 (Apr-Jun): Ends Jun 30th (3 months after Apr 1st)
            WHEN EXTRACT(MONTH FROM asd.start_date) BETWEEN 4 AND 6 THEN (asd.start_date + '3 months'::interval - '1 day'::interval)::date
            ELSE NULL
        END AS end_date
    FROM all_start_dates asd
),
final_semesters AS (
    SELECT
        semester_num,
        sc.start_date,
        sc.end_date,
        -- 3. Conditional BE Year Adjustment (BE Year logic remains correct)
        CASE
            WHEN start_month BETWEEN 7 AND 12 THEN base_be_year    -- Sem 1 & Sem 2 start (Nov-Dec)
            WHEN start_month BETWEEN 1 AND 6 THEN base_be_year - 1  -- Sem 2 start (Jan-Mar) & Sem 3 (Apr-Jun)
            ELSE base_be_year
        END AS final_be_year
    FROM semester_calc sc
)
-- 4. Final INSERT statement
INSERT INTO semester (semester_value, start_date, end_date)
SELECT
    final_semesters.semester_num::TEXT || '/' || final_semesters.final_be_year::TEXT AS semester_value,
    final_semesters.start_date,
    final_semesters.end_date
FROM final_semesters
WHERE final_semesters.semester_num > 0
AND final_semesters.start_date < '2030-07-01'::date;

-- status
INSERT INTO status (status_value) VALUES
    ('Active'),
    ('InActive'),
    ('Pending');


-- class_day
INSERT INTO class_days (class_day_value) VALUES
    ('Sunday'),
    ('Monday'),
    ('Tuesday'),
    ('Wednesday'),
    ('Thursday'),
    ('Friday'),
    ('Saturday');

-- course_programs
INSERT INTO course_programs (course_program_value) VALUES
    ('General'),
    ('International');
