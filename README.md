# Proglog

A toy project demonstrating how to serve an in memory log.

## Usage

Launch the server by doing:

```
$ make
```

Try adding some records:

```
$ curl -X POST localhost:8080 -d '{"record":{"value":"ujfespkbitg6re3w"}}'
{"offset":0}
$ curl -X POST localhost:8080 -d '{"record":{"value":"ujfespkb"}}'        
{"offset":1}
```

Read them:

```
$ curl -X GET localhost:8080 -d '{"offset":0}'                            
{"record":{"value":"ujfespkbitg6re3w","offset":0}}
$ curl -X GET localhost:8080 -d '{"offset":1}'                    
{"record":{"value":"ujfespkb","offset":1}}
```

## Testing

It should be enough to run:

```
$ make test
```

In order to check coverage:

```
$ make cover
```