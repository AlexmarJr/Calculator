import type { Action, CalculatorState } from '../models/calculator.ts'
import { formatExpression, formatNumber } from './formatting.ts'

const keys: { label: string; name: string; action: Action; kind?: string }[] = [
  { label: 'AC', name: 'All clear', action: { type: 'clear' }, kind: 'utility' },
  { label: '+/−', name: 'Change sign', action: { type: 'sign' }, kind: 'utility' },
  { label: '%', name: 'Percentage', action: { type: 'percent' }, kind: 'utility' },
  { label: '÷', name: 'Divide', action: { type: 'operator', value: '/' }, kind: 'operator' },
  { label: '7', name: '7', action: { type: 'digit', value: '7' } },
  { label: '8', name: '8', action: { type: 'digit', value: '8' } },
  { label: '9', name: '9', action: { type: 'digit', value: '9' } },
  { label: '×', name: 'Multiply', action: { type: 'operator', value: '*' }, kind: 'operator' },
  { label: '4', name: '4', action: { type: 'digit', value: '4' } },
  { label: '5', name: '5', action: { type: 'digit', value: '5' } },
  { label: '6', name: '6', action: { type: 'digit', value: '6' } },
  { label: '−', name: 'Subtract', action: { type: 'operator', value: '-' }, kind: 'operator' },
  { label: '1', name: '1', action: { type: 'digit', value: '1' } },
  { label: '2', name: '2', action: { type: 'digit', value: '2' } },
  { label: '3', name: '3', action: { type: 'digit', value: '3' } },
  { label: '+', name: 'Add', action: { type: 'operator', value: '+' }, kind: 'operator' },
  { label: '0', name: '0', action: { type: 'digit', value: '0' } },
  { label: ',', name: 'Decimal separator', action: { type: 'decimal' } },
  { label: '⌫', name: 'Backspace', action: { type: 'backspace' }, kind: 'utility' },
  { label: '=', name: 'Equals', action: { type: 'equals' }, kind: 'equals' },
]

type CalculatorViewProps = {
  state: CalculatorState
  busy: boolean
  dispatch: (action: Action) => Promise<void>
}

export function CalculatorView({ state, busy, dispatch }: CalculatorViewProps) {
  const displayValue = formatNumber(state.display)
  const environment = import.meta.env.VITE_APP_ENV ?? (import.meta.env.DEV ? 'local' : 'production')

  return (
    <div className="app-shell">
      <header className="site-header">
        <a className="wordmark" href="./" aria-label="Calculator home">
          <span className="brand-icon" aria-hidden="true">
            <span />
            <span />
            <span />
            <span />
          </span>
          calculator<span className="wordmark-dot">.</span>
        </a>
        {environment === 'local' && (
          <a
            className="production-link"
            href="https://sezzle.137.131.234.28.sslip.io/"
            target="_blank"
            rel="noopener noreferrer"
          >
            APP in production
            <span aria-hidden="true">↗</span>
          </a>
        )}
        <span className="header-note">
          <span className="status-dot" /> A little more clarity
        </span>
      </header>

      <main>
        <div className="workspace">
          <section className="calculator" aria-label="Calculator">
            <div className="display">
              <div className="display-topline">
                <span>LET'S WORK IT OUT</span>
                <span className="display-indicator" aria-hidden="true" />
              </div>
              <div className="expression">{formatExpression(state.expression) || '\u00a0'}</div>
              <output
                className={
                  'result' +
                  (displayValue.length > 16
                    ? ' result-long'
                    : displayValue.length > 10
                      ? ' result-compact'
                      : '')
                }
                aria-live="polite"
                aria-atomic="true"
                aria-label="Result"
                aria-busy={busy}
              >
                {state.error ? 'Error' : displayValue}
              </output>
              <p
                className={'display-caption' + (state.error ? ' error-message' : '')}
                role={state.error ? 'alert' : undefined}
              >
                {state.error ||
                  (busy
                    ? 'Working…'
                    : state.justCalculated
                      ? 'Your answer is ready.'
                      : 'Good things start with a number.')}
              </p>
            </div>

            <div className="controls">
              <div className="science-keys">
                <button
                  type="button"
                  className={'science-key' + (state.operator === '^' ? ' selected' : '')}
                  title="Toggle exponentiation"
                  aria-label="Exponentiation"
                  aria-pressed={state.operator === '^'}
                  onClick={() => dispatch({ type: 'operator', value: '^' })}
                >
                  <span className="science-symbol">
                    x<sup>y</sup>
                  </span>
                  <span>Power</span>
                  <kbd>^</kbd>
                </button>
                <button
                  type="button"
                  className="science-key"
                  title="Calculate the square root of the current number"
                  aria-keyshortcuts="R"
                  aria-label="Square root"
                  onClick={() => dispatch({ type: 'sqrt' })}
                >
                  <span className="science-symbol">√</span>
                  <span>Square root</span>
                  <kbd>R</kbd>
                </button>
              </div>
              <div className="keypad">
                {keys.map((key) => (
                  <button
                    type="button"
                    key={key.name}
                    aria-label={key.name}
                    aria-pressed={
                      key.action.type === 'operator'
                        ? state.operator === key.action.value
                        : undefined
                    }
                    className={
                      'key ' +
                      (key.kind ?? 'number') +
                      (key.action.type === 'operator' && state.operator === key.action.value
                        ? ' selected'
                        : '')
                    }
                    onClick={() => dispatch(key.action)}
                  >
                    {key.label}
                  </button>
                ))}
              </div>
              <div className="calculator-footer">
                <span className="tiny-spark" aria-hidden="true">
                  ✦
                </span>{' '}
                A fresh perspective on the everyday.
              </div>
            </div>
          </section>

          <aside className="sidebar">
            <section className="history-panel" aria-labelledby="history-title">
              <div className="panel-heading">
                <div>
                  <p className="eyebrow">A PAPER TRAIL</p>
                  <h2 id="history-title">Your calculations</h2>
                </div>
                <button
                  type="button"
                  className="text-button"
                  disabled={!state.history.length}
                  onClick={() => dispatch({ type: 'clear-history' })}
                >
                  Clear
                </button>
              </div>
              <p className="panel-description">A little room to retrace your steps.</p>
              {state.history.length ? (
                <ol className="history-list">
                  {state.history.map((item) => (
                    <li key={item.id}>
                      <button
                        type="button"
                        className="history-entry"
                        onClick={() => dispatch({ type: 'restore', value: item })}
                        aria-label={
                          'Use result ' +
                          formatNumber(item.result) +
                          ' from ' +
                          formatExpression(item.expression)
                        }
                      >
                        <span className="history-expression">
                          {formatExpression(item.expression)} =
                        </span>
                        <span className="history-result">
                          {formatNumber(item.result)}
                          <span className="reuse-arrow" aria-hidden="true">
                            ↗
                          </span>
                        </span>
                      </button>
                    </li>
                  ))}
                </ol>
              ) : (
                <div className="history-empty">
                  <span className="empty-icon" aria-hidden="true">
                    <svg viewBox="0 0 48 48" fill="none" stroke="currentColor" strokeWidth="1.5">
                      <rect x="12" y="7" width="24" height="34" rx="4" />
                      <path d="M18 16h12M18 22h8M18 28h12M18 34h5" strokeLinecap="round" />
                    </svg>
                  </span>
                  <h3>A clean slate.</h3>
                  <p>
                    Your answers will land here.
                    <br />
                    Select one to use it again.
                  </p>
                </div>
              )}
              <p className="history-note">Your last 20 answers, kept for this session.</p>
            </section>

            <section className="tips-panel" aria-labelledby="tips-title">
              <span className="tips-star" aria-hidden="true">
                ✳
              </span>
              <h2 id="tips-title">A few handy things</h2>
              <p>
                Use <kbd>Enter</kbd> for the answer, <kbd>Esc</kbd> to reset, and <kbd>R</kbd> for
                square root.
              </p>
              <div className="tip-divider" />
              <p>
                Percentages understand the context.
                <br />
                <strong>200 + 10% = 220</strong>
              </p>
              <p className="order-note">Operations calculate in entry order.</p>
            </section>
          </aside>
        </div>
      </main>
    </div>
  )
}
