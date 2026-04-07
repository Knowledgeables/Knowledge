INSERT OR IGNORE INTO users (username, email, password_hash)
VALUES
('admin', 'admin@dev.local', '$2a$12$Zd7hfqU8jdcyK0Q8oDxNOONbiU4GxPkY2eRR9Pw0hZQyhAZaHdkmG'),
('tester', 'tester@dev.local', '$2a$12$sHphNz2knALIBeijiWoCqe3erFOhSY/Ke8kUt/09g3SkRySl0kTp2');

INSERT OR IGNORE INTO pages (title, url, language, content)
VALUES
('Docker Intro', 'https://en.wikipedia.org/wiki/Docker_(software)', 'en', 'Introduction to Docker'),
('Go Basics', 'https://en.wikipedia.org/wiki/Go_(programming_language)', 'en', 'Intro to Go'),

-- Test data
('PostgreSQL Indexing', 'https://en.wikipedia.org/wiki/PostgreSQL', 'en', 'Improve query performance with indexes https://en.wikipedia.org/wiki/PostgreSQL'),
('Redis Caching', 'https://en.wikipedia.org/wiki/Redis', 'en', 'In-memory caching strategies https://en.wikipedia.org/wiki/Redis'),
('Microservices Architecture', 'https://en.wikipedia.org/wiki/Microservices', 'en', 'Design scalable distributed systems https://en.wikipedia.org/wiki/Microservices'),
('GraphQL API Design', 'https://en.wikipedia.org/wiki/GraphQL', 'en', 'Flexible APIs with GraphQL https://en.wikipedia.org/wiki/GraphQL'),
('CI/CD Pipelines', 'https://en.wikipedia.org/wiki/CI/CD', 'en', 'Automate build, test, and deployment https://en.wikipedia.org/wiki/CI/CD'),
('TypeScript Basics', 'https://en.wikipedia.org/wiki/TypeScript', 'en', 'Typed JavaScript for safer code https://en.wikipedia.org/wiki/TypeScript'),
('Node.js Streams', 'https://en.wikipedia.org/wiki/Node.js', 'en', 'Efficient data processing in Node.js https://en.wikipedia.org/wiki/Node.js'),
('Linux Networking', 'https://en.wikipedia.org/wiki/Computer_network', 'en', 'TCP/IP, ports, and routing basics https://en.wikipedia.org/wiki/Computer_network'),
('Web Security Essentials', 'https://en.wikipedia.org/wiki/Web_security', 'en', 'XSS, CSRF, and authentication https://en.wikipedia.org/wiki/Web_security'),
('OAuth2 Explained', 'https://en.wikipedia.org/wiki/OAuth', 'en', 'Secure authorization flows https://en.wikipedia.org/wiki/OAuth'),
('JWT Authentication', 'https://en.wikipedia.org/wiki/JSON_Web_Token', 'en', 'Stateless authentication with tokens https://en.wikipedia.org/wiki/JSON_Web_Token'),
('Clean Code Principles', 'https://en.wikipedia.org/wiki/Software_maintenance', 'en', 'Write maintainable software https://en.wikipedia.org/wiki/Software_maintenance'),
('System Design Basics', 'https://en.wikipedia.org/wiki/Systems_design', 'en', 'Scalability, load balancing, caching https://en.wikipedia.org/wiki/Systems_design'),
('Distributed Systems Deep Dive', 'https://en.wikipedia.org/wiki/Distributed_computing', 'en',
'Distributed systems are complex by nature. They require careful handling of consistency, 
availability, and partition tolerance (CAP theorem https://en.wikipedia.org/wiki/CAP_theorem). In real-world systems, engineers 
must deal with issues like network latency, partial failures, retries, idempotency, and 
eventual consistency https://en.wikipedia.org/wiki/Eventual_consistency. Techniques such as replication https://en.wikipedia.org/wiki/Replication_(computing), 
sharding https://en.wikipedia.org/wiki/Shard_(database_architecture), leader election, and 
consensus algorithms (like Raft https://en.wikipedia.org/wiki/Raft_(algorithm) or Paxos https://en.wikipedia.org/wiki/Paxos_(computer_science)) are essential. 
Observability, monitoring, and tracing also play a crucial role in debugging distributed environments.');