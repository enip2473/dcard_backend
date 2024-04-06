RPS Record

4/6 00:00 Environment: local RPS: 500
4/6 22:00 Environment: GCP RPS: 500 -> database is the bottleneck
4/6 23:00 Environment: GCP RPS: 1000 -> database reach limit of 100 RPS, trying not accessing database
4/7 02:00 Environment: GCP RPS: 2000 Redis no free tier
4/7 02:45 Use Redis, PostgreSQL, and APP all in local with docker-compose, RPS: 3000