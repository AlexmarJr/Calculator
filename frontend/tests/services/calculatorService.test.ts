import { test } from 'node:test'
import assert from 'node:assert/strict'
import { createCalculatorService } from '../../src/services/calculatorService.ts'
import { initialState } from '../../src/models/calculator.ts'
import type { Action, CalculatorState } from '../../src/models/calculator.ts'

function deferred<T>() {
  let resolve!: (value: T) => void
  const promise = new Promise<T>((done) => {
    resolve = done
  })
  return { promise, resolve }
}
test('service forwards every button action without processing it', async () => {
  const sent: Action[] = []
  const backendState = {
    ...initialState,
    display: '987',
    expression: 'Backend-selected expression',
  }
  let latest = initialState
  const controller = createCalculatorService(
    async (_state, action) => {
      sent.push(action)
      return backendState
    },
    (state) => {
      latest = state
    },
  )
  const actions: Action[] = [
    { type: 'digit', value: '2' },
    { type: 'decimal' },
    { type: 'backspace' },
    { type: 'sign' },
    { type: 'operator', value: '+' },
    { type: 'operator', value: '^' },
    { type: 'sqrt' },
    { type: 'percent' },
    { type: 'equals' },
    { type: 'clear-history' },
  ]
  for (const action of actions) await controller.dispatch(action)
  assert.deepEqual(sent, actions)
  assert.deepEqual(latest, backendState)
})
test('queued requests use each previous server response as their state', async () => {
  const started = deferred<void>(),
    response = deferred<CalculatorState>()
  const calls: { state: CalculatorState; action: Action }[] = []
  const firstResponse = { ...initialState, display: '123' }
  const lastResponse = { ...initialState, display: '456' }
  let latest = initialState
  const controller = createCalculatorService(
    async (state, action) => {
      calls.push({ state, action })
      if (calls.length === 1) {
        started.resolve()
        return response.promise
      }
      return lastResponse
    },
    (state) => {
      latest = state
    },
  )
  const first = controller.dispatch({ type: 'digit', value: '9' })
  const second = controller.dispatch({ type: 'operator', value: '^' })
  await started.promise
  assert.equal(calls.length, 1)
  response.resolve(firstResponse)
  await Promise.all([first, second])
  assert.deepEqual(calls[1].state, firstResponse)
  assert.deepEqual(latest, lastResponse)
})
test('clear cancels obsolete responses and is itself processed by the backend', async () => {
  const started = deferred<void>(),
    response = deferred<CalculatorState>()
  let latest = initialState
  let requestSignal: AbortSignal | undefined
  const controller = createCalculatorService(
    async (_state, action, signal) => {
      if (action.type === 'equals') {
        requestSignal = signal
        started.resolve()
        return response.promise
      }
      return action.type === 'clear' ? initialState : { ...initialState, display: '9' }
    },
    (state) => {
      latest = state
    },
  )
  const old = controller.dispatch({ type: 'equals' })
  await started.promise
  await controller.dispatch({ type: 'clear' })
  await controller.dispatch({ type: 'digit', value: '9' })
  assert.equal(requestSignal?.aborted, true)
  response.resolve({ ...initialState, display: '5' })
  await old
  assert.equal(latest.display, '9')
})
test('service preserves the display when a request fails', async () => {
  let latest = initialState
  const service = createCalculatorService(
    async () => {
      throw new Error('Cannot reach the calculator service.')
    },
    (state) => {
      latest = state
    },
  )
  await service.dispatch({ type: 'equals' })
  assert.equal(latest.display, '0')
  assert.match(latest.error!, /Cannot reach/)
})
