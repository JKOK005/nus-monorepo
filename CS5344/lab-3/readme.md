## CS 5344 Lab 3 submission

### Directories
Sources

| Description | Path |
| --- | --- |
| Source Data | ./src/main/resources/datafiles |
| Query text  | ./src/main/resources/query.txt |
| Stop words  | ./src/main/resources/stopwords.txt |

Results (after execution of programme)

| Description | Path |
| --- | --- |
| Part 1 | ./src/main/resources/output/part1/part-xxxx.csv |
| Part 2 | ./src/main/resources/output/part2/out.txt |

### Compile and execute
Clone this repository and cd to the respective project
```sbtshell
git clone https://github.com/JKOK005/nus-monorepo.git

cd CS5344/lab-3
```

Compile the jar and install all dependencies located in `build.sbt`

```sbtshell
sbt assembly
```

Execute the programme as a spark job with the following parameters

For Part 1:
```sbtshell
spark-submit --master local[2] --deploy-mode client --conf "spark.hadoop.mapreduce.fileoutputcommitter.algorithm.version=2" \
--driver-memory 1G --executor-cores 2 --conf spark.dynamicAllocation.maxExecutors=2 --executor-memory 1G \
--class com.CS5344.lab_3.part1 ./target/scala-2.11/lab-3-assembly-0.1.jar
```

For Part 2:
```sbtshell
spark-submit --master local[2] --deploy-mode client --conf "spark.hadoop.mapreduce.fileoutputcommitter.algorithm.version=2" \
--driver-memory 1G --executor-cores 2 --conf spark.dynamicAllocation.maxExecutors=2 --executor-memory 1G \
--class com.CS5344.lab_3.part2 ./target/scala-2.11/lab-3-assembly-0.1.jar
```