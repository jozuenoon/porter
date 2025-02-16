# Porter API
Porter - the ports service.

## Tools
* **just:** To execute commands defined in a `.justfile`, you need to have the `just` tool installed on your system [Just website](https://just.systems/).
* **golangci-lint:** For linter you need to have the `golangci-lint` tool installed on your system.
* **moq:** For mocking you need to have the `moq` tool installed on your system. [Moq website](https://github.com/matryer/moq).

## NOTE
The `.justfile` should contain commands mentioned below for convenience.

## Building

Run command from `.justfile` using `just` command followed by the target name.

```bash
just build
```

## Running
Run created docker image.

```bash
docker run -it -p 8080:8080 --memory=200m q88/porter:latest
```
## API

### Upload file

* **Path:** `/api/v1/ports`
* **Method:** `POST`
* **Content-Type:** `application/json`
* **Request Body:** JSON port data.

```bash
curl -T <file>.json -X POST -H "Content-Type: application/json" http://localhost:8080/api/v1/ports
```

### Get port by UN/LOCODE

* **Path:** `/api/v1/port`
* **Method:** `GET`
* **Query Parameter:** `unloc` - Port UN/LOCODE (e.g., `AEAJM`)

```bash
curl -X GET http://localhost:8080/api/v1/port?unloc=AEAJM
```

## Testing

1. Preapare input file, eg. 1GB json file with ports, it can contain duplicate ports.
2. Build docker image with `just build`.
3. Run docker image with `docker run -it -p 8080:8080 --memory=200m q88/porter:latest`.
4. Identify container id.
5. Run stats `docker stats <container_id>` to see container stats.
6. Upload the file and see memory usage in docker stats.
