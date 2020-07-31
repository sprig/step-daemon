package com.colingodsey.stepd

object PageFormat {
  case object SP_4x4D_128 extends PageFormat {
    val Directional = true
    val BytesPerChunk = 256
    val BytesPerSegment = 2
    val StepsPerSegment = 8
    val MaxStepsPerSegment = 7 //segment is 8 ticks, but has max of 7

    val SegmentsPerChunk = BytesPerChunk / BytesPerSegment
    val StepsPerChunk = SegmentsPerChunk * StepsPerSegment
  }

  case object SP_4x2_256 extends PageFormat {
    val Directional = false
    val BytesPerChunk = 256
    val StepsPerChunk = 1024
    val StepsPerSegment = 4
    val MaxStepsPerSegment = 3
  }

  case object SP_4x1_512 extends PageFormat {
    val Directional = false
    val BytesPerChunk = 256
    val StepsPerChunk = 512
    val StepsPerSegment = 1
    val MaxStepsPerSegment = 1
  }
}
trait PageFormat {
  val BytesPerChunk: Int
  val StepsPerSegment: Int
  val MaxStepsPerSegment: Int
  val StepsPerChunk: Int
  val Directional: Boolean
}
