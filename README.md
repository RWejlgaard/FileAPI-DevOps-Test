# DevOps Coding Test

## How to Build/Run
I'm utilizing docker-compose for hosting my small file server API

In order to run execute the following command in the root folder of the project
```bash
docker-compose up
```

This will build and spawn a docker container containing the go application
and expose the containers port `8080` to a random port on our host

In order to see which port is bound we can execute `docker-compose ps`


## How to Test

For testing we can try to execute the different endpoints
```bash
curl localhost:32768/files/

> <pre>
> <a href=".idea/">.idea/</a>
> <a href="README.md">README.md</a>
> <a href="docker-compose.yml">docker-compose.yml</a>
> <a href="main.go">main.go</a>
> </pre>
```

We can upload a file like this:
```bash
curl localhost:32768 -X POST -F uploadFile=path/to/some/file

> h1r3mes00n.pdf
```

And we can delete a file like this:
```bash
curl localhost:32768/delete/h1r3mes00n.pdf

> h1r3mes00n.pdf has been deleted
```