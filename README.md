
# go-new

This repository has templates to golang projects where should be used with the command `gonew`

## Requirements

You need the `gonew` installed to use these templates. You can do it using the command bellow or following the instructions on the site <https://go.dev/blog/gonew>

``` bash
go install golang.org/x/tools/cmd/gonew@latest
```

## How to use

To create a project using the template you should use the bellow command.

``` bash
gonew github.com/leandrosa/gonew/<template_folder> <module_destine>
```

Example:
To create a project using the helloWorld template, you can use the bellow example:

``` bash
gonew github.com/leandrosa/gonew/helloWorld mytest.com/myhelloWorld
```
