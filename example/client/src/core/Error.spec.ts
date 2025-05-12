import { describe, expect, it } from 'vitest'
import ErrorPage from './Error'

describe('ErrorPage component', () => {
  it('should match the snapshot', () => {
    const result = ErrorPage({})
    expect(result).toMatchSnapshot()
  })
})
