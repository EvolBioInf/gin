# [`gin`](https://owncloud.gwdg.de/index.php/s/zxz4J1ekIrnbabv)
## Description
Annotate arbitrary genome intervals with genes and carry out
functional enrichment analysis with Monte Carlo test.
## Author
[Bernhard Haubold](http://guanine.evolbio.mpg.de/), `haubold@evolbio.mpg.de`
## Make the Programs
If you are on an Ubuntu system like Ubuntu on
[wsl](https://learn.microsoft.com/en-us/windows/wsl/install) under
MS-Windows or the [Ubuntu Docker
container](https://hub.docker.com/_/ubuntu), you can clone the
repository and change into it.

`git clone https://github.com/evolbioinf/gin`  
`cd gin`

Then install the additional dependencies by running the script
[`setup.sh`](scripts/setup.sh).

`bash scripts/setup.sh`

Make the programs.

`make`

The directory `bin` now contains the binaries.
