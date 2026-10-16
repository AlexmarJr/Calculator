import { test } from 'node:test'
import assert from 'node:assert/strict'
import { actionForKey } from '../../src/controllers/keyboard.ts'

test('shortcuts map to API actions, including both decimal separators', () => {
  assert.deepEqual(actionForKey('7'), { type: 'digit', value: '7' })
  assert.deepEqual(actionForKey('^'), { type: 'operator', value: '^' })
  assert.deepEqual(actionForKey('.'), { type: 'decimal' })
  assert.deepEqual(actionForKey(','), { type: 'decimal' })
  assert.deepEqual(actionForKey('r'), { type: 'sqrt' })
  assert.deepEqual(actionForKey('R'), { type: 'sqrt' })
  assert.deepEqual(actionForKey('%'), { type: 'percent' })
  assert.deepEqual(actionForKey('Backspace'), { type: 'backspace' })
  assert.deepEqual(actionForKey('Escape'), { type: 'clear' })
  assert.deepEqual(actionForKey('Delete'), { type: 'clear' })
  assert.equal(actionForKey('Tab'), undefined)
})

test('Enter calculates unless the focused control uses native activation', () => {
  assert.deepEqual(actionForKey('Enter'), { type: 'equals' })
  assert.equal(actionForKey('Enter', true), undefined)
  assert.deepEqual(actionForKey('=', true), { type: 'equals' })
})
