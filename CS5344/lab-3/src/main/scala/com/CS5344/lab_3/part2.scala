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

  def getTF_IDF(df: DataFrame, dictionaryDf: DataFrame, wordCount: Long, numDoc: Long): DataFrame = {
    df.join(dictionaryDf, df("col") === dictionaryDf("dictionary_col") ,"inner")
      .withColumn("tf_idf", (df("count") / wordCount) * (log(lit(numDoc) / dictionaryDf("dictionary_count"))))
  }

  implicit val spark: SparkSession = initSparkContext("CS5344-lab-3-part-2")

  val dataFileDir       = "./src/main/resources/datafiles"
  val stopWordsFilePath = "./src/main/resources/stopwords.txt"
  val queryTextPath     = "./src/main/resources/query.txt"

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
                            .count
                            .withColumnRenamed("col", "dictionary_col")
                            .withColumnRenamed("count", "dictionary_count")
                            .cache

  // Computes TF-IDF value of for each document
  val TF_IDF = singleDocDf.map(
                            eachDf => {
                              eachDf.cache
                              val totalWordCount  = eachDf.agg(sum("count")).first.getLong(0)
                              getTF_IDF(eachDf, dictionaryDf, totalWordCount, dataFilePtr.listFiles.length)
                            }
                          )

  // Normalize TF_IDF
  val TF_IDF_norm = TF_IDF.map(
                            eachDf => {
                              eachDf.cache
                              val totalWordSquared = eachDf.agg(sum(pow(col("tf_idf"), 2))).first.getDouble(0)
                              eachDf.withColumn("tf_idf_norm", col("tf_idf") / lit(math.sqrt(totalWordSquared)))
                            }
                          )

  // Compute relevance vector
  val relevanceVector = TF_IDF_norm.map(
                            eachDf => {
                              val filtered = eachDf.filter(col("col").isin(List("wonderful","story") :_*))
                                                    .select("tf_idf_norm")

                              if (filtered.count > 0) filtered.agg(sum("tf_idf_norm")).first.getDouble(0)
                              else 0
                            }
                        )

  // Output to file
  val zipped = relevanceVector zip dataFilePtr.listFiles
  zipped.foreach { println }
}
