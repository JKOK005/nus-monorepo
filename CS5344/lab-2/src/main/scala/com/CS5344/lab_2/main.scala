package com.CS5344.lab_2

import org.apache.spark.SparkConf
import org.apache.spark.sql.SparkSession

object main extends App {
  def initSparkContext(applicationName: String): SparkSession = {
    val sparkConfig = new SparkConf().setAppName(applicationName)
    SparkSession.builder().config(sparkConfig).enableHiveSupport().getOrCreate()
  }
  val spark       = initSparkContext("CS5344-lab-2")

  val reviewsRdd  = spark.read.json("./src/main/resources/reviews_Musical_Instruments_5.json")
                                  .select("asin").rdd

  val parseReviewsRdd = reviewsRdd.filter(row => row(0) != null)
                                  .map(row => (row(0).toString, 1)).reduceByKey(_ + _)

  val metadataRdd     = spark.read.json("./src/main/resources/meta_Musical_Instruments.json")
                                  .select("asin","price").rdd

  val metadataPairedRdd = metadataRdd.filter(row => row(0) != null)
                                      .map(row => (row(0).toString, row(1)))

  // Perform join here
  val mergedRdd   = parseReviewsRdd.join(metadataPairedRdd)
                                    .sortBy(x => x._2._1.toInt, ascending = false)

  // Take top 10 results
  println("Top 10 products based on review counts")
  mergedRdd.take(10).foreach{result => println(s"Product ID: ${result._1}, Reviews: ${result._2._1}, Price: ${result._2._2}")}
}
