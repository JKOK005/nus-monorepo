package com.CS5344.lab_3

import com.CS5344.lab_3.part1.initSparkContext
import org.apache.spark.SparkConf
import org.apache.spark.sql.SparkSession

object part2 extends App {
  def initSparkContext(applicationName: String): SparkSession = {
    val sparkConfig = new SparkConf().setAppName(applicationName)
    SparkSession.builder().config(sparkConfig).enableHiveSupport().getOrCreate()
  }

  implicit val spark: SparkSession = initSparkContext("CS5344-lab-3-part-1")
}
