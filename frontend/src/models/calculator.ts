export type Operator = '+' | '-' | '*' | '/' | '^'
export type HistoryEntry = { id: number; expression: string; result: string }
export type CalculatorState = {
  display: string
  stored: string | null
  operator: Operator | null
  awaitingOperand: boolean
  justCalculated: boolean
  expression: string
  operandLabel: string | null
  error: string | null
  history: HistoryEntry[]
  nextId: number
}
export type Action =
  | { type: 'digit'; value: string }
  | { type: 'operator'; value: Operator }
  | { type: 'restore'; value: HistoryEntry }
  | {
      type:
        'decimal' | 'equals' | 'clear' | 'backspace' | 'sign' | 'percent' | 'sqrt' | 'clear-history'
    }

// Initial display values only. Go processes every input and operation.
export const initialState: CalculatorState = {
  display: '0',
  stored: null,
  operator: null,
  awaitingOperand: false,
  justCalculated: false,
  expression: '',
  operandLabel: null,
  error: null,
  history: [],
  nextId: 1,
}
