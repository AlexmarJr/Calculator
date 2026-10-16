import { test } from 'node:test'
import assert from 'node:assert/strict'
import { formatExpression, formatNumber } from '../../src/views/formatting.ts'

test('Brazilian formatting preserves exact server-provided text', () => {
  for (const [raw, masked] of [
    ['0', '0'],
    ['1234', '1.234'],
    ['-1234567.50', '-1.234.567,50'],
    ['0.', '0,'],
    ['-0.05', '-0,05'],
    ['123456789012345', '123.456.789.012.345'],
    ['1.25e+21', '1,25e+21'],
    ['1e-7', '1e-7'],
  ])
    assert.equal(formatNumber(raw), masked)
  assert.equal(formatExpression('1234.50 + 10000.25 ='), '1.234,50 + 10.000,25 =')
})
