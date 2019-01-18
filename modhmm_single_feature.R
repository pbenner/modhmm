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

eval.mixture.components <- function(x, model, weights, components=NULL) {
    if (model$Name == "scalar:mixture distribution") {
        if (is.null(components)) {
            components <- 1:length(model$Distributions)
        }
        y <- rep(0, length(x))
        for (i in components) {
            y <- y + eval.mixture.components(x, model$Distributions[[i]], weights=model$Parameters[[i]])
        }
        return (y)
    } else
    if (model$Name == "scalar:poisson distribution") {
        return (weights[1]*dpois(x, model$Parameters[1]))
    } else
    if (model$Name == "scalar:geometric distribution") {
        return (weights[1]*dgeom(x, model$Parameters[1]))
    } else
    if (model$Name == "scalar:pdf translation") {
        return (eval.mixture.components(x+model$Parameters[1], model$Distributions[[1]], weights))
    } else
    if (model$Name == "scalar:delta distribution") {
        y <- rep(0, length(x))
        y[x == model$Parameters[1]] = weights[1]
        return (y)
    } else {
        stop(sprintf("could not parse: %s", model$Name))
    }
}

## -----------------------------------------------------------------------------

plot.mixture <- function(x, model, weights=NULL, col=NULL, lty=NULL, ...) {
    if (is.null(col) || is.na(col)) {
        col="black"
    }
    if (is.null(lty) || is.na(lty)) {
        lty=1
    }
    if (model$Name == "scalar:mixture distribution") {
        for (i in 1:length(model$Distributions)) {
            plot.mixture(x, model$Distributions[[i]], weights=model$Parameters[[i]], col=col[i], lty=lty[i], ...)
        }
    } else {
        y <- eval.mixture.components(x, model, weights)
        if (sum(y != 0) > 1) {
            lines(x, y, col=col, lty=lty, ...)
        } else {
            points(x, y, col=col, lty=lty, ...)
        }
    }
}

plot.mixture.joined <- function(x, model, components, col=NULL, lty=NULL, ...) {
    if (is.null(col) || is.na(col)) {
        col=rep("black", 2)
    }
    if (is.null(lty) || is.na(lty)) {
        lty=c(1,1)
    }
    n  <- length(model$Distributions)
    y1 <- eval.mixture.components(x, model, components=components)
    y2 <- eval.mixture.components(x, model, components=(1:n)[!(1:n %in% components)])
    lines(x, y1, lty=lty[1], col=col[1], ...)
    lines(x, y2, lty=lty[2], col=col[2], ...)
}

## -----------------------------------------------------------------------------

plot.model.and.counts <- function(modelFilename, countsFilename, componentsFilename=NULL, xlab="coverage", ylab="probability", log="y", main="", lty=2:100, col=NULL, add=FALSE, ...) {
    counts <- fromJSON(file=countsFilename)
    model  <- fromJSON(file= modelFilename)
    if (!add) {
        plot(Y/sum(Y) ~ X, counts, type="l", xlab=xlab, ylab=ylab, log=log, main=main, ...)
    }
    if (is.null(componentsFilename)) {
        plot.mixture(counts$X, model, lty=lty, col=col, ...)
    } else {
        components <- fromJSON(file=componentsFilename)
        if (length(components) > 0) {
            components <- components + 1
        }
        plot.mixture.joined(counts$X, model, components, lty=lty, col=col, ...)
    }
}

## -----------------------------------------------------------------------------

plot.model.dir <- function(dir,
                           open.assay = "atac",
                           features   = c(open.assay, "h3k27ac", "h3k27me3", "h3k4me1", "h3k4me3", "h3k4me3o1", "h3k9me3", "rna"),
                           xlim       = c(400, 200, 100, 100, 300, 800, 200, 200, 200),
                           ...) {
    par(mfrow=c(3,3))
    for (i in 1:length(features)) {
        plot.model.and.counts(sprintf("%s/%s.json", dir, features[i]),
                              sprintf("%s/%s.counts.json", dir, features[i]),
                              sprintf("%s/%s.components.json", dir, features[i]),
                              xlim=c(0,xlim[i]), main=features[i], ...)
        legend("topright", legend=c("foreground", "background"), lty=2:3)
    }
    plot.model.and.counts(sprintf("%s/rna.json", dir, features[8]),
                          sprintf("%s/rna.counts.json", dir, features[8]),
                          sprintf("%s/rna-low.components.json", dir, features[8]),
                          xlim=c(0,xlim[9]), main="rna-low", ...)
    legend("topright", legend=c("foreground", "background"), lty=2:3)
}
