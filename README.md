## ModHMM

ModHMM is a highly modular genome segmentation method based on a hidden Markov model that incorporates genome-wide predictions from a set of classifiers. In order to simplify usage, ModHMM implements a default set of classifiers, but also allows to use predictions from third party methods.

References:

Philipp Benner and Martin Vingron. *ModHMM: A Modular Supra-Bayesian Genome Segmentation Method*. International Conference on Research in Computational Molecular Biology (RECOMB). Springer, Cham, 2019. S. 35-50. [[Link]](https://link.springer.com/chapter/10.1007/978-3-030-17083-7_3)

<center><img src="https://raw.githubusercontent.com/pbenner/modhmm/master/README_example1.png" alt="ModHMM" width="720" height="429" /></center>

### Available Segmentations

ModHMM segmentations are available for several ENCODE data sets:

Tissue | Date |Segmentation | Posteriors | Config | Single-feature model | Comment
-------|------|-------------|------------|--------|----------------------|--------
GRCh38 ascending aorta | 2019-06-24 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-ascending-aorta/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/GRCh38-ascending-aorta/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-ascending-aorta.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-ascending-aorta-models.tar.bz2) | total RNA-seq
GRCh38 gastrocnemius medialis | 2019-06-24 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-gastrocnemius-medialis/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/GRCh38-gastrocnemius-medialis/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-gastrocnemius-medialis.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-gastrocnemius-medialis-models.tar.bz2) | total RNA-seq
GRCh38 heart left ventricle | 2019-06-24 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-heart-left-ventricle/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/GRCh38-heart-left-ventricle/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-heart-left-ventricle.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-heart-left-ventricle-models.tar.bz2) | total RNA-seq
GRCh38 lung upper lobe | 2019-06-24 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-lung-upper-lobe/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/GRCh38-lung-upper-lobe/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-lung-upper-lobe.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-lung-upper-lobe-models.tar.bz2) | total RNA-seq
GRCh38 pancreas body | 2019-06-24 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-pancreas-body/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/GRCh38-pancreas-body/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-pancreas-body.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-pancreas-body-models.tar.bz2) | total RNA-seq
GRCh38 spleen | 2019-06-24 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-spleen/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/GRCh38-spleen/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-spleen.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-spleen-models.tar.bz2) | total RNA-seq
GRCh38 stomach | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-stomach/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/GRCh38-stomach/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-stomach.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-stomach-models.tar.bz2) | total RNA-seq
GRCh38 tibial nerve | 2019-06-24 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-tibial-nerve/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/GRCh38-tibial-nerve/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-tibial-nerve.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-tibial-nerve-models.tar.bz2) | total RNA-seq
GRCh38 transverse colon | 2019-06-24 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-transverse-colon/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/GRCh38-transverse-colon/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-transverse-colon.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-transverse-colon-models.tar.bz2) | total RNA-seq
GRCh38 uterus | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-uterus/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/GRCh38-uterus/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-uterus.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/GRCh38-uterus-models.tar.bz2) | total RNA-seq
mm10 forebrain embryo day11.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day11.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-forebrain-embryo-day11.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day11.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day11.5-models.tar.bz2) | poly-A RNA-seq
mm10 forebrain embryo day12.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day12.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-forebrain-embryo-day12.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day12.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day12.5-models.tar.bz2) | poly-A RNA-seq
mm10 forebrain embryo day13.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day13.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-forebrain-embryo-day13.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day13.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day13.5-models.tar.bz2) | poly-A RNA-seq
mm10 forebrain embryo day14.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day14.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-forebrain-embryo-day14.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day14.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day14.5-models.tar.bz2) | poly-A RNA-seq
mm10 forebrain embryo day15.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day15.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-forebrain-embryo-day15.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day15.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day15.5-models.tar.bz2) | poly-A RNA-seq
mm10 forebrain embryo day16.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day16.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-forebrain-embryo-day16.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day16.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-forebrain-embryo-day16.5-models.tar.bz2) | poly-A RNA-seq
mm10 heart embryo day14.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-heart-embryo-day14.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-heart-embryo-day14.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-heart-embryo-day14.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-heart-embryo-day14.5-models.tar.bz2) | poly-A RNA-seq
mm10 heart embryo day15.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-heart-embryo-day15.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-heart-embryo-day15.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-heart-embryo-day15.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-heart-embryo-day15.5-models.tar.bz2) | poly-A RNA-seq
mm10 hindbrain embryo day11.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day11.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-hindbrain-embryo-day11.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day11.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day11.5-models.tar.bz2) | poly-A RNA-seq
mm10 hindbrain embryo day12.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day12.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-hindbrain-embryo-day12.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day12.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day12.5-models.tar.bz2) | poly-A RNA-seq
mm10 hindbrain embryo day13.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day13.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-hindbrain-embryo-day13.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day13.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day13.5-models.tar.bz2) | poly-A RNA-seq
mm10 hindbrain embryo day14.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day14.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-hindbrain-embryo-day14.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day14.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day14.5-models.tar.bz2) | poly-A RNA-seq
mm10 hindbrain embryo day15.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day15.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-hindbrain-embryo-day15.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day15.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day15.5-models.tar.bz2) | poly-A RNA-seq
mm10 hindbrain embryo day16.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day16.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-hindbrain-embryo-day16.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day16.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-hindbrain-embryo-day16.5-models.tar.bz2) | poly-A RNA-seq
mm10 kidney embryo day14.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day14.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-kidney-embryo-day14.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day14.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day14.5-models.tar.bz2) | poly-A RNA-seq
mm10 kidney embryo day15.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day15.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-kidney-embryo-day15.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day15.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day15.5-models.tar.bz2) | poly-A RNA-seq
mm10 kidney embryo day16.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day16.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-kidney-embryo-day16.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day16.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-kidney-embryo-day16.5-models.tar.bz2) | poly-A RNA-seq
mm10 limb embryo day14.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-limb-embryo-day14.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-limb-embryo-day14.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-limb-embryo-day14.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-limb-embryo-day14.5-models.tar.bz2) | poly-A RNA-seq
mm10 limb embryo day15.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-limb-embryo-day15.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-limb-embryo-day15.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-limb-embryo-day15.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-limb-embryo-day15.5-models.tar.bz2) | poly-A RNA-seq
mm10 liver embryo day11.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day11.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-liver-embryo-day11.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day11.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day11.5-models.tar.bz2) | poly-A RNA-seq
mm10 liver embryo day12.5 | 2019-06-23 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day12.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-liver-embryo-day12.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day12.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day12.5-models.tar.bz2) | poly-A RNA-seq
mm10 liver embryo day13.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day13.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-liver-embryo-day13.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day13.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day13.5-models.tar.bz2) | poly-A RNA-seq
mm10 liver embryo day14.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day14.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-liver-embryo-day14.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day14.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day14.5-models.tar.bz2) | poly-A RNA-seq
mm10 liver embryo day15.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day15.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-liver-embryo-day15.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day15.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day15.5-models.tar.bz2) | poly-A RNA-seq
mm10 liver embryo day16.5 | 2019-06-23 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day16.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-liver-embryo-day16.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day16.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-liver-embryo-day16.5-models.tar.bz2) | poly-A RNA-seq
mm10 lung embryo day14.5 | 2019-06-23 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day14.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-lung-embryo-day14.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day14.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day14.5-models.tar.bz2) | poly-A RNA-seq
mm10 lung embryo day15.5 | 2019-06-23 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day15.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-lung-embryo-day15.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day15.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day15.5-models.tar.bz2) | poly-A RNA-seq
mm10 lung embryo day16.5 | 2019-06-22 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day16.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-lung-embryo-day16.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day16.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-lung-embryo-day16.5-models.tar.bz2) | poly-A RNA-seq
mm10 midbrain embryo day11.5 | 2019-06-23 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day11.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-midbrain-embryo-day11.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day11.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day11.5-models.tar.bz2) | poly-A RNA-seq
mm10 midbrain embryo day12.5 | 2019-06-23 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day12.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-midbrain-embryo-day12.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day12.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day12.5-models.tar.bz2) | poly-A RNA-seq
mm10 midbrain embryo day13.5 | 2019-06-23 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day13.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-midbrain-embryo-day13.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day13.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day13.5-models.tar.bz2) | poly-A RNA-seq
mm10 midbrain embryo day14.5 | 2019-06-23 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day14.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-midbrain-embryo-day14.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day14.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day14.5-models.tar.bz2) | poly-A RNA-seq
mm10 midbrain embryo day15.5 | 2019-06-23 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day15.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-midbrain-embryo-day15.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day15.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day15.5-models.tar.bz2) | poly-A RNA-seq
mm10 midbrain embryo day16.5 | 2019-06-23 | [Segmentation](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day16.5/segmentation.bed.gz) |   [Posteriors](https://owww.molgen.mpg.de/~benner/pool/modhmm/mm10-midbrain-embryo-day16.5/) |       [Config](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day16.5.json) |       [Models](https://github.com/pbenner/modhmm-segmentations/raw/master/mm10-midbrain-embryo-day16.5-models.tar.bz2) | poly-A RNA-seq

### Installation

ModHMM can be installed by either downloading a binary from the [binary repository](https://github.com/pbenner/modhmm-binary) or by compiling the program from source.

To compile ModHMM you must first install the [Go compiler](https://golang.org/dl/). Afterwards, you may install ModHMM as follows:
```sh
  go get -v github.com/pbenner/modhmm
  cd $GOPATH/src/github.com/pbenner/modhmm
  make install
```

### Computing ModHMM Segmentations

Unlike most other genome segmentation methods, ModHMM so far depends on data from a fixed set of features or assays (ATAC/DNase, H3K27ac, H3K27me3, H3K4me1, H3K4me3, WCE/IgG, and RNA-seq). Open chromatin information must be either provided as ATAC-seq or DNase-seq data. Preferentially, the data should be provided as BAM files, but it is also possible to use bigWig files as input. If BAM files are provided, ModHMM first computes a coverage of each feature and in case of single-end sequencing data automatically estimates the mean fragment length. The most difficult and crucial step of computing ModHMM segmentations is the enrichment analysis of individual features. ModHMM uses a feature specific mixture model for estimating genome-wide enrichment probabilities (i.e. for each feature the probability of a peak at each genomic position). These single-feature models can be estimated by ModHMM, but there is so far no good heuristic for selecting the correct number of components. To make the software easily applicable to new data sets, ModHMM implements a default single-feature model. Genome segmentations based on this default model are slightly less accurate, but allow the user to obtain a quick annotation with very little effort.

ModHMM requires a configuration file in JSON format. The following is a simple example where data is provided as BAM files:
```R
{
    # Directory containing feature alignment files
    "Bam Directory" : ".bam",
    # Names of alignment files (for each feature a comma separated list of
    # replicates must be specified, any number of replicates is supported)
    "Bam Files"     : {
        "ATAC"      : [    "atac-rep1.bam",     "atac-rep2.bam"],
        #"DNase"    : [   "dnase-rep1.bam",    "dnase-rep2.bam"],
        "H3K27ac"   : [ "h3k27ac-rep1.bam",  "h3k27ac-rep2.bam"],
        "H3K27me3"  : ["h3k27me3-rep1.bam", "h3k27me3-rep2.bam"],
        "H3K9me3"   : [ "h3k9me3-rep1.bam",  "h3k9me3-rep2.bam"],
        "H3K4me1"   : [ "h3k4me1-rep1.bam",  "h3k4me1-rep2.bam"],
        "H3K4me3"   : [ "h3k4me3-rep1.bam",  "h3k4me3-rep2.bam"],
        "RNA"       : [     "rna-rep1.bam",      "rna-rep2.bam"],
        "Control"   : [ "control-rep1.bam",  "control-rep2.bam"]
    },
    # Number of threads used for computing coverage bigWigs (memory intense!)
    "Coverage Threads"                : 5,
    # ModHMM has several parameters that must be estimated, the single-feature mixture
    # and the HMM transition parameters. By default ModHMM won't estimate these parameters,
    # but use a fallback model with pre-estimated parameter. This option selects the type
    # of fallback model ["mm10", "GRCh38"]
    "Model Fallback"                  : "mm10",
    # Directory containing all auxiliary files and the final segmentation
    "Directory"                       : "mm10-liver-embryo-day12.5",
    "Description"                     : "liver embryo day12.5",
    # Number of threads used for evaluating classifiers and computing the segmentation
    "Threads"                         : 20,
    # Verbose level (0: no output, 1: low, 2: high)
    "Verbose"                         : 1
}
```
The following configuration can be used if data instead is given in bigWig format:
```R
{
    # Data is provided as bigWig files. Set all coverage files static!
    "Coverage Files": {
        "ATAC"   : {"Filename": "coverage-atac.bw",    "Static": true },
        #"DNase" : {"Filename": "coverage-dnase.bw",   "Static": true },
        "H3K27ac": {"Filename": "coverage-h3k27ac.bw", "Static": true },
        "H3K4me1": {"Filename": "coverage-h3k4me1.bw", "Static": true },
        "H3K4me3": {"Filename": "coverage-h3k4me3.bw", "Static": true },
        "RNA"    : {"Filename": "coverage-rna.bw",     "Static": true },
        "Control": {"Filename": "coverage-control.bw", "Static": true }
    },
    # ModHMM has several parameters that must be estimated, the single-feature mixture
    # and the HMM transition parameters. By default ModHMM won't estimate these parameters,
    # but use a fallback model with pre-estimated parameter. This option selects the type
    # of fallback model ["mm10", "GRCh38"]
    "Model Fallback"                  : "mm10",
    # Directory containing all auxiliary files and the final segmentation
    "Directory"                       : "mm10-liver-embryo-day12.5",
    "Description"                     : "liver embryo day12.5",
    # Number of threads used for evaluating classifiers and computing the segmentation
    "Threads"                         : 20,
    # Verbose level (0: no output, 1: low, 2: high)
    "Verbose"                         : 1
}
```
Coverage bigWig files must be placed in the directory `mm10-liver-embryo-day12.5`. Setting the option `Static` to true tells ModHMM that the provided bigWig files are not automatically generated and should not be overwritten. It is also possible to use BAM files for some of the features while for the remaining features bigWig files are given.

ModHMM computes segmentations in several stages. At every stage the output is saved as a bigWig file, which can be inspected in a genome browser. The location and name of each bigWig file can be configured. A full set of all options is printed with `modhmm --genconf`.

To execute ModHMM simply run (assuming the configuration file is named `config.json`):
```sh
  modhmm -c config.json segmentation
```

### State Glossary

ModHMM states are not numbered but directly linked to known chromatin states. The following is a list of state abbreviations used in the resulting segmentation files

State name | Description
-----------|------------------------------------------
PA         | active promoter
EA         | active enhancer
BI         | bivalet region
PI         | primed region
EA:tr      | active enhancer in a transcribed region
BI:tr      | bivalet region in a transcribed region
PI:tr      | primed region in a transcribed region
TR         | transcribed region
TL         | low transcription
R1         | H3K27me3 repressed
R2         | H3K9me3 repressed
NS         | no signal
CL         | control signal

### Extracting Promoter and Enhancer Predictions

Genome segmentations are discretized predictions of chromatin states. They do not contain any information about the certainty of a particular prediction. Another drawback is that the number of predicted promoters and enhancers depends on the quality of the data, in particular the sequencing depth. Especially for differential analysis the dependency on the data quality might be hindering. In addition to genome segmentations, ModHMM can compute chromatin state probabilities:
```sh
  modhmm -c config.json eval-posterior-marginals
```
This command will compute genome-wide probabilities for all chromatin states and export them as bigWig files named `posterior-marginal-STATE.bw`. By default, probabilities are on *log-scale*. With the option `--std-scale` ModHMM exports probabilities on standard scale. In this case, bigWig files are named `posterior-marginal-exp-STATE.bw`.

ModHMM also implements a simple peak-finding algorithm that can be used to call high-probability regions in chromatin state probability tracks:
```sh
  modhmm -c config.json call-posterior-marginal-peaks --threshold=0.8
```
This command outputs tables with identified peaks, i.e. all regions with probabilities higher than the given threshold.

### Estimating Single-Feature Models

For detecting peaks, i.e. for separating signal from noise, ModHMM by default uses a single-feature mixture model that was estimated on either a mouse or human data set. It uses quantile-normalization to fit the provided data to the default model before evaluating genome-wide peak probabilities. This procedure allows to easily apply ModHMM to new data sets, but is less accurate than using a model that was estimated from the actual data at hand.

The following command estimates a new single-feature model for all features:
```sh
  modhmm -c config.json estimate-single-feature
```
It uses a default set of components for the mixture model of each feature and estimates the parameters from the observed coverages. The resulting estimates can be visualized using
```sh
  modhmm -c config.json plot-single-feature
```
If the mixture model of a feature, say H3K27ac, poorly separates signal from noise, it is possible to adapt the number of mixture components, i.e.
```sh
  modhmm -c config.json estimate-single-feature --force h3k27ac 1 2 1
```
This model would use a single delta distributions for all zero counts (first 1), two Poisson distributions (2), and a single geometric distribution (second 1). Which components are used to model the signal (foreground) is specified in the file `h3k27ac.components.json`. The result should look similar to the mixture model in the Figure below.

<img src="https://raw.githubusercontent.com/pbenner/modhmm/master/README_sf.png" alt="ModHMM SF" width="790" height="262" />

To obtain a list of mixture components including the estimated parameters use
```sh
  modhmm -c config.json print-single-feature
```

### Using ModHMM as a Peak Caller

Most peak callers use a single pre-defined model for computing enrichment probabilities and detecting peaks. In most cases there is a strong model misfit, because of the strong heterogeneity of ChIP-seq data. ModHMM instead allows to fit a mixture distribution (single-feature model) to the observed coverage values with a user-defined set of components. The following command calls ATAC-seq peaks using the estimated single-feature model, if available (see previous section):
```sh
  modhmm -c config.json call-single-feature-peaks atac
```

### Use Cases
#### Example 1: Compute segmentation on ENCODE data from mouse embyonic liver at day 12.5

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
    "Model Fallback"                  : "mm10",
    "Directory"                       : "mm10-liver-embryo-day12.5",
    "Description"                     : "liver embryo day12.5",
    "Threads"                         : 20,
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

#### Example 2: Estimate custom single-feature models on ENCODE data from mouse embyonic forebrain at day 11.5

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
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature h3k4me3o1 1 4 2
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature h3k9me3   2 4 1
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature rna       1 0 4
  modhmm -c mm10-forebrain-embryo-day11.5.json estimate-single-feature control   7 2 1

  echo '[3,4]'     > mm10-forebrain-embryo-day11.5:models/dnase.components.json
  echo '[8]'       > mm10-forebrain-embryo-day11.5:models/h3k27me3.components.json
  echo '[5,6,7,8]' > mm10-forebrain-embryo-day11.5:models/h3k4me1.components.json
  echo '[3,4]'     > mm10-forebrain-embryo-day11.5:models/h3k4me3.components.json
  echo '[6]'       > mm10-forebrain-embryo-day11.5:models/h3k4me3o1.components.json
  echo '[5,6]'     > mm10-forebrain-embryo-day11.5:models/h3k9me3.components.json
  echo '[2,3,4]'   > mm10-forebrain-embryo-day11.5:models/rna.components.json
  echo '[1,2]'     > mm10-forebrain-embryo-day11.5:models/rna-low.components.json
  echo '[9]'       > mm10-forebrain-embryo-day11.5:models/control.components.json
```
