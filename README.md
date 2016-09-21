## RHM is poc / experimental cli for feedhenry written in go. This is not for production use and is being used along with the [Golang intro](https://github.com/fheng/golang-intro)


We are using github.com/urfave/cli to help with the common cli functionality 

structure:

```
|--commands
        |-- get 
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

##Developing

```
    mkdir -p $GOPATH/src/github/feedhenry 
    cd $GOPATH/src/github/feedhenry
    git clone git@github.com/feedhenry/rhm.git
    cd rhm 
    go build .
    //test it works
    ./rhm 
    run the tests go test ./... 
```

More details to come shortly
