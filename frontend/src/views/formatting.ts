// Format text instead of converting to a number: preserve typed decimals and precision.
export function formatNumber(value: string): string {
  const [mantissa, exponent] = value.split(/[eE]/)
  const [integer, fraction] = mantissa.split('.')
  const grouped = integer.replace(/\B(?=(\d{3})+(?!\d))/g, '.')
  return (
    grouped +
    (fraction === undefined ? '' : ',' + fraction) +
    (exponent === undefined ? '' : 'e' + exponent)
  )
}

export function formatExpression(expression: string): string {
  return expression.replace(/-?\d+(?:\.\d*)?(?:e[+-]?\d+)?/gi, formatNumber)
}
