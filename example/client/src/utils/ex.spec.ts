import { describe, expect, it } from 'vitest'
import { add } from './ex'

describe('add', () => {
  it('should match the snapshot for various inputs', () => {
    const results = [
      { input: [2, 3], output: add(2, 3) },
      { input: [5, -3], output: add(5, -3) },
      { input: [-4, -6], output: add(-4, -6) },
      { input: [0, 0], output: add(0, 0) },
      { input: [0, 7], output: add(0, 7) },
      { input: [9, 0], output: add(9, 0) },
      { input: [100, 200], output: add(100, 200) },
      { input: [-50, 50], output: add(-50, 50) },
    ]
    expect(results).toMatchSnapshot()
  })
})
