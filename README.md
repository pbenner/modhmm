## ModHMM

ModHMM is a highly modular genome segmentation method that incorporates genome-wide predictions from a set of classifiers. It implements a basic set of classifiers, but also allows to use predictions from third party classifiers.

### Installation



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

Create output directory
```sh
  mkdir mm10-liver-embryo-day12.5
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
```json
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

### Example 3: Estimate single-feature models on ENCODE data from mouse embyonic forebrain at day 11.5

Create a configuration file named `mm10-forebrain-embryo-day11.5.json`:
```json
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
    "Single-Feature Model Directory"  : "mm10-forebrain-embryo-day11.5:models",
    "Directory"                       : "mm10-forebrain-embryo-day11.5",
    "Description"                     : "forebrain embryo day11.5",
    "Threads"                         : 20,
    "Verbose"                         : 1
}
```

Create directories:
```ssh
  mkdir mm10-forebrain-embryo-day11.5
  mkdir mm10-forebrain-embryo-day11.5:models
```

Compute coverages and count files:
```ssh
  modhmm -c mm10-forebrain-embryo-day11.5.json compute-counts
```

Estimate a single-feature model for H3K27ac with one dirac component, two Poisson, and two geometric components:
```sh
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature h3k27ac 1 2 2
```

The resulting estimate can be easily inspected in *GNU-R* by first sourcing the file *modhmm_single_feature.R* from the ModHMM source package:
```R
  source("path/to/modhmm_single_feature.R")
```

Plot a histogram of the raw data and the estimated mixture distribution:
```R
  plot.model.and.counts("mm10-forebrain-embryo-day11.5:models/h3k27ac.json",
                        "mm10-forebrain-embryo-day11.5:models/h3k27ac.counts.json",
                         xlim=c(0,100))
  legend("topright", legend=0:4, pch=c(1,rep(NA,4)), lty=c(NA,2:5))
```

Select component 4 as foreground:
```sh
  echo '[4]' > mm10-forebrain-embryo-day11.5:models/h3k27ac.components.json
```

Visualize the foreground and background components of the mixture distribution:
```R
  plot.model.and.counts("mm10-forebrain-embryo-day11.5:models/h3k27ac.json",
                        "mm10-forebrain-embryo-day11.5:models/h3k27ac.counts.json",
                        "mm10-forebrain-embryo-day11.5:models/h3k27ac.components.json",
                        xlim=c(0,100))
  legend("topright", legend=c("foreground", "background"), lty=2:3)
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
  echo '[1,3,4,6]' > mm10-forebrain-embryo-day11.5:models/h3k4me1.components.json
  echo '[3,4]'     > mm10-forebrain-embryo-day11.5:models/h3k4me3.components.json
  echo '[2]'       > mm10-forebrain-embryo-day11.5:models/h3k4me3o1.components.json
  echo '[3,6]'     > mm10-forebrain-embryo-day11.5:models/h3k9me3.components.json
  echo '[1,3,4]'   > mm10-forebrain-embryo-day11.5:models/rna.components.json
  echo '[1,3]'     > mm10-forebrain-embryo-day11.5:models/rna-low.components.json
  echo '[9]'       > mm10-forebrain-embryo-day11.5:models/control.components.json
```
