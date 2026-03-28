--lookup table

--verified
CREATE TABLE class_days(
    class_day_ID  SERIAL PRIMARY KEY,
    class_day_value VARCHAR(20) NOT NULL,
    class_day_value_thai VARCHAR(20) NOT NULL
);

--verified
CREATE TABLE course_programs(
    course_program_ID  SERIAL PRIMARY KEY,
    course_program_value VARCHAR(20) NOT NULL,
    course_program_value_thai VARCHAR(20) NOT NULL
);

--verified
CREATE TABLE status(
    status_ID  SERIAL PRIMARY KEY,
    status_value VARCHAR(20) NOT NULL
);


CREATE TABLE grades(
    grade_ID  SERIAL PRIMARY KEY,
    grade_value VARCHAR(10) NOT NULL
);

CREATE TABLE semester(
    semester_ID  SERIAL PRIMARY KEY,
    semester_value VARCHAR(10) NOT NULL UNIQUE,
    start_date DATE ,
    end_date DATE ,
    is_active BOOLEAN DEFAULT false

   
);
--lookup table

CREATE TABLE professors(
    professor_ID SERIAL PRIMARY KEY,
    firstname VARCHAR(100) NOT NULL,
    lastname VARCHAR(150) NOT NULL,
    prefix VARCHAR(20),
    firstname_thai VARCHAR(100),
    lastname_thai VARCHAR(100),
     email VARCHAR(50),
    created_at TIMESTAMP 
);

CREATE TABLE accountants(
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(100) NOT NULL,
    lastname VARCHAR(150) NOT NULL,
     email VARCHAR(50),
    created_at TIMESTAMP 
);

CREATE TABLE students(
    student_ID INTEGER PRIMARY KEY,
    firstname VARCHAR(100),
    lastname VARCHAR(100) ,
    prefix VARCHAR(20),
    firstname_thai VARCHAR(100),
    lastname_thai VARCHAR(100),
    email VARCHAR(50),
    phone_number  VARCHAR(20) 
);

--verified
CREATE TABLE transcript_storage(
    transcript_ID SERIAL PRIMARY KEY,
    file_bytes BYTEA,
    file_name VARCHAR(100),
    student_ID INTEGER UNIQUE,
    CONSTRAINT FK_student_ID
        FOREIGN KEY (student_ID)
        REFERENCES students(student_ID)
);

--verified
CREATE TABLE bank_account_storage(
    bank_account_ID SERIAL PRIMARY KEY,
    file_bytes BYTEA,
    file_name VARCHAR(100),
    student_ID INTEGER UNIQUE,
    CONSTRAINT FK_student_ID
        FOREIGN KEY (student_ID)
        REFERENCES students(student_ID)
);

--verified
CREATE TABLE student_card_storage(
    student_card_ID SERIAL PRIMARY KEY,
    file_bytes BYTEA,
    file_name VARCHAR(100),
    student_ID INTEGER UNIQUE,
    CONSTRAINT FK_student_ID
        FOREIGN KEY (student_ID)
        REFERENCES students(student_ID)
);

--verified
CREATE TABLE courses(
    course_ID SERIAL PRIMARY KEY,
    course_code VARCHAR(20) NOT NULL,
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
    work_hour INTEGER,
    created_date TIMESTAMP ,
    deleted_date TIMESTAMP,
    CONSTRAINT FK_professor_ID
        FOREIGN KEY (professor_ID)
        REFERENCES professors(professor_ID),
    CONSTRAINT FK_course_program_ID
        FOREIGN KEY (course_program_ID)
        REFERENCES course_programs(course_program_ID),    
    CONSTRAINT FK_semester_ID
        FOREIGN KEY (semester_ID)
        REFERENCES semester(semester_ID),
    CONSTRAINT FK_class_day_ID
        FOREIGN KEY (class_day_ID)
        REFERENCES class_days(class_day_ID)
);

--verified
CREATE TABLE ta_job_posting(
    id SERIAL PRIMARY KEY,
    professor_ID INTEGER NOT NULL,
    task VARCHAR(200) NOT NULL,
    ta_allocation INTEGER NOT NULL,
    location  VARCHAR(20) NOT NULL,
    status_ID INTEGER NOT NULL,
    course_ID INTEGER ,
    grade_ID INTEGER NOT NULL,
    created_date TIMESTAMP ,
    deleted_date TIMESTAMP,
    CONSTRAINT FK_status_ID
        FOREIGN KEY (status_ID)
        REFERENCES status(status_ID),
    CONSTRAINT FK_professor_ID
        FOREIGN KEY (professor_ID)
        REFERENCES professors(professor_ID),
    CONSTRAINT FK_course_ID
        FOREIGN KEY (course_ID)
        REFERENCES courses(course_ID),
    CONSTRAINT FK_grade_ID
        FOREIGN KEY (grade_ID)
        REFERENCES grades(grade_ID)
);

CREATE TABLE ta_application(
    id SERIAL PRIMARY KEY,
    student_ID INT NOT NULL,
    status_ID INT NOT NULL,
    job_post_ID INT NOT NULL,
    grade VARCHAR(10) NOT NULL,
    purpose VARCHAR(100) NOT NULL,
      reject_reason VARCHAR(200),
    created_date TIMESTAMP,
    deleted_date TIMESTAMP,
    CONSTRAINT FK_student_ID
        FOREIGN KEY (student_ID)
        REFERENCES students(student_ID),
    CONSTRAINT FK_status_ID
        FOREIGN KEY (status_ID)
        REFERENCES status(status_ID),
    CONSTRAINT FK_job_post_ID
        FOREIGN KEY (job_post_ID)
        REFERENCES ta_job_posting(id)
);

CREATE TABLE ta_courses(
    id SERIAL PRIMARY KEY,
    student_ID INT,
    course_ID INT,
    created_date TIMESTAMP,
    deleted_date TIMESTAMP,
    CONSTRAINT FK_student_ID
        FOREIGN KEY (student_ID)
        REFERENCES students(student_ID),
    CONSTRAINT FK_course_ID
        FOREIGN KEY (course_ID)
        REFERENCES courses(course_ID)
);


CREATE TABLE ta_duty_historys(
    id SERIAL PRIMARY KEY,
    date TIMESTAMP,
    course_ID INT,
    student_ID INT,
    CONSTRAINT FK_course_ID
        FOREIGN KEY (course_ID)
        REFERENCES courses(course_ID),
    CONSTRAINT FK_student_ID
        FOREIGN KEY (student_ID)
        REFERENCES students(student_ID)
);

CREATE TABLE discord_channels(
    channel_id VARCHAR(50) UNIQUE,
    role_id VARCHAR(50),
    channel_name VARCHAR(100),
    course_ID INTEGER,
    CONSTRAINT FK_course_ID
        FOREIGN KEY (course_ID)
        REFERENCES courses(course_ID)
);

-- adding for cannot insert duplicate studentid courseid and date
-- ALTER TABLE ta_duty_historys 
-- ADD CONSTRAINT unique_attendance UNIQUE (student_ID, course_ID, date);
CREATE TABLE email_history(
    id SERIAL PRIMARY KEY,
    subject VARCHAR(200),
    body VARCHAR(2000),
    received_name  VARCHAR(100),
    n_received INT,
    status_ID INT,
    created_date TIMESTAMP,
    CONSTRAINT FK_status_ID
        FOREIGN KEY (status_ID)
        REFERENCES status(status_ID)
);

CREATE TABLE holidays(
    id SERIAL PRIMARY KEY,
    holiday_date DATE UNIQUE,
    name_eng VARCHAR(200),
    name_thai VARCHAR(200),
    category VARCHAR(20)
);
-- insert constant values

INSERT INTO professors (firstname, lastname, prefix, firstname_thai, lastname_thai, email) VALUES
('Amnach', 'Khawne', 'ผศ. ดร.', 'อำนาจ', 'ขาวเน', 'amnach.kh@kmitl.ac.th'),
('Charoen', 'Vongchumyen', 'รศ. ดร.', 'เจริญ ', 'วงษ์ชุ่มเย็น', 'charoen.vo@kmitl.ac.th'),
('Kietikul', 'Jearanaitanakij', 'รศ. ดร.', 'เกียรติกูล', 'เจียรนัยธนะกิจ', 'kietikul.je@kmitl.ac.th'),
('Orachat', 'Chitsobhuk', 'รศ. ดร.', 'อรฉัตร', 'จิตต์โสภักตร์', 'orachat.ch@kmitl.ac.th'),
('Rathachai', 'Chawuthai', 'รศ. ดร.', 'รัฐชัย', 'ชาวอุทัย', 'rathachai.ch@kmitl.ac.th'),
('Sakchai', 'Thipchaksurat', 'รศ. ดร.', 'ศักดิ์ชัย', 'ทิพย์จักษุรัตน์', 'sakchai.th@kmitl.ac.th '),
('Surin', 'Kittitornkun', 'รศ. ดร.', 'สุรินทร์', 'กิตติธรกุล', 'surin.ki@kmitl.ac.th'),
('Aranya', 'Walairacht', 'ผศ. ดร.', 'อรัญญา', 'วลัยรัชต์','kwaranya@kmitl.ac.th'),
('Chompoonuch', 'Tengcharoen', 'ผศ. ดร.', 'ชมพูนุท', 'เต็งเจริญ', 'chompoonuch.te@kmitl.ac.th'),
('Chutimet', 'Srinilta',  'ผศ. ดร.', 'ชุติเมษฏ์', 'ศรีนิลทา', 'kschutim@kmitl.ac.th'),
('Pakorn', 'Watanachaturaporn', 'ผศ. ดร.', 'ปกรณ์', 'วัฒนจตุรพร', 'pakorn.wa@KMITL.ac.th'),
('Phongsak', 'Keeratiwintakorn', 'ผศ. ดร.', 'พงษ์ศักดิ์', 'กีรติวินทกร', 'phongsak.ke@kmitl.ac.th'),
('Thanunchai', 'Threepak',  'ผศ. ดร.', 'ธนัญชัย', 'ตรีภาค', 'thanunchai.th@kmitl.ac.th'),
('Akkradach', 'Watcharapupong', 'ผศ.', 'อัครเดช', 'วัชระภูพงษ์','thanunchai.th@kmitl.ac.th'),
('Bundit', 'Pasaya', 'ผศ.', 'บัณฑิต', 'พัสยา', 'bundit.pa@kmitl.ac.th'),
('Kiatnarong', 'Tongprasert', 'ผศ.', 'เกียรติณรงค์', 'ทองประเสริฐ', 'kiatnarong.to@kmitl.ac.th'),
('Sorayut', 'Glomglome', 'ผศ.', 'สรยุทธ', 'กลมกล่อม', 'sorayut.gl@kmitl.ac.th'),
('Thana', 'Hongsuwan', 'ผศ.', 'ธนา', 'หงษ์สุวรรณ', 'khthana@kmitl.ac.th'),
('Jirayu', 'Petchhan', 'ดร.', 'จิรายุ', 'เพชรแหน', 'jirayu.pe@kmitl.ac.th'),
('Parinya', 'Ekparinya', 'ดร.', 'ปริญญา', 'เอกปริญญา', 'parinya.ek@kmitl.ac.th'),
('Kanut', 'Tangtisanon', 'อ.', 'คณัฐ', 'ตังติสาานนท์', 'ktkanut@kmitl.ac.th'),
('Watjanapong', 'Kasemsiri', 'อ.', 'วัจนพงศ์', 'เกษมศิริ', 'watjanapong.ka@kmitl.ac.th'),
('Jirasak', 'Sittigorn', 'อ.', 'จิระศักดิ์', 'สิทธิกร', 'ksjirasa@gmail.com');

-- -- semester date
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
    ('OPEN'),
    ('CLOSE'),
    ('PENDING'),
    ('REJECTED'),
    ('APPROVED'),
    ('SUCCESSFUL'),
    ('FAILED');
 


-- class_day
INSERT INTO class_days (class_day_value,class_day_value_thai) VALUES
    ('Sunday','อาทิตย์'),
    ('Monday','จันทร์'),
    ('Tuesday','อังคาร'),
    ('Wednesday','พุธ'),
    ('Thursday','พฤหัสบดี'),
    ('Friday','ศุกร์'),
    ('Saturday','เสาร์');

-- course_programs
INSERT INTO course_programs (course_program_value,course_program_value_thai) VALUES
    ('General','ทั่วไป'),
    ('International','นานาชาติ'),
    ('Continuing','ต่อเนื่อง');

--grades
INSERT INTO grades (grade_value) VALUES
    ('A'),
    ('B+'),
    ('B'),
    ('C+'),
    ('C'),
    ('D+'),
    ('D');

INSERT INTO students (student_ID,firstname,lastname) VALUES
    (1,'guest','test'),
    (2,'guest2','test2')