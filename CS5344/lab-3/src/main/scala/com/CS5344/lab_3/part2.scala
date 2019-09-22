package com.CS5344.lab_3

import org.apache.spark.SparkConf
import org.apache.spark.sql.{DataFrame, SparkSession}
import org.apache.spark.sql.functions._
import commons.Utils
import java.io.File
import scala.collection.mutable.ListBuffer

object part2 extends App {
  def initSparkContext(applicationName: String): SparkSession = {
    val sparkConfig = new SparkConf().setAppName(applicationName)
    SparkSession.builder().config(sparkConfig).enableHiveSupport().getOrCreate()
  }

  def getTF_IDF(df: DataFrame, dictionaryDf: DataFrame, wordCount: Long, numDoc: Int): DataFrame = {
    df.join(dictionaryDf, df("col") === dictionaryDf("col") ,"inner")
      .withColumn("tf_idf", (df("count") / wordCount) * (log(lit(numDoc) / dictionaryDf("count"))))
  }

  implicit val spark: SparkSession = initSparkContext("CS5344-lab-3-part-2")

  val dataFileDir       = "./src/main/resources/datafiles"
  val stopWordsFilePath = "./src/main/resources/stopwords.txt"

  val stopWordsDf = spark.read.textFile(stopWordsFilePath)
                                .toDF("words")
                                .withColumn("words", lower(col("words")))
                                .cache

  val singleDocDf   = new ListBuffer[DataFrame]()   // To track information for a single document
  val allDocDf      = new ListBuffer[DataFrame]()   // To track aggregated information across all documents
  val dataFilePtr   = new File(dataFileDir)

  dataFilePtr.listFiles.foreach(
    filePath => {
      val dataDf              = spark.read.textFile(filePath.toString).toDF("words")
      val singleWordDf        = Utils.dfToWordsVector(dataDf)
      val removeStopWordsDf   = Utils.dfRemoveStopWords(singleWordDf, stopWordsDf)

      singleDocDf.append(removeStopWordsDf.groupBy("col").count)
      allDocDf.append(removeStopWordsDf.withColumn("path", lit(filePath.toString)))
    }
  )

  // Compute dictionary of all words and the number of documents that contains them
  val dictionaryDf = allDocDf.reduce(_.union(_))
                            .distinct
                            .groupBy("col")
                            .count.cache

  // Computes TF-IDF value of for each document


}
