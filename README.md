# [`gin`](https://owncloud.gwdg.de/index.php/s/f7jx4E87boc9QU4)
## Description
Annotate arbitrary genome intervals with genes and carry out
functional enrichment analysis with Monte Carlo test.
## Authors
[Bernhard Haubold](http://guanine.evolbio.mpg.de/)
and [Beatriz Vieira Mourato](https://beatrizvm.github.io/) 
## Make the Programs
If you are on an Ubuntu system like Ubuntu on
[wsl](https://learn.microsoft.com/en-us/windows/wsl/install) under
MS-Windows or the [Ubuntu Docker
container](https://hub.docker.com/_/ubuntu), you can clone the
repository and change into it.

```
git clone https://github.com/evolbioinf/gin
cd gin
```

Then install the additional dependencies by running the script
[`setup.sh`](scripts/setup.sh).

```
bash scripts/setup.sh
```

Make the programs.

```
make
```

The directory `bin` now contains the binaries. You can also download
example data into the directory `data` and test the programs on it.

```
make test
```


