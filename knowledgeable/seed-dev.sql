INSERT OR IGNORE INTO users (username, email, password_hash)
VALUES
('admin', 'admin@dev.local', '$2a$12$Zd7hfqU8jdcyK0Q8oDxNOONbiU4GxPkY2eRR9Pw0hZQyhAZaHdkmG'),
('tester', 'tester@dev.local', '$2a$12$sHphNz2knALIBeijiWoCqe3erFOhSY/Ke8kUt/09g3SkRySl0kTp2');

INSERT OR IGNORE INTO pages (title, url, language, content)
VALUES
('Docker Intro', '/docker', 'en', 'Introduction to Docker'),
('Go Basics', '/go', 'en', 'Intro to Go'),
('Docker Intro DK', '/docker-da', 'da', 'Introduktion til Docker');