package com.CS5344.lab_3

import org.apache.spark.SparkConf
import org.apache.spark.sql.{DataFrame, SparkSession}
import org.apache.spark.sql.functions._
import commons.Utils
import java.io.File

object part2 extends App {
  def initSparkContext(applicationName: String): SparkSession = {
    val sparkConfig = new SparkConf().setAppName(applicationName)
    SparkSession.builder().config(sparkConfig).enableHiveSupport().getOrCreate()
  }

  implicit val spark: SparkSession = initSparkContext("CS5344-lab-3-part-2")

  val dataFileDir       = "./src/main/resources/datafiles"
  val stopWordsFilePath = "./src/main/resources/stopwords.txt"

  val stopWordsDf       = spark.read.textFile(stopWordsFilePath)
                                .toDF("words")
                                .withColumn("words", lower(col("words")))
                                .cache
}
