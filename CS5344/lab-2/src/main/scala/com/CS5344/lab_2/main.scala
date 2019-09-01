package com.CS5344.lab_2

import org.apache.spark.SparkConf
import org.apache.spark.sql.SparkSession

object main extends App {
  def initSparkContext(applicationName: String): SparkSession = {
    val sparkConfig = new SparkConf().setAppName(applicationName)
    SparkSession.builder().config(sparkConfig).enableHiveSupport().getOrCreate()
  }

  val spark       = initSparkContext("CS5344-lab-2")
}
