# RSS Reader

- [x] fetch & parse rss
- [x] save rss to db
- [x] add cron job (worker) to fetch rss
- [x] fetch rss from multiple sources 
```sql
INSERT INTO sources (id,name,url,created_at,updated_at,deleted_at) VALUES (
1,'espn nba','https://www.espn.com/espn/rss/nba/news','2020-04-26 06:08:10','2020-04-26 06:08:10',NULL);
INSERT INTO sources (id,name,url,created_at,updated_at,deleted_at) VALUES (
2,'hn newest golang','https://hnrss.org/newest?q=golang','2020-04-26 06:08:32','2020-04-26 06:08:32',NULL);
INSERT INTO sources (id,name,url,created_at,updated_at,deleted_at) VALUES (
3,'hn job','https://hnrss.org/jobs','2020-04-26 06:08:53','2020-04-26 06:08:53',NULL);
INSERT INTO sources (id,name,url,created_at,updated_at,deleted_at) VALUES (
4,'hn polls','https://hnrss.org/polls','2020-04-26 06:08:53','2020-04-26 06:08:53',NULL);
```