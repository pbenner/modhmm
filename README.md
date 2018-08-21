## ModHMM

ModHMM is a highly modular genome segmentation method that incorporates genome-wide predictions from a set of classifiers. It implements a basic set of classifiers, but also allows to use predictions from third party classifiers.

### Installation



### Example 1: ENCODE mouse embyonic liver at day 12.5

Download BAM files from ENCODE and store them in a directory called `.bam`:
```sh
  # ATAC-seq
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF929LOH.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF848NLJ.bam
  # H3K27ac
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF524ZFV.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF322QGS.bam
  # H3K27me3
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF811DWT.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF171KAM.bam
  # H3K9me3
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF293UCG.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF777XFH.bam
  # H3K4me1
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF788JMC.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF340ACH.bam
  # H3K4me3
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF211WGC.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF587PZE.bam
  # RNA-seq
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF405LEY.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF627PCS.bam
  # Control
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF865QGZ.bam
  wget --directory-prefix=.bam http://www.encodeproject.org/files/ENCFF929LOH/@@download/ENCFF438RYK.bam
```

Create a configuration file named `mm10-liver-embryo-day12.5.conf`:
```json
{
    "Bam Directory" : ".bam",
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
    "Coverage Threads"                : 5,
    "Single-Feature Model Directory"  : "mm10-liver-embryo-day12.5:models",
    "Directory"                       : "mm10-liver-embryo-day12.5",
    "Description"                     : "liver embryo day12.5",
    "Threads"                         : 20,
    "Verbose"                         : 1
}
```

Execute ModHMM:
```sh
  mkdir mm10-liver-embryo-day12.5
  modhmm -c mm10-liver-embryo-day12.5.conf segmentation
```
