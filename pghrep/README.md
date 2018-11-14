# GO-Scratch

This simple layout for projects on Go language

Just clone this repository:

    $ git clone https://github.com/Juev/go-scratch new-project
    $ cd new-project
    $ rm -rf .git
    $ git init

Then change Makefile.

Or you can use my bash script:

    $ curl -o /usr/local/bin/go-scratch https://gist.githubusercontent.com/Juev/648da9739319c4514b2556ecb97e919f/raw/05ad9d488044a9bf6c72f6fb807eed81812b4ca0/go-scratch
    $ chmod +x /usr/local/bin/go-scratch
    $ go-scratch new-project

After this you should change Makefile.

## Makefile

You should change values of the following variables

* `BINARY` change from `go-scratch` to the name of your project
* `GITHUB_USERNAME` change from `Juev` to your username
* `VERSION` set to your project version

That should be it.
