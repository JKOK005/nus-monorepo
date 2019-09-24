package com.CS5344.lab_3

import org.apache.spark.SparkConf
import org.apache.spark.sql.{DataFrame, SparkSession}
import org.apache.spark.sql.functions._
import commons.Utils
import java.io.{File, PrintWriter}
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
  val outputFilePath    = "./src/main/resources/output/part2/out.txt"

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

  /**
    * Compute dictionary of all words and the number of documents that contains them
    * All documents are read individually. Stop words and punctuations are removed
    * We first take the distinct collection of words for each document
    * We then union all documents together and generate the document count of word W, groupBy word W
    */
  val dictionaryDf = allDocDf.reduce(_.union(_))
                            .distinct
                            .groupBy("col")
                            .count
                            .withColumnRenamed("col", "dictionary_col")
                            .withColumnRenamed("count", "dictionary_count")
                            .cache

  /**
    * TF-IDF is defined as the produce of term frequency and IDF value
    * term frequency (TF) of word W : (count of W in document) / (total words in document)
    * IDF for word W : log((size of entire document set) / (number of documents containing W))
    *
    * TF-IDF = TF * IDF
    */
  val TF_IDF = singleDocDf.map(
                            eachDf => {
                              eachDf.cache
                              val totalWordCount  = eachDf.agg(sum("count")).first.getLong(0)
                              getTF_IDF(eachDf, dictionaryDf, totalWordCount, dataFilePtr.listFiles.length)
                            }
                          )

  /**
    * Normalization of all TF-IDF value for a single document
    * For a given word W:
    *   Compute S = sum(All [TF-IDF]^2 of words)
    *   normalized TF-IDF = TF-IDF / sqrt(S)
    */
  val TF_IDF_norm = TF_IDF.map(
                            eachDf => {
                              eachDf.cache
                              val totalWordSquared = eachDf.agg(sum(pow(col("tf_idf"), 2))).first.getDouble(0)
                              eachDf.withColumn("tf_idf_norm", col("tf_idf") / lit(math.sqrt(totalWordSquared)))
                            }
                          )

  /**
    * Compute relevance vector for a list of query words Q
    * This is essentially filtering the normalized TF_IDF dataframe for all rows containing words in Q
    * We then sum up all the normalized TF_IDF values to get the relevance for the document
    */
  val relevanceVector = TF_IDF_norm.map(
                            eachDf => {
                              val filtered = eachDf.filter(col("col").isin(List("wonderful","story") :_*))
                                                    .select("tf_idf_norm")

                              if (filtered.count > 0) filtered.agg(sum("tf_idf_norm")).first.getDouble(0)
                              else 0
                            }
                        )

  /**
    * Write results to file
    */
  val zipped = relevanceVector zip dataFilePtr.listFiles
  val file   = new File(outputFilePath)
  file.getParentFile.mkdirs
  file.createNewFile()
  val writer = new PrintWriter(file)
  zipped.foreach { i => writer.write(i.toString); writer.write("\n") }
  writer.close
}
