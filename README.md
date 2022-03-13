# Secret Scanner

This project can add repository and detect secret key on each repository.

## Design Overview
- secret-scanner serves CRUD for repository (gofiber framework and postgres database)
- trigger a scan, and update the result to "Queued"
- register job to the server, until the server is ready to process the scan request
- there are configurable max worker in <<.env>> to do concurency scan
- on each worker
  - receive job and start download repository (go-getter)
  - do recursively scan to find SECRET_KEY
  - trigger result back to server
- server update findings and scan status

## Usage
- Download and run
```
git clone https://github.com/treewai/secret-scanner`

cd secret-scanner

docker-compose up
```



- Add repository
```
curl -X POST "http://localhost:8080/v1/repositories" -H "Content-type: application/json" -d '{"name": "treewai", "url": "github.com/treewai/secret-scanner"}'
```

- Get repositories
```
curl -X GET "http://localhost:8080/v1/repositories"

[{"repoId":"1","name":"treewai","url":"github.com/treewai/secret-scanner"}]
```

- Get repository
```
curl -X GET "http://localhost:8080/v1/repositories/1"

{"repoId":"1","name":"treewai","url":"github.com/treewai/secret-scanner"}
```

- Update repository
```
curl -X PATCH "http://localhost:8080/v1/repositories/1" -H "Content-type: application/json" -d '{"name": "yutthana", "url": "github.com/treewai/secret-scanner"}'
```

- Delete repository
```
curl -X DELETE "http://localhost:8080/v1/repositories/1"
```

- Trigger a scan
```
curl -X POST "http://localhost:8080/v1/scans/1"
```

- Get scan results
```
curl -X GET "http://localhost:8080/v1/scans"

[{"id":"1","status":"Success","repositoryName":"yutthana","repositoryUrl":"github.com/treewai/secret-scanner","findings":{"findings":[{"locaton":{"path":"README.md","position":{"begin":{"line":15}}},"metadata":{"description":"Define secret key","severity":"HIGN"},"ruleId":"1","type":"sast"},{"locaton":{"path":"README.md","position":{"begin":{"line":15}}},"metadata":{"description":"Define secret key","severity":"HIGN"},"ruleId":"1","type":"sast"}]},"queuedAt":"0001-01-01T00:00:00Z","scanningAt":"2022-03-13T19:27:50.784136Z","finishedAt":"2022-03-13T19:27:52.743708Z"}]
```

**Remark:**

- Only has OpenAPI 3.0 spec (not swagger UI interface)
- Unit tests (not completed)
- Exported function comment (not completed)
