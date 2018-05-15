## Deprecation Notice
This repository has been deprecated and is not being maintained. It should not be used. If you have any questions, please get in touch with the collaborators.

## RHM is poc / experimental cli for feedhenry written in go. This is not for production use and is being used along with the [Golang intro](https://github.com/fheng/golang-intro)


We are using github.com/urfave/cli to help with the common cli functionality 

Propsed structure:

```
|--commands (main location for all of the domain logic)
        |-- get (read and list)
        |   |-- projects
        |   |-- apps
        |-- delete 
        |   |-- projects
        |   |-- apps
        |-- update 
        |   |-- projects
        |   |-- apps
        |-- use 
        |   |-- projects
        |   |-- apps
        |-- create 
        |   |-- projects
        |   |-- apps
       login
|--request //request helpers
|--storage //storing data to disk
|--ui //user interface (cli only)
|--vendor //dependencies
| 
main //where it all starts         
```    

New commands should go under commands.

## Developing

Ensure you are running go 1.7 or later.

Install glide first [it's on GitHub](https://github.com/Masterminds/glide)
```
go get github.com/Masterminds/glide
```

```
    go get github.com/feedhenry/rhm
    cd $GOPATH/src/github.com/feedhenry/rhm
    glide install
    go install
    # test it works
    rhm 
    # run the tests 
    go test ./... 
```

## Building 

You can run ```make ci``` to run a full build or you can run ```make test``` just to run the tests.

## Releasing

There is a release target in the make file. Before releasing ensure you update the version in the Makefile then run 
```
make release
```
