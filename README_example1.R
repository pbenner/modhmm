#! /usr/bin/env Rscript

library(Gviz)
library(GenomicFeatures)
library(rtracklayer)
library(org.Hs.eg.db)
library(TxDb.Hsapiens.UCSC.hg38.knownGene)

## -----------------------------------------------------------------------------

RGB <- function(a,b,c) rgb(a/255,b/255,c/255,1.0)

## coordinates
## -----------------------------------------------------------------------------

regions <- GRanges(seqnames="chr1", ranges=IRanges(start=12035661, end=12215660))

## genes
## -----------------------------------------------------------------------------
if (TRUE) {
    genes <- GeneRegionTrack(TxDb.Hsapiens.UCSC.hg38.knownGene, name = "ensGene", transcriptAnnotation = "symbol", just.group="right",
                             background.title="white", fontcolor.title = "black", col.axis="black", cex.group=1.0)

    symbols <- unlist(mapIds(org.Hs.eg.db, gene(genes), "SYMBOL", "ENTREZID", multiVals = "first"))
    symbol(genes) <- symbols[gene(genes)]
}
## data
## -----------------------------------------------------------------------------
options(ucscChromosomeNames=FALSE)

dir <- "/project/modhmm-data/philipp/modhmm-encode-1.2.2/GRCh38-gastrocnemius-medialis"

d1  <- DataTrack(import(sprintf("%s/coverage-atac.bw", dir), format="bw", selection=regions),
                 name = "ATAC", showAxis=TRUE, cex.title=1.0,
                 background.title="transparent", col.border.title="transparent", background.panel="transparent", fontcolor.title = "black", col.axis="black", col="dodgerblue3")
d2  <- DataTrack(import(sprintf("%s/coverage-h3k27ac.bw", dir), format="bw", selection=regions),
                 name = "H3K27ac", showAxis=TRUE, cex.title=1.0,
                 background.title="transparent", col.border.title="transparent", background.panel="transparent", fontcolor.title = "black", col.axis="black", col="dodgerblue3")
d3  <- DataTrack(import(sprintf("%s/coverage-h3k4me1.bw", dir), format="bw", selection=regions),
                 name = "H3K4me1", showAxis=TRUE, cex.title=1.0,
                 background.title="transparent", col.border.title="transparent", background.panel="transparent", fontcolor.title = "black", col.axis="black", col="dodgerblue3")
d4  <- DataTrack(import(sprintf("%s/coverage-h3k4me3.bw", dir), format="bw", selection=regions),
                 name = "H3K4me3", showAxis=TRUE, cex.title=1.0,
                 background.title="transparent", col.border.title="transparent", background.panel="transparent", fontcolor.title = "black", col.axis="black", col="dodgerblue3")
d5  <- DataTrack(import(sprintf("%s/coverage-h3k27me3.bw", dir), format="bw", selection=regions),
                 name = "H3K27me3", showAxis=TRUE, cex.title=1.0,
                 background.title="transparent", col.border.title="transparent", background.panel="transparent", fontcolor.title = "black", col.axis="black", col="dodgerblue3")
d6 <- DataTrack(import(sprintf("%s/coverage-h3k9me3.bw", dir), format="bw", selection=regions),
                name = "H3K9me3", showAxis=TRUE, cex.title=1.0,
                background.title="transparent", col.border.title="transparent", background.panel="transparent", fontcolor.title = "black", col.axis="black", col="dodgerblue3")
d7 <- DataTrack(import(sprintf("%s/coverage-rna.bw", dir), format="bw", selection=regions),
                name = "log RNA", showAxis=TRUE, cex.title=1.0,
                transformation = function(x) { log(x+1) },
                background.title="transparent", col.border.title="transparent", background.panel="transparent", fontcolor.title = "black", col.axis="black", col="dodgerblue3")

## -------------------------------------------------------------------------
s1 <- import(sprintf("%s/segmentation.bed.gz", dir), format="bed")
## annotation
t1 <- AnnotationTrack(s1, width = 10, stacking = "full", fontsize.item=12, feature = s1$name, id=s1$name, featureAnnotation="id",
                      background.title="transparent", col.border.title="transparent", background.panel="transparent", fontcolor.title = "black", col.axis="black", col="transparent", rotation.title = 0, size=1.0,
                      fontcolor.item="black", cex=1.0, stackHeight=2,
                      "PA"    = "transparent",
                      "EA"    = "transparent",
                      "EA:tr" = "transparent",
                      "PR"    = "transparent",
                      "PR:tr" = "transparent",
                      "BI"    = "transparent",
                      "BI:tr" = "transparent",
                      "TR"    = "transparent",
                      "R1"    = "transparent",
                      "R2"    = "transparent",
                      "CL"    = "transparent")
## colored bars
t2 <- AnnotationTrack(s1, width = 10, stacking = "dense", fontsize.item=8, feature = s1$name, 
                      background.title="transparent", col.border.title="transparent", background.panel="transparent", fontcolor.title = "black", col.axis="black", col="transparent", rotation.title = 0, size=1.5, fontcolor.item="black", stackHeight=0.5,
                      "PA"    = RGB(30,144,255),
                      "EA"    = RGB(0,100,0),
                      "EA:tr" = RGB(0,100,0),
                      "PR"    = RGB(255,140,0),
                      "PR:tr" = RGB(255,140,0),
                      "BI"    = RGB(178,34,34),
                      "BI:tr" = RGB(178,34,34),
                      "TR"    = RGB(143,188,143),
                      "R1"    = RGB(255,69,0),
                      "R2"    = RGB(255,69,0),
                      "CL"    = RGB(255,0,255))

## plot
## -------------------------------------------------------------------------

chromosome <- as.character(seqnames(regions)[1])
from       <- start(regions)[1]
to         <- end  (regions)[1]

axis <- GenomeAxisTrack(name=chromosome, showTitle=TRUE, cex=1.0, cex.title=1.0,
                        background.title="white", fontcolor.title = "black", col.axis="black", rotation.title = 0, fontsize=12)

png("README_example1.png", height=700, width=1000)
plotTracks(c(axis, t1, t2, d1, d2, d3, d4, d5, d6, d7, genes),
           chromosome = chromosome, from = from, to = to, rotate.title=90,
           type="h", genome="hg38", window=1000)
dev.off()
