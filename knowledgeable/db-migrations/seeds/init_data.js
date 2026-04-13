export async function seed(knex) {
  // ryd først (rigtig rækkefølge pga. FK hvis de findes)
  await knex('pages').del();
  await knex('users').del();

  // USERS
  await knex('users')
    .insert([
      {
        username: 'admin',
        email: 'admin@dev.local',
        password_hash: '$2a$12$8fhiDpxQAlrTsZfQYLiKUeELzXABImIynuZSx3nwyNfIOofEg1W3i'       
      },
      {
        username: 'tester',
        email: 'tester@dev.local',
        password_hash: '$2a$12$sHphNz2knALIBeijiWoCqe3erFOhSY/Ke8kUt/09g3SkRySl0kTp2'
      }
    ])
    .onConflict(['email'])
    .ignore();

  // PAGES
  await knex('pages')
    .insert([
      {
        title: 'Docker Intro',
        url: 'https://en.wikipedia.org/wiki/Docker_(software)',
        language: 'en',
        content: 'Introduction to Docker'
      },
      {
        title: 'Go Basics',
        url: 'https://en.wikipedia.org/wiki/Go_(programming_language)',
        language: 'en',
        content: 'Intro to Go'
      },
      {
        title: 'PostgreSQL Indexing',
        url: 'https://en.wikipedia.org/wiki/PostgreSQL',
        language: 'en',
        content: 'Improve query performance with indexes https://en.wikipedia.org/wiki/PostgreSQL'
      },
      {
        title: 'Redis Caching',
        url: 'https://en.wikipedia.org/wiki/Redis',
        language: 'en',
        content: 'In-memory caching strategies https://en.wikipedia.org/wiki/Redis'
      },
      {
        title: 'Microservices Architecture',
        url: 'https://en.wikipedia.org/wiki/Microservices',
        language: 'en',
        content: 'Design scalable distributed systems https://en.wikipedia.org/wiki/Microservices'
      },
      {
        title: 'GraphQL API Design',
        url: 'https://en.wikipedia.org/wiki/GraphQL',
        language: 'en',
        content: 'Flexible APIs with GraphQL https://en.wikipedia.org/wiki/GraphQL'
      },
      {
        title: 'CI/CD Pipelines',
        url: 'https://en.wikipedia.org/wiki/CI/CD',
        language: 'en',
        content: 'Automate build, test, and deployment https://en.wikipedia.org/wiki/CI/CD'
      },
      {
        title: 'TypeScript Basics',
        url: 'https://en.wikipedia.org/wiki/TypeScript',
        language: 'en',
        content: 'Typed JavaScript for safer code https://en.wikipedia.org/wiki/TypeScript'
      },
      {
        title: 'Node.js Streams',
        url: 'https://en.wikipedia.org/wiki/Node.js',
        language: 'en',
        content: 'Efficient data processing in Node.js https://en.wikipedia.org/wiki/Node.js'
      },
      {
        title: 'Linux Networking',
        url: 'https://en.wikipedia.org/wiki/Computer_network',
        language: 'en',
        content: 'TCP/IP, ports, and routing basics https://en.wikipedia.org/wiki/Computer_network'
      },
      {
        title: 'Web Security Essentials',
        url: 'https://en.wikipedia.org/wiki/Web_security',
        language: 'en',
        content: 'XSS, CSRF, and authentication https://en.wikipedia.org/wiki/Web_security'
      },
      {
        title: 'OAuth2 Explained',
        url: 'https://en.wikipedia.org/wiki/OAuth',
        language: 'en',
        content: 'Secure authorization flows https://en.wikipedia.org/wiki/OAuth'
      },
      {
        title: 'JWT Authentication',
        url: 'https://en.wikipedia.org/wiki/JSON_Web_Token',
        language: 'en',
        content: 'Stateless authentication with tokens https://en.wikipedia.org/wiki/JSON_Web_Token'
      },
      {
        title: 'Clean Code Principles',
        url: 'https://en.wikipedia.org/wiki/Software_maintenance',
        language: 'en',
        content: 'Write maintainable software https://en.wikipedia.org/wiki/Software_maintenance'
      },
      {
        title: 'System Design Basics',
        url: 'https://en.wikipedia.org/wiki/Systems_design',
        language: 'en',
        content: 'Scalability, load balancing, caching https://en.wikipedia.org/wiki/Systems_design'
      },
      {
        title: 'Distributed Systems Deep Dive',
        url: 'https://en.wikipedia.org/wiki/Distributed_computing',
        language: 'en',
        content: `Distributed systems are complex by nature. They require careful handling of consistency, 
availability, and partition tolerance (CAP theorem https://en.wikipedia.org/wiki/CAP_theorem). In real-world systems, engineers 
must deal with issues like network latency, partial failures, retries, idempotency, and 
eventual consistency https://en.wikipedia.org/wiki/Eventual_consistency. Techniques such as replication https://en.wikipedia.org/wiki/Replication_(computing), 
sharding https://en.wikipedia.org/wiki/Shard_(database_architecture), leader election, and 
consensus algorithms (like Raft https://en.wikipedia.org/wiki/Raft_(algorithm) or Paxos https://en.wikipedia.org/wiki/Paxos_(computer_science)) are essential. 
Observability, monitoring, and tracing also play a crucial role in debugging distributed environments.`
      }
    ])
    .onConflict(['url'])
    .ignore();
}