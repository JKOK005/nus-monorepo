## CS 5344 Lab 1 submission

### Compile and execute
First, compile the jar and install all dependencies located in `build.sbt`

```scala
sbt assembly
```

Execute the programme as a spark job with the following parameters
```scala
spark-submit --master local[2] --deploy-mode client --conf "spark.hadoop.mapreduce.fileoutputcommitter.algorithm.version=2" \
--driver-memory 1G --executor-cores 2 --conf spark.dynamicAllocation.maxExecutors=5 --executor-memory 1G \
--class com.CS5344.lab_1.main ./target/scala-2.11/lab-1-assembly-0.1.jar
```

This would create the resultant file in `./output` directory