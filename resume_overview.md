Cisco thousandEyes
- Developed a scalable integration dashboard using React, TypeScript, and Tailwind CSS to visualize multi-source telemetry, enhancing team efficiency
and reducing analysis time by 30%.

Situation: During my internship on the Integrations team at Cisco ThousandEyes, the team owned the infrastructure used by multiple engineering teams across the org. Each team has its own telemetry sources, AWS endpoints, and health metrics which were scattered across different tools. This made debugging especially hard to handle by the integrations team, and the engineers had hard time understand what went wrong where.
Task: I was given the task to build a centralized dashboard that could visualize multi-source telemetry and infrastructure metrics in one place. Goal was to improve especially debugging time, and be flexible such that if more integrations were to be added, the system would scale
Action: Started with a single team which had only a few resources deployed. Shadowed senior engineer in the initial meeting requirement gathering meetings to understand their design and service, and later internal discussion to understand what metrics would be useful.
Gathering the AWS endpoints and authentication was the pain point. The access control across the accounts were initial hiccups, but once we figured it out, after that the flow became similar across the teams, it was just adding the additional service as required.
Result: The dashboard became the single source of truth for the team, reducing jumping across tools.


- Built and optimized serverless backend services and GraphQL APIs with Node.js, Express, and AWS Lambda, improving data ingestion speed by 40%
and reducing manual workflows by 70%.

This is also in sync with the dashboard application. When the user logins into the application, instead of fetching all the services and getting overwhelmed
with data, I used graphql apis to fetch and process only the data that was necessary.


- Automated webhook processing and CI/CD pipelines with Docker, Kubernetes, and GitHub Actions, achieving 99.9% uptime and reducing
synchronization errors


Electric Hydrogen
- Developed a Python-based telemetry and monitoring system using AWS Lambda, CDK, and Slack Webhooks, improving real-time visibility and
reducing cloud monitoring costs by 30%.


- Enhanced ETL reliability by 40% through DLQ-enabled Lambda → S3 ingestion workflows and automated alerting mechanisms for proactive data
management.

- Optimized database performance and analytics speed by 35% via Python automation, PostgreSQL schema tuning, and modular microservice design
for scalable backend operations.


- Streamlined CI/CD and infrastructure provisioning using Terraform, Boto3, and GitHub Actions, cutting deployment lead times by 25% and ensuring
consistent, reproducible environments

JPMorgan Chase
- Engineered cloud migration of legacy financial services to AWS using Terraform and CDK, reducing provisioning time by 20% and ensuring consistent
multi-environment deployments.

- Boosted service reliability by 35% through Lambda/ECS optimization with ALB health checks, Step Functions, and DLQ-based recovery for high-load
transaction systems.

- Developed CI/CD pipelines with GitHub Actions and SonarQube, cutting failed deployments by 30% while improving code quality and release stability.

- Refactored monolithic applications into Node.js microservices and containerized them via Docker and ECR, enhancing scalability and accelerating
feature delivery

Intern JPMorgan Chase & Co
- Developed a financial reporting platform using GoLang and AngularJS, automating onboarding analysis workflows and reducing manual reporting
effort by 30%.

- Automated containerized deployments on AWS EKS with Docker and Jenkins CI/CD pipelines, accelerating release cycles by 40% and improving
environment consistency.

- Enhanced system reliability by building Go-based RESTful APIs and implementing unit/integration testing with Jest and Mocha, improving overall
application stability.









Serverless Data Processing Pipeline
- Built a highly available ETL pipeline using AWS Lambda, S3, and DynamoDB, processing streaming data from multiple sources with automated
scalability and ensuring 99.9% data reliability.
- Automated infrastructure deployment with Terraform and AWS CDK and designed CloudWatch dashboards to monitor metrics and detect anomalies
in real-time, reducing manual intervention.


Distributed Cache with Consistent Hashing
- Developed a fault-tolerant distributed key-value store in Go with consistent hashing and replication, achieving sub-10ms read/write latency under
benchmark tests.
- Implemented write-ahead logging and LRU eviction policies, containerized nodes with Docker, and orchestrated clusters using Kubernetes, optimizing
memory usage and ensuring high durability.



