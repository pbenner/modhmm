## Copyright (C) 2018 Philipp Benner
##
## This program is free software: you can redistribute it and/or modify
## it under the terms of the GNU General Public License as published by
## the Free Software Foundation, either version 3 of the License, or
## (at your option) any later version.
##
## This program is distributed in the hope that it will be useful,
## but WITHOUT ANY WARRANTY; without even the implied warranty of
## MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
## GNU General Public License for more details.
##
## You should have received a copy of the GNU General Public License
## along with this program.  If not, see <http://www.gnu.org/licenses/>.
##

library(rjson)

## -----------------------------------------------------------------------------

plot.distribution <- function(x, json, weights=NULL, col=NULL, lty=NULL, ...) {
    if (is.null(col) || is.na(col)) {
        col="black"
    }
    if (is.null(lty) || is.na(lty)) {
        lty=1
    }
    if (json$Name == "scalar:mixture distribution") {
        for (i in 1:length(json$Distributions)) {
            plot.distribution(x, json$Distributions[[i]], weights=json$Parameters[[i]], col=col[i], lty=lty[i], ...)
        }
    } else
    if (json$Name == "scalar:poisson distribution") {
        lines(x, weights[1]*dpois(x, json$Parameters[1]), col=col, lty=lty, ...)
    } else
    if (json$Name == "scalar:geometric distribution") {
        lines(x, weights[1]*dgeom(x, json$Parameters[1]), col=col, lty=lty, ...)
    } else
    if (json$Name == "scalar:pdf translation") {
        plot.distribution(x+json$Parameters[1], json$Distributions[[1]], weights=weights, col=col, lty=lty, ...)
    } else
    if (json$Name == "scalar:delta distribution") {
        points(json$Parameters[1], weights[1], col=col, ...)
    } else {
        stop(sprintf("could not parse: %s", json$Name))
    }
}

## -----------------------------------------------------------------------------

plot.distribution.and.counts <- function(modelFilename, countsFilename, xlab="coverage", ylab="probability", log="y", main="", ...) {
    counts <- fromJSON(file=countsFilename)
    model  <- fromJSON(file= modelFilename)
    plot(Y/sum(Y) ~ X, counts, type="l", xlab=xlab, ylab=ylab, log=log, main=main, ...)
    plot.distribution(counts$X, model, lty=1:100, ...)
}
