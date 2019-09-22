package com.CS5344.lab_3

import org.apache.spark.SparkConf
import org.apache.spark.sql.{DataFrame, SparkSession}
import org.apache.spark.sql.functions._
import commons.Utils
import java.io.File
import scala.collection.mutable.ListBuffer

object part1 extends App {
  def initSparkContext(applicationName: String): SparkSession = {
    val sparkConfig = new SparkConf().setAppName(applicationName)
    SparkSession.builder().config(sparkConfig).enableHiveSupport().getOrCreate()
  }

  def getMin = udf((arr: Seq[Long]) => arr.min)

  implicit val spark: SparkSession = initSparkContext("CS5344-lab-3-part-1")

  val dataFileDir       = "./src/main/resources/datafiles"
  val stopWordsFilePath = "./src/main/resources/stopwords.txt"

  val stopWordsDf   = spark.read.textFile(stopWordsFilePath)
                          .toDF("words")
                          .withColumn("words", lower(col("words")))
                          .cache

  val allAggregatedDfLst  = new ListBuffer[DataFrame]()
  val dataFilePtr   = new File(dataFileDir)

  dataFilePtr.listFiles.foreach(
    filePath => {
      val dataDf              = spark.read.textFile(filePath.toString).toDF("words")
      val singleWordDf        = Utils.dfToWordsVector(dataDf)
      val removeStopWordsDf   = Utils.dfRemoveStopWords(singleWordDf, stopWordsDf)
      val aggregatedDf        = removeStopWordsDf.groupBy("col").count.cache  // Aggregate by word and count
      allAggregatedDfLst.append(aggregatedDf)
    }
  )

  val finalDf = allAggregatedDfLst.reduce(_.union(_))
                                  .groupBy("col")
                                  .agg(collect_list("count") as "counts")
                                  .filter(size(col("counts")) === dataFilePtr.listFiles.length)
                                  .withColumn("min_count", getMin(col("counts")))
                                  .orderBy(desc("min_count")).select("col", "min_count")

  finalDf.coalesce(1)
          .write
          .mode("Overwrite")
          .option("header","true")
          .csv("./src/main/resources/output/part1") // Order by count and write to resources/output
}