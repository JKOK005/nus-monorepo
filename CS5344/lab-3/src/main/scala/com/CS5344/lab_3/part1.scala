package com.CS5344.lab_3

import org.apache.spark.SparkConf
import org.apache.spark.sql.SparkSession
import org.apache.spark.sql.functions.explode
import org.apache.spark.sql.functions._

object part1 extends App {
  def initSparkContext(applicationName: String): SparkSession = {
    val sparkConfig = new SparkConf().setAppName(applicationName)
    SparkSession.builder().config(sparkConfig).enableHiveSupport().getOrCreate()
  }

  val spark: SparkSession = initSparkContext("CS5344-lab-3-part-1")
  import spark.implicits._

  val dataDf      = spark.read.textFile("./src/main/resources/datafiles")
  val stopWordsDf = spark.read.textFile("./src/main/resources/stopwords.txt")

  val wordsDf = dataDf.filter(row => row.length != 0)
                      .map(
                        row => {
                          row.replaceAll("[\\p{Punct}&&[^'-]]+", " ") // Replaces all punctuations (except ') with space
                              .replaceAll("( )+", " ")                // Replaces all consecutive spaces with a single space
                              .toLowerCase.trim                                           // Cast to lower case and trim left / right space
                              .split(" ")                                         // Split at space
                        }
                      ).select(explode($"value"))                                         // Flatmap to single row

  val filteredDf = wordsDf.join(stopWordsDf, wordsDf("col") === stopWordsDf("value"), "leftanti") // Removes all stop words from main DF

  val aggregatedDf = filteredDf.groupBy("col").count()                              // Aggregate by word and count

  aggregatedDf.orderBy(desc("count"))
              .coalesce(1)
              .write
              .mode("Overwrite")
              .option("header","true")
              .csv("./src/main/resources/output/part1")                            // Order by count and write to resources/output
}
