export const timeDurationEnum = {
  minDuration: -9223372036854776000,
  maxDuration: 9223372036854776000,
  Nanosecond: 1,
  Microsecond: 1000,
  Millisecond: 1000000,
  Second: 1000000000,
  Minute: 60000000000,
  Hour: 3600000000000
} as const

export type TimeDurationEnum = (typeof timeDurationEnum)[keyof typeof timeDurationEnum]

export type TimeDuration = TimeDurationEnum