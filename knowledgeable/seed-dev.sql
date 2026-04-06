INSERT OR IGNORE INTO users (username, email, password_hash)
VALUES
('admin', 'admin@dev.local', '$2a$12$Zd7hfqU8jdcyK0Q8oDxNOONbiU4GxPkY2eRR9Pw0hZQyhAZaHdkmG'),
('tester', 'tester@dev.local', '$2a$12$sHphNz2knALIBeijiWoCqe3erFOhSY/Ke8kUt/09g3SkRySl0kTp2');

INSERT OR IGNORE INTO pages (title, url, language, content)
VALUES
('Docker Intro', '/docker', 'en', 'Introduction to Docker'),
('Go Basics', '/go', 'en', 'Intro to Go'),

-- Test data
('Spring Boot Guide', '/spring-boot', 'en', 'Build REST APIs with Spring Boot'),
('React Fundamentals', '/react', 'en', 'Learn components, hooks, and state'),
('Kubernetes 101', '/kubernetes', 'en', 'Container orchestration basics'),
('PostgreSQL Indexing', '/postgres-indexing', 'en', 'Improve query performance with indexes'),
('Redis Caching', '/redis', 'en', 'In-memory caching strategies'),
('Microservices Architecture', '/microservices', 'en', 'Design scalable distributed systems'),
('GraphQL API Design', '/graphql', 'en', 'Flexible APIs with GraphQL'),
('CI/CD Pipelines', '/cicd', 'en', 'Automate build, test, and deployment'),
('TypeScript Basics', '/typescript', 'en', 'Typed JavaScript for safer code'),
('Node.js Streams', '/node-streams', 'en', 'Efficient data processing in Node.js'),
('Linux Networking', '/linux-networking', 'en', 'TCP/IP, ports, and routing basics'),
('Web Security Essentials', '/web-security', 'en', 'XSS, CSRF, and authentication'),
('OAuth2 Explained', '/oauth2', 'en', 'Secure authorization flows'),
('JWT Authentication', '/jwt', 'en', 'Stateless authentication with tokens'),
('Clean Code Principles', '/clean-code', 'en', 'Write maintainable software'),
('System Design Basics', '/system-design', 'en', 'Scalability, load balancing, caching'),
('Distributed Systems Deep Dive', '/distributed-systems', 'en', 
'Distributed systems are complex by nature. They require careful handling of consistency, 
availability, and partition tolerance (CAP theorem). In real-world systems, engineers 
must deal with issues like network latency, partial failures, retries, idempotency, and 
eventual consistency. Techniques such as replication, sharding, leader election, and 
consensus algorithms (like Raft or Paxos) are essential. Observability, monitoring, 
and tracing also play a crucial role in debugging distributed environments.');