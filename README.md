# go-english-test-platform

backend platform for english proficiency testing built with go, featuring asynchronous scoring, file storage, caching, and monitoring.

## tech stack

| category           | technologies                                          |
| ------------------ | ----------------------------------------------------- |
| backend            | go, gin                                               |
| database           | postgresql                                            |
| cache              | redis                                                 |
| messaging          | rabbitmq                                              |
| object storage     | minio                                                 |
| authentication     | jwt                                                   |
| database migration | goose                                                 |
| monitoring         | prometheus, grafana                                   |
| architecture       | clean architecture, repository pattern, service layer |
| processing         | background worker, event-driven architecture          |
| devops             | docker, docker compose                                |

## features

- authentication & authorization
- test management
- section management
- question management
- test attempts
- answer submission
- asynchronous scoring
- results & analytics
- file upload management
- dashboard analytics
- redis caching
- prometheus metrics
- grafana dashboards
