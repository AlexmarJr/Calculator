import type { Action, CalculatorState } from '../models/calculator.ts'

export type SendAction = (
  state: CalculatorState,
  action: Action,
  signal?: AbortSignal,
) => Promise<CalculatorState>

function isState(value: unknown): value is CalculatorState {
  if (!value || typeof value !== 'object') return false
  const state = value as Record<string, unknown>
  const optionalText = (item: unknown) => item === null || typeof item === 'string'
  return (
    typeof state.display === 'string' &&
    typeof state.expression === 'string' &&
    optionalText(state.stored) &&
    optionalText(state.operator) &&
    optionalText(state.operandLabel) &&
    optionalText(state.error) &&
    typeof state.awaitingOperand === 'boolean' &&
    typeof state.justCalculated === 'boolean' &&
    typeof state.nextId === 'number' &&
    Array.isArray(state.history) &&
    state.history.every(
      (item) =>
        item &&
        typeof item === 'object' &&
        typeof item.id === 'number' &&
        typeof item.expression === 'string' &&
        typeof item.result === 'string',
    )
  )
}

export async function sendAction(
  state: CalculatorState,
  action: Action,
  signal?: AbortSignal,
  fetcher: typeof fetch = fetch,
): Promise<CalculatorState> {
  const controller = new AbortController()
  const cancel = () => controller.abort()
  signal?.addEventListener('abort', cancel, { once: true })
  if (signal?.aborted) controller.abort()
  let timedOut = false
  const timeout = setTimeout(() => {
    timedOut = true
    controller.abort()
  }, 10000)
  try {
    const response = await fetcher('/api/calculator/action', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ state, action }),
      signal: controller.signal,
    })
    const data = (await response.json()) as { state?: unknown; error?: unknown }
    if (!response.ok)
      throw new Error(
        typeof data?.error === 'string' ? data.error : 'Unable to process input. Try again.',
      )
    if (!isState(data?.state))
      throw new Error('The calculator service returned an invalid response.')
    return data.state
  } catch (error) {
    if (signal?.aborted) throw error
    if (timedOut) throw new Error('The calculator service timed out. Try again.')
    if (error instanceof TypeError)
      throw new Error('Cannot reach the calculator service. Try again.')
    if (error instanceof SyntaxError)
      throw new Error('The calculator service returned an invalid response.')
    throw error
  } finally {
    clearTimeout(timeout)
    signal?.removeEventListener('abort', cancel)
  }
}
