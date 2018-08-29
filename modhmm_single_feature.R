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

eval.mixture.components <- function(x, json, weights, components=NULL) {
    if (json$Name == "scalar:mixture distribution") {
        if (is.null(components)) {
            components <- 1:length(json$Distributions)
        }
        y <- rep(0, length(x))
        for (i in components) {
            y <- y + eval.mixture.components(x, json$Distributions[[i]], weights=json$Parameters[[i]])
        }
        return (y)
    } else
    if (json$Name == "scalar:poisson distribution") {
        return (weights[1]*dpois(x, json$Parameters[1]))
    } else
    if (json$Name == "scalar:geometric distribution") {
        return (weights[1]*dgeom(x, json$Parameters[1]))
    } else
    if (json$Name == "scalar:pdf translation") {
        return (eval.mixture.components(x+json$Parameters[1], json$Distributions[[1]], weights))
    } else
    if (json$Name == "scalar:delta distribution") {
        y <- rep(0, length(x))
        y[x == json$Parameters[1]] = weights[1]
        return (y)
    } else {
        stop(sprintf("could not parse: %s", json$Name))
    }
}

## -----------------------------------------------------------------------------

plot.mixture <- function(x, json, weights=NULL, col=NULL, lty=NULL, ...) {
    if (is.null(col) || is.na(col)) {
        col="black"
    }
    if (is.null(lty) || is.na(lty)) {
        lty=1
    }
    if (json$Name == "scalar:mixture distribution") {
        for (i in 1:length(json$Distributions)) {
            plot.mixture(x, json$Distributions[[i]], weights=json$Parameters[[i]], col=col[i], lty=lty[i], ...)
        }
    } else {
        y <- eval.mixture.components(x, json, weights)
        if (sum(y != 0) > 1) {
            lines(x, y, col=col, lty=lty, ...)
        } else {
            points(x, y, col=col, lty=lty, ...)
        }
    }
}

## -----------------------------------------------------------------------------

plot.model.and.counts <- function(modelFilename, countsFilename, components.fg=NULL, xlab="coverage", ylab="probability", log="y", main="", lty=2:100, col=NULL, ...) {
    counts <- fromJSON(file=countsFilename)
    model  <- fromJSON(file= modelFilename)
    plot(Y/sum(Y) ~ X, counts, type="l", xlab=xlab, ylab=ylab, log=log, main=main, ...)
    if (is.null(components.fg)) {
        plot.mixture(counts$X, model, lty=lty, col=col, ...)
    } else {
        if (is.null(col) || is.na(col)) {
            col=rep("black", 2)
        }
        if (is.null(lty) || is.na(lty)) {
            lty=c(1,1)
        }
        n  <- length(model$Distributions)
        y1 <- eval.mixture.components(counts$X, model, components=components.fg)
        y2 <- eval.mixture.components(counts$X, model, components=(1:n)[!(1:n %in% components.fg)])
        lines(counts$X, y1, lty=lty[1], col=col[1], ...)
        lines(counts$X, y2, lty=lty[2], col=col[2], ...)
    }
}
