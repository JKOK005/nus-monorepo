package com.CS5344.lab_3.commons

import org.apache.spark.sql.{DataFrame, SparkSession}
import org.apache.spark.sql.functions.explode

object Utils {

  def dfToWordsVector(df: DataFrame)(implicit spark: SparkSession): DataFrame = {
    import spark.implicits._

    df.filter(row => !(row.mkString("").isEmpty && row.length>0))
      .map(
        row => {
          row.toString
            .replaceAll("[\\p{Punct}&&[^'-]]+", " ")    // Replaces all punctuations (except ') with space
            .replaceAll("( )+", " ")                    // Replaces all consecutive spaces with a single space
            .toLowerCase.trim                                               // Cast to lower case and trim left / right space
            .split(" ")                                             // Split at space
        }
      ).select(explode($"value"))                                           // Flatmap to single row
  }

  def dfRemoveStopWords(df: DataFrame, stopWordsDf: DataFrame): DataFrame = {
    // Removes all stop words from given df
    df.join(stopWordsDf, df("col") === stopWordsDf("words"), "leftanti")
  }
}
