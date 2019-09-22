package com.CS5344.lab_3

import org.apache.spark.SparkConf
import org.apache.spark.sql.SparkSession
import org.apache.spark.sql.functions._
import commons.Utils

object part1 extends App {
  def initSparkContext(applicationName: String): SparkSession = {
    val sparkConfig = new SparkConf().setAppName(applicationName)
    SparkSession.builder().config(sparkConfig).enableHiveSupport().getOrCreate()
  }

  implicit val spark: SparkSession = initSparkContext("CS5344-lab-3-part-1")

  val dataDf      = spark.read.textFile("./src/main/resources/datafiles").toDF("words")
  val stopWordsDf = spark.read.textFile("./src/main/resources/stopwords.txt").toDF("words")

  val wordsDf     = Utils.dfToWordsVector(dataDf)
  val filteredDf  = Utils.dfRemoveStopWords(wordsDf, stopWordsDf)

  val aggregatedDf  = filteredDf.groupBy("col").count() // Aggregate by word and count
  aggregatedDf.orderBy(desc("count"))
              .coalesce(1)
              .write
              .mode("Overwrite")
              .option("header","true")
              .csv("./src/main/resources/output/part1") // Order by count and write to resources/output
}