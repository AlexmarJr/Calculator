import type { SendAction } from '../repositories/calculatorRepository.ts'
import type { Action, CalculatorState } from '../models/calculator.ts'
import { initialState } from '../models/calculator.ts'

// Transport only: Go returns the complete state for every button action.
// Keep requests ordered; reset and history reuse cancel obsolete responses.
export function createCalculatorService(
  send: SendAction,
  publish: (state: CalculatorState) => void,
  publishBusy: (busy: boolean) => void = () => {},
) {
  let state = initialState
  let queue = Promise.resolve()
  let generation = 0
  let active: AbortController | null = null
  function cancel() {
    generation += 1
    active?.abort()
    active = null
    queue = Promise.resolve()
    publishBusy(false)
  }
  function dispatch(action: Action): Promise<void> {
    if (action.type === 'clear' || action.type === 'restore') cancel()
    const current = generation
    queue = queue.then(async () => {
      if (generation !== current) return
      active = new AbortController()
      publishBusy(true)
      try {
        const response = await send(state, action, active.signal)
        if (generation !== current) return
        state = response
        publish(state)
      } catch (error) {
        if (generation !== current) return
        state = {
          ...state,
          error: error instanceof Error ? error.message : 'Unable to process input.',
        }
        publish(state)
      } finally {
        if (generation === current) {
          active = null
          publishBusy(false)
        }
      }
    })
    return queue
  }
  return { dispatch, cancel }
}
