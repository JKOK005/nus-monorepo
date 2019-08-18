package com.CS5344.lab_1

import org.apache.spark.SparkConf
import org.apache.spark.sql.SparkSession

object main extends App {
  def initSparkContext(applicationName: String): SparkSession = {
    val sparkConfig = new SparkConf().setAppName(applicationName)
    SparkSession.builder().config(sparkConfig).enableHiveSupport().getOrCreate()
  }

  val spark       = initSparkContext("CS5344-lab-1")
  val wordsRdd    = spark.sparkContext.textFile("./src/main/resources/word_count.txt")
  val finalRdd    = wordsRdd.map(_.trim)
                            .map(each_line => each_line.replaceAll("[\\.$|,]", ""))
                            .flatMap(each_line => each_line.split(" "))
                            .map(each_word => (each_word, 1))
                            .reduceByKey(_+_)

  finalRdd.coalesce(1).saveAsTextFile("./output")
}
