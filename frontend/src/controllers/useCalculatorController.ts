import { useEffect, useState } from 'react'
import { sendAction } from '../repositories/calculatorRepository.ts'
import { createCalculatorService } from '../services/calculatorService.ts'
import { initialState } from '../models/calculator.ts'
import { actionForKey } from './keyboard.ts'

export function useCalculatorController() {
  const [state, setState] = useState(initialState)
  const [busy, setBusy] = useState(false)
  const [service] = useState(() => createCalculatorService(sendAction, setState, setBusy))
  useEffect(() => () => service.cancel(), [service])

  useEffect(() => {
    const handleKey = (event: KeyboardEvent) => {
      const target = event.target as HTMLElement | null
      if (
        event.ctrlKey ||
        event.metaKey ||
        event.altKey ||
        target?.isContentEditable ||
        target?.matches('input, textarea, select')
      )
        return
      const action = actionForKey(
        event.key,
        Boolean(target?.closest('.history-entry, .text-button, a')),
      )
      if (action) {
        event.preventDefault()
        void service.dispatch(action)
      }
    }
    window.addEventListener('keydown', handleKey)
    return () => window.removeEventListener('keydown', handleKey)
  }, [service])

  return { state, busy, dispatch: service.dispatch }
}
