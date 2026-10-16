import { test } from 'node:test'
import assert from 'node:assert/strict'
import { sendAction } from '../../src/repositories/calculatorRepository.ts'
import { initialState } from '../../src/models/calculator.ts'

test('REST client sends the complete state and button action, then returns server state unchanged', async () => {
  let sent: RequestInit | undefined
  const backendState = {
    ...initialState,
    display: '42',
    expression: 'Server expression',
    justCalculated: true,
  }
  const received = await sendAction(
    initialState,
    { type: 'equals' },
    undefined,
    async (url, init) => {
      assert.equal(url, '/api/calculator/action')
      sent = init
      return Response.json({ state: backendState })
    },
  )
  assert.equal(sent?.method, 'POST')
  assert.deepEqual(JSON.parse(sent!.body as string), {
    state: initialState,
    action: { type: 'equals' },
  })
  assert.deepEqual(received, backendState)
})
test('repository reports backend, invalid-response, and network errors', async () => {
  await assert.rejects(
    sendAction(initialState, { type: 'digit', value: 'x' }, undefined, async () =>
      Response.json({ error: 'Invalid input' }, { status: 400 }),
    ),
    /Invalid input/,
  )
  await assert.rejects(
    sendAction(initialState, { type: 'equals' }, undefined, async () =>
      Response.json({ state: { display: 42 } }),
    ),
    /invalid response/,
  )
  await assert.rejects(
    sendAction(initialState, { type: 'equals' }, undefined, async () => {
      throw new TypeError('Network down')
    }),
    /Cannot reach/,
  )
})
