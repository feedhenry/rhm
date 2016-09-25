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

```
    go get github.com/feedhenry/rhm
    cd $GOPATH/src/github.com/feedhenry/rhm
    go build .
    # test it works
    ./rhm 
    # run the tests 
    go test ./... 
```

More details to come shortly
