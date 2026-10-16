import type { Action, Operator } from '../models/calculator.ts'

// Map shortcuts to actions only; Go interprets the action.
export function actionForKey(key: string, nativeEnter = false): Action | undefined {
  if (/^[0-9]$/.test(key)) return { type: 'digit', value: key }
  if (['+', '-', '*', '/', '^'].includes(key)) return { type: 'operator', value: key as Operator }
  if (key === '.' || key === ',') return { type: 'decimal' }
  if (key === '%') return { type: 'percent' }
  if (key.toLowerCase() === 'r') return { type: 'sqrt' }
  if (key === 'Backspace') return { type: 'backspace' }
  if (key === 'Escape' || key === 'Delete') return { type: 'clear' }
  if (key === '=' || (key === 'Enter' && !nativeEnter)) return { type: 'equals' }
}
