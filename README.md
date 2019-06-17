## ModHMM

ModHMM is a highly modular genome segmentation method based on hidden Markov models that incorporates genome-wide predictions from a set of classifiers. In order to simplify usage, ModHMM implements a default set of classifiers, but also allows to use predictions from third party methods.

### Installation

ModHMM can be installed by either downloading a binary from the [binary repository](https://github.com/pbenner/modhmm-binary) or by compiling the program from source.

To compile ModHMM you must first install the [Go compiler](https://golang.org/dl/). Afterwards, you may install ModHMM as follows:
```sh
  go get -v github.com/pbenner/modhmm
  cd $GOPATH/src/github.com/pbenner/modhmm
  make install
```

### Available Segmentations

ModHMM segmentations are available for several ENCODE data sets:

Tissue                         | Single-feature model | Segmentation
-------------------------------|----------------------|---------------
mm10 forebrain embryo day 11.5 | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day11.5-models.tar.bz2) | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day11.5/segmentation.bed.gz) |
mm10 forebrain embryo day 12.5 | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day12.5-models.tar.bz2) | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day12.5/segmentation.bed.gz) |
mm10 forebrain embryo day 13.5 | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day13.5-models.tar.bz2) | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day13.5/segmentation.bed.gz) |
mm10 forebrain embryo day 14.5 | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day14.5-models.tar.bz2) | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day14.5/segmentation.bed.gz) |
mm10 forebrain embryo day 15.5 | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day15.5-models.tar.bz2) | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day15.5/segmentation.bed.gz) |
mm10 forebrain embryo day 16.5 | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day16.5-models.tar.bz2) | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day16.5/segmentation.bed.gz) |
mm10 heart embryo day 14.5     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-heart-embryo-day14.5-models.tar.bz2)     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-heart-embryo-day14.5/segmentation.bed.gz)     |
mm10 heart embryo day 15.5     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-heart-embryo-day15.5-models.tar.bz2)     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-heart-embryo-day15.5/segmentation.bed.gz)     |
mm10 hindbrain embryo day 11.5 | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day11.5-models.tar.bz2) | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day11.5/segmentation.bed.gz) |
mm10 hindbrain embryo day 12.5 | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day12.5-models.tar.bz2) | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day12.5/segmentation.bed.gz) |
mm10 hindbrain embryo day 13.5 | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day13.5-models.tar.bz2) | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day13.5/segmentation.bed.gz) |
mm10 hindbrain embryo day 14.5 | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day14.5-models.tar.bz2) | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day14.5/segmentation.bed.gz) |
mm10 hindbrain embryo day 15.5 | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day15.5-models.tar.bz2) | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day15.5/segmentation.bed.gz) |
mm10 hindbrain embryo day 16.5 | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day16.5-models.tar.bz2) | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day16.5/segmentation.bed.gz) |
mm10 kidney embryo day 14.5    | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day14.5-models.tar.bz2)    | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day14.5/segmentation.bed.gz)    |
mm10 kidney embryo day 15.5    | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day15.5-models.tar.bz2)    | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day15.5/segmentation.bed.gz)    |
mm10 kidney embryo day 16.5    | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day16.5-models.tar.bz2)    | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day16.5/segmentation.bed.gz)    |
mm10 limb embryo day 14.5      | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-limb-embryo-day14.5-models.tar.bz2)      | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-limb-embryo-day14.5/segmentation.bed.gz)      |
mm10 limb embryo day 15.5      | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-limb-embryo-day15.5-models.tar.bz2)      | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-limb-embryo-day15.5/segmentation.bed.gz)      |
mm10 liver embryo day 11.5     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day11.5-models.tar.bz2)     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day11.5/segmentation.bed.gz)     |
mm10 liver embryo day 12.5     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day12.5-models.tar.bz2)     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day12.5/segmentation.bed.gz)     |
mm10 liver embryo day 13.5     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day13.5-models.tar.bz2)     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day13.5/segmentation.bed.gz)     |
mm10 liver embryo day 14.5     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day14.5-models.tar.bz2)     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day14.5/segmentation.bed.gz)     |
mm10 liver embryo day 15.5     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day15.5-models.tar.bz2)     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day15.5/segmentation.bed.gz)     |
mm10 liver embryo day 16.5     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day16.5-models.tar.bz2)     | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day16.5/segmentation.bed.gz)     |
mm10 lung embryo day 14.5      | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day14.5-models.tar.bz2)      | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day14.5/segmentation.bed.gz)      |
mm10 lung embryo day 15.5      | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day15.5-models.tar.bz2)      | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day15.5/segmentation.bed.gz)      |
mm10 lung embryo day 16.5      | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day16.5-models.tar.bz2)      | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day16.5/segmentation.bed.gz)      |
mm10 midbrain embryo day 11.5  | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day11.5-models.tar.bz2)  | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day11.5/segmentation.bed.gz)  |
mm10 midbrain embryo day 12.5  | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day12.5-models.tar.bz2)  | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day12.5/segmentation.bed.gz)  |
mm10 midbrain embryo day 13.5  | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day13.5-models.tar.bz2)  | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day13.5/segmentation.bed.gz)  |
mm10 midbrain embryo day 14.5  | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day14.5-models.tar.bz2)  | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day14.5/segmentation.bed.gz)  |
mm10 midbrain embryo day 15.5  | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day15.5-models.tar.bz2)  | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day15.5/segmentation.bed.gz)  |
mm10 midbrain embryo day 16.5  | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day16.5-models.tar.bz2)  | [Download](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day16.5/segmentation.bed.gz)  |

### Computing ModHMM Segmentations

Unlike most other genome segmentation methods, ModHMM so far depends on data from a fixed set of assays (ATAC/DNase, H3K27ac, H3K27me3, H3K4me1, H3K4me3, WCE/IgG, and RNA-seq). Preferentially, the data should be provided as BAM files, but it is also possible to use bigWig files as input. If BAM files are provided, ModHMM computes coverages and in case of single-end sequencing data automatically estimates the mean fragment length.

### Example 1: Compute segmentation on ENCODE data from mouse embyonic liver at day 12.5

Download BAM files from ENCODE and store them in a directory called `.bam`:
```sh
  # ATAC-seq
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF929LOH.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF848NLJ/@@download/ENCFF848NLJ.bam
  # H3K27ac
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF524ZFV/@@download/ENCFF524ZFV.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF322QGS/@@download/ENCFF322QGS.bam
  # H3K27me3
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF811DWT/@@download/ENCFF811DWT.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF171KAM/@@download/ENCFF171KAM.bam
  # H3K9me3
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF293UCG/@@download/ENCFF293UCG.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF777XFH/@@download/ENCFF777XFH.bam
  # H3K4me1
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF788JMC/@@download/ENCFF788JMC.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF340ACH/@@download/ENCFF340ACH.bam
  # H3K4me3
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF211WGC/@@download/ENCFF211WGC.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF587PZE/@@download/ENCFF587PZE.bam
  # RNA-seq
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF405LEY/@@download/ENCFF405LEY.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF627PCS/@@download/ENCFF627PCS.bam
  # Control
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF865QGZ/@@download/ENCFF865QGZ.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF438RYK/@@download/ENCFF438RYK.bam
```

Create a configuration file named `mm10-liver-embryo-day12.5.json` (ModHMM accepts an extended JSON format that allows comments):
```R
{
    # Directory containing feature alignment files
    "Bam Directory" : ".bam",
    # Names of alignment files
    "Bam Files"     : {
        "ATAC"      : ["ENCFF929LOH.bam", "ENCFF848NLJ.bam"],
        "H3K27ac"   : ["ENCFF524ZFV.bam", "ENCFF322QGS.bam"],
        "H3K27me3"  : ["ENCFF811DWT.bam", "ENCFF171KAM.bam"],
        "H3K9me3"   : ["ENCFF293UCG.bam", "ENCFF777XFH.bam"],
        "H3K4me1"   : ["ENCFF788JMC.bam", "ENCFF340ACH.bam"],
        "H3K4me3"   : ["ENCFF211WGC.bam", "ENCFF587PZE.bam"],
        "RNA"       : ["ENCFF405LEY.bam", "ENCFF627PCS.bam"],
        "Control"   : ["ENCFF865QGZ.bam", "ENCFF438RYK.bam"]
    },
    # Number of threads used for computing coverage bigWigs (memory intense!)
    "Coverage Threads"                : 5,
    # Directory containing the single-feature models that must either be
    # estimated by hand or downloaded from the ModHMM repository
    "Single-Feature Model Directory"  : "mm10-liver-embryo-day12.5:models",
    # Directory containing all auxiliary files and the final segmentation
    "Directory"                       : "mm10-liver-embryo-day12.5",
    "Description"                     : "liver embryo day12.5",
    # Number of threads used for evaluating classifiers and computing the segmentation
    "Threads"                         : 20,
    # Verbose level (0: no output, 1: low, 2: high)
    "Verbose"                         : 1
}
```

Create output directories
```sh
  mkdir mm10-liver-embryo-day12.5
  mkdir mm10-liver-embryo-day12.5:models
```

Execute ModHMM:
```sh
  modhmm -c mm10-liver-embryo-day12.5.json segmentation
```

### Example 2: Compute segmentation on ENCODE data from mouse embyonic forebrain at day 11.5 with single-feature models estimated on embryonic liver at day 12.5

Download BAM files from ENCODE and store them in a directory called `.bam`:
```sh
  # DNase-seq
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF546SVK/@@download/ENCFF546SVK.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF358BLW/@@download/ENCFF358BLW.bam
  # H3K27ac
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF439HJF/@@download/ENCFF439HJF.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF393PYK/@@download/ENCFF393PYK.bam
  # H3K27me3
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF854DPK/@@download/ENCFF854DPK.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF330LCP/@@download/ENCFF330LCP.bam
  # H3K9me3
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF828ITY/@@download/ENCFF828ITY.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF670JTV/@@download/ENCFF670JTV.bam
  # H3K4me1
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF528ZVN/@@download/ENCFF528ZVN.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF695PCS/@@download/ENCFF695PCS.bam
  # H3K4me3
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF437KKV/@@download/ENCFF437KKV.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF354JHH/@@download/ENCFF354JHH.bam
  # RNA-seq
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF625THA/@@download/ENCFF625THA.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF177AZU/@@download/ENCFF177AZU.bam
  # Control
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF631YQS/@@download/ENCFF631YQS.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF658BBR/@@download/ENCFF658BBR.bam
```

Create a configuration file named `mm10-forebrain-embryo-day11.5.json`:
```R
{
    "Bam Directory" : ".bam",
    "Bam Files"     : {
        "ATAC"      : ["ENCFF426VDN.bam", "ENCFF275OKU.bam"],
        #"DNase"    : ["ENCFF546SVK.bam", "ENCFF358BLW.bam"],
        "H3K27ac"   : ["ENCFF439HJF.bam", "ENCFF393PYK.bam"],
        "H3K27me3"  : ["ENCFF854DPK.bam", "ENCFF330LCP.bam"],
        "H3K9me3"   : ["ENCFF828ITY.bam", "ENCFF670JTV.bam"],
        "H3K4me1"   : ["ENCFF528ZVN.bam", "ENCFF695PCS.bam"],
        "H3K4me3"   : ["ENCFF437KKV.bam", "ENCFF354JHH.bam"],
        "RNA"       : ["ENCFF625THA.bam", "ENCFF177AZU.bam"],
        "Control"   : ["ENCFF631YQS.bam", "ENCFF658BBR.bam"]
    },
    "Coverage Threads"                : 5,
    "Single-Feature Model Directory"  : "mm10-liver-embryo-day12.5:models",
    "Directory"                       : "mm10-forebrain-embryo-day11.5",
    "Description"                     : "forebrain embryo day11.5",
    "Threads"                         : 20,
    "Verbose"                         : 1
}
```

Create output directories
```sh
  mkdir mm10-forebrain-embryo-day11.5
  mkdir mm10-forebrain-embryo-day11.5:models
```

Execute ModHMM:
```sh
  modhmm -c mm10-forebrain-embryo-day11.5.json segmentation
```

### Example 3: Estimate custom single-feature models on ENCODE data from mouse embyonic forebrain at day 11.5

Create a configuration file named `mm10-forebrain-embryo-day11.5.json` and set model files static to prevent automatic updates:
```R
{
    "Bam Directory" : ".bam",
    "Bam Files"     : {
        #"ATAC"     : ["ENCFF426VDN.bam", "ENCFF275OKU.bam"],
        "DNase"     : ["ENCFF546SVK.bam", "ENCFF358BLW.bam"],
        "H3K27ac"   : ["ENCFF439HJF.bam", "ENCFF393PYK.bam"],
        "H3K27me3"  : ["ENCFF854DPK.bam", "ENCFF330LCP.bam"],
        "H3K9me3"   : ["ENCFF828ITY.bam", "ENCFF670JTV.bam"],
        "H3K4me1"   : ["ENCFF528ZVN.bam", "ENCFF695PCS.bam"],
        "H3K4me3"   : ["ENCFF437KKV.bam", "ENCFF354JHH.bam"],
        "RNA"       : ["ENCFF625THA.bam", "ENCFF177AZU.bam"],
        "Control"   : ["ENCFF631YQS.bam", "ENCFF658BBR.bam"]
    },
    "Coverage Threads"                : 5,
    "Single-Feature Model Static"     : true,
    "Single-Feature Model Directory"  : "mm10-forebrain-embryo-day11.5:models",
    "Directory"                       : "mm10-forebrain-embryo-day11.5",
    "Description"                     : "forebrain embryo day11.5",
    "Threads"                         : 20,
    "Verbose"                         : 1
}
```

Create directories:
```sh
  mkdir mm10-forebrain-embryo-day11.5
  mkdir mm10-forebrain-embryo-day11.5:models
```

Compute coverages and count files:
```sh
  modhmm -c mm10-forebrain-embryo-day11.5.json compute-counts
```

Estimate a single-feature model for H3K27ac with one dirac component, two Poisson, and two geometric components:
```sh
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature h3k27ac 1 2 2
```

The resulting estimate can be easily inspected:
```sh
  modhmm -c mm10-forebrain-embryo-day11.5.json plot-single-feature --xlim=0-200 h3k27ac
```

Select component 4 as foreground:
```sh
  echo '[4]' > mm10-forebrain-embryo-day11.5:models/h3k27ac.components.json
```

Visualize the merged foreground and background components of the mixture distribution:
```R
  modhmm -c mm10-forebrain-embryo-day11.5.json plot-single-feature --xlim=0-200 h3k27ac
```

Repeat these steps for all remaining features:
```sh
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature dnase     1 1 3
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature h3k27me3  4 4 1
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature h3k4me1   1 8 0
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature h3k4me3   1 1 3
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature h3k4me3o1 0 1 2
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature h3k9me3   2 4 1
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature rna       1 0 4
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature control   7 2 1

  echo '[3,4]'     > mm10-forebrain-embryo-day11.5:models/dnase.components.json
  echo '[8]'       > mm10-forebrain-embryo-day11.5:models/h3k27me3.components.json
  echo '[5,6,7,8]' > mm10-forebrain-embryo-day11.5:models/h3k4me1.components.json
  echo '[3,4]'     > mm10-forebrain-embryo-day11.5:models/h3k4me3.components.json
  echo '[2]'       > mm10-forebrain-embryo-day11.5:models/h3k4me3o1.components.json
  echo '[5,6]'     > mm10-forebrain-embryo-day11.5:models/h3k9me3.components.json
  echo '[2,3,4]'   > mm10-forebrain-embryo-day11.5:models/rna.components.json
  echo '[1,2]'     > mm10-forebrain-embryo-day11.5:models/rna-low.components.json
  echo '[9]'       > mm10-forebrain-embryo-day11.5:models/control.components.json
```
